package gateclient

import (
	"crypto/tls"
	"github.com/codilime/floodgate/config"
	"github.com/codilime/floodgate/config/auth"
	"net/http"
	"testing"
	"time"
)

func TestX509Authenticate(t *testing.T) {
	type args struct {
		config config.Config
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "invalid cert",
			args: args{
				config: config.Config{
					Auth: auth.Config{
						X509: auth.X509{
							Cert: "-----BEGIN CERTIFICATE-----\nHelloWorldCert==-----END RSA PRIVATE KEY-----",
							Key:  "-----BEGIN RSA PRIVATE KEY-----\nHelloWorldPrivate==\n-----END RSA PRIVATE KEY-----",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "cert not configured",
			args: args{
				config: config.Config{},
			},
			wantErr: true,
		},
		{
			name: "valid cert",
			args: args{
				config: config.Config{
					Auth: auth.Config{
						X509: auth.X509{
							Cert: "-----BEGIN CERTIFICATE-----\nMIIFazCCA1OgAwIBAgIUYQNloiIKBdaBUD8dvaUDbPDa8wMwDQYJKoZIhvcNAQEL\nBQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM\nGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yMDA1MjcxMDA0MDZaFw0yMTA1\nMjcxMDA0MDZaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw\nHwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwggIiMA0GCSqGSIb3DQEB\nAQUAA4ICDwAwggIKAoICAQCnlB10FqOQCDp+0IQ7vPeBXD3O2WnW4mQ8gUN0c0au\niF3Ldvh+qoeyoULzRV8TCG3hroIYZykHKINqJEaV5cPHtURcqGQ373MMMCKc9K6A\n4bIwDFEhPZHKl0HoOGx4ZkJRsMiZ0u7zbprz72gteXv3CLiuosXRLM1BTDKrNTiz\npXfmXonYkmuDjZ/nQ4HR9e96nH5lfNaXIzq63/BEocCXncR2AiBpGCuy7hxJZ5+v\nEGUeH7a4EENFrTKOIWDkq2dRbmySEtio1fXAFtKX/caryGo4YLt6hyT5H1Ds33N7\nQw5Szyt0ouFgCc/wk7aXVO9A4zm++1y7be7axI8EWhPdONhwxzaykiN1MY/LnyvZ\nsXuRKZWzwISYJ2kPAuM0QQwhRhnVL8oMrEpAJjL0gIfykAmKUAOTJ/cjQgHI2xIp\nr3pKVtjxhrp21hXxRLmmHV6pwrgJ0PVkoQHmbfk0vF51U9GP5BThOytX/LnaGZUc\nuTSH1OUKM6hvxqH3RZAmFENbR3gzwJqQ9yzxlj5YmbIO3JjC39o96OY0m3zx4FxL\noeyAn7mVB8Gf/voOy/0VsRc3SeSzn+lI4nXalY0s7vC0EqLCKF/FA4pYlWRZaFlz\nHk8nrD7g0ZnpcO+a8PPlGqdYCNeSEaViXJNHt03Wv/noNN98cgwZXGqz+sNnSRsd\n0wIDAQABo1MwUTAdBgNVHQ4EFgQUFfPP3nOeKy8V3VduLZkIOlJA6pkwHwYDVR0j\nBBgwFoAUFfPP3nOeKy8V3VduLZkIOlJA6pkwDwYDVR0TAQH/BAUwAwEB/zANBgkq\nhkiG9w0BAQsFAAOCAgEAgWHZN7dAfsr9EbBpxBrNIIkua8iJYcsAdPAwAxK6UV45\napPPsBc5czUcqt4WFSye8ghWooT77BOYR4qVMfa85Ten7esnKBP7cruXZVDotVyN\nFn7vn0b3e3pY1VEhEbxOvYOd8BSYDaKpKrSkgMjDQnoDAUdHF6AMcLIZb+l3OpsL\nBm7GPoNOuyKN8/YDQXZZAL/oGbQdjroNsiQOM/Oemk9VlFGkdXrgqhaOBqgDv+3O\n/pTn3553J/cPwsslygGWCLYHQeRmfNbxzOCW2gtSXvI+iGPJpb2md/TUrZPOIcKp\nQAD3/Q2srSJWKcQPsANWiA1wPYm4atitbEskGea+C1+YrjO7y6JD3sms71iTgGMA\nD1AkvBXb0UTDKEKyyGuWBNaiEQwuj4EsTzADs4IPH1jxIwPt2PKTtrFyUkPCBhuM\nN2l19HE19DRefxARKzaCN27ZbUKqyTrnnX5b/hOupPpi5cpB9f9mCKvhvCLv3XjA\noEcegBbhz4Sa9rYmDpSXSlJ8UJ+8FR/b9cscldG1F6DHWm48qJA1rRylejLLNAqy\nO3Lbzg/8wXWmSc5K+7y3phlpFPiY0sAyq93cVwKpZ2iEAIT/gIVQHmnws0pxbGAY\nmECwzPMrL+0+aTYDz4J/CjEpKisaooiXwTj8WtmNLFHoEiT0ngoi37JgNEpROUQ=\n-----END CERTIFICATE-----",
							Key:  "-----BEGIN PRIVATE KEY-----\n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\nrYHC77OD6rKzNPO1lwvG8JMrRsTgWBY=\n-----END PRIVATE KEY-----",
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var httpClient = &http.Client{
				Timeout: 5 * time.Second,
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: tt.args.config.Insecure},
				},
			}

			_, err := X509Authenticate(httpClient, &tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("X509Authenticate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
