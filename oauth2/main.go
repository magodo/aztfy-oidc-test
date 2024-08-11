package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google/externalaccount"
)

const (
	EnvClientID       = "ARM_CLIENT_ID"
	EnvSubscriptionID = "ARM_SUBSCRIPTION_ID"
)

func main() {
	clientId, ok := os.LookupEnv(EnvClientID)
	if !ok {
		log.Fatalf("%q not defined", EnvClientID)
	}
	subId, ok := os.LookupEnv(EnvSubscriptionID)
	if !ok {
		log.Fatalf("%q not defined", EnvSubscriptionID)
	}
	ctx := context.Background()
	conf := externalaccount.Config{
		Audience:         "api://AzureADTokenExchange",
		SubjectTokenType: "urn:ietf:params:oauth:token-type:jwt",
		TokenURL:         "https://login.microsoftonline.com/oauth2/v2.0/token",
		ClientID:         clientId,
		CredentialSource: &externalaccount.CredentialSource{
			URL: os.Getenv("ACTIONS_ID_TOKEN_REQUEST_URL"),
			Headers: map[string]string{
				"Accept":        "application/json",
				"Authorization": fmt.Sprintf("Bearer %s", os.Getenv("ACTIONS_ID_TOKEN_REQUEST_TOKEN")),
				"Content-Type":  "application/x-www-form-urlencoded",
			},
			Format: externalaccount.Format{
				Type:                  "json",
				SubjectTokenFieldName: "value",
			},
		},
		Scopes: []string{"https://management.azure.com"},
	}

	ts, err := externalaccount.NewTokenSource(ctx, conf)
	if err != nil {
		log.Fatal(err)
	}

	c := oauth2.NewClient(ctx, ts)
	resp, err := c.Get(fmt.Sprintf("https://management.azure.com/subscriptions/%s/resourcegroups/magodo-test?api-version=2020-06-01", subId))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Status)
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
