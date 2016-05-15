package netapi

import (
	"fmt"
	"log"
	"os"

	"github.com/cloudfoundry/cli/plugin"
)

type Client struct{}

func (c *Client) Allow(guid1, guid2, token string) error {
	fmt.Printf("client allow: %s %s %s\n", guid1, guid2, token)
	return nil
}

//go:generate counterfeiter -o ../fakes/Client.go --fake-name Client . client
type client interface {
	Allow(guid1, guid2, token string) error
}

//go:generate counterfeiter -o ../fakes/UserLogger.go --fake-name UserLogger . userLogger
type userLogger interface {
	Printf(format string, v ...interface{})
}

type Runner struct {
	Client     client
	UserLogger userLogger
}

func (r *Runner) Run(cliConnection plugin.CliConnection, args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("missing required arguments, try -h")
	}
	name1 := args[1]
	name2 := args[2]
	app1, err := cliConnection.GetApp(name1)
	if err != nil {
		return fmt.Errorf("getting app %s: %s", name1, err)
	}
	app2, err := cliConnection.GetApp(name2)
	if err != nil {
		return fmt.Errorf("getting app %s: %s", name2, err)
	}
	token, err := cliConnection.AccessToken()
	if err != nil {
		return fmt.Errorf("getting token: %s", err)
	}
	err = r.Client.Allow(app1.Guid, app2.Guid, token)
	if err != nil {
		return fmt.Errorf("allow: %s", err)
	}
	r.UserLogger.Printf("allowed %s <--> %s\n", name1, name2)
	return nil
}

const (
	CommandAllow    = "net-allow"
	CommandDisallow = "net-disallow"
	CommandList     = "net-list"
)

type Plugin struct{}

func (p *Plugin) Run(cliConnection plugin.CliConnection, args []string) {
	logger := log.New(os.Stdout, "", 0)

	runner := &Runner{
		Client:     &Client{},
		UserLogger: logger,
	}

	if err := runner.Run(cliConnection, args); err != nil {
		logger.Fatalf("%s", err)
	}
}

func (p *Plugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "connet",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 12,
		},
		Commands: []plugin.Command{
			plugin.Command{
				Name:     CommandAllow,
				HelpText: "Allow direct network traffic between two apps",
				UsageDetails: plugin.Usage{
					Usage: fmt.Sprintf("cf %s APP_ONE APP_TWO", CommandAllow),
				},
			},
			plugin.Command{
				Name:     CommandDisallow,
				HelpText: "Remove an existing net-allow rule",
				UsageDetails: plugin.Usage{
					Usage: fmt.Sprintf("cf %s APP_ONE APP_TWO", CommandDisallow),
				},
			},
			plugin.Command{
				Name:     CommandList,
				HelpText: "List all network allow rules",
				UsageDetails: plugin.Usage{
					Usage: fmt.Sprintf("cf %s", CommandList),
				},
			},
		},
	}
}
