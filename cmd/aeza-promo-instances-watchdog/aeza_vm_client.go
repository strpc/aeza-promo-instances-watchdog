package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const contentTypeJSON = "application/json"

type AezaVMClient struct {
	apiKey        string
	token         string
	ses6          string
	aezaVmBaseURL string
	client        *http.Client
}

func NewAezaVMClient(timeout time.Duration, aezaVmBaseURL string, apiKey string) *AezaVMClient {
	if timeout == 0 {
		timeout = defaultTimeout
	}
	defaultClient := &http.Client{
		Timeout: timeout,
	}
	return &AezaVMClient{
		apiKey:        apiKey,
		client:        defaultClient,
		aezaVmBaseURL: aezaVmBaseURL,
	}
}

func (ae *AezaVMClient) Auth() error {
	method := "auth/v3/auth_by_key"
	Url, err := ae.generateVMURL(method)
	if err != nil {
		return fmt.Errorf("generate url: %w", err)
	}

	body := map[string]string{
		"key": ae.apiKey,
	}

	bodyBytes, err := json.Marshal(&body)
	if err != nil {
		return fmt.Errorf("marshal auth body: %w", err)
	}

	resp, err := ae.client.Post(Url, contentTypeJSON, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("send auth request: %w", err)
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error read body: %w", err)
	}

	var resultSchema VMAuthSchema
	if err := json.Unmarshal(bodyResp, &resultSchema); err != nil {
		return fmt.Errorf("unmarshal body: %w", err)
	}
	ae.token = resultSchema.Token
	ae.ses6 = resultSchema.Session
	return nil
}

func (ae *AezaVMClient) GetAllInstances() (VMListInstansSchema, error) {
	method := "vm/v3/host"
	var result VMListInstansSchema
	//var result interface{}

	Url, err := ae.generateVMURL(method)
	if err != nil {
		return result, fmt.Errorf("generate url: %w", err)
	}
	err = ae.sendVMRequest(http.MethodGet, Url, &result)
	return result, err
}

func (ae *AezaVMClient) StartInstance(instanceID int) error {
	method := fmt.Sprintf("vm/v3/host/%d/start", instanceID)
	Url, err := ae.generateVMURL(method)
	if err != nil {
		return fmt.Errorf("generate url: %w", err)
	}
	var result StartInstanceSchema
	if err = ae.sendVMRequest(http.MethodPost, Url, &result); err != nil {
		fmt.Errorf("start instance error: %w", err)
	}
	if result.Id != instanceID {
		return fmt.Errorf("error start instance: %d", instanceID)
	}
	return nil
}

func (ae *AezaVMClient) sendVMRequest(method string, Url string, resultSchemaPointer interface{}) error {
	req, err := http.NewRequest(method, Url, nil)
	if err != nil {
		return fmt.Errorf("generate new request: %w", err)
	}

	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: ae.token,
	})
	req.AddCookie(&http.Cookie{
		Name:  "ses6",
		Value: ae.ses6,
	})
	req.Header.Add("x-xsrf-token", ae.token)

	resp, err := ae.client.Do(req)
	if err != nil {
		fmt.Errorf("create request: %w", err)
	}

	defer resp.Body.Close()

	bodyRaw, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}
	err = json.Unmarshal(bodyRaw, &resultSchemaPointer)
	if err != nil {
		return fmt.Errorf("unmarshall data: %w", err)
	}
	return nil
}

func (ae *AezaVMClient) generateVMURL(method string) (string, error) {
	u, err := url.Parse(ae.aezaVmBaseURL)
	if err != nil {
		return "", fmt.Errorf("parse base api url: %w", err)
	}
	u.Path = method
	return u.String(), nil
}
