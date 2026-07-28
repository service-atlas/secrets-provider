package secretsprovider

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-secretsmanager-caching-go/v2/secretcache"
)

type DatabaseInfo struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Provider interface {
	GetSecret(ctx context.Context, name string) (string, error)
	GetDatabaseInfo(ctx context.Context) (DatabaseInfo, error)
}

var newSecretCache = func() (secretCache, error) {
	return secretcache.New()
}

func NewProvider(ctx context.Context) (Provider, error) {
	switch os.Getenv("SECRETS_PROVIDER") {
	case "aws":
		cache, err := newSecretCache()
		if err != nil {
			return nil, fmt.Errorf("creating secret cache: %w", err)
		}
		return &AWSProvider{cache: cache}, nil
	default:
		return &EnvProvider{}, nil
	}
}
