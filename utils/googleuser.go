package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"user-service/config"
)

type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func GetGoogleUserInfo(code string) (*GoogleUser, error) {
	clientID := config.AppConfig.GoogleClientID
	clientSecret := config.AppConfig.GoogleClientSecret
	redirectURI := config.AppConfig.GoogleRedirectURL

	tokenURL := "https://oauth2.googleapis.com/token"
	resp, err := http.PostForm(tokenURL, url.Values{
		"code":          {code},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"redirect_uri":  {redirectURI},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %v", err)
	}
	defer resp.Body.Close()

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		IdToken     string `json:"id_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %v", err)
	}

	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + tokenResponse.AccessToken
	resp, err = http.Get(userInfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	userInfo, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user info: %v", err)
	}

	var googleUser GoogleUser
	if err := json.Unmarshal(userInfo, &googleUser); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %v", err)
	}

	return &googleUser, nil
}
