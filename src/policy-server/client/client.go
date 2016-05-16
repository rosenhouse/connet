package client

import (
	"fmt"
	"net/http"
	"policy-server/models"

	"github.com/dghubble/sling"
)

func NewOuterClient(baseURL string, httpClient *http.Client) *OuterClient {
	return &OuterClient{}
}

type OuterClient struct {
	slingClient *sling.Sling
}

func (c *OuterClient) ListRules() ([]models.Rule, error) {
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

func (c *OuterClient) AddRule(rule models.Rule) error {
	return nil
}

func (c *OuterClient) DeleteRule(rule models.Rule) error {
	return nil
}
