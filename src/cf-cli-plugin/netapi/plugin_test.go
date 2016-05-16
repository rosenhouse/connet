package netapi_test

import (
	"cf-cli-plugin/fakes"
	"cf-cli-plugin/netapi"
	"errors"
	"log"

	"github.com/cloudfoundry/cli/plugin/models"
	"github.com/cloudfoundry/cli/plugin/pluginfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("plugin", func() {
	var (
		runner           *netapi.Runner
		client           *fakes.Client
		userLoggerBuffer *gbytes.Buffer
		cliConnection    *pluginfakes.FakeCliConnection
	)

	BeforeEach(func() {
		client = &fakes.Client{}
		userLoggerBuffer = gbytes.NewBuffer()
		cliConnection = &pluginfakes.FakeCliConnection{}
		cliConnection.GetAppStub = func(name string) (plugin_models.GetAppModel, error) {
			if name == "app-banana" {
				return plugin_models.GetAppModel{Guid: "guid-banana"}, nil
			} else if name == "app-apple" {
				return plugin_models.GetAppModel{Guid: "guid-apple"}, nil
			} else {
				return plugin_models.GetAppModel{}, errors.New("bad app")
			}
		}
		cliConnection.AccessTokenReturns("some-access-token", nil)
		runner = &netapi.Runner{
			Client:        client,
			UserLogger:    log.New(userLoggerBuffer, "", 0),
			CliConnection: cliConnection,
		}
	})

	Describe("net-allow", func() {
		It("calls netapi.Allow() with the app guids and token", func() {
			err := runner.Run([]string{"net-allow", "app-apple", "app-banana"})
			Expect(err).NotTo(HaveOccurred())
			Expect(client.AllowCallCount()).To(Equal(1))
			rule, token := client.AllowArgsForCall(0)
			Expect(rule.AppGuid1).To(Equal("guid-apple"))
			Expect(rule.AppGuid2).To(Equal("guid-banana"))
			Expect(token).To(Equal("some-access-token"))
			Expect(userLoggerBuffer).To(gbytes.Say("allowed app-apple <--> app-banana"))
		})

		Context("when missing required arguments", func() {
			It("returns a helpful error", func() {
				err := runner.Run([]string{"net-allow"})
				Expect(err).To(MatchError("missing required arguments, try -h"))

				err = runner.Run([]string{"net-allow", "app-apple"})
				Expect(err).To(MatchError("missing required arguments, try -h"))
			})
		})
	})

	Describe("net-disallow", func() {
		It("calls netapi.Disallow() with the app guids and token", func() {
			err := runner.Run([]string{"net-disallow", "app-apple", "app-banana"})
			Expect(err).NotTo(HaveOccurred())
			Expect(client.DisallowCallCount()).To(Equal(1))
			rule, token := client.DisallowArgsForCall(0)
			Expect(rule.AppGuid1).To(Equal("guid-apple"))
			Expect(rule.AppGuid2).To(Equal("guid-banana"))
			Expect(token).To(Equal("some-access-token"))
			Expect(userLoggerBuffer).To(gbytes.Say("disallowed app-apple <--> app-banana"))
		})

		Context("when missing required arguments", func() {
			It("returns a helpful error", func() {
				err := runner.Run([]string{"net-disallow"})
				Expect(err).To(MatchError("missing required arguments, try -h"))

				err = runner.Run([]string{"net-disallow", "app-apple"})
				Expect(err).To(MatchError("missing required arguments, try -h"))
			})
		})
	})

	Describe("net-list", func() {
		BeforeEach(func() {
			client.ListReturns([]netapi.Rule{
				{AppGuid1: "apple", AppGuid2: "banana"},
				{AppGuid1: "plum", AppGuid2: "peach"},
			}, nil)
		})

		It("calls netapi.List() with the token", func() {
			err := runner.Run([]string{"net-list"})
			Expect(err).NotTo(HaveOccurred())
			Expect(client.ListCallCount()).To(Equal(1))
			token := client.ListArgsForCall(0)
			Expect(token).To(Equal("some-access-token"))
			Expect(userLoggerBuffer).To(gbytes.Say("net-allow rules:\n"))
			Expect(userLoggerBuffer).To(gbytes.Say("apple <--> banana\n"))
			Expect(userLoggerBuffer).To(gbytes.Say("plum <--> peach"))
		})
	})

})
