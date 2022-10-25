package ontap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SVM struct {
	UUID *string `json:"uuid,omitempty"`

	Name                string          `json:"name,omitempty"`
	Aggregates          []UUIDRef       `json:"aggregates,omitempty"`
	AggregatesDelegated bool            `json:"aggregates_delegated,omitempty"`
	Certificate         UUIDRef         `json:"certificate,omitempty"`
	CIFS                *SVMCIFS        `json:"cifs,omitempty"`
	Comment             string          `json:"comment,omitempty"`
	DNS                 *SVMDNS         `json:"dns,omitempty"`
	FCInterfaces        []FCInterface   `json:"fc_interfaces,omitempty"`
	FCP                 *FCP            `json:"fcp,omitempty"`
	IPInterfaces        []IPInterface   `json:"ip_interfaces,omitempty"`
	IPSpace             UUIDRef         `json:"ipspace,omitempty"`
	ISCSI               ISCSI           `json:"iscsi,omitempty"`
	Language            string          `json:"language,omitempty"`
	LDAP                *LDAP           `json:"ldap,omitempty"`
	NFS                 *NFS            `json:"nfs,omitempty"`
	NIS                 *NIS            `json:"nis,omitempty"`
	NSSwitch            *NSSwitch       `json:"nsswitch,omitempty"`
	NVME                *NVME           `json:"nvme,omitempty"`
	Routes              []Route         `json:"routes,omitempty"`
	S3                  *S3             `json:"s3,omitempty"`
	Snapmirror          *Snapmirror     `json:"snapmirror,omitempty"`
	SnapshotPolicy      *SnapshotPolicy `json:"snapshot_policy,omitempty"`
	State               string          `json:"state,omitempty"`
	Subtype             string          `json:"subtype,omitempty"`
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

type S3 struct {
	Enabled bool    `json:"enabled,omitempty"`
	Name    *string `json:"name,omitempty"`
}

type Snapmirror struct {
	IsProtected           bool  `json:"is_protected,omitempty"`
	ProtectedVolumesCount int64 `json:"protected_volulumes_count,omitempty"`
}

type SnapshotPolicy struct {
	Name string `json:"name,omitempty"`
	UUID string `json:"uuid,omitempty"`
}

type SVMSearchResult struct {
	NumRecords int64    `json:"num_records,omitempty"`
	Records    []Record `json:"records,omitempty"`
}
type Record struct {
	UUID string `json:"uuid,omitempty"`
	Name string `json:"name,omitempty"`
}

func (c *Client) CreateSVM(svm *SVM) (*SVM, error) {

	req_SVMJSON, err := json.Marshal(svm)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/api/svm/svms?return_records=true", c.HostURL), bytes.NewBuffer(req_SVMJSON))

	if err != nil {
		return nil, err
	}

	_, err = c.doRequest(req)

	if err != nil {
		return nil, err
	}

	new_svm, err := c.GetSVM(nil, &svm.Name)

	if err != nil {
		return nil, err
	}

	return new_svm, nil
}

func (c *Client) GetSVM(uuid *string, name *string) (*SVM, error) {

	if name != nil {
		req_id, err := http.NewRequest("GET", fmt.Sprintf("https://%s/api/svm/svms?name=%s", c.HostURL, *name), nil)

		if err != nil {
			return nil, err
		}

		body, err := c.doRequest(req_id)

		if err != nil {
			return nil, err
		}

		svm_result := SVMSearchResult{}
		json.Unmarshal(body, &svm_result)
		uuid = &svm_result.Records[0].UUID
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/api/svm/svms/%s", c.HostURL, *uuid), nil)

	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)

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

func (c *Client) UpdateSVM(svm *SVM) (*SVM, error) {

	uuid := svm.UUID
	svm.UUID = nil

	svm.FCP = nil
	req_body, err := json.Marshal(svm)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("https://%s/api/svm/svms/%s", c.HostURL, *uuid), bytes.NewBuffer(req_body))

	if err != nil {
		return nil, err
	}

	_, err = c.doRequest(req)

	if err != nil {
		return nil, err
	}

	svm_result, err := c.GetSVM(uuid, nil)

	if err != nil {
		return nil, err
	}

	return svm_result, nil
}

func (c *Client) DeleteSVM(svm *SVM) error {

	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://%s/api/svm/svms/%s", c.HostURL, *svm.UUID), nil)

	if err != nil {
		return err
	}

	_, err = c.doRequest(req)

	if err != nil {
		return err
	}

	return nil
}
