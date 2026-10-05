// SPDX-License-Identifier: MPL-2.0

package googleworkspace

import (
	"context"
	"os"
	"testing"
)

func TestConfigLoadAndValidate_credsInvalidJSON(t *testing.T) {
	config := &apiClient{
		Credentials:           "{this is not json}",
		ImpersonatedUserEmail: "my-fake-email@example.com",
	}

	diags := config.loadAndValidate(context.Background())
	if !diags.HasError() {
		t.Fatalf("expected error, but got nil")
	}
}

func TestConfigLoadAndValidate_credsJSON(t *testing.T) {
	contents, err := os.ReadFile(testFakeCredentialsPath)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	config := &apiClient{
		Credentials:           string(contents),
		ImpersonatedUserEmail: "my-fake-email@example.com",
	}

	diags := config.loadAndValidate(context.Background())
	err = checkDiags(diags)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}

func TestConfigLoadAndValidate_credsFromFile(t *testing.T) {
	config := &apiClient{
		Credentials:           testFakeCredentialsPath,
		ImpersonatedUserEmail: "my-fake-email@example.com",
	}

	diags := config.loadAndValidate(context.Background())
	err := checkDiags(diags)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}

func TestConfigLoadAndValidate_credsNoImpersonation(t *testing.T) {
	config := &apiClient{
		Credentials: testFakeCredentialsPath,
	}

	diags := config.loadAndValidate(context.Background())
	err := checkDiags(diags)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}

func TestConfigOauthScopes_custom(t *testing.T) {
	config := &apiClient{
		Credentials:           testFakeCredentialsPath,
		ClientScopes:          []string{"https://www.googleapis.com/auth/admin/directory"},
		ImpersonatedUserEmail: "my-fake-email@example.com",
	}

	diags := config.loadAndValidate(context.Background())
	err := checkDiags(diags)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	if len(config.ClientScopes) != 1 {
		t.Fatalf("expected 1 scope, got %d scopes: %v", len(config.ClientScopes), config.ClientScopes)
	}
	if config.ClientScopes[0] != "https://www.googleapis.com/auth/admin/directory" {
		t.Fatalf("expected scope to be %q, got %q", "https://www.googleapis.com/auth/admin/directory", config.ClientScopes[0])
	}
}

func TestConfigLoadAndValidate_accessTokenInvalid(t *testing.T) {
	config := &apiClient{
		AccessToken:           "abcdefghijklmnopqrstuvwxyz",
		Customer:              os.Getenv("GOOGLEWORKSPACE_CUSTOMER_ID"),
		ImpersonatedUserEmail: os.Getenv("GOOGLEWORKSPACE_IMPERSONATED_USER_EMAIL"),
		ClientScopes:          []string{"https://www.googleapis.com/auth/admin.directory.domain"},
	}

	config.loadAndValidate(context.Background())
	diags := checkValidCreds(config)
	err := checkDiags(diags)
	if err == nil {
		t.Fatalf("expected error, but got nil")
	}
}
