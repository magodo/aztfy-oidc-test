package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google/externalaccount"
)

var (
	clientId = flag.String("client-id", "", "arm client id")
	subId    = flag.String("sub-id", "", "arm sub id")
)

func main() {
	flag.Parse()
	if *clientId == "" {
		log.Fatal("client id not specified")
	}
	if *subId == "" {
		log.Fatal("subscription id not specified")
	}
	ctx := context.Background()
	conf := externalaccount.Config{
		Audience:         "api://AzureADTokenExchange",
		SubjectTokenType: "urn:ietf:params:oauth:token-type:jwt",
		TokenURL:         "https://login.microsoftonline.com/oauth2/v2.0/token",
		ClientID:         *clientId,
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
	resp, err := c.Get(fmt.Sprintf("https://management.azure.com/subscriptions/%s/resourcegroups/aztfy?api-version=2020-06-01", *subId))
	log.Println(resp)
	if resp != nil {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(string(b))
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Status)
}
