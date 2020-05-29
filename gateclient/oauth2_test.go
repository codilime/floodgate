package gateclient

import (
	"github.com/codilime/floodgate/config/auth"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newConf(url string) auth.OAuth2 {
	var oauth2Conf = auth.OAuth2{
		Enabled:      true,
		TokenURL:     url + "/token",
		AuthURL:      url + "/auth",
		ClientID:     "clientId",
		ClientSecret: "clientSecret",
		Scopes:       []string{"profile"},
	}

	return oauth2Conf
}

func TestOAuth2Authenticate_Setup(t *testing.T) {
	conf := newConf("server")

	oauth := &OAuth2Authentication{}
	_, err := oauth.setup(&conf)
	if err != nil {
		t.Errorf("OAuth2Authenticate.setup() error = %v, wantErr false", err)
	}
}

func TestOAuth2Authenticate_SetupIncorrect(t *testing.T) {
	conf := auth.OAuth2{}
	oauth := &OAuth2Authentication{}
	_, err := oauth.setup(&conf)

	if (err != nil) != true {
		t.Errorf("OAuth2Authenticate.setup() error = %v, wantErr true", err)
	}
}

func TestOAuth2Authenticate_GenerateCodeIdentifier(t *testing.T) {
	oauth := &OAuth2Authentication{}
	verifier, code, err := oauth.generateCodeVerifier()
	if err != nil {
		t.Errorf("OAuth2Authenticate.generateCodeVerifier() error = %v, wantErr false", err)
	}

	if verifier == "" {
		t.Errorf("OAuth2Authenticate.generateCodeVerifier() verifier should not be empty")
	}

	if code == "" {
		t.Errorf("OAuth2Authenticate.generateCodeVerifier() code should not be empty")
	}
}

func TestOAuth2Authenticate_GenerateAuthCodeUrl(t *testing.T) {
	conf := newConf("server")

	oauth := &OAuth2Authentication{}
	oauth, _ = oauth.setup(&conf)

	url, err := oauth.generateAuthCodeUrl()
	if err != nil {
		t.Errorf("OAuth2Authenticate.generateAuthCodeUrl() error = %v, wantErr false", err)
	}

	if url == "" {
		t.Errorf("OAuth2Authenticate.generateAuthCodeUrl() url should not be empty")
	}
}

func TestOAuth2Authenticate_GenerateToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != "/token" {
			t.Errorf("OAuth2Authenticate.GenerateToken() unexpected exchange request URL, %v is found.", r.URL)
		}

		w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		w.Write([]byte("access_token=ACCESS_TOKEN&scope=user&token_type=bearer")) //nolint:errcheck
	}))
	defer ts.Close()

	conf := newConf(ts.URL)

	oauth := &OAuth2Authentication{}
	oauth, _ = oauth.setup(&conf)
	_, _ = oauth.generateAuthCodeUrl()

	err := oauth.generateToken("test")
	if err != nil {
		t.Errorf("OAuth2Authenticate.GenerateToken() unexpected error %v", err)
	}

	if oauth.Token != nil && oauth.Token.AccessToken != "ACCESS_TOKEN" {
		t.Errorf("OAuth2Authenticate.GenerateToken() access token is invalid got %v want ACCESS_TOKEN", oauth.Token.AccessToken)
	}
}

func TestOAuth2Authenticate_RefreshToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != "/token" {
			t.Errorf("OAuth2Authenticate.RefreshToken() unexpected exchange request URL, %v is found.", r.URL)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"access_token":"NEW_ACCESS_TOKEN",  "scope": "user", "token_type": "bearer"}`)) //nolint:errcheck
	}))
	defer ts.Close()

	conf := newConf(ts.URL)

	oauth := &OAuth2Authentication{}
	oauth, _ = oauth.setup(&conf)

	oauth.CachedToken = oauth2.Token{
		AccessToken:  "ACCESS_TOKEN",
		RefreshToken: "REFRESH_TOKEN",
		TokenType:    "bearer",
		Expiry:       time.Now(),
	}

	err := oauth.refreshToken()
	if err != nil {
		t.Errorf("OAuth2Authenticate.RefreshToken() unexpected error %v", err)
	}

	expected := "NEW_ACCESS_TOKEN"
	if oauth.Token != nil && oauth.Token.AccessToken != expected {
		t.Errorf("OAuth2Authenticate.RefreshToken() access token is invalid got %v want %v", oauth.Token.AccessToken, expected)
	}
}
