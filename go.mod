module github.com/runtimeracer/kajitool

go 1.16

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/runtimeracer/go-graphql-client v0.2.3
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.1
	golang.org/x/oauth2 v0.0.0-20210402161424-2e8d93401602 // indirect
)

replace github.com/runtimeracer/go-graphql-client v0.2.3 => ../go-graphql-client
