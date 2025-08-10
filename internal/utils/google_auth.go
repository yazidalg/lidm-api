package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

type GoogleTokenInfo struct {
	Aud           string `json:"aud"`
	Azp           string `json:"azp"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Exp           string `json:"exp"`
	FamilyName    string `json:"family_name"`
	GivenName     string `json:"given_name"`
	Iat           string `json:"iat"`
	Iss           string `json:"iss"`
	Locale        string `json:"locale"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Sub           string `json:"sub"`
}

type GoogleAuth struct {
	httpClient *http.Client
	clientID   string
}

func NewGoogleAuth(clientID string) *GoogleAuth {
	return &GoogleAuth{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		clientID: clientID,
	}
}

// VerifyGoogleToken verifies Google ID token using Google's tokeninfo endpoint
func (g *GoogleAuth) VerifyGoogleToken(idToken string) (*GoogleUserInfo, error) {
	url := "https://oauth2.googleapis.com/tokeninfo?id_token=" + idToken

	resp, err := g.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid token, status: %d", resp.StatusCode)
	}

	var tokenInfo GoogleTokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&tokenInfo); err != nil {
		return nil, fmt.Errorf("failed to decode token info: %v", err)
	}

	// Validate audience (client ID) if provided
	if g.clientID != "" && tokenInfo.Aud != g.clientID {
		return nil, fmt.Errorf("invalid audience: expected %s, got %s", g.clientID, tokenInfo.Aud)
	}

	// Check if email is verified
	if tokenInfo.EmailVerified != "true" {
		return nil, fmt.Errorf("email not verified")
	}

	// Validate issuer
	if tokenInfo.Iss != "https://accounts.google.com" && tokenInfo.Iss != "accounts.google.com" {
		return nil, fmt.Errorf("invalid issuer: %s", tokenInfo.Iss)
	}

	// Convert to our UserInfo struct
	userInfo := &GoogleUserInfo{
		ID:            tokenInfo.Sub,
		Email:         tokenInfo.Email,
		VerifiedEmail: tokenInfo.EmailVerified == "true",
		Name:          tokenInfo.Name,
		GivenName:     tokenInfo.GivenName,
		FamilyName:    tokenInfo.FamilyName,
		Picture:       tokenInfo.Picture,
		Locale:        tokenInfo.Locale,
	}

	return userInfo, nil
}

// GetUserInfoFromAccessToken gets user info using access token
func (g *GoogleAuth) GetUserInfoFromAccessToken(accessToken string) (*GoogleUserInfo, error) {
	url := "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken

	resp, err := g.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info, status: %d", resp.StatusCode)
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %v", err)
	}

	return &userInfo, nil
}
