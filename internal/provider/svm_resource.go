package ontap

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ontap "github.com/ybizeul/terraform-provider-ontap/ontap_client_go"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &SVMResource{}
var _ resource.ResourceWithImportState = &SVMResource{}

func NewSVMResource() resource.Resource {
	return &SVMResource{}
}

// SVMResource defines the resource implementation.
type SVMResource struct {
	client *ontap.Client
}

// SVMResourceModel describes the resource data model.
type SVMResourceModel struct {
	UUID types.String `tfsdk:"uuid"`
	// Aggregates          []AggregateResourceModel     `tfsdk:"aggregates"`
	// AggregatesDelegated types.Bool                     `tfsdk:"aggregates_delegated"`
	// Certificate         types.String                   `tfsdk:"certificate"`
	CIFS    *CIFSResourceModel `tfsdk:"cifs"`
	Comment types.String       `tfsdk:"comment"`
	DNS     *DNSResourceModel  `tfsdk:"dns"`
	// FCInterfaces        []FCInterfaceResourceModel   `tfsdk:"fc_interfaces"`
	// FCP                 types.Bool                     `tfsdk:"fcp"`
	IPInterfaces []IPInterfaceResourceModel `tfsdk:"ip_interfaces"`
	IPSpace      *IPSpaceResourceModel      `tfsdk:"ipspace"`
	// ISCSI               types.Bool                     `tfsdk:"iscsi"`
	// Language            types.String                   `tfsdk:"language"`
	// LDAP                *LDAPResourceModel           `tfsdk:"ldap"`
	Name types.String `tfsdk:"name"`
	NFS  types.Bool   `tfsdk:"nfs"`
	// NIS                 *NISResourceModel            `tfsdk:"nis"`
	// NVME                types.Bool                     `tfsdk:"nvme"`
	// NSSwitch            *NSSwitchResourceModel       `tfsdk:"nsswitch"`
	Routes []RouteResourceModel `tfsdk:"routes"`
	// S3                  *S3ResourceModel             `tfsdk:"s3"`
	// Snapmirror          *SnapmirrorResourceModel     `tfsdk:"snapmirror"`
	// SnapshotPolicy      *SnapshotPolicyResourceModel `tfsdk:"snapshot_policy"`
	// State               types.String                   `tfsdk:"state"`
	Subtype types.String `tfsdk:"subtype"`
}

/*
****************************

	cifs

*****************************
*/

func NewCIFSResourceModel() CIFSResourceModel {
	return CIFSResourceModel{
		ADDomain: nil,
		Enabled:  types.Bool{Null: true},
		Name:     &types.String{Null: true},
	}
}

type CIFSResourceModel struct {
	ADDomain *ADDomainResourceModel `tfsdk:"ad_domain"`
	Enabled  types.Bool             `tfsdk:"enabled"`
	Name     *types.String          `tfsdk:"name"`
}

type ADDomainResourceModel struct {
	FQDN               types.String `tfsdk:"fqdn"`
	OrganizationalUnit types.String `tfsdk:"organizational_unit"`
}

type DNSResourceModel struct {
	Domains []types.String `tfsdk:"domains"`
	Servers []types.String `tfsdk:"servers"`
}

/*
****************************

	ip_interfaces

*****************************
*/

func NewIPInterfaceResourceModel() IPInterfaceResourceModel {
	return IPInterfaceResourceModel{
		IP: IPInterfaceIPResourceModel{
			Address: types.String{Null: true},
			Netmask: types.String{Null: true},
		},
		Name:          types.String{Null: true},
		ServicePolicy: types.String{Null: true},
		UUID:          types.String{Null: true},
	}
}

type IPInterfaceResourceModel struct {
	IP            IPInterfaceIPResourceModel `tfsdk:"ip"`
	Name          types.String               `tfsdk:"name"`
	ServicePolicy types.String               `tfsdk:"service_policy"`
	Services      []types.String             `tfsdk:"services"`
	UUID          types.String               `tfsdk:"uuid"`
}

type IPInterfaceIPResourceModel struct {
	Address types.String `tfsdk:"address"`
	Netmask types.String `tfsdk:"netmask"`
}

type IPSpaceResourceModel struct {
	Name types.String `tfsdk:"name"`
	UUID types.String `tfsdk:"uuid"`
}

/*
****************************

	routes

*****************************
*/
type RouteResourceModel struct {
	Destination RouteDestinationResourceModel `tfsdk:"destination"`
	Gateway     types.String                  `tfsdk:"gateway"`
}
type RouteDestinationResourceModel struct {
	Address types.String `tfsdk:"address"`
	Family  types.String `tfsdk:"family"`
	Netmask types.String `tfsdk:"netmask"`
}

func (r *SVMResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_svm"
}

func (r *SVMResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "An ONTAP SVM",
		Attributes: map[string]tfsdk.Attribute{
			"uuid": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			// "aggregates": {
			// 	Computed: true,
			// 	Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
			// 		"name": {
			// 			Type:     types.StringType,
			// 			Required: true,
			// 		},
			// 		"uuid": {
			// 			Type:     types.StringType,
			// 			Required: true,
			// 		},
			// 	}),
			// },
			// "aggregates_delegated": {
			// 	Type:     types.BoolType,
			// 	Optional: true,
			// },
			// "certificate": {
			// 	Type:     types.StringType,
			// 	Optional: true,
			// },
			"cifs": {
				Optional: true,
				Type: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"ad_domain": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"fqdn":                types.StringType,
								"organizational_unit": types.StringType,
							},
						},
						"enabled": types.BoolType,
						"name":    types.StringType,
					},
				},
			},
			"comment": {
				Type:     types.StringType,
				Optional: true,
			},
			"dns": {
				Optional: true,
				Type: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"domains": types.ListType{ElemType: types.StringType},
						"servers": types.ListType{ElemType: types.StringType},
					},
				},
			},
			// "fc_interfaces": {
			// 	Optional: true,
			// 	Type: types.ListType{
			// 		ElemType: types.ObjectType{
			// 			AttrTypes: map[string]attr.Type{
			// 				"data_protocol": types.StringType,
			// 				"name":          types.StringType,
			// 				"uuid":          types.StringType,
			// 			},
			// 		},
			// 	},
			// },
			// "fcp": {
			// 	Type:     types.BoolType,
			// 	Optional: true,
			// },
			"ip_interfaces": {
				Optional: true,
				Type: types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"ip": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"address": types.StringType,
									"netmask": types.StringType,
								},
							},
							"name":           types.StringType,
							"service_policy": types.StringType,
							"services": types.ListType{
								ElemType: types.StringType,
							},
							"uuid": types.StringType,
						},
					},
				},
			},
			"ipspace": {
				Optional: true,
				Type: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"name": types.StringType,
						"uuid": types.StringType,
					},
				},
			},
			// "iscsi": {
			// 	Type:     types.BoolType,
			// 	Optional: true,
			// },
			// "language": {
			// 	Type:     types.StringType,
			// 	Optional: true,
			// },
			// "ldap": {
			// 	Optional: true,
			// 	Type: types.ObjectType{
			// 		AttrTypes: map[string]attr.Type{
			// 			"ad_domain": types.StringType,
			// 			"base_dn":   types.StringType,
			// 			"bind_dn":   types.StringType,
			// 			"enabled":   types.BoolType,
			// 			"servers": types.ListType{
			// 				ElemType: types.StringType,
			// 			},
			// 		},
			// 	},
			// },
			"name": {
				Type:     types.StringType,
				Optional: true,
			},
			"nfs": {
				Type:     types.BoolType,
				Optional: true,
			},
			// "nis": {
			// 	Optional: true,
			// 	Type: types.ObjectType{
			// 		AttrTypes: map[string]attr.Type{
			// 			"domain":  types.StringType,
			// 			"enabled": types.BoolType,
			// 			"servers": types.ListType{
			// 				ElemType: types.StringType,
			// 			},
			// 		},
			// 	},
			// },
			// "nvme": {
			// 	Type:     types.BoolType,
			// 	Optional: true,
			// },
			// "nsswitch": {
			// 	Optional: true,
			// 	Type: types.ObjectType{
			// 		AttrTypes: map[string]attr.Type{
			// 			"group": types.ListType{
			// 				ElemType: types.StringType,
			// 			},
			// 			"hosts": types.ListType{
			// 				ElemType: types.StringType,
			// 			},
			// 			"namemap": types.ListType{
			// 				ElemType: types.StringType,
			// 			},
			// 			"netgroup": types.ListType{
			// 				ElemType: types.StringType,
			// 			},
			// 			"passwd": types.ListType{
			// 				ElemType: types.StringType,
			// 			},
			// 		},
			// 	},
			// },
			"routes": {
				Optional: true,
				Type: types.ListType{
					ElemType: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"destination": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"address": types.StringType,
									"family":  types.StringType,
									"netmask": types.StringType,
								},
							},
							"gateway": types.StringType,
						},
					},
				},
			},
			// "s3": {
			// 	Optional: true,
			// 	Type: types.ObjectType{
			// 		AttrTypes: map[string]attr.Type{
			// 			"enabled": types.BoolType,
			// 			"name":    types.StringType,
			// 		},
			// 	},
			// },
			// "snapmirror": {
			// 	Optional: true,
			// 	Type: types.ObjectType{
			// 		AttrTypes: map[string]attr.Type{
			// 			"is_protected":            types.BoolType,
			// 			"protected_volumes_count": types.Int64Type,
			// 		},
			// 	},
			// },
			// "snapshot_policy": {
			// 	Optional: true,
			// 	Type: types.ObjectType{
			// 		AttrTypes: map[string]attr.Type{
			// 			"name": types.StringType,
			// 			"uuid": types.StringType,
			// 		},
			// 	},
			// },
			// "state": {
			// 	Type:     types.StringType,
			// 	Optional: true,
			// },
			"subtype": {
				Type:     types.StringType,
				Optional: true,
			},
		},
	}, nil
}

func (r *SVMResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SVMResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *SVMResourceModel

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

	svm := ontap.SVM{}

	svm.Name = data.Name.Value

	if data.CIFS != nil {
		cifs := ontap.SVMCIFS{}
		cifs.ADDomain = &ontap.ADDomain{
			FQDN:               data.CIFS.ADDomain.FQDN.Value,
			OrganizationalUnit: data.CIFS.ADDomain.OrganizationalUnit.Value,
		}
		cifs.Enabled = data.CIFS.Enabled.Value
		cifs.Name = &data.CIFS.Name.Value
	}

	svm.Comment = data.Comment.Value

	if data.DNS != nil {
		dns := ontap.SVMDNS{}
		dns.Domains = []string{}
		for _, d := range data.DNS.Domains {
			dns.Domains = append(dns.Domains, d.Value)
		}
		dns.Servers = []string{}
		for _, d := range data.DNS.Servers {
			dns.Servers = append(dns.Servers, d.Value)
		}
		svm.DNS = &dns
	}

	created_svm, err := r.client.CreateSVM(&svm)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create SVM, got error: %s", err))
		return
	}
	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log

	// Save data into Terraform state

	data.UUID = types.String{Value: string(*created_svm.UUID)}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SVMResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *SVMResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	SVM, err := r.client.GetSVM(&data.UUID.Value, nil)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read SVM, got error: %s", err))
		return
	}
	data.UUID = types.String{Value: *SVM.UUID}
	data.Name = types.String{Value: SVM.Name}

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

func (r *SVMResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *SVMResourceModel
	var state *SVMResourceModel

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
	svm := ontap.SVM{}
	svm.UUID = &state.UUID.Value
	svm.Name = plan.Name.Value

	// tflog.Trace(ctx, "creating a SVM +%v", map[string]interface{}{
	// 	"data":   data,
	// 	"SVM":  SVM,
	// 	"__uuid": data.Id.Value,
	// })

	updated_SVM, err := r.client.UpdateSVM(&svm)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update SVM, got error: %s", err))
		return
	}
	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	plan.UUID = types.String{Value: *updated_SVM.UUID}
	plan.Name = types.String{Value: updated_SVM.Name}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SVMResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *SVMResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	svm := ontap.SVM{}
	svm.UUID = &data.UUID.Value

	r.client.DeleteSVM(&svm)

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := d.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *SVMResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
