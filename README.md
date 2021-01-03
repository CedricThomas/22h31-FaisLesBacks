# FaisLesBacks

A note application backend

This application use [Auth0](https://auth0.com/) for the users management, [FCM](https://firebase.google.com/docs/cloud-messaging) for notifications and [Airtable](https://airtable.com) as a storage.

## Getting Started

```
go build ./cmd/server
./server
```

## Configuration

This application can be configured in the file `internal/pkg/config.go`:
```
type Config struct {
	Port         string   `env:"PORT" envDefault:"9090"`
	Certificate  string   `env:"CERTIFICATE,file" envDefault:"../../certificates/auth0.pem" json:"-"`
	Issuer       string   `env:"ISSUER" envDefault:"https://dev-dgoly5h6.eu.auth0.com/"`
	Audience     []string `env:"AUDIENCE" envDefault:"casseur_flutter"`
	ApiKey       string   `env:"API_KEY,required"` // AirTable API Key
	BaseID       string   `env:"BASE_ID,required"` // AirTable Base ID
	FcmServerKey string   `env:"FCM_SERVER_KEY,required"`
}
```
