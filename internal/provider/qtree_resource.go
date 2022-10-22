package ontap

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ontap "github.com/ybizeul/terraform-provider-ontap/ontap_client_go"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &QtreeResource{}
var _ resource.ResourceWithImportState = &QtreeResource{}

func NewQtreeResource() resource.Resource {
	return &QtreeResource{}
}

// QtreeResource defines the resource implementation.
type QtreeResource struct {
	client *ontap.Client
}

// QtreeResourceModel describes the resource data model.
type QtreeResourceModel struct {
	Id types.String `tfsdk:"uuid"`

	SVMUUID    types.String `tfsdk:"svm_uuid"`
	VolumeUUID types.String `tfsdk:"volume_uuid"`

	Name           types.String `tfsdk:"name"`
	QtreeID        types.Int64  `tfsdk:"id"`
	Path           types.String `tfsdk:"path"`
	SecurityStyle  types.String `tfsdk:"security_style"`
	UnixPermission types.Int64  `tfsdk:"unix_permissions"`
}

func (r *QtreeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_qtree"
}

func (r *QtreeResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example resource",

		Attributes: map[string]tfsdk.Attribute{
			"uuid": {
				MarkdownDescription: "Qtree UUID, which is <VolumeUUID>/<QtreeID>",
				Optional:            true,
				Computed:            true,
				Type:                types.StringType,
			},
			"id": {
				MarkdownDescription: "Example identifier",
				Type:                types.Int64Type,
				Computed:            true,
			},
			"svm_uuid": {
				MarkdownDescription: "Example identifier",
				Type:                types.StringType,
				Required:            true,
			},
			"volume_uuid": {
				MarkdownDescription: "Example identifier",
				Type:                types.StringType,
				Required:            true,
			},
			"name": {
				MarkdownDescription: "Example identifier",
				Type:                types.StringType,
				Required:            true,
			},
			"path": {
				MarkdownDescription: "Example identifier",
				Type:                types.StringType,
				Computed:            true,
			},
			"security_style": {
				MarkdownDescription: "Example identifier",
				Type:                types.StringType,
				Optional:            true,
			},
			"unix_permissions": {
				MarkdownDescription: "Example identifier",
				Type:                types.Int64Type,
				Optional:            true,
				Computed:            true,
			},
		},
	}, nil
}

func (r *QtreeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ontap.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *QtreeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *QtreeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := d.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	//     return
	// }

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	qtree := ontap.Qtree{}
	qtree.Name = data.Name.Value
	qtree.SVMUUID = data.SVMUUID.Value
	qtree.VolumeUUID = data.VolumeUUID.Value
	qtree.SecurityStyle = data.SecurityStyle.Value
	qtree.UnixPermission = data.UnixPermission.Value

	created_qtree, err := r.client.CreateQtree(&qtree)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create qtree, got error: %s", err))
		return
	}
	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a QTREE +%v", map[string]interface{}{
		"uuid": qtree.UUID,
	})

	data.Name = types.String{Value: created_qtree.Name}
	data.Id = types.String{Value: created_qtree.UUID}
	data.SVMUUID = types.String{Value: created_qtree.SVMUUID}
	data.VolumeUUID = types.String{Value: created_qtree.VolumeUUID}
	data.SecurityStyle = types.String{Value: created_qtree.SecurityStyle}
	data.UnixPermission = types.Int64{Value: created_qtree.UnixPermission}
	data.Path = types.String{Value: created_qtree.Path}
	data.QtreeID = types.Int64{Value: created_qtree.Id}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QtreeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *QtreeResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	qtree, err := r.client.GetQtree(data.Id.Value, "")
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read qtree, got error: %s", err))
		return
	}
	data.Id = types.String{Value: qtree.UUID}
	data.QtreeID = types.Int64{Value: qtree.Id}
	data.Name = types.String{Value: qtree.Name}
	data.VolumeUUID = types.String{Value: qtree.VolumeUUID}
	data.Path = types.String{Value: qtree.Path}
	data.SecurityStyle = types.String{Value: qtree.SecurityStyle}
	data.UnixPermission = types.Int64{Value: qtree.UnixPermission}
	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := d.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *QtreeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *QtreeResourceModel
	var state *QtreeResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := d.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }
	qtree := ontap.Qtree{}
	qtree.UUID = state.Id.Value
	qtree.Name = plan.Name.Value
	qtree.SecurityStyle = plan.SecurityStyle.Value
	qtree.UnixPermission = plan.UnixPermission.Value

	// tflog.Trace(ctx, "creating a QTREE +%v", map[string]interface{}{
	// 	"data":   data,
	// 	"qtree":  qtree,
	// 	"__uuid": data.Id.Value,
	// })

	updated_qtree, err := r.client.UpdateQtree(&qtree)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update qtree, got error: %s", err))
		return
	}
	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	plan.Id = types.String{Value: updated_qtree.UUID}
	plan.Name = types.String{Value: updated_qtree.Name}
	plan.SVMUUID = types.String{Value: updated_qtree.SVMUUID}
	plan.VolumeUUID = types.String{Value: updated_qtree.VolumeUUID}
	plan.SecurityStyle = types.String{Value: updated_qtree.SecurityStyle}
	plan.UnixPermission = types.Int64{Value: updated_qtree.UnixPermission}
	plan.Path = types.String{Value: updated_qtree.Path}
	plan.QtreeID = types.Int64{Value: updated_qtree.Id}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *QtreeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *QtreeResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	qtree := ontap.Qtree{}
	qtree.UUID = data.Id.Value
	qtree.VolumeUUID = data.VolumeUUID.Value

	r.client.DeleteQtree(&qtree)

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := d.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *QtreeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
