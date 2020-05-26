package gateclient

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/codilime/floodgate/config"
	"io/ioutil"
	"net/http"
)

// x509Authenticate is used to authenticate using x509
func x509Authenticate(httpClient *http.Client, floodgateConfig *config.Config) (*http.Client, error) {
	x509Config := floodgateConfig.Auth.X509

	if x509Config.CertPath == "" || x509Config.KeyPath == "" {
		return nil, fmt.Errorf("incorrect x509 configuration")
	}

	cert, err := tls.LoadX509KeyPair(x509Config.CertPath, x509Config.KeyPath)
	if err != nil {
		return nil, err
	}

	clientCA, err := ioutil.ReadFile(x509Config.CertPath)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
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
