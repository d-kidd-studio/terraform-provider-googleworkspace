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

	creds, err := googleoauth.CredentialsFromJSONWithParams(context.Background(), []byte(credentials), googleoauth.CredentialsParams{
		Scopes: []string{
			"https://www.googleapis.com/auth/admin.directory.group",
		},
		Subject: impersonatedUser,
	})
	if err != nil {
		t.Fatalf("creating DWD credentials: %v", err)
	}

	if _, err := creds.TokenSource.Token(); err != nil {
		t.Fatalf("obtaining DWD token for %s: %v", impersonatedUser, err)
	}
}
