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
	UUID                types.String                   `tfsdk:"uuid"`
	Aggregates          []AggregateDataSourceModel     `tfsdk:"aggregates"`
	AggregatesDelegated types.Bool                     `tfsdk:"aggregates_delegated"`
	Certificate         types.String                   `tfsdk:"certificate"`
	CIFS                *CIFSDataSourceModel           `tfsdk:"cifs"`
	Comment             types.String                   `tfsdk:"comment"`
	DNS                 *DNSDataSourceModel            `tfsdk:"dns"`
	FCInterfaces        []FCInterfaceDataSourceModel   `tfsdk:"fc_interfaces"`
	FCP                 types.Bool                     `tfsdk:"fcp"`
	IPInterfaces        []IPInterfaceDataSourceModel   `tfsdk:"ip_interfaces"`
	IPSpace             *IPSpaceDataModel              `tfsdk:"ipspace"`
	ISCSI               types.Bool                     `tfsdk:"iscsi"`
	Language            types.String                   `tfsdk:"language"`
	LDAP                *LDAPDataSourceModel           `tfsdk:"ldap"`
	Name                types.String                   `tfsdk:"name"`
	NFS                 types.Bool                     `tfsdk:"nfs"`
	NIS                 *NISDataSourceModel            `tfsdk:"nis"`
	NVME                types.Bool                     `tfsdk:"nvme"`
	NSSwitch            *NSSwitchDataSourceModel       `tfsdk:"nsswitch"`
	Routes              []RouteDataSourceModel         `tfsdk:"routes"`
	S3                  *S3DataSourceModel             `tfsdk:"s3"`
	Snapmirror          *SnapmirrorDataSourceModel     `tfsdk:"snapmirror"`
	SnapshotPolicy      *SnapshotPolicyDataSourceModel `tfsdk:"snapshot_policy"`
	State               types.String                   `tfsdk:"state"`
	Subtype             types.String                   `tfsdk:"subtype"`
}

/*
****************************

	aggregates

*****************************
*/
type AggregateDataSourceModel struct {
	Name types.String `tfsdk:"name"`
	UUID types.String `tfsdk:"uuid"`
}

/*
****************************

	cifs

*****************************
*/

func NewCIFSDataSourceModel() CIFSDataSourceModel {
	return CIFSDataSourceModel{
		ADDomain: nil,
		Enabled:  types.Bool{Null: true},
		Name:     &types.String{Null: true},
	}
}

type CIFSDataSourceModel struct {
	ADDomain *ADDomainDataSourceModel `tfsdk:"ad_domain"`
	Enabled  types.Bool               `tfsdk:"enabled"`
	Name     *types.String            `tfsdk:"name"`
}

type ADDomainDataSourceModel struct {
	FQDN               types.String `tfsdk:"fqdn"`
	OrganizationalUnit types.String `tfsdk:"organizational_unit"`
}

type DNSDataSourceModel struct {
	Domains []types.String `tfsdk:"domains"`
	Servers []types.String `tfsdk:"servers"`
}

/*
****************************

	fc_interfaces

*****************************
*/
func NewFCInterfaceDataSourceModel() FCInterfaceDataSourceModel {
	return FCInterfaceDataSourceModel{
		Name: types.String{Null: true},
		UUID: types.String{Null: true},
	}
}

type FCInterfaceDataSourceModel struct {
	DataProtocal types.String `tfsdk:"data_protocol"`
	Name         types.String `tfsdk:"name"`
	UUID         types.String `tfsdk:"uuid"`
}

/*
****************************

	ip_interfaces

*****************************
*/

func NewIPInterfaceDataSourceModel() IPInterfaceDataSourceModel {
	return IPInterfaceDataSourceModel{
		IP: IPInterfaceIPDataSourceModel{
			Address: types.String{Null: true},
			Netmask: types.String{Null: true},
		},
		Name:          types.String{Null: true},
		ServicePolicy: types.String{Null: true},
		UUID:          types.String{Null: true},
	}
}

type IPInterfaceDataSourceModel struct {
	IP            IPInterfaceIPDataSourceModel `tfsdk:"ip"`
	Name          types.String                 `tfsdk:"name"`
	ServicePolicy types.String                 `tfsdk:"service_policy"`
	Services      []types.String               `tfsdk:"services"`
	UUID          types.String                 `tfsdk:"uuid"`
}

type IPInterfaceIPDataSourceModel struct {
	Address types.String `tfsdk:"address"`
	Netmask types.String `tfsdk:"netmask"`
}

type IPSpaceDataModel struct {
	Name types.String `tfsdk:"name"`
	UUID types.String `tfsdk:"uuid"`
}

/*
****************************

	ldap

*****************************
*/
func NewLDAPDataSourceModel() LDAPDataSourceModel {
	return LDAPDataSourceModel{
		ADDomain: types.String{Null: true},
		BaseDN:   types.String{Null: true},
		BindDN:   types.String{Null: true},
		Servers:  nil,
	}
}

type LDAPDataSourceModel struct {
	ADDomain types.String   `tfsdk:"ad_domain"`
	BaseDN   types.String   `tfsdk:"base_dn"`
	BindDN   types.String   `tfsdk:"bind_dn"`
	Enabled  types.Bool     `tfsdk:"enabled"`
	Servers  []types.String `tfsdk:"servers"`
}

/*
****************************

	nis

*****************************
*/
func NewNISDataSourceModel() NISDataSourceModel {
	return NISDataSourceModel{
		Domain:  types.String{Null: true},
		Servers: nil,
	}
}

type NISDataSourceModel struct {
	Domain  types.String   `tfsdk:"domain"`
	Enabled types.Bool     `tfsdk:"enabled"`
	Servers []types.String `tfsdk:"servers"`
}

/*
****************************

	nsswitch

*****************************
*/
type NSSwitchDataSourceModel struct {
	Group    []types.String `tfsdk:"group"`
	Hosts    []types.String `tfsdk:"hosts"`
	Namemap  []types.String `tfsdk:"namemap"`
	Netgroup []types.String `tfsdk:"netgroup"`
	Passwd   []types.String `tfsdk:"passwd"`
}

/*
****************************

	routes

*****************************
*/
type RouteDataSourceModel struct {
	Destination RouteDestinationDataSourceModel `tfsdk:"destination"`
	Gateway     types.String                    `tfsdk:"gateway"`
}
type RouteDestinationDataSourceModel struct {
	Address types.String `tfsdk:"address"`
	Family  types.String `tfsdk:"family"`
	Netmask types.String `tfsdk:"netmask"`
}

type S3DataSourceModel struct{}

type SnapmirrorDataSourceModel struct{}

type SnapshotPolicyDataSourceModel struct{}

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
				"fc_interfaces": {
					Optional: true,
					Type: types.ListType{
						ElemType: types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"data_protocol": types.StringType,
								"name":          types.StringType,
								"uuid":          types.StringType,
							},
						},
					},
				},
				"fcp": {
					Type:     types.BoolType,
					Optional: true,
				},
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
					Type: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"ad_domain": types.StringType,
							"base_dn":   types.StringType,
							"bind_dn":   types.StringType,
							"enabled":   types.BoolType,
							"servers": types.ListType{
								ElemType: types.StringType,
							},
						},
					},
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
					Type: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"domain":  types.StringType,
							"enabled": types.BoolType,
							"servers": types.ListType{
								ElemType: types.StringType,
							},
						},
					},
				},
				"nvme": {
					Type:     types.BoolType,
					Optional: true,
				},
				"nsswitch": {
					Optional: true,
					Type: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"group": types.ListType{
								ElemType: types.StringType,
							},
							"hosts": types.ListType{
								ElemType: types.StringType,
							},
							"namemap": types.ListType{
								ElemType: types.StringType,
							},
							"netgroup": types.ListType{
								ElemType: types.StringType,
							},
							"passwd": types.ListType{
								ElemType: types.StringType,
							},
						},
					},
				},
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
				"s3": {
					Optional: true,
					Type: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"enabled": types.BoolType,
							"name":    types.StringType,
						},
					},
				},
				"snapmirror": {
					Optional: true,
					Type: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"is_protected":            types.BoolType,
							"protected_volumes_count": types.Int64Type,
						},
					},
				},
				"snapshot_policy": {
					Optional: true,
					Type: types.ObjectType{
						AttrTypes: map[string]attr.Type{
							"name": types.StringType,
							"uuid": types.StringType,
						},
					},
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

	// UUID
	data.UUID = types.String{Value: SVM.UUID}

	// Aggregates
	data.Aggregates = []AggregateDataSourceModel{}
	for _, aggr := range SVM.Aggregates {
		data.Aggregates = append(data.Aggregates, AggregateDataSourceModel{
			Name: types.String{Value: aggr.Name},
			UUID: types.String{Value: aggr.UUID},
		})
	}

	// AggregatesDelegated
	data.AggregatesDelegated = types.Bool{Value: SVM.AggregatesDelegated}

	// CIFS
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
		cifs.Enabled = types.Bool{Value: SVM.CIFS.Enabled}
		data.CIFS = &cifs
	}

	// Certificate

	data.Certificate = types.String{Value: SVM.Certificate.UUID}
	// This code commented implements CIFS settings with types.Object but the syntax
	// of repeated AttrTypes didn't look like a good pattern.
	// Instead, we replaced :
	//
	// CIFS                types.Object         `tfsdk:"cifs"`
	//
	// by
	//
	// CIFS                *CIFSDataSourceModel `tfsdk:"cifs"`
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

	// Comment
	data.Comment = types.String{Value: SVM.Comment}

	// DNS
	data.DNS = &DNSDataSourceModel{}
	for _, d := range SVM.DNS.Domains {
		data.DNS.Domains = append(data.DNS.Domains, types.String{Value: d})
	}
	for _, d := range SVM.DNS.Servers {
		data.DNS.Servers = append(data.DNS.Servers, types.String{Value: d})
	}

	// FC Interfaces
	for _, i := range SVM.FCInterfaces {
		iface := FCInterfaceDataSourceModel{
			DataProtocal: types.String{Value: i.DataProtocal},
			Name:         types.String{Value: i.Name},
			UUID:         types.String{Value: i.UUID},
		}

		data.FCInterfaces = append(data.FCInterfaces, iface)
	}

	// FCP
	data.FCP = types.Bool{Value: SVM.FCP.Enabled}

	// Name
	data.Name = types.String{Value: SVM.Name}

	// IP Interfaces
	for _, i := range SVM.IPInterfaces {
		iface := NewIPInterfaceDataSourceModel()

		iface.IP.Address = types.String{Value: i.IP.Address}

		if i.IP.Netmask != nil {
			iface.IP.Netmask = types.String{Value: *i.IP.Netmask}
		}

		iface.Name = types.String{Value: i.Name}

		if i.ServicePolicy != nil {
			iface.ServicePolicy = types.String{Value: *i.ServicePolicy}
		}

		iface.UUID = types.String{Value: i.UUID}

		for _, s := range i.Services {
			iface.Services = append(iface.Services, types.String{Value: s})
		}

		data.IPInterfaces = append(data.IPInterfaces, iface)
	}

	// IPSpace

	data.IPSpace = &IPSpaceDataModel{
		Name: types.String{Value: SVM.IPSpace.Name},
		UUID: types.String{Value: SVM.IPSpace.UUID},
	}

	// ISCSI
	data.ISCSI = types.Bool{Value: SVM.ISCSI.Enabled}

	// Language
	data.Language = types.String{Value: SVM.Language}

	// LDAP
	ldap := NewLDAPDataSourceModel()
	if SVM.LDAP.ADDomain != nil {
		ldap.ADDomain = types.String{Value: *SVM.LDAP.ADDomain}
	}
	if SVM.LDAP.BaseDN != nil {
		ldap.BaseDN = types.String{Value: *SVM.LDAP.BaseDN}
	}
	if SVM.LDAP.BindDN != nil {
		ldap.BindDN = types.String{Value: *SVM.LDAP.BindDN}
	}
	ldap.Enabled = types.Bool{Value: SVM.LDAP.Enabled}
	data.LDAP = &ldap

	for _, s := range SVM.LDAP.Servers {
		data.LDAP.Servers = append(data.LDAP.Servers, types.String{Value: s})
	}

	// NFS
	data.NFS = types.Bool{Value: SVM.NFS.Enabled}

	// NIS
	nis := NewNISDataSourceModel()
	if SVM.NIS.Domain != nil {
		nis.Domain = types.String{Value: *SVM.NIS.Domain}
	}
	data.NIS = &nis
	for _, s := range SVM.LDAP.Servers {
		data.NIS.Servers = append(data.NIS.Servers, types.String{Value: s})
	}

	// NSSwitch
	nsswitch := NSSwitchDataSourceModel{}
	for _, a := range SVM.NSSwitch.Group {
		nsswitch.Group = append(nsswitch.Group, types.String{Value: a})
	}
	for _, a := range SVM.NSSwitch.Hosts {
		nsswitch.Hosts = append(nsswitch.Hosts, types.String{Value: a})
	}
	for _, a := range SVM.NSSwitch.Namemap {
		nsswitch.Namemap = append(nsswitch.Namemap, types.String{Value: a})
	}
	for _, a := range SVM.NSSwitch.Netgroup {
		nsswitch.Netgroup = append(nsswitch.Netgroup, types.String{Value: a})
	}
	for _, a := range SVM.NSSwitch.Passwd {
		nsswitch.Passwd = append(nsswitch.Passwd, types.String{Value: a})
	}
	data.NSSwitch = &nsswitch

	// NVME
	data.NVME = types.Bool{Value: SVM.NVME.Enabled}

	// Routes
	routes := []RouteDataSourceModel{}

	for _, r := range SVM.Routes {
		route := RouteDataSourceModel{
			Destination: RouteDestinationDataSourceModel{
				Address: types.String{Value: r.Destination.Address},
				Family:  types.String{Value: r.Destination.Family},
				Netmask: types.String{Value: r.Destination.Netmask},
			},
		}
		routes = append(routes, route)
	}
	data.Routes = routes

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
