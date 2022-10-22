package ontap

import "encoding/json"

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
