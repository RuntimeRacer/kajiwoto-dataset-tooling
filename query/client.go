package query

import (
	"context"
	"fmt"
	"github.com/runtimeracer/go-graphql-client"
	"net/http"
)

// headerTransport is used to add custom headers to the request
// shootout to tgwizard; https://github.com/shurcooL/graphql/issues/28
type headerTransport struct {
	base    http.RoundTripper
	headers map[string]string
}

func (h *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := CloneRequest(req)
	for key, val := range h.headers {
		req2.Header.Set(key, val)
	}
	return h.base.RoundTrip(req2)
}

func (h *headerTransport) GetHeaders() map[string]string {
	return h.headers
}

func (h *headerTransport) AddHeaders(newHeaders map[string]string) {
	for k, v := range newHeaders {
		h.headers[k] = v
	}
}

// kajiwotoClient is a custom graphql client for kajiwoto reqeusts
type kajiwotoClient struct {
	client          *graphql.Client
	transportClient *http.Client
}

func GetKajiwotoClient(endpoint string) *kajiwotoClient {
	// Init HTTP Client
	transportClient := &http.Client{
		Transport: &headerTransport{
			base:    http.DefaultTransport,
			headers: map[string]string{},
		},
	}

	return &kajiwotoClient{
		client:          graphql.NewClient(endpoint, transportClient),
		transportClient: transportClient,
	}
}

func (c *kajiwotoClient) GetHeaders() map[string]string {
	return c.transportClient.Transport.(*headerTransport).GetHeaders()
}

func (c *kajiwotoClient) AddHeaders(newHeaders map[string]string) {
	c.transportClient.Transport.(*headerTransport).AddHeaders(newHeaders)
}

// DoLoginUserPW performs login via user / pw combination
func (c *kajiwotoClient) DoLoginUserPW(username, password string) (result LoginResult, err error) {
	// Sanity check
	if username == "" || password == "" {
		return result, fmt.Errorf("invalid login credentials")
	}

	vars := map[string]interface{}{
		"usernameOrEmail": graphql.String(username),
		"password":        graphql.String(password),
	}

	loginResult := kajiwotoLoginUserPWMutation{}
	if errLogin := c.performGraphMutation(vars, &loginResult); errLogin != nil {
		return result, fmt.Errorf("unable to login, response: %q", errLogin)
	}

	// Build generic Result object
	result = LoginResult{
		Login:   loginResult.Login,
		Welcome: loginResult.Welcome,
	}

	return result, nil
}

// DoLoginAuthToken performs login via session key if available
func (c *kajiwotoClient) DoLoginAuthToken(authToken string) (result LoginResult, err error) {
	// Sanity check
	if authToken == "" {
		return result, fmt.Errorf("invalid login credentials")
	}

	vars := map[string]interface{}{
		"authToken": graphql.Token(authToken),
		"action":    graphql.String(""),
	}

	// Add Auth-Token header
	headers := map[string]string{
		"auth_token": authToken,
	}
	c.AddHeaders(headers)

	loginResult := kajiwotoLoginAuthTokenMutation{}
	if errLogin := c.performGraphMutation(vars, &loginResult); errLogin != nil {
		return result, fmt.Errorf("unable to login, response: %q", errLogin)
	}

	// Build generic Result object
	result = LoginResult{
		Login:   loginResult.Login,
		Welcome: loginResult.Welcome,
	}

	return result, nil
}

func (c *kajiwotoClient) performGraphMutation(vars map[string]interface{}, mutation interface{}) error {
	return c.client.Mutate(context.Background(), mutation, vars)
}

// CloneRequest creates a shallow copy of the request along with a deep copy of the Headers.
func CloneRequest(req *http.Request) *http.Request {
	r := new(http.Request)

	// shallow clone
	*r = *req

	// deep copy headers
	r.Header = CloneHeader(req.Header)

	return r
}

// CloneHeader creates a deep copy of an http.Header.
func CloneHeader(in http.Header) http.Header {
	out := make(http.Header, len(in))
	for key, values := range in {
		newValues := make([]string, len(values))
		copy(newValues, values)
		out[key] = newValues
	}
	return out
}
