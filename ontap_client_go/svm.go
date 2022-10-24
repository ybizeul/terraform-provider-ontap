package ontap

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SVM struct {
	UUID string `json:"uuid,omitempty"`

	Name                string        `json:"name,omitempty"`
	Aggregates          []UUIDRef     `json:"aggregates,omitempty"`
	AggregatesDelegated bool          `json:"aggregates_delegated,omitempty"`
	Certificate         UUIDRef       `json:"certificate,omitempty"`
	CIFS                *SVMCIFS      `json:"cifs,omitempty"`
	Comment             string        `json:"comment"`
	DNS                 SVMDNS        `json:"dns"`
	FCInterfaces        []FCInterface `json:"fc_interfaces"`
	FCP                 FCP           `json:"fcp"`
	IPInterfaces        []IPInterface `json:"ip_interfaces"`
	IPSpace             UUIDRef       `json:"ipspace,omitempty"`
	ISCSI               ISCSI         `json:"iscsi,omitempty"`
	Language            string        `json:"language,omitempty"`
	LDAP                *LDAP         `json:"ldap,omitempty"`
	NFS                 *NFS          `json:"nfs,omitempty"`
	NIS                 *NIS          `json:"nis,omitempty"`
	NSSwitch            *NSSwitch     `json:"nsswitch,omitempty"`
	NVME                *NVME         `json:"nvme,omitempty"`
	Routes              []Route       `json:"routes,omitempty"`
}

type SVMCIFS struct {
	ADDomain *ADDomain `json:"ad_domain,omitempty"`
	Enabled  bool      `json:"enabled,omitempty"`
	Name     *string   `json:"name,omitempty"`
}

type ADDomain struct {
	FQDN               string `json:"fqdn,omitempty"`
	OrganizationalUnit string `json:"organizational_unit,omitempty"`
}

type SVMDNS struct {
	Domains []string `json:"domains,omitempty"`
	Servers []string `json:"servers,omitempty"`
}

type FCInterface struct {
	DataProtocal string `json:"data_protocol,omitempty"`
	Name         string `json:"name,omitempty"`
	UUID         string `json:"uuid,omitempty"`
}

type FCP struct {
	Enabled bool `json:"enabled"`
}
type IPInterface struct {
	IP            IPInterfaceIP `json:"ip,omitempty"`
	Name          string        `json:"name,omitempty"`
	ServicePolicy *string       `json:"service_policy,omitempty"`
	Services      []string      `json:"services,omitempty"`
	UUID          string        `json:"uuid,omitempty"`
}

type IPInterfaceIP struct {
	Address string  `json:"address,omitempty"`
	Netmask *string `json:"netmask,omitempty"`
}
type ISCSI struct {
	Enabled bool `json:"enabled,omitempty"`
}

type LDAP struct {
	ADDomain *string  `json:"ad_domain,omitempty"`
	BaseDN   *string  `json:"base_dn,omitempty"`
	BindDN   *string  `json:"bind_dn,omitempty"`
	Enabled  bool     `json:"enabled,omitempty"`
	Servers  []string `json:"servers,omitempty"`
}

type NFS struct {
	Enabled bool `json:"enabled,omitempty"`
}

type NIS struct {
	Domain  *string  `json:"domain,omitempty"`
	Enabled bool     `json:"enabled,omitempty"`
	Servers []string `json:"servers,omitempty"`
}

type NSSwitch struct {
	Group    []string `json:"group,omitempty"`
	Hosts    []string `json:"hosts,omitempty"`
	Namemap  []string `json:"namemap,omitempty"`
	Netgroup []string `json:"netgroup,omitempty"`
	Passwd   []string `json:"passwd,omitempty"`
}

type NVME struct {
	Enabled bool `json:"enabled"`
}

type Route struct {
	Destination RouteDestination `json:"destination,omitempty"`
	Gateway     string           `json:"gateway,omitempty"`
}
type RouteDestination struct {
	Address string `json:"address,omitempty"`
	Family  string `json:"family,omitempty"`
	Netmask string `json:"netmask,omitempty"`
}

func (c *Client) GetSVM(uuid string) (*SVM, error) {
	// s := strings.Split(uuid, "/")

	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/api/svm/svms/%s", c.HostURL, uuid), nil)

	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	fmt.Print(body)
	if err != nil {
		return nil, err
	}

	svm := SVM{}

	err = json.Unmarshal(body, &svm)

	if err != nil {
		return nil, err
	}

	return &svm, nil
}
