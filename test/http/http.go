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
	key := "/opt/kubernetes/ssl/kubernetes.key"
	cert := "/opt/kubernetes/ssl/kubernetes.crt"
	ca := "/opt/kubernetes/ssl/ca.crt"

	clientCert, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		log.Fatal(err)
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(ca)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      caCertPool,
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	response, err := client.Get(url)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	fmt.Println(response)
	fmt.Println(response.Body)
	fmt.Println(response.Header.Get("Date"))

}
