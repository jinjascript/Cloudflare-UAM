// Jinja Cloudflare API Client
// Part of the Jinja CPU Monitoring System
// Created by Jinja (2025)
// This module handles all Cloudflare API interactions for DDoS protection

package cloudflare

import (
        "bytes"
        "encoding/json"
        "fmt"
        "io/ioutil"
        "net/http"
        "time"
)

// Client represents a Cloudflare API client
type Client struct {
        APIKey  string
        Email   string
        ZoneID  string
        BaseURL string
        client  *http.Client
}

// SecurityLevelSetting represents the security level setting for Cloudflare
type SecurityLevelSetting struct {
        Value string `json:"value"`
}

// SecurityLevelPayload is the payload for the security level API request
type SecurityLevelPayload struct {
        Value string `json:"value"`
}

// SecurityLevelResponse represents the response from Cloudflare API
type SecurityLevelResponse struct {
        Success bool            `json:"success"`
        Errors  []CloudflareError `json:"errors"`
        Result  json.RawMessage `json:"result"`
}

// CloudflareError represents an error from the Cloudflare API
type CloudflareError struct {
        Code    int    `json:"code"`
        Message string `json:"message"`
}

// NewClient creates a new Cloudflare API client
func NewClient(apiKey, email, zoneID string) *Client {
        return &Client{
                APIKey:  apiKey,
                Email:   email,
                ZoneID:  zoneID,
                BaseURL: "https://api.cloudflare.com/client/v4",
                client: &http.Client{
                        Timeout: 10 * time.Second,
                },
        }
}

// EnableUnderAttackMode enables the "Under Attack Mode" in Cloudflare
func (c *Client) EnableUnderAttackMode() error {
        return c.setSecurityLevel("under_attack")
}

// DisableUnderAttackMode disables the "Under Attack Mode" in Cloudflare
func (c *Client) DisableUnderAttackMode() error {
        return c.setSecurityLevel("medium")
}

// setSecurityLevel sets the security level in Cloudflare
// Possible values: off, essentially_off, low, medium, high, under_attack
func (c *Client) setSecurityLevel(level string) error {
        url := fmt.Sprintf("%s/zones/%s/settings/security_level", c.BaseURL, c.ZoneID)

        payload := SecurityLevelPayload{
                Value: level,
        }

        data, err := json.Marshal(payload)
        if err != nil {
                return fmt.Errorf("error marshaling request payload: %w", err)
        }

        req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(data))
        if err != nil {
                return fmt.Errorf("error creating request: %w", err)
        }

        req.Header.Set("X-Auth-Email", c.Email)
        req.Header.Set("X-Auth-Key", c.APIKey)
        req.Header.Set("Content-Type", "application/json")

        resp, err := c.client.Do(req)
        if err != nil {
                return fmt.Errorf("error executing request: %w", err)
        }
        defer resp.Body.Close()

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                return fmt.Errorf("error reading response body: %w", err)
        }

        var response SecurityLevelResponse
        if err := json.Unmarshal(body, &response); err != nil {
                return fmt.Errorf("error unmarshalling response: %w", err)
        }

        if !response.Success {
                if len(response.Errors) > 0 {
                        return fmt.Errorf("Cloudflare API error: %s (code: %d)", 
                                response.Errors[0].Message, response.Errors[0].Code)
                }
                return fmt.Errorf("unknown Cloudflare API error")
        }

        return nil
}

// GetCurrentSecurityLevel retrieves the current security level from Cloudflare
func (c *Client) GetCurrentSecurityLevel() (string, error) {
        url := fmt.Sprintf("%s/zones/%s/settings/security_level", c.BaseURL, c.ZoneID)

        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
                return "", fmt.Errorf("error creating request: %w", err)
        }

        req.Header.Set("X-Auth-Email", c.Email)
        req.Header.Set("X-Auth-Key", c.APIKey)
        req.Header.Set("Content-Type", "application/json")

        resp, err := c.client.Do(req)
        if err != nil {
                return "", fmt.Errorf("error executing request: %w", err)
        }
        defer resp.Body.Close()

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                return "", fmt.Errorf("error reading response body: %w", err)
        }

        var response SecurityLevelResponse
        if err := json.Unmarshal(body, &response); err != nil {
                return "", fmt.Errorf("error unmarshalling response: %w", err)
        }

        if !response.Success {
                if len(response.Errors) > 0 {
                        return "", fmt.Errorf("Cloudflare API error: %s (code: %d)", 
                                response.Errors[0].Message, response.Errors[0].Code)
                }
                return "", fmt.Errorf("unknown Cloudflare API error")
        }

        var setting SecurityLevelSetting
        if err := json.Unmarshal(response.Result, &setting); err != nil {
                return "", fmt.Errorf("error unmarshalling result: %w", err)
        }

        return setting.Value, nil
}
