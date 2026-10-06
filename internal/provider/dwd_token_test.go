// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package googleworkspace

import (
	"context"
	"os"
	"testing"

	googleoauth "golang.org/x/oauth2/google"
)

func TestDWDToken(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("set TF_ACC=1 to run the Domain-Wide Delegation token test")
	}

	credentials := os.Getenv("GOOGLEWORKSPACE_CREDENTIALS")
	if credentials == "" {
		t.Fatal("GOOGLEWORKSPACE_CREDENTIALS must be set")
	}

	impersonatedUser := os.Getenv("GOOGLEWORKSPACE_IMPERSONATED_USER_EMAIL")
	if impersonatedUser == "" {
		t.Fatal("GOOGLEWORKSPACE_IMPERSONATED_USER_EMAIL must be set")
	}

	contents, _, err := pathOrContents(credentials)
	if err != nil {
		t.Fatalf("reading GOOGLEWORKSPACE_CREDENTIALS: %v", err)
	}

	creds, err := googleoauth.CredentialsFromJSONWithTypeAndParams(
		context.Background(),
		[]byte(contents),
		googleoauth.ServiceAccount,
		googleoauth.CredentialsParams{
			Scopes: []string{
				"https://www.googleapis.com/auth/admin.directory.group",
			},
			Subject: impersonatedUser,
		},
	)
	if err != nil {
		t.Fatalf("creating DWD credentials: %v", err)
	}

	if _, err := creds.TokenSource.Token(); err != nil {
		t.Fatalf("obtaining DWD token for %s: %v", impersonatedUser, err)
	}
}

func TestDWDTokenDefaultScopes(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("set TF_ACC=1 to run the Domain-Wide Delegation token test")
	}

	credentials := os.Getenv("GOOGLEWORKSPACE_CREDENTIALS")
	if credentials == "" {
		t.Fatal("GOOGLEWORKSPACE_CREDENTIALS must be set")
	}

	impersonatedUser := os.Getenv("GOOGLEWORKSPACE_IMPERSONATED_USER_EMAIL")
	if impersonatedUser == "" {
		t.Fatal("GOOGLEWORKSPACE_IMPERSONATED_USER_EMAIL must be set")
	}

	contents, _, err := pathOrContents(credentials)
	if err != nil {
		t.Fatalf("reading GOOGLEWORKSPACE_CREDENTIALS: %v", err)
	}

	creds, err := googleoauth.CredentialsFromJSONWithTypeAndParams(
		context.Background(),
		[]byte(contents),
		googleoauth.ServiceAccount,
		googleoauth.CredentialsParams{
			Scopes:  DefaultClientScopes,
			Subject: impersonatedUser,
		},
	)
	if err != nil {
		t.Fatalf("creating DWD credentials with provider default scopes: %v", err)
	}

	if _, err := creds.TokenSource.Token(); err != nil {
		t.Fatalf("obtaining DWD token for %s with provider default scopes: %v", impersonatedUser, err)
	}
}

func TestDWDTokenDefaultScopesIndividually(t *testing.T) {
	if os.Getenv("TF_ACC") != "1" {
		t.Skip("set TF_ACC=1 to run the Domain-Wide Delegation token test")
	}

	credentials := os.Getenv("GOOGLEWORKSPACE_CREDENTIALS")
	if credentials == "" {
		t.Fatal("GOOGLEWORKSPACE_CREDENTIALS must be set")
	}

	impersonatedUser := os.Getenv("GOOGLEWORKSPACE_IMPERSONATED_USER_EMAIL")
	if impersonatedUser == "" {
		t.Fatal("GOOGLEWORKSPACE_IMPERSONATED_USER_EMAIL must be set")
	}

	contents, _, err := pathOrContents(credentials)
	if err != nil {
		t.Fatalf("reading GOOGLEWORKSPACE_CREDENTIALS: %v", err)
	}

	for _, scope := range DefaultClientScopes {
		t.Run(scope, func(t *testing.T) {
			creds, err := googleoauth.CredentialsFromJSONWithTypeAndParams(
				context.Background(),
				[]byte(contents),
				googleoauth.ServiceAccount,
				googleoauth.CredentialsParams{
					Scopes:  []string{scope},
					Subject: impersonatedUser,
				},
			)
			if err != nil {
				t.Fatalf("creating DWD credentials for scope %s: %v", scope, err)
			}

			if _, err := creds.TokenSource.Token(); err != nil {
				t.Fatalf("obtaining DWD token for scope %s: %v", scope, err)
			}
		})
	}
}
