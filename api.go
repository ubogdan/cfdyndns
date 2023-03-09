package cfdyndns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type API struct {
	APIToken string
	client   http.Client
}

func New(apiToken string) *API {
	return &API{
		APIToken: apiToken,
		client: http.Client{
			Timeout: 3 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:    10,
				IdleConnTimeout: 30 * time.Second,
			},
		},
	}
}

func (a *API) GetZone(domain string) (*Zone, error) {
	query, err := url.Parse("https://api.cloudflare.com/client/v4/zones")
	if err != nil {
		return nil, err
	}
	params := query.Query()
	params.Set("name", domain)
	params.Set("status", "active")
	query.RawQuery = params.Encode()

	req, err := http.NewRequest(http.MethodGet, query.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Authorization", "Bearer "+a.APIToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http.Do %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll %w", err)
	}

	var zoneQuery ZoneQuery

	err = json.Unmarshal(data, &zoneQuery)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal %w", err)
	}

	if len(zoneQuery.Result) == 0 {
		return nil, fmt.Errorf("no zone found for domain %s", domain)
	}

	return &zoneQuery.Result[0], nil
}

func (a *API) GetRecord(zoneID, record string) (*Record, error) {
	query, err := url.Parse("https://api.cloudflare.com/client/v4/zones/" + zoneID + "/dns_records")
	if err != nil {
		return nil, fmt.Errorf("url.Parse %w", err)
	}
	params := query.Query()
	params.Set("name", record)
	params.Set("type", "A")
	query.RawQuery = params.Encode()

	req, err := http.NewRequest(http.MethodGet, query.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Authorization", "Bearer "+a.APIToken)
	req.Header.Set("Content-Type", jsonMime)

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http.Do %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll %w", err)
	}

	var recordQuery RecordQuery

	err = json.Unmarshal(data, &recordQuery)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal %w", err)
	}

	if len(recordQuery.Result) == 0 {
		return nil, fmt.Errorf("no zone found for domain %s", record)
	}

	return &recordQuery.Result[0], nil
}

func (a *API) UpdateRecord(record Record) error {
	query, err := url.Parse("https://api.cloudflare.com/client/v4/zones/" + record.ZoneId + "/dns_records/" + record.Id)
	if err != nil {
		return fmt.Errorf("url.Parse %w", err)
	}

	payload, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("json.Marshal %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, query.String(), bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("http.NewRequest %w", err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Authorization", "Bearer "+a.APIToken)
	req.Header.Set("Content-Type", jsonMime)

	resp, err := a.client.Do(req)
	if err != nil {
		return fmt.Errorf("http.Do %w", err)
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("io.ReadAll %w", err)
	}

	return nil
}
