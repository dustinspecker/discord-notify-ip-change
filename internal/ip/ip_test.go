package ip_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/onsi/gomega"

	"github.com/dustinspecker/discord-notify-ip-change/internal/ip"
)

func TestGetReturnsErrorWhenURLInvalid(t *testing.T) {
	g := gomega.NewWithT(t)

	_, err := ip.Get("")
	g.Expect(err).ToNot(gomega.BeNil(), "error should be returned when invalid URL")

	g.Expect(err.Error()).To(gomega.Equal(`error getting URL "": Get "": unsupported protocol scheme ""`))
}

func TestGetReturnsErrorWhenURLNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	g := gomega.NewWithT(t)

	_, err := ip.Get(server.URL)
	g.Expect(err).ToNot(gomega.BeNil(), "error should be returned when URL not found")

	g.Expect(err.Error()).To(gomega.Equal(fmt.Sprintf(`error getting URL %q: encountered bad status: 404`, server.URL)))
}

func TestGetReturnsErrorWhenInvalidResponseFormat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"ip"`)
	}))
	defer server.Close()

	g := gomega.NewWithT(t)

	_, err := ip.Get(server.URL)
	g.Expect(err).ToNot(gomega.BeNil(), "error should be returned when URL has invalid response format")

	g.Expect(err.Error()).To(gomega.Equal("error unmarshalling response: unexpected EOF"))
}

func TestGetReturnsIP(t *testing.T) {
	expectedIP := "192.168.0.1"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, fmt.Sprintf(`{"ip": "%s"}`, expectedIP))
	}))
	defer server.Close()

	g := gomega.NewWithT(t)

	output, err := ip.Get(server.URL)
	g.Expect(err).To(gomega.BeNil(), "no error should be returned when URL is retrieved")

	g.Expect(output).To(gomega.Equal(expectedIP), "Get should return output from server")
}
