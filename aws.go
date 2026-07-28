package secretsprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type secretCache interface {
	GetSecretStringWithContext(ctx context.Context, secretID string) (string, error)
}

// AWSProvider is a secrets provider that fetches and caches secrets from AWS Secrets Manager.
type AWSProvider struct {
	cache secretCache
}

// GetSecret retrieves the value of a secret by its name from AWS Secrets Manager or its cache.
// It returns an error if the secret is not found, has no string value, or another issue occurs.
func (p *AWSProvider) GetSecret(ctx context.Context, name string) (string, error) {
	out, err := p.cache.GetSecretStringWithContext(ctx, name)
	if err != nil {
		return "", fmt.Errorf("fetching secret %q from AWS: %w", name, err)
	}
	if out == "" {
		return "", fmt.Errorf("secret %q has no string value", name)
	}
	return out, nil
}

// GetDatabaseInfo retrieves database connection information from an AWS Secrets Manager secret.
// It looks up the secret name from the DATABASE_SECRET_NAME environment variable.
// Returns an error if the environment variable is not set, if the secret is not found, or if the secret is invalid JSON.
func (p *AWSProvider) GetDatabaseInfo(ctx context.Context) (DatabaseInfo, error) {
	secretName, found := os.LookupEnv("DATABASE_SECRET_NAME")
	if !found {
		return DatabaseInfo{}, fmt.Errorf("DATABASE_SECRET_NAME not set")
	}
	jsonSecret, err := p.GetSecret(ctx, secretName)
	if err != nil {
		return DatabaseInfo{}, err
	}
	secret := DatabaseInfo{}
	err = json.Unmarshal([]byte(jsonSecret), &secret)
	if err != nil {
		return DatabaseInfo{}, err
	}
	return secret, nil
}
