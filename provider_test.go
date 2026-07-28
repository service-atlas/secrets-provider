package secretsprovider

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
)

func TestNewProvider_EnvProviderDefault(t *testing.T) {
	t.Setenv("SECRETS_PROVIDER", "")
	provider, err := NewProvider(t.Context())
	if err != nil {
		t.Fatalf("error creating provider. err: %v", err)
	}
	_, ok := provider.(*EnvProvider)
	if !ok {
		t.Errorf("expected *EnvProvider, got %T", provider)
	}
}

func TestNewProvider_AWS(t *testing.T) {
	t.Setenv("SECRETS_PROVIDER", "aws")
	// We use a mock loader to avoid side effects and ensure it returns a valid config
	oldLoader := loadAWSConfig
	loadAWSConfig = func(ctx context.Context, optFns ...func(*awsconfig.LoadOptions) error) (aws.Config, error) {
		return aws.Config{}, nil
	}
	defer func() { loadAWSConfig = oldLoader }()

	provider, err := NewProvider(t.Context())
	if err != nil {
		t.Fatalf("error creating provider. err: %v", err)
	}
	_, ok := provider.(*AWSProvider)
	if !ok {
		t.Errorf("expected *AWSProvider, got %T", provider)
	}
}

func TestNewProvider_ExplicitEnv(t *testing.T) {
	t.Setenv("SECRETS_PROVIDER", "env")
	provider, err := NewProvider(t.Context())
	if err != nil {
		t.Fatalf("error creating provider. err: %v", err)
	}
	_, ok := provider.(*EnvProvider)
	if !ok {
		t.Errorf("expected *EnvProvider, got %T", provider)
	}
}

func TestNewProvider_AWSError(t *testing.T) {
	t.Setenv("SECRETS_PROVIDER", "aws")
	oldLoader := loadAWSConfig
	loadAWSConfig = func(ctx context.Context, optFns ...func(*awsconfig.LoadOptions) error) (aws.Config, error) {
		return aws.Config{}, fmt.Errorf("forced error")
	}
	defer func() { loadAWSConfig = oldLoader }()

	_, err := NewProvider(t.Context())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "loading AWS config: forced error" {
		t.Errorf("expected error %q, got %q", "loading AWS config: forced error", err.Error())
	}
}
