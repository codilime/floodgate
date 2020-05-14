package gateclient

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/codilime/floodgate/config"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
)

var (
	conf *oauth2.Config
)

// Source: https://github.com/spinnaker/spin/blob/master/cmd/gateclient/client.go
func oAuth2Authenticate(floodgateConfig *config.Config) (*oauth2.Token, error) {
	oauth2Config := floodgateConfig.Auth.OAuth2

	if oauth2Config.TokenUrl == "" || oauth2Config.AuthUrl == "" || len(oauth2Config.Scopes) == 0 {
		return nil, fmt.Errorf("incorrect oauth2 configuration")
	}

	conf = &oauth2.Config{
		ClientID:     oauth2Config.ClientId,
		ClientSecret: oauth2Config.ClientSecret,
		Scopes:       oauth2Config.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  oauth2Config.AuthUrl,
			TokenURL: oauth2Config.TokenUrl,
		},
		RedirectURL: "http://localhost:8085",
	}

	var token *oauth2.Token

	if oauth2Config.CachedToken.Valid() {
		_, _ = oAuth2RefreshToken()
	} else {
		t, err := oAuth2GetToken()
		if err != nil {
			return nil, err
		}

		token = t
	}

	return token, nil
}

func oAuth2RefreshToken() (*oauth2.Token, error) {
	return nil, nil
}

func oAuth2GetToken() (*oauth2.Token, error) {
	//Setup HTTP server to get callback
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		fmt.Fprintln(w, code)
	})
	go http.ListenAndServe(":8085", nil)

	verifier, verifierCode, err := generateCodeVerifier()
	if err != nil {
		return nil, err
	}

	codeVerifier := oauth2.SetAuthURLParam("code_verifier", verifier)
	codeChallenge := oauth2.SetAuthURLParam("code_challenge", verifierCode)
	challengeMethod := oauth2.SetAuthURLParam("code_challenge_method", "S256")

	url := conf.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce, challengeMethod, codeChallenge)

	log.Infof("Go to url and authenticate:\n%s\n", url)
	log.Infof("Paste verification code: ")

	var code string
	fmt.Scanf("%s", &code)

	token, err := conf.Exchange(context.Background(), code, codeVerifier)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func generateCodeVerifier() (verifier string, code string, err error) {
	randomBytes := make([]byte, 64)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", err
	}

	verifier = base64.RawURLEncoding.EncodeToString(randomBytes)
	verifierHash := sha256.Sum256([]byte(verifier))
	code = base64.RawURLEncoding.EncodeToString(verifierHash[:]) // Slice for type conversion

	return verifier, code, nil
}
