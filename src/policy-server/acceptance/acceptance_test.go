package acceptance_test

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"policy-server/client"
	"policy-server/config"
	"policy-server/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"github.com/pivotal-golang/lager/lagertest"
)

var _ = Describe("Policy server", func() {
	var (
		session        *gexec.Session
		address        string
		configFilePath string
		outerClient    *client.OuterClient

		logger *lagertest.TestLogger
	)

	BeforeEach(func() {
		address = fmt.Sprintf("127.0.0.1:%d", 4001+GinkgoParallelNode())

		logger = lagertest.NewTestLogger("test")
		configFilePath = WriteConfigFile(&config.ServerConfig{
			ListenAddress: address,
		})

		serverCmd := exec.Command(serverBinPath, "-configFile", configFilePath)
		var err error
		session, err = gexec.Start(serverCmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		outerClient = client.NewOuterClient("http://"+address, http.DefaultClient)
	})

	AfterEach(func() {
		session.Interrupt()
		Eventually(session, DEFAULT_TIMEOUT).Should(gexec.Exit(0))
		Expect(os.Remove(configFilePath)).To(Succeed())
	})

	var serverIsAvailable = func() error {
		return VerifyTCPConnection(address)
	}

	It("should boot and gracefully terminate", func() {
		Eventually(serverIsAvailable, DEFAULT_TIMEOUT).Should(Succeed())

		Consistently(session).ShouldNot(gexec.Exit())

		session.Interrupt()
		Eventually(session, DEFAULT_TIMEOUT).Should(gexec.Exit(0))
	})

	Describe("rule lifecycle", func() {
		It("should support list, add and delete on the set of rules", func() {
			Eventually(serverIsAvailable, DEFAULT_TIMEOUT).Should(Succeed())

			By("listing the rules")
			rules, err := outerClient.ListRules()
			Expect(err).NotTo(HaveOccurred())
			Expect(rules).To(BeEmpty())

			By("adding a new rule")
			Expect(outerClient.AddRule(models.Rule{
				Group1: "group1",
				Group2: "group2",
			})).To(Succeed())

			By("adding a second rule")
			Expect(outerClient.AddRule(models.Rule{
				Group1: "group2",
				Group2: "group3",
			})).To(Succeed())

			By("adding a third rule")
			Expect(outerClient.AddRule(models.Rule{
				Group1: "group2",
				Group2: "group2",
			})).To(Succeed())

			By("listing the rules")
			rules, err = outerClient.ListRules()
			Expect(err).NotTo(HaveOccurred())
			Expect(rules).To(HaveLen(3))
			Expect(rules).To(ConsistOf([]models.Rule{
				{Group1: "group1", Group2: "group2"},
				{Group1: "group2", Group2: "group3"},
				{Group1: "group2", Group2: "group2"},
			}))

			By("removing the second rule")
			Expect(outerClient.DeleteRule(models.Rule{
				Group1: "group2",
				Group2: "group3",
			})).To(Succeed())

			By("listing the rules")
			rules, err = outerClient.ListRules()
			Expect(err).NotTo(HaveOccurred())
			Expect(rules).To(HaveLen(2))
			Expect(rules).To(ConsistOf([]models.Rule{
				{Group1: "group1", Group2: "group2"},
				{Group1: "group2", Group2: "group2"},
			}))
		})
	})
})
