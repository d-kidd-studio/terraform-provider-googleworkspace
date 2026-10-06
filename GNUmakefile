SWEEP?=us-central1
TEST?=$$(go list ./...)

default: build

build: fmtcheck
	go install

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -w -s ./internal/provider

# Currently required by tf-deploy compile
fmtcheck:
	@echo "==> Checking source code against gofmt..."
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

generate: build
	go generate  ./...

lint:
	@echo "==> Checking source code against linters..."
	@golangci-lint run ./internal/provider

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test ./internal/provider -v -sweep=$(SWEEP) -sweep-run=$(SWEEPARGS) -timeout 60m

.PHONY: test test-unit testacc testacc-gmail

# Run fast local unit tests. Acceptance tests are skipped by the SDK unless TF_ACC=1.
test: test-unit

test-unit: fmtcheck
	TF_ACC=0 go test -count=1 $(TESTARGS) -skip '^(TestAcc|TestDWD)' -timeout=30s $(TEST)

# Run acceptance tests against a real Google Workspace. Gmail mailbox tests are
# excluded; use testacc-gmail when GOOGLEWORKSPACE_TEST_GMAIL_USER is available.
testacc: fmtcheck
	TF_ACC=1 go test -count=1 $(TEST) -v $(TESTARGS) -run '^(TestAcc|TestDWD)' -skip '^TestAccResourceGmailSendAsAlias_.*$' -timeout 120m

# Run only the Gmail acceptance tests. These require a real licensed Gmail user.
testacc-gmail: fmtcheck
	TF_ACC=1 go test -count=1 $(TEST) -v $(TESTARGS) -run '^TestAccResourceGmailSendAsAlias_.*$$' -timeout 120m
