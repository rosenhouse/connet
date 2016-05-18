package client

import (
	"fmt"
	"net/http"
	"policy-server/models"

	"github.com/dghubble/sling"
)

func NewInnerClient(baseURL string, httpClient *http.Client) *InnerClient {
	slingClient := sling.New().Client(httpClient).Base(baseURL).Set("Accept", "application/json")
	return &InnerClient{
		slingClient: slingClient,
	}
}

type InnerClient struct {
	slingClient *sling.Sling
}

func (c *InnerClient) Poll() ([]models.Rule, error) {
	var rules []models.Rule

	resp, err := c.slingClient.New().Get("/rules").Receive(&rules, nil)
	if err != nil {
		return nil, fmt.Errorf("list rules: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list rules: unexpected status code: %s", resp.Status)
	}

	return rules, nil
}
