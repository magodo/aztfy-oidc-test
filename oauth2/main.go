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
	EnvResourceId = "TEST_RESOURCE_ID"
	EnvAPIVersion = "TEST_API_VERSION"
)

func main() {
	ctx := context.Background()
	conf := externalaccount.Config{
		Audience:         "api://AzureADTokenExchange",
		SubjectTokenType: "urn:ietf:params:oauth:token-type:jwt",
		TokenURL:         "https://login.microsoftonline.com/oauth2/v2.0/token",
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
	resourceId, ok := os.LookupEnv(EnvResourceId)
	if !ok {
		log.Fatalf("%q not defined", EnvResourceId)
	}
	apiVersion, ok := os.LookupEnv(EnvAPIVersion)
	if !ok {
		log.Fatalf("%q not defined", EnvAPIVersion)
	}
	resp, err := c.Get(fmt.Sprintf("https://management.azure.com%s?api-version=%s", resourceId, apiVersion))
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
