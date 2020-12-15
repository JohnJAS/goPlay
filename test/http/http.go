package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	server := "shcCDFRH75vm01-0.hpeswlab.net"
	port := "8443"

	url := fmt.Sprintf("https://%s:%s", server, port)

	// Load client cert
	key := "/opt/kubernetes/ssl/kubectl-kube-api-client.key"
	cert := "/opt/kubernetes/ssl/kubectl-kube-api-client.crt"

	clientCert, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		log.Fatal(err)
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile("/opt/kubernetes/ssl/ca.crt")
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{clientCert},
				RootCAs:      caCertPool,
			},
		},
	}

	response, err := client.Get(url)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	fmt.Println(response)

}
