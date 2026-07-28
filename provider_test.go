package secretsprovider

import (
	"context"
	"fmt"
	"testing"
)

type mockCache struct {
	err error
}

func (m *mockCache) GetSecretStringWithContext(_ context.Context, _ string) (string, error) {
	return "", m.err
}

func TestNewProvider_EnvProviderDefault(t *testing.T) {
	t.Setenv("SECRETS_PROVIDER", "")
	provider, err := NewProvider()
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
	// We use a mock to avoid side effects
	oldNewCache := newSecretCache
	newSecretCache = func() (secretCache, error) {
		return &mockCache{}, nil
	}
	defer func() { newSecretCache = oldNewCache }()

	provider, err := NewProvider()
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
	provider, err := NewProvider()
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
	oldNewCache := newSecretCache
	newSecretCache = func() (secretCache, error) {
		return nil, fmt.Errorf("forced error")
	}
	defer func() { newSecretCache = oldNewCache }()

	_, err := NewProvider()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "creating secret cache: forced error" {
		t.Errorf("expected error %q, got %q", "creating secret cache: forced error", err.Error())
	}
}
