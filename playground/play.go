// Example application that uses all of the available API options.
package main

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

//Refer to : https://www.alvestrand.no/objectid/2.5.29.17.html
//
//OID value: 2.5.29.17
//
//OID description: id-ce-subjectAltName
//
//This extension contains one or more alternative names, using any of a variety of name forms, for the entity that is bound by the CA to the certified public key.
//
//This extension may, at the option of the certificate issuer, be either critical or non-critical. An implementation which supports this extension is not required to be able to process all name forms. If the extension is flagged critical, at least one of the name forms that is present shall be recognized and processed, otherwise the certificate shall be considered invalid.
var oidExtensionSubjectAltName = []int{2, 5, 29, 17}

type Certificate struct {
	*tls.Certificate
	x509cert *x509.Certificate
}

func (c *Certificate) ReadFromFile(pem string) (block []byte, err error) {
	return ioutil.ReadFile(pem)
}

func (c *Certificate) InitCert(certPEMBlock []byte) (err error) {
	//resolve the first leaf cert block  -----BEGIN CERTIFICATE-----  -----END CERTIFICATE-----
	certDERBlock, restPEMBlock := pem.Decode(certPEMBlock)
	if certDERBlock == nil {
		return errors.New("failed to decode certPEMBlock")
	}
	if c.Certificate == nil {
		c.Certificate = new(tls.Certificate)
	}
	//append cert block to cert instance
	c.Certificate.Certificate = append(c.Certificate.Certificate, certDERBlock.Bytes)
	//continue to resolve certificate chain
	certDERBlockChain, _ := pem.Decode(restPEMBlock)
	if certDERBlockChain != nil {
		//append cert chain to cert instance
		c.Certificate.Certificate = append(c.Certificate.Certificate, certDERBlockChain.Bytes)
	}
	c.x509cert, err = x509.ParseCertificate(certDERBlock.Bytes)

	return
}

func (c *Certificate) InitKey(keyPEMBlock []byte) (err error) {
	//resolve private key ------BEGIN RSA PRIVATE KEY----- -----END RSA PRIVATE KEY-----
	keyDERBlock, _ := pem.Decode(keyPEMBlock)
	if keyDERBlock == nil {
		return errors.New("failed to decode keyDERBlock")
	}

	var key interface{}
	var errParsePK error
	if keyDERBlock.Type == "RSA PRIVATE KEY" {
		//RSA PKCS1
		key, errParsePK = x509.ParsePKCS1PrivateKey(keyDERBlock.Bytes)
	} else if keyDERBlock.Type == "PRIVATE KEY" {
		//pkcs8
		key, errParsePK = x509.ParsePKCS8PrivateKey(keyDERBlock.Bytes)
	}
	if errParsePK != nil {
		return errParsePK
	} else {
		c.PrivateKey = key
	}
	return
}

func (c *Certificate) HasSANExtension() bool {
	return oidInExtensions(oidExtensionSubjectAltName, c.x509cert.Extensions)
}

func (c *Certificate) CheckHostnameInSAN(hostname string) (err error) {
	if c.HasSANExtension() {
		return c.x509cert.VerifyHostname(hostname)
	} else {
		return errors.New("No SAN")
	}
	return
}

func oidInExtensions(oid asn1.ObjectIdentifier, extensions []pkix.Extension) bool {
	for _, e := range extensions {
		if e.Id.Equal(oid) {
			return true
		}
	}
	return false
}
