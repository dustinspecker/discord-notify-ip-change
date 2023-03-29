package integration_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var discordNotifyIPChangeCLI string

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = BeforeSuite(func() {
	var err error
	discordNotifyIPChangeCLI, err = gexec.Build("github.com/dustinspecker/discord-notify-ip-change/cmd/discord-notify-ip-change")
	Expect(err).To(BeNil(), "failed to build discord-notify-ip-change")
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
