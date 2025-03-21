package config

import (
	"context"
	"fmt"
	"log"
	"os"

	vault "github.com/hashicorp/vault/api"
)

var (
	ENV     string
	DB_HOST string
	DB_NAME string
	DB_PORT string
	DB_PASS string
	DB_USER string
)

func LoadConfig() {
	if os.Getenv("VAULT_TOKEN") != "" {
		vaultSecrets, err := loadVaultSecrets()
		if err != nil {
			log.Printf("Error loading secrets from Vault: %v", err)
		} else {
			DB_HOST = getString(vaultSecrets, "DB_HOST", os.Getenv("DB_HOST"))
			DB_NAME = getString(vaultSecrets, "DB_NAME", os.Getenv("DB_NAME"))
			DB_PORT = getString(vaultSecrets, "DB_PORT", os.Getenv("DB_PORT"))
			DB_PASS = getString(vaultSecrets, "DB_PASS", os.Getenv("DB_PASS"))
			DB_USER = getString(vaultSecrets, "DB_USER", os.Getenv("DB_USER"))
			ENV = getString(vaultSecrets, "ENV", os.Getenv("ENV"))
		}
		log.Printf("DEBUG - All vault secrets : %v", vaultSecrets)
	}
	if DB_HOST == "" {
		DB_HOST = os.Getenv("DB_HOST")
	}
	if DB_NAME == "" {
		DB_NAME = os.Getenv("DB_NAME")
	}
	if DB_PORT == "" {
		DB_PORT = os.Getenv("DB_PORT")
	}
	if DB_PASS == "" {
		DB_PASS = os.Getenv("DB_PASS")
	}
	if DB_USER == "" {
		DB_USER = os.Getenv("DB_USER")
	}
	if ENV == "" {
		ENV = os.Getenv("ENV")
	}
	if ENV == "" || DB_HOST == "" || DB_NAME == "" || DB_PORT == "" || DB_PASS == "" || DB_USER == "" {
		log.Fatal("Missing environment variables")
	}
}

func loadVaultSecrets() (map[string]interface{}, error) {
	config := vault.DefaultConfig()
	if addr := os.Getenv("VAULT_ADDR"); addr != "" {
		config.Address = addr
	}
	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	token := os.Getenv("VAULT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("missing vault_token")
	}
	client.SetToken(token)

	secret, err := client.KVv2("secret").Get(context.Background(), "gpsd/user-mgmt")
	if err != nil {
		return nil, fmt.Errorf("failed to read secret from vault: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no secret found in vault")
	}
	return secret.Data, nil
}

func getString(secrets map[string]interface{}, key, def string) string {
	if val, ok := secrets[key].(string); ok && val != "" {
		return val
	}
	return def
}
