package config

import (
	"context"
	"fmt"
	"log"
	"os"

	vault "github.com/hashicorp/vault/api"
)

var (
	ENV             string
	DB_URL          string
	DB_HOST         string
	DB_PORT         string
	DB_NAME         string
	DB_PASS         string
	DB_USERNAME     string
	JWT_SECRET      string
	GRAPHHOPPER_KEY string
	PORT            string
	TOMTOM_API_KEY  string
	GRAPHHOPPER_URL string
	TOMTOM_URL      string
)

func LoadConfig() {
	if os.Getenv("VAULT_TOKEN") != "" {
		vaultSecrets, err := loadVaultSecrets()
		if err != nil {
			log.Printf("Error loading secrets from Vault: %v", err)
		} else {
			DB_URL = getString(vaultSecrets, "DATABASE_URL", os.Getenv("DATABASE_URL"))
			JWT_SECRET = getString(vaultSecrets, "JWT_SECRET", os.Getenv("JWT_SECRET"))
			DB_USERNAME = getString(vaultSecrets, "DATABASE_USERNAME", os.Getenv("DATABASE_USERNAME"))
			DB_PASS = getString(vaultSecrets, "DATABASE_PASS", os.Getenv("DATABASE_PASS"))
			GRAPHHOPPER_KEY = getString(vaultSecrets, "GRAPHHOPPER_KEY", os.Getenv("GRAPHHOPPER_KEY"))
			PORT = getString(vaultSecrets, "PORT", os.Getenv("PORT"))
			TOMTOM_API_KEY = getString(vaultSecrets, "TOMTOM_API_KEY", os.Getenv("TOMTOM_API_KEY"))
			GRAPHHOPPER_URL = getString(vaultSecrets, "GRAPHHOPPER_URL", os.Getenv("GRAPHHOPPER_URL"))
			TOMTOM_URL = getString(vaultSecrets, "TOMTOM_URL", os.Getenv("TOMTOM_URL"))
		}
		log.Printf("DEUBG - All vault secrets : %v", vaultSecrets)
	}
	if DB_URL == "" {
		DB_URL = os.Getenv("DATABASE_URL")
	}
	if DB_PASS == "" {
		DB_PASS = os.Getenv("DATABASE_PASS")
	}
	if DB_USERNAME == "" {
		DB_USERNAME = os.Getenv("DATABASE_USERNAME")
	}
	if JWT_SECRET == "" {
		JWT_SECRET = os.Getenv("JWT_SECRET")
	}
	if GRAPHHOPPER_KEY == "" {
		GRAPHHOPPER_KEY = os.Getenv("GRAPHHOPPER_KEY")
	}
	if PORT == "" {
		PORT = os.Getenv("PORT")
	}
	if TOMTOM_API_KEY == "" {
		TOMTOM_API_KEY = os.Getenv("TOMTOM_API_KEY")
	}
	if GRAPHHOPPER_URL == "" {
		GRAPHHOPPER_URL = os.Getenv("GRAPHHOPPER_URL")
	}
	if TOMTOM_URL == "" {
		TOMTOM_URL = os.Getenv("TOMTOM_URL")
	}
	if ENV == "" {
		ENV = os.Getenv("ENV")
	}
	if ENV == "" || DB_URL == "" || DB_PASS == "" || DB_USERNAME == "" || JWT_SECRET == "" || GRAPHHOPPER_KEY == "" || PORT == "" || TOMTOM_API_KEY == "" || GRAPHHOPPER_URL == "" || TOMTOM_URL == "" {
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

	secret, err := client.KVv2("secret").Get(context.Background(), "gpsd/map-mgmt")
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
