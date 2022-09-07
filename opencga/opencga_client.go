package opencga

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

/*
This module provides an OpenCGA client that can be reused in the provider read and create modules
for querying and posting to the OpenCGA API.

It contains common functionality for detecting errors and parsing JSON payloads.
It uses the type structs from opencga_types.go to return data to the caller.

*/

type APIClient struct {
	BaseUrl    string
	Token      string
	HttpClient *http.Client
}

func newClient(baseUrl string) *APIClient {
	c := &APIClient{}
	c.BaseUrl = baseUrl
	c.HttpClient = http.DefaultClient
	log.Printf("created api client for: %s\n", c.BaseUrl)
	return c
}

func buildRequest(client *APIClient, path string, body interface{}, params map[string]string) (*http.Request, error) {
	url := fmt.Sprintf("%s/opencga/webservices/rest/v1/%s", client.BaseUrl, path)

	var req *http.Request
	var err error
	if body != nil {
		reqBody, err := json.Marshal(body)
		if err != nil {
			return nil, errors.New("Failed to convert body to json")
		}
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	} else {
		req, err = http.NewRequest("GET", url, nil)
	}
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	buildQuery(client, req, params)
	return req, nil
}

func buildQuery(client *APIClient, req *http.Request, params map[string]string) {
	q := req.URL.Query()
	if client.Token != "" {
		q.Add("sid", client.Token)
	}
	for key, val := range params {
		q.Add(key, val)
	}
	req.URL.RawQuery = q.Encode()
}

func (c *APIClient) Login(user string, password string) error {
	path := fmt.Sprintf("users/%s/login", user)
	body := map[string]string{
		"password": password,
	}
	req, err := buildRequest(c, path, body, nil)
	if err != nil {
		return err
	}
	resp, err := c.Call(req)
	if err != nil {
		return err
	}

	var login Login
	err = mapstructure.Decode(resp.Results[0], &login)
	if err != nil {
		return err
	}

	c.Token = login.Token
	return nil
}

func (c *APIClient) Call(req *http.Request) (*Response, error) {
	log.Printf("calling: %s %s", req.Method, req.URL)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawdata, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jsondata map[string]interface{}
	err = json.Unmarshal(rawdata, &jsondata)
	if err != nil {
		log.Printf("Failed to unmarshall data: %s", rawdata)
		return nil, err
	}
	log.Printf("received: %v", jsondata)

	var api_response ApiResponse
	err = mapstructure.Decode(jsondata, &api_response)
	if err != nil {
		return nil, err
	}

	// Check response for errors
	if api_response.Error != "" {
		err = fmt.Errorf("API Error: %s", api_response.Error)
	}
	if len(api_response.Responses) != 1 {
		err = fmt.Errorf("API Error: expecting 1 response, got %d", len(api_response.Responses))
	}
	if api_response.Responses[0].ErrorMsg != "" {
		err = fmt.Errorf("API Error: %s", api_response.Responses[0].ErrorMsg)
	}
	if err != nil {
		return nil, err
	}

	return &api_response.Responses[0], nil
}
