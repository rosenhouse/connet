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

var _ = Describe("net-allow", func() {
	var (
		runner           *netapi.Runner
		client           *fakes.Client
		userLoggerBuffer *gbytes.Buffer
		cliConnection    *pluginfakes.FakeCliConnection
	)

	BeforeEach(func() {
		client = &fakes.Client{}
		userLoggerBuffer = gbytes.NewBuffer()
		runner = &netapi.Runner{
			Client:     client,
			UserLogger: log.New(userLoggerBuffer, "", 0),
		}
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
	})

	It("calls netapi.Allow() with the app guids and token", func() {
		err := runner.Run(cliConnection, []string{"net-allow", "app-apple", "app-banana"})
		Expect(err).NotTo(HaveOccurred())
		Expect(client.AllowCallCount()).To(Equal(1))
		g1, g2, token := client.AllowArgsForCall(0)
		Expect(g1).To(Equal("guid-apple"))
		Expect(g2).To(Equal("guid-banana"))
		Expect(token).To(Equal("some-access-token"))
		Expect(userLoggerBuffer).To(gbytes.Say("allowed app-apple <--> app-banana"))
	})

	Context("when missing required arguments", func() {
		It("returns a helpful error", func() {
			err := runner.Run(cliConnection, []string{"net-allow"})
			Expect(err).To(MatchError("missing required arguments, try -h"))

			err = runner.Run(cliConnection, []string{"net-allow", "app-apple"})
			Expect(err).To(MatchError("missing required arguments, try -h"))
		})
	})

})
