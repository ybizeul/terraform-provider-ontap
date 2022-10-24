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
	DataProtocal string `json:"data_protocol"`
	Name         string `json:"name"`
	UUID         string `json:"uuid"`
}

type FCP struct {
	Enabled bool `json:"enabled"`
}
type IPInterface struct {
	IP            IPInterfaceIP `json:"ip"`
	Name          string        `json:"name"`
	ServicePolicy *string       `json:"service_policy"`
	Services      []string      `json:"services"`
	UUID          string        `json:"uuid"`
}

type IPInterfaceIP struct {
	Address string  `json:"address"`
	Netmask *string `json:"netmask"`
}

// This is the JSON representation of a Qtree for REST Create / Update

// func (svm SVM) RestMarshall() ([]byte, error) {
// 	qtree_json := qtree

// 	if qtree_json.SVMUUID != "" {
// 		qtree_json.SVM = UUIDRef{UUID: qtree_json.SVMUUID}
// 	}
// 	qtree_json.SVMUUID = ""

// 	if qtree_json.SVMUUID != "" {
// 		qtree_json.Volume = UUIDRef{UUID: qtree_json.VolumeUUID}
// 	}
// 	qtree_json.VolumeUUID = ""

// 	qtree_json.UUID = ""

// 	return json.Marshal(qtree_json)
// }

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
