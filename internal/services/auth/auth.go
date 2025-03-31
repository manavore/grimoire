package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/manavore/grimoire/internal/config"
)

type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	IDToken          string `json:"id_token"`
	NotBeforePolicy  int    `json:"not-before-policy"`
	SessionState     string `json:"session_state"`
	Scope            string `json:"scope"`
}

type UserInfo struct {
	Sub               string `json:"sub"`
	EmailVerified     bool   `json:"email_verified"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	Email             string `json:"email"`
	RealmAccess       struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
}

type AuthService struct {
	config *config.KeycloakConfig
	client *http.Client
}

func NewAuthService(config *config.KeycloakConfig) *AuthService {
	return &AuthService{
		config: config,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (a *AuthService) GenerateAuthURL(state string) string {
	authURL, _ := url.Parse(a.config.GetAuthURL())

	q := authURL.Query()
	q.Add("client_id", a.config.ClientID)
	q.Add("response_type", "code")
	q.Add("redirect_uri", a.config.RedirectURL)
	q.Add("scope", "openid profile email")
	q.Add("state", state)
	authURL.RawQuery = q.Encode()

	return authURL.String()
}

func (a *AuthService) ExchangeCodeForToken(code string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", a.config.ClientID)
	if a.config.ClientSecret != "" {
		data.Set("client_secret", a.config.ClientSecret)
	}
	data.Set("redirect_uri", a.config.RedirectURL)

	req, err := http.NewRequest("POST", a.config.GetTokenURL(), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token exchange failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

func (a *AuthService) GetUserInfo(accessToken string) (*UserInfo, error) {
	req, err := http.NewRequest("GET", a.config.GetUserInfoURL(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("userinfo request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func (a *AuthService) VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Implémentation simplifiée - en production, vous devriez récupérer les clés publiques
		// du serveur Keycloak via l'endpoint JWKS
		return nil, errors.New("not implemented - fetch public key from Keycloak JWKS endpoint")
	})

	return token, err
}

func (a *AuthService) RefreshToken(refreshToken string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", a.config.ClientID)
	if a.config.ClientSecret != "" {
		data.Set("client_secret", a.config.ClientSecret)
	}

	req, err := http.NewRequest("POST", a.config.GetTokenURL(), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token refresh failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

func (a *AuthService) CreateLogoutURL(idToken string) string {
	logoutURL, _ := url.Parse(a.config.GetLogoutURL())
	q := logoutURL.Query()
	q.Add("id_token_hint", idToken)
	if a.config.PostLogoutURL != "" {
		q.Add("post_logout_redirect_uri", a.config.PostLogoutURL)
	}
	logoutURL.RawQuery = q.Encode()

	return logoutURL.String()
}
