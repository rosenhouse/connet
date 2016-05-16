package main

import (
	"encoding/json"
	"flag"
	"lib/marshal"
	"log"
	"os"
	"policy-server/config"
	"policy-server/handlers"

	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/grouper"
	"github.com/tedsuo/ifrit/http_server"
	"github.com/tedsuo/ifrit/sigmon"
	"github.com/tedsuo/rata"
)

func main() {
	var configFilePath string
	const configFileFlag = "configFile"

	flag.StringVar(&configFilePath, configFileFlag, "", "")
	flag.Parse()

	conf, err := config.ParseConfigFile(configFilePath)

	logger := lager.NewLogger("policy-server")
	marshaler := marshal.MarshalFunc(json.Marshal)

	rataHandlers := rata.Handlers{}
	rataHandlers["rules_list"] = &handlers.RulesList{
		Logger:    logger,
		Marshaler: marshaler,
	}

	routes := rata.Routes{
		{Name: "rules_list", Method: "GET", Path: "/rules"},
	}

	rataRouter, err := rata.NewRouter(routes, rataHandlers)
	if err != nil {
		log.Fatalf("unable to create rata router: %s", err) // not tested
	}

	httpServer := http_server.New(conf.ListenAddress, rataRouter)

	members := grouper.Members{
		{"http_server", httpServer},
	}

	group := grouper.NewOrdered(os.Interrupt, members)

	monitor := ifrit.Invoke(sigmon.New(group))

	err = <-monitor.Wait()
	if err != nil {
		log.Fatalf("terminated: %s", err)
	}
}
