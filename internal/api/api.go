package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"framer/internal/core"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/httplog/v2"
)

var ApiClient = NewClient("http://localhost:" + os.Getenv("PORT"))

// HTTPClient is an interface to allow for easy testing and swapping implementations
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client is a wrapper for the http.Client to add custom functionality
type Client struct {
	httpClient HTTPClient
	baseURL    string
}

type LoginResponseDto struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// NewClient creates a new instance of Client with the specified baseURL
func NewClient(baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{},
		baseURL:    baseURL,
	}
}

// SetHTTPClient sets a custom HTTP client for the wrapper
func (c *Client) SetHTTPClient(client HTTPClient) {
	c.httpClient = client
}

func isHttpError(resp int) bool {
	return resp < http.StatusOK || resp >= http.StatusBadRequest
}

// Request makes an HTTP request with the given method, path, request payload, and response struct
func (c *Client) Request(method, path string, reqBody, resBody interface{}, getTokens ...func() (string, string)) (int, error) {
	// Create a new request URL using the base URL and path
	url := c.baseURL + path

	// Serialize the request body if it is not nil
	var body []byte
	var err error
	if reqBody != nil {
		body, err = json.Marshal(reqBody)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	// Create a new HTTP request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to create request: %w", err)
	}

	refreshToken := ""
	accessToken := ""
	if len(getTokens) > 0 {
		accessToken, refreshToken = getTokens[0]()
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	// Set the content type header for JSON
	req.Header.Set("Content-Type", "application/json")

	// Send the request using the client
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read and deserialize the response body into the response struct
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
	}

	if isHttpError(resp.StatusCode) {
		refreshResponse := &LoginResponseDto{}
		if resp.StatusCode == http.StatusUnauthorized {
			_, err := c.Request("POST", RefreshApiPath, &LoginResponseDto{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
			}, refreshResponse)
			if err != nil || isHttpError(resp.StatusCode) {
				return http.StatusUnauthorized, fmt.Errorf("unauthorized: %w", err)
			}
		}
		// Return an error with the status and response body.
		return resp.StatusCode, fmt.Errorf("request failed with status: %s, response: %s", resp.Status, respBody)
	}

	if resBody != nil {
		val := string(respBody)
		fmt.Println(val)
		err = json.Unmarshal(respBody, resBody)
		if err != nil {
			return resp.StatusCode, fmt.Errorf("failed to unmarshal response body: %w", err)
		}
	}

	return resp.StatusCode, nil
}
func WriteJSONError(w http.ResponseWriter, statusCode int, err error) {
	var errorResponse interface{}
	if errs, ok := err.(interface{ Unwrap() []error }); ok {
		var messages []string
		for _, e := range errs.Unwrap() {
			messages = append(messages, e.Error())
		}
		errorResponse = map[string][]string{"errors": messages}
	} else {
		errorResponse = map[string]string{"error": err.Error()}
	}
	errorJSON, _ := json.Marshal(errorResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(errorJSON)
}

func HandleError(r *http.Request, w http.ResponseWriter, err error, statusCode int) {

	switch statusCode {
	case http.StatusInternalServerError:
		GetLogger(r).Error(err.Error())
		if os.Getenv("APP_ENV") == string(core.Development) {
			WriteJSONError(w, statusCode, err)
		} else {
			WriteJSONError(w, statusCode, errors.New("internal server error"))
		}
	default:
		WriteJSONError(w, statusCode, err)
	}
}

func GetLogger(r *http.Request) *slog.Logger {
	return httplog.LogEntry(r.Context())
}
