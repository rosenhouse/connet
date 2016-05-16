package netapi

import "fmt"

type Client struct{}

func (c *Client) Allow(rule Rule, token string) error {
	fmt.Printf("%s\n", rule)
	return nil
}
func (c *Client) Disallow(rule Rule, token string) error {
	fmt.Printf("%s\n", rule)
	return nil
}

func (c *Client) List(token string) ([]Rule, error) {
	return nil, nil
}
