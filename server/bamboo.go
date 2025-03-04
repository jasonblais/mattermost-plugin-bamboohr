package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	baseUrl               = "https://api.bamboohr.com/api/gateway.php/%s"
	apiVersion            = "v1"
	employeeDirectoryLink = "/v1/employees/directory"
)

// type employee struct {
// 	ID          string `json: "id`
// 	DisplayName string `json: displayName`
// 	JobTitle    string `json: jobTitle`
// 	Location    string `json: location`
// }

// type employeeDirectory struct {
// 	employees []employee
// }

type Client struct {
	client  *http.Client
	BaseUrl string
	// Directory *employeeDirectory
}

// Takes in your company's bamboo subdomain
func NewClient(httpClient *http.Client, subdomain string) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	builtBambooUrl := buildBambooURL(subdomain, baseUrl)

	c := &Client{
		client:  httpClient,
		BaseUrl: builtBambooUrl,
	}

	return c
}

func (c *Client) do(key string, req *http.Request) ([]byte, error) {
	req.SetBasicAuth(key, "")
	req.Header.Set("Accept", "application/json")

	// Might not be necessary
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}

	return body, nil
}

func (c *Client) buildEmployeeDirectory(key string) error {
	directoryUrl := buildUrlToDirectory(c.BaseUrl, employeeDirectoryLink)
	req, err := http.NewRequest("GET", directoryUrl, nil)

	if err != nil {
		return err
	}
	bytes, err := c.do(key, req)

	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &c.Directory)

	if err != nil {
		return err
	}

	return nil
}

func buildBambooURL(subdomain string, baseUrl string) string {
	return fmt.Sprintf(baseUrl, subdomain)
}

func buildUrlToDirectory(b string, d string) string {
	return b + d
}
