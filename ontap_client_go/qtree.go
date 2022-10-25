package ontap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Qtree struct {
	UUID string `json:"uuid,omitempty"`

	SVM        UUIDRef `json:"svm"`
	Volume     UUIDRef `json:"volume"`
	VolumeUUID string  `json:"volume_uuid,omitempty"`
	SVMUUID    string  `json:"svm_uuid,omitempty"`

	Name           string `json:"name,omitempty"`
	Id             int64  `json:"id,omitempty"`
	Path           string `json:"path,omitempty"`
	SecurityStyle  string `json:"security_style,omitempty"`
	UnixPermission int64  `json:"unix_permissions,omitempty"`
}

// This is the JSON representation of a Qtree for REST Create / Update

func (qtree Qtree) RestMarshall() ([]byte, error) {
	qtree_json := qtree

	if qtree_json.SVMUUID != "" {
		qtree_json.SVM = UUIDRef{UUID: qtree_json.SVMUUID}
	}
	qtree_json.SVMUUID = ""

	if qtree_json.VolumeUUID != "" {
		qtree_json.Volume = UUIDRef{UUID: qtree_json.VolumeUUID}
	}
	qtree_json.VolumeUUID = ""

	qtree_json.UUID = ""

	return json.Marshal(qtree_json)
}

func (c *Client) CreateQtree(qtree *Qtree) (*Qtree, error) {
	qtree_copy := *qtree

	qtree_copy.SVM = UUIDRef{
		UUID: qtree_copy.SVMUUID,
	}
	qtree_copy.Volume = UUIDRef{
		UUID: qtree_copy.VolumeUUID,
	}
	qtree_copy.VolumeUUID = ""
	qtree_copy.SVMUUID = ""

	req_qtreeJSON, err := json.Marshal(qtree_copy)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://%s/api/storage/qtrees?return_records=true", c.HostURL), bytes.NewBuffer(req_qtreeJSON))

	if err != nil {
		return nil, err
	}

	_, err = c.doRequest(req)

	if err != nil {
		return nil, err
	}

	result_qtree, err := c.GetQtreeInVolume(qtree_copy.Volume.UUID, qtree.Name)

	if err != nil {
		return nil, err
	}

	return result_qtree, nil
}

func (c *Client) GetQtree(uuid string, qtreeName string) (*Qtree, error) {
	// s := strings.Split(uuid, "/")

	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/api/storage/qtrees/%s", c.HostURL, uuid), nil)

	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	qtree := Qtree{}

	err = json.Unmarshal(body, &qtree)

	if err != nil {
		return nil, err
	}

	qtree.VolumeUUID = qtree.Volume.UUID
	qtree.SVMUUID = qtree.SVM.UUID

	qtree.SVM = UUIDRef{}
	qtree.Volume = UUIDRef{}

	qtree.UUID = uuid
	return &qtree, nil
}
func (c *Client) GetQtreeInVolume(volume_uuid string, name string) (*Qtree, error) {

	req_id, err := http.NewRequest("GET", fmt.Sprintf("https://%s/api/storage/qtrees?volume.uuid=%s&name=%s", c.HostURL, volume_uuid, name), nil)

	if err != nil {
		return nil, err
	}

	body_id, err := c.doRequest(req_id)
	if err != nil {
		return nil, err
	}

	qtree_result := QtreeJSONRecordsResponse{}
	err = json.Unmarshal(body_id, &qtree_result)
	if err != nil {
		return nil, err
	}
	qtreeUUID := fmt.Sprintf("%s/%d", volume_uuid, qtree_result.Records[0].ID)

	req, err := http.NewRequest("GET", fmt.Sprintf("https://%s/api/storage/qtrees/%s", c.HostURL, qtreeUUID), nil)

	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	qtree := Qtree{}
	err = json.Unmarshal(body, &qtree)
	if err != nil {
		return nil, err
	}
	qtree.UUID = qtreeUUID
	qtree.VolumeUUID = qtree.Volume.UUID
	qtree.SVMUUID = qtree.SVM.UUID

	qtree.SVM = UUIDRef{}
	qtree.Volume = UUIDRef{}

	if err != nil {
		return nil, err
	}

	return &qtree, nil
}
func (c *Client) UpdateQtree(qtree *Qtree) (*Qtree, error) {

	req_body, err := qtree.RestMarshall()

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", fmt.Sprintf("https://%s/api/storage/qtrees/%s", c.HostURL, qtree.UUID), bytes.NewBuffer(req_body))

	if err != nil {
		return nil, err
	}

	_, err = c.doRequest(req)

	if err != nil {
		return nil, err
	}

	qtree_result, err := c.GetQtree(qtree.UUID, "")

	if err != nil {
		return nil, err
	}

	return qtree_result, nil
}

func (c *Client) DeleteQtree(qtree *Qtree) error {

	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://%s/api/storage/qtrees/%s", c.HostURL, qtree.UUID), nil)

	if err != nil {
		return err
	}

	_, err = c.doRequest(req)

	if err != nil {
		return err
	}

	return nil
}
