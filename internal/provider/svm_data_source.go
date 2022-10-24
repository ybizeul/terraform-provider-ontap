package ontap

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ontap "github.com/ybizeul/terraform-provider-ontap/ontap_client_go"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &SVMDataSource{}

func NewSVMDataSource() datasource.DataSource {
	return &SVMDataSource{}
}

// ExampleDataSource defines the data source implementation.
type SVMDataSource struct {
	client *ontap.Client
}

// ExampleDataSourceModel describes the data source data model.
type SVMDataSourceModel struct {
	UUID                types.String               `tfsdk:"uuid"`
	Aggregates          []AggregateDataSourceModel `tfsdk:"aggregates"`
	AggregatesDelegated types.Bool                 `tfsdk:"aggregates_delegated"`
	Certificate         types.String               `tfsdk:"certificate"`
	CIFS                *CIFSDataSourceModel       `tfsdk:"cifs"`
	Comment             types.String               `tfsdk:"comment"`
	DNS                 types.Map                  `tfsdk:"dns"`
	FCInterfaces        types.List                 `tfsdk:"fc_interfaces"`
	FCP                 types.Bool                 `tfsdk:"fcp"`
	IPInterfaces        types.List                 `tfsdk:"ip_interfaces"`
	IPSpace             types.Map                  `tfsdk:"ipspace"`
	ISCSI               types.Bool                 `tfsdk:"iscsi"`
	Language            types.String               `tfsdk:"language"`
	LDAP                types.Map                  `tfsdk:"ldap"`
	Name                types.String               `tfsdk:"name"`
	NFS                 types.Bool                 `tfsdk:"nfs"`
	NIS                 types.Map                  `tfsdk:"nis"`
	NVME                types.Bool                 `tfsdk:"nvme"`
	NSSwitch            types.Map                  `tfsdk:"nsswitch"`
	Routes              types.List                 `tfsdk:"routes"`
	S3                  types.Map                  `tfsdk:"s3"`
	Snapmirror          types.Map                  `tfsdk:"snapmirror"`
	SnapshotPolicy      types.Map                  `tfsdk:"snapshot_policy"`
	State               types.String               `tfsdk:"state"`
	Subtype             types.String               `tfsdk:"subtype"`
}

type AggregateDataSourceModel struct {
	Name types.String `tfsdk:"name"`
	UUID types.String `tfsdk:"uuid"`
}

type CIFSDataSourceModel struct {
	ADDomain *ADDomainDataSourceModel `tfsdk:"ad_domain"`
	Enabled  types.Bool               `tfsdk:"enabled"`
	Name     *types.String            `tfsdk:"name"`
}

func NewCIFSDataSourceModel() CIFSDataSourceModel {
	return CIFSDataSourceModel{
		ADDomain: nil,
		Enabled:  types.Bool{Null: true},
		Name:     &types.String{Null: true},
	}
}

type ADDomainDataSourceModel struct {
	FQDN               types.String `tfsdk:"fqdn"`
	OrganizationalUnit types.String `tfsdk:"organizational_unit"`
}

func (d *SVMDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_svm"
}

func (d *SVMDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
			// This description is used by the documentation generator and the language server.
			MarkdownDescription: "Example data source",
			Attributes: map[string]tfsdk.Attribute{
				"uuid": {
					Type:     types.StringType,
					Required: true,
				},
				"aggregates": {
					Computed: true,
					Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
						"name": {
							Type:     types.StringType,
							Required: true,
						},
						"uuid": {
							Type:     types.StringType,
							Required: true,
						},
					}),
				},
				"aggregates_delegated": {
					Type:     types.BoolType,
					Optional: true,
				},
				"certificate": {
					Type:     types.StringType,
					Optional: true,
				},
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
				// "cifs": {
				// 	Optional: true,
				// 	Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
				// 		"ad_domain": {
				// 			Required: true,
				// 			Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
				// 				"fqdn": {
				// 					Type:     types.StringType,
				// 					Required: true,
				// 				},
				// 				"organizational_unit": {
				// 					Type:     types.StringType,
				// 					Required: true,
				// 				},
				// 			}),
				// 		},
				// 		"enabled": {
				// 			Type:     types.BoolType,
				// 			Required: true,
				// 		},
				// 		"name": {
				// 			Type:     types.StringType,
				// 			Required: true,
				// 		},
				// 	}),
				// },
				"comment": {
					Type:     types.StringType,
					Optional: true,
				},
				"dns": {
					Optional: true,
					Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
						"domains": {
							Required: true,
							Type: types.ListType{
								ElemType: types.StringType,
							},
						},
						"servers": {
							Required: true,
							Type: types.ListType{
								ElemType: types.StringType,
							},
						},
					}),
				},
				"fc_interfaces": {
					Optional: true,
					Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
						"data_protocol": {
							Type:     types.StringType,
							Required: true,
						},
						"location": {
							Required: true,
							Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
								"port": {
									Required: true,
									Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
										"name": {
											Type:     types.StringType,
											Required: true,
										},
										"node": {
											Type:     types.StringType,
											Required: true,
										},
										"uuid": {
											Type:     types.StringType,
											Required: true,
										},
									}),
								},
							}),
						},
						"name": {
							Type:     types.StringType,
							Required: true,
						},
						"uuid": {
							Type:     types.StringType,
							Required: true,
						},
					}),
				},
				"fcp": {
					Type:     types.BoolType,
					Optional: true,
				},
				"ip_interfaces": {
					Optional: true,
					Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
						"ip": {
							Required: true,
							Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
								"address": {
									Type:     types.StringType,
									Required: true,
								},
								"netmask": {
									Type:     types.StringType,
									Required: true,
								},
							}),
						},
						"location": {
							Required: true,
							Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
								"broadcast_domain": {
									Required: true,
									Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
										"name": {
											Type:     types.StringType,
											Required: true,
										},
										"uuid": {
											Type:     types.StringType,
											Required: true,
										},
									}),
								},
								"home_node": {
									Required: true,
									Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
										"name": {
											Type:     types.StringType,
											Required: true,
										},
										"uuid": {
											Type:     types.StringType,
											Required: true,
										},
									}),
								},
							}),
						},
						"name": {
							Type:     types.StringType,
							Optional: true,
						},
						"service_policy": {
							Type:     types.StringType,
							Optional: true,
						},
						"services": {
							Type: types.ListType{
								ElemType: types.StringType,
							},
							Optional: true,
						},
						"uuid": {
							Type:     types.StringType,
							Required: true,
						},
					}),
				},
				"ipspace": {
					Optional: true,
					Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
						"name": {
							Type:     types.StringType,
							Required: true,
						},
						"uuid": {
							Type:     types.StringType,
							Required: true,
						},
					}),
				},
				"iscsi": {
					Type:     types.BoolType,
					Optional: true,
				},
				"language": {
					Type:     types.StringType,
					Optional: true,
				},
				"ldap": {
					Optional: true,
					Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
						"ad_domain": {
							Type:     types.StringType,
							Required: true,
						},
						"base_dn": {
							Type:     types.StringType,
							Required: true,
						},
						"bind_dn": {
							Type:     types.StringType,
							Required: true,
						},
						"enabled": {
							Type:     types.BoolType,
							Required: true,
						},
						"servers": {
							Type: types.ListType{
								ElemType: types.StringType,
							},
							Required: true,
						},
					}),
				},
				"name": {
					Type:     types.StringType,
					Optional: true,
				},
				"nfs": {
					Type:     types.BoolType,
					Optional: true,
				},
				"nis": {
					Optional: true,
					Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
						"ad_domain": {
							Type:     types.StringType,
							Required: true,
						},
						"enabled": {
							Type:     types.BoolType,
							Required: true,
						},
						"servers": {
							Type: types.ListType{
								ElemType: types.StringType,
							},
							Required: true,
						},
					}),
				},
				"nvme": {
					Type:     types.BoolType,
					Optional: true,
				},
				"nsswitch": {
					Optional: true,
					Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
						"group": {
							Type: types.ListType{
								ElemType: types.StringType,
							},
							Required: true,
						},
						"hosts": {
							Type: types.ListType{
								ElemType: types.StringType,
							},
							Required: true,
						},
						"namemap": {
							Type: types.ListType{
								ElemType: types.StringType,
							},
							Required: true,
						},
						"netgroup": {
							Type: types.ListType{
								ElemType: types.StringType,
							},
							Required: true,
						},
						"passwd": {
							Type: types.ListType{
								ElemType: types.StringType,
							},
							Required: true,
						},
					}),
				},
				"routes": {
					Optional: true,
					Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
						"destination": {
							Required: true,
							Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
								"address": {
									Type:     types.StringType,
									Required: true,
								},
								"family": {
									Type:     types.StringType,
									Required: true,
								},
								"netmask": {
									Type:     types.StringType,
									Required: true,
								},
							}),
						},
						"gateway": {
							Type:     types.StringType,
							Required: true,
						},
					}),
				},
				"s3": {
					Optional: true,
					Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
						"enabled": {
							Type:     types.BoolType,
							Required: true,
						},
						"name": {
							Type:     types.StringType,
							Required: true,
						},
					}),
				},
				"snapmirror": {
					Optional: true,
					Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
						"is_protected": {
							Type:     types.BoolType,
							Required: true,
						},
						"protected_volumes_count": {
							Type:     types.Int64Type,
							Required: true,
						},
					}),
				},
				"snapshot_policy": {
					Optional: true,
					Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
						"name": {
							Type:     types.StringType,
							Required: true,
						},
						"uuid": {
							Type:     types.StringType,
							Required: true,
						},
					}),
				},
				"state": {
					Type:     types.StringType,
					Optional: true,
				},
				"subtype": {
					Type:     types.StringType,
					Optional: true,
				},
			},
		},
		nil
}

func (d *SVMDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SVMDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SVMDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	SVM, err := d.client.GetSVM(data.UUID.Value)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read SVM, got error: %s", err))
		return
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.

	data.UUID = types.String{Value: SVM.UUID}

	data.Aggregates = []AggregateDataSourceModel{}
	for _, aggr := range SVM.Aggregates {
		data.Aggregates = append(data.Aggregates, AggregateDataSourceModel{
			Name: types.String{Value: aggr.Name},
			UUID: types.String{Value: aggr.UUID},
		})
	}

	if SVM.CIFS != nil {
		cifs := NewCIFSDataSourceModel()

		if SVM.CIFS.Name != nil {
			cifs.Name = &types.String{Value: *SVM.CIFS.Name}
		}

		if SVM.CIFS.ADDomain != nil {
			cifs.ADDomain = &ADDomainDataSourceModel{
				FQDN:               types.String{Value: SVM.CIFS.ADDomain.FQDN},
				OrganizationalUnit: types.String{Value: SVM.CIFS.ADDomain.OrganizationalUnit},
			}
		}

		data.CIFS = &cifs
	}

	// This code commented implements CIFS settings with types.Object but the syntax
	// of repeated AttrTypes didn't look like a good pattern.
	// Instead, we replaced :
	// CIFS                types.Object       `tfsdk:"cifs"`
	//
	// by
	//
	// CIFS                *CIFSDataSourceModel       `tfsdk:"cifs"`
	//
	// in the model
	//
	// if SVM.CIFS != nil {
	// 	cifs := types.Object{
	// 		Attrs: map[string]attr.Value{
	// 			"ad_domain": types.Object{
	// 				Attrs: map[string]attr.Value{
	// 					"fqdn":                types.String{Value: SVM.CIFS.ADDomain.FQDN},
	// 					"organizational_unit": types.String{Value: SVM.CIFS.ADDomain.OrganizationalUnit},
	// 				},
	// 				AttrTypes: map[string]attr.Type{
	// 					"fqdn":                types.StringType,
	// 					"organizational_unit": types.StringType,
	// 				},
	// 			},
	// 			"enabled": types.Bool{Value: SVM.CIFS.Enabled},
	// 			"name":    types.String{Value: SVM.CIFS.Name},
	// 		},
	// 		AttrTypes: map[string]attr.Type{
	// 			"ad_domain": types.ObjectType{
	// 				AttrTypes: map[string]attr.Type{
	// 					"fqdn":                types.StringType,
	// 					"organizational_unit": types.StringType,
	// 				},
	// 			},
	// 			"enabled": types.BoolType,
	// 			"name":    types.StringType,
	// 		},
	// 	}
	// 	data.CIFS = cifs
	// }

	data.AggregatesDelegated = types.Bool{Value: SVM.AggregatesDelegated}

	data.Name = types.String{Value: SVM.Name}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
