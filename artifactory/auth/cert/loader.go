package cert

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/jfrog/jfrog-client-go/utils/errorutils"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func loadCertificates(caCertPool *x509.CertPool, certificatesDirPath string) error {
	if !fileutils.IsPathExists(certificatesDirPath, false) {
		return nil
	}
	files, err := ioutil.ReadDir(certificatesDirPath)
	if err != nil {
		return errorutils.WrapError(err)
	}
	for _, file := range files {
		caCert, err := ioutil.ReadFile(filepath.Join(certificatesDirPath, file.Name()))
		if err != nil {
			return errorutils.WrapError(err)
		}
		caCertPool.AppendCertsFromPEM(caCert)
	}
	return nil
}

func GetTransportWithLoadedCert(certificatesDirPath string, insecureTls bool, transport *http.Transport) (*http.Transport, error) {
	// Remove once SystemCertPool supports windows
	caCertPool, err := loadSystemRoots()
	if err != nil {
		return nil, errorutils.WrapError(err)
	}
	err = loadCertificates(caCertPool, certificatesDirPath)
	if err != nil {
		return nil, err
	}
	// Setup HTTPS client
	tlsConfig := &tls.Config{
		RootCAs:            caCertPool,
		ClientSessionCache: tls.NewLRUClientSessionCache(1),
		InsecureSkipVerify: insecureTls,
	}
	tlsConfig.BuildNameToCertificate()
	transport.TLSClientConfig = tlsConfig

	return transport, nil
}
