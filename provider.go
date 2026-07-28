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

	// GetSecret retrieves the secret value associated with the provided name from the underlying storage or environment.
	GetSecret(ctx context.Context, name string) (string, error)

	// GetDatabaseInfo retrieves database connection details such as URL, username, and password from the underlying provider.
	GetDatabaseInfo(ctx context.Context) (DatabaseInfo, error)
}

var newSecretCache = func() (secretCache, error) {
	return secretcache.New()
}

// NewProvider creates a new secrets Provider based on the SECRETS_PROVIDER environment variable.
// Returns an AWSProvider if SECRETS_PROVIDER is set to "aws", otherwise defaults to an EnvProvider.
func NewProvider() (Provider, error) {
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
