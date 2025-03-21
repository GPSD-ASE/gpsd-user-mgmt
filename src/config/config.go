package config

import (
	"context"
	"fmt"
	"log"
	"os"

	vault "github.com/hashicorp/vault/api"
)

var (
	USER_MGMT_ENV      string
	USER_MGMT_DB_HOST  string
	USER_MGMT_DB_NAME  string
	USER_MGMT_DB_PORT  string
	USER_MGMT_DB_PASS  string
	USER_MGMT_DB_USER  string
	USER_MGMT_APP_PORT string
)

func LoadConfig() {
	if os.Getenv("VAULT_TOKEN") != "" {
		vaultSecrets, err := loadVaultSecrets()
		if err != nil {
			log.Printf("Error loading secrets from Vault: %v", err)
		} else {
			USER_MGMT_ENV = getString(vaultSecrets, "USER_MGMT_ENV", os.Getenv("USER_MGMT_ENV"))
			USER_MGMT_DB_HOST = getString(vaultSecrets, "USER_MGMT_DB_HOST", os.Getenv("USER_MGMT_DB_HOST"))
			USER_MGMT_DB_NAME = getString(vaultSecrets, "USER_MGMT_DB_NAME", os.Getenv("USER_MGMT_DB_NAME"))
			USER_MGMT_DB_PORT = getString(vaultSecrets, "USER_MGMT_DB_PORT", os.Getenv("USER_MGMT_DB_PORT"))
			USER_MGMT_DB_PASS = getString(vaultSecrets, "USER_MGMT_DB_PASS", os.Getenv("USER_MGMT_DB_PASS"))
			USER_MGMT_DB_USER = getString(vaultSecrets, "USER_MGMT_DB_USER", os.Getenv("USER_MGMT_DB_USER"))
			USER_MGMT_APP_PORT = getString(vaultSecrets, "USER_MGMT_APP_PORT", os.Getenv("USER_MGMT_APP_PORT"))
		}
		log.Printf("DEBUG - All vault secrets : %v", vaultSecrets)
	}
	if USER_MGMT_DB_HOST == "" {
		USER_MGMT_DB_HOST = os.Getenv("USER_MGMT_DB_HOST")
	}
	if USER_MGMT_DB_NAME == "" {
		USER_MGMT_DB_NAME = os.Getenv("USER_MGMT_DB_NAME")
	}
	if USER_MGMT_DB_PORT == "" {
		USER_MGMT_DB_PORT = os.Getenv("USER_MGMT_DB_PORT")
	}
	if USER_MGMT_DB_PASS == "" {
		USER_MGMT_DB_PASS = os.Getenv("USER_MGMT_DB_PASS")
	}
	if USER_MGMT_DB_USER == "" {
		USER_MGMT_DB_USER = os.Getenv("USER_MGMT_DB_USER")
	}
	if USER_MGMT_ENV == "" {
		USER_MGMT_ENV = os.Getenv("USER_MGMT_ENV")
	}
	if USER_MGMT_APP_PORT == "" {
		USER_MGMT_APP_PORT = os.Getenv("USER_MGMT_APP_PORT")
	}
	if USER_MGMT_APP_PORT == "" || USER_MGMT_ENV == "" || USER_MGMT_DB_HOST == "" || USER_MGMT_DB_NAME == "" || USER_MGMT_DB_PORT == "" || USER_MGMT_DB_PASS == "" || USER_MGMT_DB_USER == "" {
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
