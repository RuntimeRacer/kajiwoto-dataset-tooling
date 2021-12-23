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
	req2 := cloneRequest(req)
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
			base: http.DefaultTransport,
			headers: map[string]string{
				// Default Headers
				"Content-Type": "application/json",
			},
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
		return result, errLogin
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

func (c *kajiwotoClient) GetAITrainerGroup(aiTrainerGroupID, authToken string) (result AITrainerGroup, err error) {
	// Sanity check
	if authToken == "" {
		return result, fmt.Errorf("invalid auth token")
	}
	if aiTrainerGroupID == "" {
		return result, fmt.Errorf("invalid trainer group ID")
	}

	vars := map[string]interface{}{
		"aiTrainerGroupId": graphql.String(aiTrainerGroupID),
	}

	// Add Auth-Token header
	headers := map[string]string{
		"auth_token": authToken,
	}
	c.AddHeaders(headers)

	// Execute Query
	aiTrainerGroupResult := kajiwotoDatasetAITrainerGroupQuery{}
	if errLogin := c.performGraphQuery(vars, &aiTrainerGroupResult); errLogin != nil {
		return result, fmt.Errorf("unable to fetch AI trainer group, response: %q", errLogin)
	}

	// Build generic Result object
	result = aiTrainerGroupResult.AITrainerGroup
	return result, nil
}

func (c *kajiwotoClient) GetAITrainedList(aiTrainerGroupID, searchQuery, authToken string, limit, page int) (result []AITrained, err error) {
	// Sanity check
	if authToken == "" {
		return result, fmt.Errorf("invalid auth token")
	}
	if aiTrainerGroupID == "" {
		return result, fmt.Errorf("invalid trainer group ID")
	}
	if limit < 1 || limit > 100 {
		return result, fmt.Errorf("limit exceeds allowed range")
	}
	if page < 0 {
		return result, fmt.Errorf("page cannot be negative")
	}

	vars := map[string]interface{}{
		"aiTrainerGroupId": graphql.String(aiTrainerGroupID),
		"limit":            graphql.Int(limit),
		"page":             graphql.Int(page),
		"searchQuery":      graphql.String(searchQuery),
	}

	// Add Auth-Token header
	headers := map[string]string{
		"auth_token": authToken,
	}
	c.AddHeaders(headers)

	// Execute Query
	aiTrainedListResult := kajiwotoDatasetAITrainedListQuery{}
	if errLogin := c.performGraphQuery(vars, &aiTrainedListResult); errLogin != nil {
		return result, fmt.Errorf("unable to fetch AI trainer group, response: %q", errLogin)
	}

	// Build generic Result object
	result = aiTrainedListResult.AITrainedList
	return result, nil
}

// DoTrainDataset performs login via session key if available
func (c *kajiwotoClient) DoTrainDataset(aiTrainerGroupID, authToken string, training []AITraining) (result TrainDatasetResult, err error) {
	// Sanity check
	if authToken == "" {
		return result, fmt.Errorf("invalid login credentials")
	}

	// Convert AITraining data to form data
	questions := make([]graphql.String, 0)
	form := make([][]graphql.String, 0)
	for _, item := range training {
		questions = append(questions, item.UserMessage)
		form = append(form, []graphql.String{item.Condition, item.Message})
	}

	vars := map[string]interface{}{
		"aiTrainerGroupId": graphql.String(aiTrainerGroupID),
		"questions":        questions,
		"form":             form,
		"editorType":       graphql.String("web-list"),
		"detailed":         graphql.Boolean(true),
		"multi":            graphql.Boolean(false), // FIXME: Does not seem to be used for web-ui. Maybe for batch?
	}

	// Add Auth-Token header
	headers := map[string]string{
		"auth_token": authToken,
	}
	c.AddHeaders(headers)

	trainingResult := kajiwotoDatasetTrainDatasetMutation{}
	if errTrain := c.performGraphMutation(vars, &trainingResult); errTrain != nil {
		return result, fmt.Errorf("unable to train dataset, response: %q", errTrain)
	}

	// Build generic Result object
	result = trainingResult.TrainDataset
	return result, nil
}

func (c *kajiwotoClient) performGraphMutation(vars map[string]interface{}, mutation interface{}) error {
	return c.client.Mutate(context.Background(), mutation, vars)
}

func (c *kajiwotoClient) performGraphQuery(vars map[string]interface{}, query interface{}) error {
	return c.client.Query(context.Background(), query, vars)
}

// cloneRequest creates a shallow copy of the request along with a deep copy of the Headers.
func cloneRequest(req *http.Request) *http.Request {
	r := new(http.Request)

	// shallow clone
	*r = *req

	// deep copy headers
	r.Header = cloneHeader(req.Header)

	return r
}

// cloneHeader creates a deep copy of an http.Header.
func cloneHeader(in http.Header) http.Header {
	out := make(http.Header, len(in))
	for key, values := range in {
		newValues := make([]string, len(values))
		copy(newValues, values)
		out[key] = newValues
	}
	return out
}
