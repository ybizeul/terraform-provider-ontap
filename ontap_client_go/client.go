package ontap

import (
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
			time.Sleep(time.Second)
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
