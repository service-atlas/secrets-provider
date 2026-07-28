package secretsprovider

import (
	"context"
	"fmt"
	"os"
)

// EnvProvider retrieves secrets from environment variables instead of external secret management systems.
type EnvProvider struct{}

// GetSecret retrieves the value of a secret identified by `name` from environment variables or returns an error if not found.
func (p *EnvProvider) GetSecret(_ context.Context, name string) (string, error) {
	if val, ok := os.LookupEnv(name); ok {
		return val, nil
	}
	return "", fmt.Errorf("secret %q not found in environment", name)
}

// GetDatabaseInfo retrieves database connection information, including URL, username, and password, from environment secrets.
// Returns a DatabaseInfo struct and an error if any required secret is missing or could not be retrieved.
func (p *EnvProvider) GetDatabaseInfo(ctx context.Context) (DatabaseInfo, error) {
	url, err := p.GetSecret(ctx, "DB_URL")
	if err != nil {
		return DatabaseInfo{}, err
	}
	username, err := p.GetSecret(ctx, "DB_USERNAME")
	if err != nil {
		return DatabaseInfo{}, err
	}
	password, err := p.GetSecret(ctx, "DB_PASSWORD")
	if err != nil {
		return DatabaseInfo{}, err
	}
	return DatabaseInfo{URL: url, Username: username, Password: password}, nil
}
