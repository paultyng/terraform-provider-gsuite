package provider

import (
	"context"
	"io/ioutil"
	"testing"
)

const testFakeCredentialsPath = "./testdata/fake_account.json"

func TestConfigLoadAndValidate_accountFilePath(t *testing.T) {
	ctx := context.Background()

	config := config{
		Credentials: testFakeCredentialsPath,
	}

	err := config.loadAndValidate(ctx)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestConfigLoadAndValidate_accountFileJSON(t *testing.T) {
	ctx := context.Background()

	contents, err := ioutil.ReadFile(testFakeCredentialsPath)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	config := config{
		Credentials: string(contents),
	}

	err = config.loadAndValidate(ctx)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestConfigLoadAndValidate_accountFileJSONInvalid(t *testing.T) {
	ctx := context.Background()

	config := config{
		Credentials: "{this is not json}",
	}

	if config.loadAndValidate(ctx) == nil {
		t.Fatalf("expected error, but got nil")
	}
}
