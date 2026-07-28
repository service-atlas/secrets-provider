# secrets-provider

A Go library for retrieving secrets and database configuration from various providers (Environment variables or AWS Secrets Manager).

## Features

- **Abstracted Provider Interface**: Switch between local environment and cloud providers without changing application logic.
- **In-Memory Caching (AWS)**: The AWS provider automatically caches secrets in memory to reduce latency and API costs.
- **Database Configuration Support**: Helper to retrieve and parse database connection details (URL, Username, Password).

## Installation

```bash
go get github.com/service-atlas/secrets-provider
```

## Usage

### 1. Initialize the Provider

The library uses the `SECRETS_PROVIDER` environment variable to determine which provider to use.

```go
import (
    "context"
    "log"
    "github.com/service-atlas/secrets-provider"
)

func main() {
    ctx := context.Background()
    provider, err := secretsprovider.NewProvider(ctx)
    if err != nil {
        log.Fatalf("failed to create provider: %v", err)
    }

    // Use the provider...
}
```

### 2. Retrieve a Simple Secret

```go
secret, err := provider.GetSecret(ctx, "MY_SECRET_NAME")
if err != nil {
    log.Fatalf("failed to get secret: %v", err)
}
fmt.Printf("Secret: %s\n", secret)
```

### 3. Retrieve Database Configuration

The `GetDatabaseInfo` method returns a `DatabaseInfo` struct containing `URL`, `Username`, and `Password`.

```go
dbInfo, err := provider.GetDatabaseInfo(ctx)
if err != nil {
    log.Fatalf("failed to get db info: %v", err)
}
fmt.Printf("Connecting to %s as %s\n", dbInfo.URL, dbInfo.Username)
```

## Configuration

### Environment Provider (Default)

Used when `SECRETS_PROVIDER` is empty or set to `env`.

- `GetSecret(name)`: Looks up the environment variable `name`.
- `GetDatabaseInfo()`: Looks up `DB_URL`, `DB_USERNAME`, and `DB_PASSWORD`.

### AWS Secrets Manager Provider

Used when `SECRETS_PROVIDER=aws`.

- **Caching**: Uses the [AWS Secrets Manager Caching Library](https://github.com/aws/aws-secretsmanager-caching-go) to store secrets in memory, reducing API calls and improving performance.
- `GetSecret(name)`: Fetches the secret with ID `name` from the cache (or AWS Secrets Manager if not cached).
- `GetDatabaseInfo()`: 
    1. Reads the environment variable `DATABASE_SECRET_NAME`.
    2. Fetches that secret from AWS (expected to be a JSON string).
    3. Unmarshals it into a `DatabaseInfo` struct.

**Expected JSON format for AWS Database Secret:**
```json
{
  "url": "bolt://localhost:7687",
  "username": "neo4j",
  "password": "password"
}
```

## Development

The project includes a `justfile` for common tasks:

- `just test`: Run short tests.
- `just test-full`: Run all tests (including AWS integration tests using Testcontainers).
- `just lint`: Run linter.