package netapi

import (
	"fmt"
	"log"
	"os"

	"github.com/cloudfoundry/cli/plugin"
)

type Plugin struct{}

func (p *Plugin) isValidCommand(command string) bool {
	cmds := p.GetMetadata().Commands
	for _, c := range cmds {
		if c.Name == command {
			return true
		}
	}
	return false
}

func (p *Plugin) Run(cliConnection plugin.CliConnection, args []string) {
	logger := log.New(os.Stdout, "", 0)

	if !p.isValidCommand(args[0]) {
		return // may be CLI-MESSAGE-UNINSTALL, just silently return
	}

	runner := &Runner{
		Client:        &Client{},
		UserLogger:    logger,
		CliConnection: cliConnection,
	}

	if err := runner.Run(args); err != nil {
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
			Minor: 15,
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
