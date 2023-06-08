package main

import (
	"fmt"
	"net/url"
	"strings"
)

type ServicesSchema struct {
	Data struct {
		Items []struct {
			Id int `json:"id"`
		} `json:"items"`
	} `json:"data"`
}

type VMGotoSchema struct {
	Data string `json:"data"`
}

func (vmgts VMGotoSchema) Token() (string, error) {
	u, err := url.Parse(aezaVmBaseURL)
	if err != nil {
		return "", fmt.Errorf("parse base url %s: %w", aezaVmBaseURL, err)
	}
	u.Path = "auth/key/"
	return strings.TrimLeft(vmgts.Data, u.String()), nil
}

type VMAuthSchema struct {
	Confirmed bool        `json:"confirmed"`
	ExpiresAt interface{} `json:"expires_at"`
	Id        int         `json:"id"`
	Session   string      `json:"session"`
	Token     string      `json:"token"`
}

type VMListInstansSchema struct {
	List []struct {
		Id  int `json:"id"`
		Ip4 []struct {
			Interface string `json:"interface"`
			Ip        string `json:"ip"`
		} `json:"ip4"`
		Ip6   []interface{} `json:"ip6"`
		State string        `json:"state"`
	}
}

type StartInstanceSchema struct {
	Id   int `json:"id"`
	Task int `json:"task"`
}
