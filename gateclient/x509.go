package gateclient

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/codilime/floodgate/config"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"net/http"
)

// X509Authenticate is used to authenticate using x509
func X509Authenticate(httpClient *http.Client, floodgateConfig *config.Config) (*http.Client, error) {
	x509Config := floodgateConfig.Auth.X509

	var cert tls.Certificate
	var clientCA []byte
	var err error
	certPool := x509.NewCertPool()

	if x509Config.CertPath != "" || x509Config.KeyPath != "" {
		certFullPath, err := homedir.Expand(x509Config.CertPath)
		if err != nil {
			return nil, err
		}

		keyFullPath, err := homedir.Expand(x509Config.KeyPath)
		if err != nil {
			return nil, err
		}

		cert, err = tls.LoadX509KeyPair(certFullPath, keyFullPath)
		if err != nil {
			return nil, err
		}

		clientCA, err = ioutil.ReadFile(certFullPath)
		if err != nil {
			return nil, err
		}
	}

	if x509Config.Cert != "" || x509Config.Key != "" {
		clientCA = []byte(x509Config.Cert)

		cert, err = tls.X509KeyPair(clientCA, []byte(x509Config.Key))
		if err != nil {
			return nil, err
		}
	}

	ok := certPool.AppendCertsFromPEM(clientCA)
	if !ok {
		return nil, errors.New("certificate is not valid")
	}

	clientTransport := httpClient.Transport.(*http.Transport)
	clientTransport.TLSClientConfig.MinVersion = tls.VersionTLS12
	clientTransport.TLSClientConfig.PreferServerCipherSuites = true
	clientTransport.TLSClientConfig.Certificates = []tls.Certificate{cert}

	return httpClient, nil
}
