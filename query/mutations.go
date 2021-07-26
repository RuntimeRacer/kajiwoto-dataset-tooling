package query

// GraphQL for requests
type kajiwotoLoginUserPWMutation struct {
	Login   Login `graphql:"login (usernameOrEmail: $usernameOrEmail, password: $password, deviceType: WEB)"`
	Welcome Welcome
}

type kajiwotoLoginAuthTokenMutation struct {
	Login   Login `graphql:"loginWithToken (authToken: $authToken, action: $action, deviceType: WEB)"`
	Welcome Welcome
}
