package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrInvalidStatusCode = errors.New("invalid status code")
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// GetRedirect returns the final URL after redirection.
func GetRedirect(url string) (string, error) {
	const op = "api.GetRedirect"

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("%s: %w: %d", op, ErrInvalidStatusCode, resp.StatusCode)
	}

	return resp.Header.Get("Location"), nil
}

// DeleteURL sends a DELETE request to the specified URL and returns the response.
func DeleteURL(url string) (*http.Response, error) {
	const op = "api.DeleteURL"

	client := &http.Client{}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return nil, fmt.Errorf("%s: %w: %d", op, ErrInvalidStatusCode, resp.StatusCode)
	}

	return resp, nil
}

// ParseErrorResponse parses the response body and returns an error response.
func ParseErrorResponse(resp *http.Response) ErrorResponse {
	var errorResponse ErrorResponse

	err := json.NewDecoder(resp.Body).Decode(&errorResponse)
	if err != nil {
		errorResponse = ErrorResponse{Error: err.Error()}
	}

	return errorResponse
}
