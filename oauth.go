package slack

import (
	"context"
	"errors"
	"net/url"
)

type OAuthResponseIncomingWebhook struct {
	URL              string `json:"url"`
	Channel          string `json:"channel"`
	ChannelID        string `json:"channel_id,omitempty"`
	ConfigurationURL string `json:"configuration_url"`
}

type OAuthResponseBot struct {
	BotUserID      string `json:"bot_user_id"`
	BotAccessToken string `json:"bot_access_token"`
}

type OAuthResponse struct {
	AccessToken     string                       `json:"access_token"`
	Scope           string                       `json:"scope"`
	TeamName        string                       `json:"team_name"`
	TeamID          string                       `json:"team_id"`
	IncomingWebhook OAuthResponseIncomingWebhook `json:"incoming_webhook"`
	Bot             OAuthResponseBot             `json:"bot"`
	UserID          string                       `json:"user_id,omitempty"`
	SlackResponse
}

type OAuthV2Response struct {
	Ok          bool   `json:"ok"`
  Error       string `json:"error"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	BotUserID   string `json:"bot_user_id"`
	AppID       string `json:"app_id"`
	Team        struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"team"`
	Enterprise struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"enterprise"`
	AuthedUser struct {
		ID          string `json:"id"`
		Scope       string `json:"scope"`
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	} `json:"authed_user"`
}

// GetOAuthToken retrieves an AccessToken
func GetOAuthToken(clientID, clientSecret, code, redirectURI string, debug bool) (accessToken string, scope string, err error) {
	return GetOAuthTokenContext(context.Background(), clientID, clientSecret, code, redirectURI, debug)
}

// GetOAuthTokenContext retrieves an AccessToken with a custom context
func GetOAuthTokenContext(ctx context.Context, clientID, clientSecret, code, redirectURI string, debug bool) (accessToken string, scope string, err error) {
	response, err := GetOAuthResponseContext(ctx, clientID, clientSecret, code, redirectURI, debug)
	if err != nil {
		return "", "", err
	}
	return response.AccessToken, response.Scope, nil
}

func GetOAuthResponse(clientID, clientSecret, code, redirectURI string, debug bool) (resp *OAuthResponse, err error) {
	return GetOAuthResponseContext(context.Background(), clientID, clientSecret, code, redirectURI, debug)
}

func GetOAuthResponseContext(ctx context.Context, clientID, clientSecret, code, redirectURI string, debug bool) (resp *OAuthResponse, err error) {
	values := url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"code":          {code},
		"redirect_uri":  {redirectURI},
	}
	response := &OAuthResponse{}
	err = post(ctx, "oauth.access", values, response, debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response, nil
}

func GetV2OAuthResponseContext(ctx context.Context, clientID, clientSecret, code, redirectURI string, debug bool) (resp *OAuthV2Response, err error) {
	values := url.Values{
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"code":          {code},
		"redirect_uri":  {redirectURI},
	}
	response := &OAuthV2Response{}
	err = post(ctx, "oauth.v2.access", values, response, debug)
	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, errors.New(response.Error)
	}
	return response, nil
}
