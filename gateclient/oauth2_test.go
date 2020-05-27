package gateclient

import (
	"testing"
)

func TestOAuth2Authenticate_GenerateCodeIdentifier(t *testing.T) {
	oauth2 := OAuth2Authentication{}

	_, _, err := oauth2.generateCodeVerifier()
	if err != nil {
		t.Errorf("OAuth2Authenticate.generateCodeVerifier() error = %v, wantErr false", err)
	}
}
