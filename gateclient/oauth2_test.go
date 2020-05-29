package gateclient

import (
	"github.com/codilime/floodgate/config/auth"
	"testing"
)

var oauth2Conf = auth.OAuth2{
	Enabled:      true,
	TokenURL:     "http://localhost/token",
	AuthURL:      "http://localhost/auth",
	ClientID:     "clientId",
	ClientSecret: "clientSecret",
	Scopes:       []string{"profile"},
}

var oauth2Test = &OAuth2Authentication{}

func TestOAuth2Authenticate_Setup(t *testing.T) {
	oauth2, err := oauth2Test.setup(&oauth2Conf)
	if err != nil {
		t.Errorf("OAuth2Authenticate.setup() error = %v, wantErr false", err)
	}

	oauth2Test = oauth2
}

func TestOAuth2Authenticate_GenerateCodeIdentifier(t *testing.T) {
	_, _, err := oauth2Test.generateCodeVerifier()
	if err != nil {
		t.Errorf("OAuth2Authenticate.generateCodeVerifier() error = %v, wantErr false", err)
	}
}

func TestOAuth2Authenticate_GenerateAuthCodeUrl(t *testing.T) {
	_, err := oauth2Test.generateAuthCodeUrl()
	if err != nil {
		t.Errorf("OAuth2Authenticate.generateAuthCodeUrl() error = %v, wantErr false", err)
	}
}
