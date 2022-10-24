package ontap

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Auth       AuthStruct
}

// AuthStruct -
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type QtreeJSONRecordsResponse struct {
	Records []QtreeJSONRecords `json:"records"`
}

type QtreeJSONRecords struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type JobResponseStruct struct {
	Job JobResponseJob `json:"job"`
}
type JobResponseJob struct {
	Links JSONResponseLinks `json:"_links"`
}
type JSONResponseLinks struct {
	Self JobResponseLinksSelf `json:"self"`
}
type JobResponseLinksSelf struct {
	HREF string `json:"href"`
}

type JobStatus struct {
	UUID        string `json:"uuid"`
	Description string `json:"description"`
	State       string `json:"state"`
	Message     string `json:"message"`
	Code        string `json:"code"`
	Start_time  string `json:"start_time"`
	End_time    string `json:"end_time"`
}

type ErrorJSON struct {
	Error ErrorJSONError `json:"error"`
}
type ErrorJSONError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type Error404 struct {
	Message string
	Code    string
}

func (e *Error404) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewClient -
func NewClient(host, username, password *string, ignoreSSLErrors bool) (*Client, error) {
	if ignoreSSLErrors {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		_, err := http.Get("https://golang.org/")
		if err != nil {
			fmt.Println(err)
		}
	}
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		// Default Hashicups URL
	}

	if host != nil {
		c.HostURL = *host
	}

	// If username or password not provided, return empty client
	if username == nil || password == nil {
		return &c, nil
	}

	c.Auth = AuthStruct{
		Username: *username,
		Password: *password,
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.SetBasicAuth(c.Auth.Username, c.Auth.Password)
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusAccepted {
		jobResponse := JobResponseStruct{}
		err = json.Unmarshal(body, &jobResponse)
		if err != nil {
			return nil, err
		}

	jobLoop:
		for {
			req, err := http.NewRequest("GET", fmt.Sprintf("https://%s%s", c.HostURL, jobResponse.Job.Links.Self.HREF), nil)
			if err != nil {
				return nil, err
			}
			body, err = c.doRequest(req)
			if err != nil {
				return nil, err
			}
			jobStatus := JobStatus{}
			err = json.Unmarshal(body, &jobStatus)

			switch jobStatus.State {
			case "error":
			case "failure":
				return nil, fmt.Errorf(jobStatus.Message)
			case "success":
				break jobLoop
			}
		}
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusAccepted {
		if res.StatusCode == 404 {
			errDescription := ErrorJSON{}

			err := json.Unmarshal(body, &errDescription)
			if err != nil {
				return nil, err
			}
			return nil, &Error404{
				Code:    errDescription.Error.Code,
				Message: errDescription.Error.Message,
			}
		}
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
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
