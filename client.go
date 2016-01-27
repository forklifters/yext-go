package yext

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	SandboxHost    string = "api-sandbox.yext.com"
	ProductionHost string = "api.yext.com"
	CUSTOMERS_PATH string = "customers"
)

var ResourceNotFound = errors.New("Resource not found")

type Client struct {
	client        *http.Client
	username      string
	password      string
	customerId    string
	baseUrl       string
	retryAttempts int

	ShowRequest bool

	LocationService    *LocationService
	ECLService         *ECLService
	CustomFieldService *CustomFieldService
	FolderService      *FolderService
	LicenseService     *LicenseService
	UserService        *UserService
}

type Config struct {
	Host string
}

func NewClient(username string, password string, customerId string, config Config) *Client {
	c := &Client{
		client:        http.DefaultClient,
		username:      username,
		password:      password,
		customerId:    customerId,
		retryAttempts: 3,
	}

	host := SandboxHost
	if config.Host != "" {
		host = config.Host
	}
	c.baseUrl = "https://" + host + "/v1"

	c.LocationService = &LocationService{client: c}
	c.ECLService = &ECLService{client: c}
	c.CustomFieldService = &CustomFieldService{client: c}
	c.FolderService = &FolderService{client: c}
	c.LicenseService = &LicenseService{client: c}
	c.UserService = &UserService{client: c}

	return c
}

func (c *Client) customerRequestUrl(path string) string {
	return fmt.Sprintf("%s/%s/%s/%s", c.baseUrl, CUSTOMERS_PATH, c.customerId, path)
}

func (c *Client) rawRequestURL(path string) string {
	return fmt.Sprintf("%s/%s", c.baseUrl, path)
}

func (c *Client) NewRequest(method string, path string) (*http.Request, error) {
	return c.NewCustomerRequestBody(method, path, nil)
}

func (c *Client) NewRawRequest(method string, path string) (*http.Request, error) {
	return c.NewRawRequestBody(method, path, nil)
}

func (c *Client) NewRequestJSON(method string, path string, obj interface{}) (*http.Request, error) {
	json, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	return c.NewCustomerRequestBody(method, path, json)
}

func (c *Client) NewRawRequestJSON(method string, path string, obj interface{}) (*http.Request, error) {
	json, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	return c.NewRawRequestBody(method, path, json)
}

func (c *Client) NewRequestBody(method string, fullPath string, data []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, fullPath, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	rawAuth := []byte(fmt.Sprintf("%v:%v", c.username, c.password))
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString(rawAuth))
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *Client) DoRequest(method string, path string, v interface{}) error {
	req, err := c.NewRequest(method, path)
	if err != nil {
		return err
	}

	return c.Do(req, v)
}

func (c *Client) DoRawRequest(method string, path string, v interface{}) error {
	req, err := c.NewRawRequest(method, path)
	if err != nil {
		return err
	}

	return c.Do(req, v)
}

func (c *Client) NewCustomerRequestBody(method string, path string, data []byte) (*http.Request, error) {
	return c.NewRequestBody(method, c.customerRequestUrl(path), data)
}

func (c *Client) NewRawRequestBody(method string, path string, data []byte) (*http.Request, error) {
	return c.NewRequestBody(method, c.rawRequestURL(path), data)
}

func (c *Client) DoRequestJSON(method string, path string, obj interface{}, v interface{}) error {
	req, err := c.NewRequestJSON(method, path, obj)
	if err != nil {
		return err
	}

	return c.Do(req, v)
}

func (c *Client) DoRawRequestJSON(method string, path string, obj interface{}, v interface{}) error {
	req, err := c.NewRawRequestJSON(method, path, obj)
	if err != nil {
		return err
	}

	return c.Do(req, v)
}

func (c *Client) Do(req *http.Request, v interface{}) error {
	// drain and cache the request body
	originalRequestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	var resultError error
	for attempt := 0; attempt <= c.retryAttempts; attempt++ {
		resultError = nil
		time.Sleep(DefaultBackoffPolicy.Duration(attempt))

		// Rehydrate the request body since it might have been drained by the previous attempt
		req.Body = ioutil.NopCloser(bytes.NewBuffer(originalRequestBody))

		if c.ShowRequest {
			fmt.Printf("%+v\n", req)
		}

		resp, err := c.client.Do(req)
		if err != nil {
			resultError = err
			continue
		}

		defer resp.Body.Close()

		if err := CheckResponseError(resp); err != nil {
			resultError = err
			continue
		}

		if v != nil {
			if w, ok := v.(io.Writer); ok {
				io.Copy(w, resp.Body)
			} else {
				resultError = json.NewDecoder(resp.Body).Decode(v)
			}
		}

		if resultError == nil {
			return nil
		}
	}
	return resultError
}

func CheckResponseError(res *http.Response) error {
	if sc := res.StatusCode; 200 <= sc && sc <= 299 {
		return nil
	} else if sc == 404 {
		return ResourceNotFound
	} else {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		errorResponse := &ErrorResponse{Response: res}
		if err := json.Unmarshal(data, errorResponse); err != nil {
			return errors.New(fmt.Sprintf("unable to unmarshal error from: %s : %s", string(data), err))
		}
		return errorResponse
	}
	return nil
}
