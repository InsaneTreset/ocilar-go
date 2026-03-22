// Package ocilar provides a Go client for the Ocilar CAPTCHA solving and document extraction API.
package ocilar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const defaultBaseURL = "https://api.ocilar.com/api/v1"

type SolveResult struct {
	Text      string `json:"text"`
	TaskID    string `json:"task_id"`
	LatencyMs int    `json:"latency_ms"`
	Status    string `json:"status"`
	Token     string `json:"token,omitempty"`
}

type ExtractResult struct {
	Data      map[string]interface{} `json:"data"`
	TaskID    string                 `json:"task_id"`
	LatencyMs int                    `json:"latency_ms"`
	Status    string                 `json:"status"`
}

type Error struct {
	StatusCode int
	Detail     string
}

func (e *Error) Error() string {
	return fmt.Sprintf("ocilar: %d %s", e.StatusCode, e.Detail)
}

type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client
}

// NewClient creates a new Ocilar API client.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: defaultBaseURL,
		http:    &http.Client{},
	}
}

// WithBaseURL sets a custom base URL.
func (c *Client) WithBaseURL(url string) *Client {
	c.baseURL = url
	return c
}

func (c *Client) post(path string, body interface{}) ([]byte, error) {
	payload, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		var errResp struct{ Detail string `json:"detail"` }
		json.Unmarshal(data, &errResp)
		return nil, &Error{StatusCode: resp.StatusCode, Detail: errResp.Detail}
	}
	return data, nil
}

func (c *Client) get(path string) ([]byte, error) {
	req, err := http.NewRequest("GET", c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-Key", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		var errResp struct{ Detail string `json:"detail"` }
		json.Unmarshal(data, &errResp)
		return nil, &Error{StatusCode: resp.StatusCode, Detail: errResp.Detail}
	}
	return data, nil
}

func parseSolve(data []byte) (*SolveResult, error) {
	var r SolveResult
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	if r.Text == "" && r.Token != "" {
		r.Text = r.Token
	}
	return &r, nil
}

func parseExtract(data []byte) (*ExtractResult, error) {
	var r ExtractResult
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// ── CAPTCHA Solving ──

func (c *Client) SolveSAT(imageBase64 string) (*SolveResult, error) {
	data, err := c.post("/solve/sat", map[string]string{"image": imageBase64})
	if err != nil { return nil, err }
	return parseSolve(data)
}

func (c *Client) SolveIMSS(imageBase64 string) (*SolveResult, error) {
	data, err := c.post("/solve/imss", map[string]string{"image": imageBase64})
	if err != nil { return nil, err }
	return parseSolve(data)
}

func (c *Client) SolveImage(imageBase64 string) (*SolveResult, error) {
	data, err := c.post("/solve/image", map[string]string{"image": imageBase64})
	if err != nil { return nil, err }
	return parseSolve(data)
}

func (c *Client) SolveRecaptchaV2(siteKey, siteURL string) (*SolveResult, error) {
	data, err := c.post("/solve/recaptcha-v2", map[string]string{"site_key": siteKey, "site_url": siteURL})
	if err != nil { return nil, err }
	return parseSolve(data)
}

func (c *Client) SolveRecaptchaV3(siteKey, siteURL, action string) (*SolveResult, error) {
	data, err := c.post("/solve/recaptcha-v3", map[string]string{"site_key": siteKey, "site_url": siteURL, "action": action})
	if err != nil { return nil, err }
	return parseSolve(data)
}

func (c *Client) SolveHCaptcha(siteKey, siteURL string) (*SolveResult, error) {
	data, err := c.post("/solve/hcaptcha", map[string]string{"site_key": siteKey, "site_url": siteURL})
	if err != nil { return nil, err }
	return parseSolve(data)
}

func (c *Client) SolveCloudflare(siteURL string) (*SolveResult, error) {
	data, err := c.post("/solve/cloudflare", map[string]string{"site_url": siteURL})
	if err != nil { return nil, err }
	return parseSolve(data)
}

func (c *Client) SolveAudio(audioBase64 string) (*SolveResult, error) {
	data, err := c.post("/solve/audio", map[string]string{"audio": audioBase64})
	if err != nil { return nil, err }
	return parseSolve(data)
}

// ── Document AI ──

func (c *Client) ExtractCSF(documentBase64 string) (*ExtractResult, error) {
	data, err := c.post("/solve/extract_csf", map[string]string{"document": documentBase64})
	if err != nil { return nil, err }
	return parseExtract(data)
}

func (c *Client) ExtractINE(documentBase64 string) (*ExtractResult, error) {
	data, err := c.post("/solve/extract_ine", map[string]string{"document": documentBase64})
	if err != nil { return nil, err }
	return parseExtract(data)
}

func (c *Client) ExtractCFDI(documentBase64 string) (*ExtractResult, error) {
	data, err := c.post("/solve/extract_cfdi", map[string]string{"document": documentBase64})
	if err != nil { return nil, err }
	return parseExtract(data)
}

func (c *Client) ExtractCURP(documentBase64 string) (*ExtractResult, error) {
	data, err := c.post("/solve/extract_curp", map[string]string{"document": documentBase64})
	if err != nil { return nil, err }
	return parseExtract(data)
}

func (c *Client) ExtractDomicilio(documentBase64 string) (*ExtractResult, error) {
	data, err := c.post("/solve/extract_domicilio", map[string]string{"document": documentBase64})
	if err != nil { return nil, err }
	return parseExtract(data)
}

func (c *Client) ExtractNomina(documentBase64 string) (*ExtractResult, error) {
	data, err := c.post("/solve/extract_nomina", map[string]string{"document": documentBase64})
	if err != nil { return nil, err }
	return parseExtract(data)
}

func (c *Client) ExtractGeneric(documentBase64 string) (*ExtractResult, error) {
	data, err := c.post("/solve/extract_generic", map[string]string{"document": documentBase64})
	if err != nil { return nil, err }
	return parseExtract(data)
}

// ── Utilities ──

func (c *Client) Hello() (map[string]interface{}, error) {
	data, err := c.post("/hello", map[string]string{})
	if err != nil { return nil, err }
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	return result, nil
}

func (c *Client) GetBalance() (map[string]interface{}, error) {
	data, err := c.get("/balance")
	if err != nil { return nil, err }
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	return result, nil
}
