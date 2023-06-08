package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const defaultTimeout = 60 * time.Second

type AezaAPIClient struct {
	apiToken       string
	aezaApiBaseURL string
	client         *http.Client
}

func NewAezaAPIClient(timeout time.Duration, apiToken string, aezaApiBaseURL string) *AezaAPIClient {
	if timeout == 0 {
		timeout = defaultTimeout
	}
	defaultClient := &http.Client{
		Timeout: timeout,
	}
	return &AezaAPIClient{
		client:         defaultClient,
		apiToken:       apiToken,
		aezaApiBaseURL: aezaApiBaseURL,
	}
}

func (a *AezaAPIClient) GetAllServices() (ServicesSchema, error) {
	method := "api/services"
	var result ServicesSchema
	Url, err := a.generateApiURL(method)
	if err != nil {
		return result, fmt.Errorf("generate api url: %w", err)
	}
	err = a.sendApiGetRequest(Url, &result)
	return result, err
}

func (a *AezaAPIClient) GetVMToken(serviceID int) (string, error) {
	method := fmt.Sprintf("api/services/%d/goto", serviceID)

	Url, err := a.generateApiURL(method)
	if err != nil {
		return "", fmt.Errorf("generate apu url: %w", err)
	}
	var respSchema VMGotoSchema
	if err = a.sendApiGetRequest(Url, &respSchema); err != nil {
		return "", err
	}
	token, err := respSchema.Token()
	if err != nil {
		return "", fmt.Errorf("vm token not find")
	}
	return token, nil
}

func (a *AezaAPIClient) sendApiGetRequest(url string, resultSchemaPointer interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("X-API-KEY", a.apiToken)
	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("send request to %s: %w", url, err)
	}
	defer resp.Body.Close()
	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error read body: %w", err)
	}
	if err = json.Unmarshal(bodyResp, resultSchemaPointer); err != nil {
		return fmt.Errorf("unmarshal response body get all services: %w", err)
	}
	return nil
}

func (a *AezaAPIClient) generateApiURL(method string) (string, error) {
	u, err := url.Parse(a.aezaApiBaseURL)
	if err != nil {
		return "", fmt.Errorf("parse base api url: %w", err)
	}
	u.Path = method
	return u.String(), nil
}
