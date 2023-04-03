package adp

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type ADPAuthenticationSystem interface {
	Authenticate() error
	NewHttpClient() (*http.Client, error)
	SetRequestAuthorizationHeader(*http.Request)
}

type ADPOAuthAuthentication struct {
	KeyFilePath         string
	CertificateFilePath string
	Credential          string
	Token               *ADPOAuthOutput
}

type ADPOAuthOutput struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func NewOAuthAuthenticationSystem(certificateFilePath, keyFilePath, credential string) (*ADPOAuthAuthentication, error) {
	err := []error{}

	if !fileExists(keyFilePath) {
		err = append(err, fmt.Errorf("key file path does not exists: %s", keyFilePath))
	}

	if !fileExists(certificateFilePath) {
		err = append(err, fmt.Errorf("certficatefile path does not exsits: %s", certificateFilePath))
	}

	if len(err) > 0 {
		return nil, errors.Join(err...)
	}

	return &ADPOAuthAuthentication{
		KeyFilePath:         keyFilePath,
		CertificateFilePath: certificateFilePath,
		Credential:          credential,
	}, nil
}

func (a *ADPOAuthAuthentication) Authenticate() error {
	client, err := a.NewHttpClient()
	if err != nil {
		return err
	}
	request, err := http.NewRequest(
		http.MethodPost,
		"https://accounts.adp.com/auth/oauth/v2/token?grant_type=client_credentials",
		nil,
	)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", a.Credential))

	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	if !isValidResponseStatusCode(resp) {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("unable to retrieve access token: %s", body)
	}

	token := ADPOAuthOutput{}
	if err = json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return err
	}
	a.Token = &token
	return nil
}

func (a *ADPOAuthAuthentication) SetRequestAuthorizationHeader(request *http.Request) {
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.Token.AccessToken))
}

func (a *ADPOAuthAuthentication) NewHttpClient() (*http.Client, error) {
	certificates, err := tls.LoadX509KeyPair(a.CertificateFilePath, a.KeyFilePath)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{certificates},
			},
		},
	}
	return client, nil
}
