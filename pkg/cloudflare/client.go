package cloudflare

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const cloudflareApiBase = "https://api.cloudflare.com/client/v4/"

type Client interface {
	List(ctx context.Context, request ListRecordRequest, zoneID string) []*Record
	Create(ctx context.Context, request CreateRecordRequest, zoneID string) *CreateResponse
	Update(ctx context.Context, request UpdateRecordRequest, zoneID, recordID string) *UpdateResponse
}

type client struct {
	Token   string
	timeout time.Duration
}

func NewClient(token string, timeout time.Duration) Client {
	return newClient(token, timeout)
}

func newClient(token string, timeout time.Duration) *client {
	return &client{
		Token:   token,
		timeout: timeout,
	}
}

type ListResponse struct {
	Success  bool          `json:"success"`
	Errors   []interface{} `json:"errors"`
	Messages []interface{} `json:"messages"`
	Result   []*Record
}

type CreateResponse struct {
	Success  bool          `json:"success"`
	Errors   []interface{} `json:"errors"`
	Messages []interface{} `json:"messages"`
	Result   *Record
}

type UpdateResponse struct {
	Success  bool          `json:"success"`
	Errors   []interface{} `json:"errors"`
	Messages []interface{} `json:"messages"`
	Result   *Record
}

type Record struct {
	Id         string    `json:"id"`
	Type       string    `json:"type"`
	Name       string    `json:"name"`
	Content    string    `json:"content"`
	Proxiable  bool      `json:"proxiable"`
	Proxied    bool      `json:"proxied"`
	Ttl        int       `json:"ttl"`
	Locked     bool      `json:"locked"`
	ZoneId     string    `json:"zone_id"`
	ZoneName   string    `json:"zone_name"`
	CreatedOn  time.Time `json:"created_on"`
	ModifiedOn time.Time `json:"modified_on"`
	Data       struct {
	} `json:"data"`
	Meta struct {
		AutoAdded bool   `json:"auto_added"`
		Source    string `json:"source"`
	} `json:"meta"`
}

type ListRecordRequest struct {
	Variant string `json:"type"`
	Name    string `json:"name"`
}

type CreateRecordRequest struct {
	Variant  string `json:"type"`
	Name     string `json:"name"`
	Content  string `json:"content"`
	TTL      int    `json:"ttl"`
	Priority int    `json:"priority"`
	Proxied  bool   `json:"proxied,omitempty"`
}

type UpdateRecordRequest struct {
	Variant string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
	Proxied bool   `json:"proxied,omitempty"`
}

// List GET zones/:zone_identifier/dns_records
func (c *client) List(ctx context.Context, in ListRecordRequest, zoneID string) []*Record {
	h := c.createHTTPClient()

	url := fmt.Sprintf(cloudflareApiBase+"zones/%s/dns_records", zoneID)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil
	}

	request.Header.Set("Authorization", "Bearer "+c.Token)
	request.Header.Set("Content-Type", "application/json")

	requestWithQuery := request.URL.Query()
	requestWithQuery.Set("type", in.Variant)
	requestWithQuery.Set("name", in.Name)

	request.URL.RawQuery = requestWithQuery.Encode()

	response, err := h.Do(request)
	if err != nil {
		return nil
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil
	}

	var r *ListResponse
	err = json.Unmarshal(body, r)
	if err != nil {
		return nil
	}

	return r.Result
}

// Create POST zones/:zone_identifier/dns_records
func (c *client) Create(ctx context.Context, in CreateRecordRequest, zoneID string) *CreateResponse {
	h := c.createHTTPClient()

	requestBody, err := json.Marshal(in)
	if err != nil {
		return nil
	}

	url := fmt.Sprintf(cloudflareApiBase+"zones/%s/dns_records", zoneID)
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil
	}

	request.Header.Set("Authorization", "Bearer "+c.Token)

	response, err := h.Do(request)
	if err != nil {
		return nil
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil
	}

	var r *CreateResponse
	err = json.Unmarshal(body, r)
	if err != nil {
		return nil
	}

	return r
}

// Update PUT zones/:zone_identifier/dns_records/:identifier
func (c *client) Update(ctx context.Context, in UpdateRecordRequest, zoneID, recordID string) *UpdateResponse {
	h := c.createHTTPClient()

	requestBody, err := json.Marshal(in)
	if err != nil {
		return nil
	}

	url := fmt.Sprintf(cloudflareApiBase+"zones/%s/dns_records/%s", zoneID, recordID)
	request, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil
	}

	request.Header.Set("Authorization", "Bearer "+c.Token)

	response, err := h.Do(request)
	if err != nil {
		return nil
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil
	}

	var r *UpdateResponse
	err = json.Unmarshal(body, r)
	if err != nil {
		return nil
	}

	return r
}

func (c *client) createHTTPClient() *http.Client {
	return &http.Client{
		Timeout: c.timeout,
	}
}
