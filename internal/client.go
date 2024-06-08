package internal

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pkg/errors"
)

func RunClient(caFile, certFile, keyFile, connectAddr string) {
	// Create a new HTTP client
	httpsClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	var err error
	httpsClient.Transport, err = NewTlsTransport(caFile, certFile, keyFile)
	if err != nil {
		log.Fatalf("failed to create TLS transport: %v", err)
	}

	// Set up the request
	endpoint := connectAddr + "/api/v1/example"
	urlObj, err := url.Parse(endpoint)
	if err != nil {
		log.Fatalf("failed to parse URL: %v", err)
	}
	queryParam := urlObj.Query()
	queryParam.Set("key", "value")
	queryParam.Set("time", time.Now().Format(time.RFC3339Nano))
	urlObj.RawQuery = queryParam.Encode()

	header := http.Header{}
	header.Set("Content-Type", "application/json")

	httpsRequest := &http.Request{
		Method: http.MethodGet,
		URL:    urlObj,
		Header: header,
	}

	for {
		// Send the request
		response, err := httpsClient.Do(httpsRequest)
		if err != nil {
			log.Fatalf("failed to send request: %v", err)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			log.Fatalf("unexpected status code: %d %s", response.StatusCode, response.Status)
		}

		// Read the response
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatalf("failed to read response body: %v", err)
		}
		log.Printf("response: %s", body)

		// Try to transfer the response to a json object
		var responseJson map[string]interface{}
		if err := json.Unmarshal(body, &responseJson); err != nil {
			log.Fatalf("failed to unmarshal response: %v", err)
		}
		log.Printf("response to json: %v", responseJson)

		// Sleep for 5 seconds
		time.Sleep(5 * time.Second)
	}
}

func NewTlsTransport(caFile, certFile, keyFile string) (*http.Transport, error) {
	// Load client certificate
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load client certificate")
	}

	// Load CA certificate
	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read CA certificate")
	}

	// Create a certificate pool and add CA certificate
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create a TLS configuration
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	return &http.Transport{
		TLSClientConfig: tlsConfig,
	}, nil
}
