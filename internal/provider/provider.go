package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ provider.Provider = &cs3Provider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &cs3Provider{
			version: version,
		}
	}
}

// cs3Provider is the provider implementation.
type cs3Provider struct {
	version string
}

// cs3ProviderModel maps provider schema data to a Go type.
type cs3ProviderModel struct {
	Region    types.String `tfsdk:"region"`
	AccessKey types.String `tfsdk:"access_key"`
	SecretKey types.String `tfsdk:"secret_key"`
}

// Metadata returns the provider type name.
func (p *cs3Provider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "customs3"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *cs3Provider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Optional: true,
			},
			"access_key": schema.StringAttribute{
				Optional: true,
			},
			"secret_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *cs3Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	// Retrieve provider data from configuration
	var config cs3ProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Region.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("region"),
			"Unknown Region",
			"The provider cannot create the Custom S3 client as there is an unknown configuration value for the AWS Region. ",
		)
	}
	if config.AccessKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("access_key"),
			"Unknown Access Key value",
			"The provider cannot create the Custom S3 client as there is an unknown configuration value for the AWS Access Key. ",
		)
	}
	if config.SecretKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("secret_key"),
			"Unknown Secret Key value",
			"The provider cannot create the Custom S3 client as there is an unknown configuration value for the AWS Secret Key. ",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	region := os.Getenv("AWS_REGION")
	access_key := os.Getenv("AWS_ACCESS_KEY_ID")
	secret_key := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if !config.Region.IsNull() {
		region = config.Region.ValueString()
	}

	if !config.AccessKey.IsNull() {
		access_key = config.AccessKey.ValueString()
	}

	if !config.SecretKey.IsNull() {
		secret_key = config.SecretKey.ValueString()
	}

	if region == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("region"),
			"Missing Region",
			"The provider cannot create the AWS client as there is a missing or empty value for the Region. ",
		)
	}

	if access_key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("access_key"),
			"Missing Access Key",
			"The provider cannot create the AWS client as there is a missing or empty value for the Access Key. ",
		)
	}

	if secret_key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("secret_key"),
			"Missing Secret Key",
			"The provider cannot create the AWS client as there is a missing or empty value for the Secret Key. ",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}
	// Create AWS client
	client, err := session.NewSession(&aws.Config{
		Region:      aws.String(region), // Specify the AWS region
		Credentials: credentials.NewStaticCredentials(access_key, secret_key, ""),
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *cs3Provider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewBucketDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *cs3Provider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewOrderResource,
	}
}
