// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package googleworkspace

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	googleoauth "golang.org/x/oauth2/google"

	"google.golang.org/api/iamcredentials/v1"
	"google.golang.org/api/option"
)

func TestAccConfigLoadAndValidate_credsFromEnv(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip(fmt.Sprintf("Network access not allowed; use TF_ACC=1 to enable"))
	}

	testAccPreCheck(t)

	creds := getTestCredsFromEnv()
	config := &apiClient{
		Credentials:           creds,
		Customer:              os.Getenv("GOOGLEWORKSPACE_CUSTOMER_ID"),
		ImpersonatedUserEmail: os.Getenv("GOOGLEWORKSPACE_IMPERSONATED_USER_EMAIL"),
	}

	diags := config.loadAndValidate(context.Background())
	err := checkDiags(diags)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	diags = checkValidCreds(config)
	err = checkDiags(diags)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}

func TestAccConfigLoadAndValidate_accessTokenInvalid(t *testing.T) {
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

func TestAccConfigLoadAndValidate_accessToken(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip(fmt.Sprintf("Network access not allowed; use TF_ACC=1 to enable"))
	}

	testAccPreCheck(t)

	creds := getTestCredsFromEnv()
	gcpConfig := &apiClient{
		Credentials:  creds,
		ClientScopes: []string{"https://www.googleapis.com/auth/cloud-platform"},
	}

	diags := gcpConfig.loadAndValidate(context.Background())
	err := checkDiags(diags)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	iamCredsService, err := iamcredentials.NewService(context.Background(), option.WithHTTPClient(gcpConfig.client))
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	serviceAccount := fmt.Sprintf("projects/-/serviceAccounts/%s", os.Getenv("GOOGLEWORKSPACE_IMPERSONATED_SERVICE_ACCOUNT"))
	tokenRequest := &iamcredentials.GenerateAccessTokenRequest{
		Lifetime: "300s",
		Scope:    []string{"https://www.googleapis.com/auth/cloud-platform"},
	}
	at, err := iamCredsService.Projects.ServiceAccounts.GenerateAccessToken(serviceAccount, tokenRequest).Do()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	config := &apiClient{
		AccessToken:           at.AccessToken,
		Customer:              os.Getenv("GOOGLEWORKSPACE_CUSTOMER_ID"),
		ImpersonatedUserEmail: os.Getenv("GOOGLEWORKSPACE_IMPERSONATED_USER_EMAIL"),
		ServiceAccount:        os.Getenv("GOOGLEWORKSPACE_IMPERSONATED_SERVICE_ACCOUNT"),
	}

	diags = config.loadAndValidate(context.Background())
	err = checkDiags(diags)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	diags = checkValidCreds(config)
	err = checkDiags(diags)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}

func TestAccConfigLoadAndValidate_accessTokenOnly(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip(fmt.Sprintf("Network access not allowed; use TF_ACC=1 to enable"))
	}

	testAccPreCheck(t)

	credsFile := getTestCredsFromEnv()

	contents, _, err := pathOrContents(credsFile)
	if err != nil {
		t.Fatalf("could not get credentials: %s", err.Error())
	}

	credParams := googleoauth.CredentialsParams{
		Scopes: []string{"https://www.googleapis.com/auth/admin.directory.group"},
	}

	creds, err := googleoauth.CredentialsFromJSONWithParams(context.Background(), []byte(contents), credParams)
	if err != nil {
		t.Fatalf("could not get oauth2 credentials: %s", err.Error())
	}

	at, err := creds.TokenSource.Token()
	if err != nil {
		t.Fatalf("could not get token from oauth2 credentials: %s", err.Error())
	}

	config := &apiClient{
		AccessToken: at.AccessToken,
		Customer:    os.Getenv("GOOGLEWORKSPACE_CUSTOMER_ID"),
	}

	diags := config.loadAndValidate(context.Background())
	err = checkDiags(diags)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	diags = checkValidCredsGroupAdmin(config)
	err = checkDiags(diags)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}

func checkValidCreds(config *apiClient) diag.Diagnostics {
	var diags diag.Diagnostics

	directoryService, diags := config.NewDirectoryService()
	if diags.HasError() {
		return diags
	}

	_, err := directoryService.Customers.Get(config.Customer).Do()
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func checkValidCredsGroupAdmin(config *apiClient) diag.Diagnostics {
	var diags diag.Diagnostics

	directoryService, diags := config.NewDirectoryService()
	if diags.HasError() {
		return diags
	}
	groupsService := directoryService.Groups
	_, err := groupsService.List().Customer(config.Customer).Do()
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
