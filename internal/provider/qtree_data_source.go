package ontap

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ontap "github.com/ybizeul/terraform-provider-ontap/ontap_client_go"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &QtreeDataSource{}

func NewQtreeDataSource() datasource.DataSource {
	return &QtreeDataSource{}
}

// ExampleDataSource defines the data source implementation.
type QtreeDataSource struct {
	client *ontap.Client
}

// ExampleDataSourceModel describes the data source data model.
type QtreeDataSourceModel struct {
	UUID types.String `tfsdk:"uuid"`

	VolumeUUID types.String `tfsdk:"volume_uuid"`

	Name           types.String `tfsdk:"name"`
	Id             types.Int64  `tfsdk:"id"`
	Path           types.String `tfsdk:"path"`
	SecurityStyle  types.String `tfsdk:"security_style"`
	UnixPermission types.Int64  `tfsdk:"unix_permissions"`
}

func (d *QtreeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_qtree"
}

func (d *QtreeDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example data source",

		Attributes: map[string]tfsdk.Attribute{
			"uuid": {
				MarkdownDescription: "Qtree UUID, which is <VolumeUUID>/<QtreeID>",
				Optional:            true,
				Type:                types.StringType,
			},
			"id": {
				MarkdownDescription: "Example identifier",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"volume_uuid": {
				MarkdownDescription: "Example identifier",
				Type:                types.StringType,
				Computed:            true,
			},
			"name": {
				MarkdownDescription: "Example identifier",
				Type:                types.StringType,
				Computed:            true,
			},
			"path": {
				MarkdownDescription: "Example identifier",
				Type:                types.StringType,
				Computed:            true,
			},
			"security_style": {
				MarkdownDescription: "Example identifier",
				Type:                types.StringType,
				Computed:            true,
			},
			"unix_permissions": {
				MarkdownDescription: "Example identifier",
				Type:                types.Int64Type,
				Computed:            true,
			},
		},
	}, nil
}

func (d *QtreeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ontap.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *QtreeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data QtreeDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	qtree, err := d.client.GetQtree(data.UUID.Value, "")
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
		return
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.UUID = types.String{Value: qtree.UUID}
	data.Id = types.Int64{Value: int64(qtree.Id)}
	data.Name = types.String{Value: qtree.Name}
	data.VolumeUUID = types.String{Value: qtree.VolumeUUID}
	data.Path = types.String{Value: qtree.Path}
	data.SecurityStyle = types.String{Value: qtree.SecurityStyle}
	data.UnixPermission = types.Int64{Value: int64(qtree.UnixPermission)}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
