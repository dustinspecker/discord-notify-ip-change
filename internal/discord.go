package internal

import (
	"fmt"
	"io"
	"net/http"
)

func SendMessage(url string, body io.Reader) error {
	_, err := http.Post(url, "application/json", body)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}

	return nil
}
