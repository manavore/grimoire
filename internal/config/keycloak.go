package config

import (
	"fmt"
	"os"
)

type KeycloakConfig struct {
	ServerURL     string
	Realm         string
	ClientID      string
	ClientSecret  string
	RedirectURL   string
	PostLogoutURL string
}

func LoadKeycloakConfig() (*KeycloakConfig, error) {
	requiredEndVars := []string{
		"KEYCLOAK_SERVER_URL",
		"KEYCLOAK_REAL",
		"KEYCLOAK_CLIENT_ID",
	}

	for _, envVar := range requiredEndVars {
		if os.Getenv(envVar) == "" {
			return nil, fmt.Errorf("Environment variable %s not definded", envVar)
		}
	}

	config := &KeycloakConfig{
		ServerURL:     os.Getenv("KEYCLOAK_SERVER_URL"),
		Realm:         os.Getenv("KEYCLOAK_REALM"),
		ClientID:      os.Getenv("KEYCLOAD_CLIENT_ID"),
		ClientSecret:  os.Getenv("KEYCLOAK_CLIENT_SECRET"),
		RedirectURL:   os.Getenv("KEYCLOAK_REDIRECT_URL"),
		PostLogoutURL: os.Getenv("KEYCLOAL_POST_LOGOUT_URL"),
	}

	return config, nil
}

func (kc *KeycloakConfig) GetWellKnownEndpoint() string {
	return fmt.Sprintf("%s/realms/%s/.well-known/openid-configuration", kc.ServerURL, kc.Realm)
}

func (kc *KeycloakConfig) GetAuthURL() string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/auth", kc.ServerURL, kc.Realm)
}

func (kc *KeycloakConfig) GetTokenURL() string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", kc.ServerURL, kc.Realm)
}

func (kc *KeycloakConfig) GetLogoutURL() string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/logout", kc.ServerURL, kc.Realm)
}

func (kc *KeycloakConfig) GetUserInfoURL() string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/userinfo", kc.ServerURL, kc.Realm)
}
