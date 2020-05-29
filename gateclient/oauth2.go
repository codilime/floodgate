package gateclient

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/codilime/floodgate/config"
	configAuth "github.com/codilime/floodgate/config/auth"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
)

// OAuth2Authentication struct is used to authenticate using oauth2
type OAuth2Authentication struct {
	Config       *oauth2.Config
	CodeVerifier oauth2.AuthCodeOption
	Token        *oauth2.Token
	CachedToken  oauth2.Token
	Done         chan bool
}

// OAuth2Authenticate is used to authenticate using oauth2
func OAuth2Authenticate(floodgateConfig *config.Config) (*oauth2.Token, error) {
	oauth2Config := floodgateConfig.Auth.OAuth2

	auth := &OAuth2Authentication{}
	auth, err := auth.setup(&oauth2Config)
	if err != nil {
		return nil, err
	}

	if oauth2Config.CachedToken.AccessToken != "" {
		auth.CachedToken = oauth2Config.CachedToken
		err := auth.refreshToken()
		if err != nil {
			return nil, err
		}
	} else {
		err := auth.getToken()
		if err != nil {
			return nil, err
		}
	}

	floodgateConfig.Auth.OAuth2.CachedToken = *auth.Token
	err = config.SaveConfig(floodgateConfig)
	if err != nil {
		return nil, err
	}

	return auth.Token, nil
}

func (a *OAuth2Authentication) setup(oauth2Config *configAuth.OAuth2) (*OAuth2Authentication, error) {
	if oauth2Config.TokenURL == "" || oauth2Config.AuthURL == "" || len(oauth2Config.Scopes) == 0 {
		return nil, fmt.Errorf("incorrect oauth2 configuration")
	}

	auth := OAuth2Authentication{
		Config: &oauth2.Config{
			ClientID:     oauth2Config.ClientID,
			ClientSecret: oauth2Config.ClientSecret,
			Scopes:       oauth2Config.Scopes,
			Endpoint: oauth2.Endpoint{
				AuthURL:  oauth2Config.AuthURL,
				TokenURL: oauth2Config.TokenURL,
			},
			RedirectURL: "http://localhost:8085/callback",
		},
		Done: make(chan bool),
	}

	return &auth, nil
}

func (a *OAuth2Authentication) refreshToken() error {
	tokenSource := a.Config.TokenSource(context.Background(), &a.CachedToken)
	t, err := tokenSource.Token()
	if err != nil {
		return err
	}

	a.Token = t

	return nil
}

func (a *OAuth2Authentication) getToken() error {
	http.HandleFunc("/callback", a.httpCallback)
	go http.ListenAndServe(":8085", nil)

	url, err := a.generateAuthCodeUrl()
	if err != nil {
		return err
	}

	log.Infof("Go to url and authenticate:\n%s\n", url)

	<-a.Done

	log.Info("Successfully authenticated")

	return nil
}

func (a *OAuth2Authentication) generateToken(code string) error {
	token, err := a.Config.Exchange(context.Background(), code, a.CodeVerifier)
	if err != nil {
		return err
	}

	a.Token = token
	return nil
}

func (a *OAuth2Authentication) generateAuthCodeUrl() (string, error) {
	verifier, verifierCode, err := a.generateCodeVerifier()
	if err != nil {
		return "", err
	}

	a.CodeVerifier = oauth2.SetAuthURLParam("code_verifier", verifier)
	codeChallenge := oauth2.SetAuthURLParam("code_challenge", verifierCode)
	challengeMethod := oauth2.SetAuthURLParam("code_challenge_method", "S256")

	url := a.Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce, challengeMethod, codeChallenge)

	return url, nil
}

func (a *OAuth2Authentication) httpCallback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")

	err := a.generateToken(code)
	if err != nil {
		log.Fatalf("Callback error: %v", err)
	}

	a.Done <- true

	fmt.Fprintf(w, "You can go back to CLI")
}

func (a *OAuth2Authentication) generateCodeVerifier() (verifier string, code string, err error) {
	randomBytes := make([]byte, 64)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", err
	}

	verifier = base64.RawURLEncoding.EncodeToString(randomBytes)
	verifierHash := sha256.Sum256([]byte(verifier))
	code = base64.RawURLEncoding.EncodeToString(verifierHash[:]) // Slice for type conversion

	return verifier, code, nil
}
