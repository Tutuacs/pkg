package uploader

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Tutuacs/pkg/config"
	"github.com/Tutuacs/pkg/types"
)

type UploadThingClient struct {
	ApiKey string
	Host   string
}

var uploader *UploadThingClient

func init() {
	uploader = nil
}

func UseUploader() (*UploadThingClient, error) {
	if uploader == nil {
		conf := config.GetUpload()
		if conf.ApiKey == "sk_live_***" {
			return nil, fmt.Errorf("miss configuration ApiKey from of UploadThing, verify your .env file or config.go")
		}
		uploader = &UploadThingClient{
			ApiKey: conf.ApiKey,
			Host:   conf.Host,
		}
	}

	return uploader, nil
}

// Use FileTypes from types package
func (u *UploadThingClient) PrepareUpload(fileConfig types.PrepareUpload) (*types.PrepareUploadResponse, error) {

	apiURL := fmt.Sprintf("%s/v7/prepareUpload", u.Host)
	fileConfig.ContentDisposition = "inline"

	jsonBody, err := json.Marshal(fileConfig)
	if err != nil {
		return nil, fmt.Errorf("error marshaling BodyJSON: %v", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("x-uploadthing-api-key", u.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response Body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %v\n Body: %s", resp.StatusCode, respBody)
	}

	var response types.PrepareUploadResponse

	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	fmt.Println("v7/prepareUpload response received.")
	return &response, nil
}

func (u *UploadThingClient) UploadFile(url string, file types.UploadFile) error {

	jsonBody, err := json.Marshal(file)
	if err != nil {
		return fmt.Errorf("error marshaling BodyJSON: %v", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("x-uploadthing-api-key", u.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("error reading response of response status: %d\n Body: %v", resp.StatusCode, err)
		}
		return fmt.Errorf("unexpected status: %v\n Body: %s", resp.StatusCode, respBody)
	}

	return nil

}
