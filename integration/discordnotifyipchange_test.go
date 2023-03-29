package integration_test

import (
	"net/http"
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("discord-notify-ip-change", func() {
	var ipServer *ghttp.Server

	BeforeEach(func() {
		ipServer = ghttp.NewServer()

		ipServer.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, `{"ip": "192.168.0.1"}"`),
			),
		)
	})

	AfterEach(func() {
		ipServer.Close()
	})

	It("sends message to Discord webhook", func() {
		discordServer := ghttp.NewServer()
		defer discordServer.Close()

		discordServer.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyContentType("application/json"),
				ghttp.VerifyRequest(http.MethodPost, "/"),
				ghttp.VerifyBody([]byte(`{"content": "192.168.0.1"}`)),
			),
		)

		command := exec.Command(discordNotifyIPChangeCLI, "-ip-url", ipServer.URL(), "-discord-webhook-url", discordServer.URL(), "-timeout", "5s")
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).To(BeNil(), "error while running command")

		Eventually(session).Should(gexec.Exit(0), "command should exit with 0 return code")

		Expect(ipServer.ReceivedRequests()).Should(HaveLen(1), "expected only 1 request made to server")
		Expect(discordServer.ReceivedRequests()).Should(HaveLen(1), "expected only 1 message to be sent to discord")
	})

	It("returns error when unable to get IP", func() {
		command := exec.Command(discordNotifyIPChangeCLI, "-ip-url", "", "-timeout", "5s")
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).To(BeNil(), "error while running command")

		Eventually(session).Should(gexec.Exit(1), "command should exit with non-zero return code")

		Expect(session.Err).To(gbytes.Say(`.*error getting public IP: error getting URL.*`), "should print error")
	})

	It("returns error when unable to send message", func() {
		command := exec.Command(discordNotifyIPChangeCLI, "-ip-url", ipServer.URL(), "-discord-webhook-url", "", "-timeout", "5s")
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).To(BeNil(), "error while running command")

		Eventually(session).Should(gexec.Exit(1), "command should exit with 1 return code")
		Expect(session.Err).To(gbytes.Say(`.*error sending message to discord: error sending message: Post "": unsupported protocol scheme "".*`), "should print error")

		Expect(ipServer.ReceivedRequests()).Should(HaveLen(1), "expected only 1 request made to server")
	})

	It("returns error when unable to render message", func() {
		command := exec.Command(discordNotifyIPChangeCLI, "-ip-url", ipServer.URL(), "-discord-webhook-url", "", "-timeout", "5s", "-format", "{{ .}")
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).To(BeNil(), "error while running command")

		Eventually(session).Should(gexec.Exit(1), "command should exit with 1 return code")
		Expect(session.Err).To(gbytes.Say(`.*error rendering message: error parsing template: template: message:1: bad character.*`), "should print error")

		Expect(ipServer.ReceivedRequests()).Should(HaveLen(1), "expected only 1 request made to server")
	})

	It("returns error when unable to parse timeout", func() {
		command := exec.Command(discordNotifyIPChangeCLI, "-ip-url", "", "-timeout", "10f")
		session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).To(BeNil(), "error while running command")

		Eventually(session).Should(gexec.Exit(1), "command should exit with non-zero return code")

		Expect(session.Err).To(gbytes.Say(`.*unable to parse timeout: time: unknown unit "f" in duration "10f".*`), "should print error")
	})
})
