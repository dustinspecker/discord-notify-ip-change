package discord

import (
	"bytes"
	"fmt"
	"net/http"
)

func SendMessage(url string) error {
	body := bytes.NewReader([]byte("hey"))

	_, err := http.Post(url, "application/json", body)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}

	return nil
}
