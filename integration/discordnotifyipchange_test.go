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

		session := runCommand("-ip-url", ipServer.URL(), "-discord-webhook-url", discordServer.URL(), "-timeout", "5s")

		Eventually(ipServer.ReceivedRequests).Should(HaveLen(1), "expected only 1 request made to server")
		Eventually(discordServer.ReceivedRequests).Should(HaveLen(1), "expected only 1 message to be sent to discord")

		Expect(session.ExitCode()).To(Equal(-1), "should continously run and not exit")
		session.Terminate().Wait()
	})

	It("continuosly sends message to Discord webhook", func() {
		ipServer.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.RespondWith(http.StatusOK, `{"ip": "192.168.0.1"}"`),
			),
		)

		discordServer := ghttp.NewServer()
		defer discordServer.Close()

		discordServer.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyContentType("application/json"),
				ghttp.VerifyRequest(http.MethodPost, "/"),
				ghttp.VerifyBody([]byte(`{"content": "192.168.0.1"}`)),
			),
			ghttp.CombineHandlers(
				ghttp.VerifyContentType("application/json"),
				ghttp.VerifyRequest(http.MethodPost, "/"),
				ghttp.VerifyBody([]byte(`{"content": "192.168.0.1"}`)),
			),
		)

		session := runCommand("-ip-url", ipServer.URL(), "-discord-webhook-url", discordServer.URL(), "-interval", "1s")

		Eventually(ipServer.ReceivedRequests).Should(HaveLen(2), "expected requests to continuously made to server")
		Eventually(discordServer.ReceivedRequests).Should(HaveLen(2), "expected message to continuously be sent to discord")

		Expect(session.ExitCode()).To(Equal(-1), "should continously run and not exit")
		session.Terminate().Wait()
	})

	It("returns error when unable to get IP", func() {
		session := runCommand("-ip-url", "", "-timeout", "5s")

		Eventually(session.Err).Should(gbytes.Say(`.*error getting public IP: error getting URL.*`), "should print error")

		Expect(session.ExitCode()).To(Equal(-1), "should continously run and not exit")
		session.Terminate().Wait()
	})

	It("returns error when unable to send message", func() {
		session := runCommand("-ip-url", ipServer.URL(), "-discord-webhook-url", "", "-timeout", "5s")

		Eventually(session.Err).Should(gbytes.Say(`.*error sending message to discord: error sending message: Post "": unsupported protocol scheme "".*`), "should print error")

		Eventually(ipServer.ReceivedRequests).Should(HaveLen(1), "expected only 1 request made to server")

		Expect(session.ExitCode()).To(Equal(-1), "should continously run and not exit")
		session.Terminate().Wait()
	})

	It("returns error when unable to render message", func() {
		session := runCommand("-ip-url", ipServer.URL(), "-discord-webhook-url", "", "-timeout", "5s", "-format", "{{ .}")

		Eventually(session.Err).Should(gbytes.Say(`.*error rendering message: error parsing template: template: message:1: bad character.*`), "should print error")

		Eventually(ipServer.ReceivedRequests).Should(HaveLen(1), "expected only 1 request made to server")

		Expect(session.ExitCode()).To(Equal(-1), "should continously run and not exit")
		session.Terminate().Wait()
	})

	It("returns error when unable to parse interval", func() {
		session := runCommand("-ip-url", "", "-interval", "10f")

		session.Wait()

		Expect(session).To(gexec.Exit(1), "command should exit with non-zero return code")

		Expect(session.Err).To(gbytes.Say(`.*unable to parse interval: time: unknown unit "f" in duration "10f".*`), "should print error")
	})

	It("returns error when unable to parse timeout", func() {
		session := runCommand("-ip-url", "", "-timeout", "10f")

		session.Wait()

		Expect(session).To(gexec.Exit(1), "command should exit with non-zero return code")

		Expect(session.Err).To(gbytes.Say(`.*unable to parse timeout: time: unknown unit "f" in duration "10f".*`), "should print error")
	})
})

func runCommand(args ...string) *gexec.Session {
	command := exec.Command(discordNotifyIPChangeCLI, args...)
	session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
	ExpectWithOffset(1, err).To(BeNil(), "error while running command")

	return session
}
