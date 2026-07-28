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

type AWSProvider struct {
	cache secretCache
}

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
