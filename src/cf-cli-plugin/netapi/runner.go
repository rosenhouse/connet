package netapi

import (
	"fmt"

	"github.com/cloudfoundry/cli/plugin"
)

const (
	CommandAllow    = "net-allow"
	CommandDisallow = "net-disallow"
	CommandList     = "net-list"
)

type Rule struct {
	AppGuid1 string
	AppGuid2 string
}

func (r Rule) String() string {
	return fmt.Sprintf("%s <--> %s", r.AppGuid1, r.AppGuid2)
}

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

//go:generate counterfeiter -o ../fakes/Client.go --fake-name Client . client
type client interface {
	Allow(rule Rule, token string) error
	Disallow(rule Rule, token string) error
	List(token string) ([]Rule, error)
}

//go:generate counterfeiter -o ../fakes/UserLogger.go --fake-name UserLogger . userLogger
type userLogger interface {
	Printf(format string, v ...interface{})
}

type Runner struct {
	Client        client
	UserLogger    userLogger
	CliConnection plugin.CliConnection
}

func (r *Runner) getRule(name1, name2 string) (Rule, error) {
	app1, err := r.CliConnection.GetApp(name1)
	if err != nil {
		return Rule{}, fmt.Errorf("getting app %s: %s", name1, err)
	}
	app2, err := r.CliConnection.GetApp(name2)
	if err != nil {
		return Rule{}, fmt.Errorf("getting app %s: %s", name2, err)
	}
	return Rule{AppGuid1: app1.Guid, AppGuid2: app2.Guid}, nil
}

func (r *Runner) Run(args []string) error {
	command := args[0]
	token, err := r.CliConnection.AccessToken()
	if err != nil {
		return fmt.Errorf("getting token: %s", err)
	}

	switch command {
	case CommandList:
		rules, err := r.Client.List(token)
		if err != nil {
			return fmt.Errorf("list: %s", err)
		}
		r.UserLogger.Printf("net-allow rules:")
		for _, rule := range rules {
			r.UserLogger.Printf("%s\n", rule)
		}
	case CommandAllow, CommandDisallow:
		if len(args) != 3 {
			return fmt.Errorf("missing required arguments, try -h")
		}
		name1 := args[1]
		name2 := args[2]
		rule, err := r.getRule(name1, name2)
		if err != nil {
			return fmt.Errorf("%s", err)
		}
		switch command {
		case CommandAllow:
			err = r.Client.Allow(rule, token)
			if err != nil {
				return fmt.Errorf("allow: %s", err)
			}
			r.UserLogger.Printf("allowed %s <--> %s\n", name1, name2)
		case CommandDisallow:
			err = r.Client.Disallow(rule, token)
			if err != nil {
				return fmt.Errorf("disallow: %s", err)
			}
			r.UserLogger.Printf("disallowed %s <--> %s\n", name1, name2)
		}
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
	return nil
}