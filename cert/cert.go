package main

import (
	// "crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"errors"
	"fmt"

	"io/ioutil"
)

var oidExtensionSubjectAltName = []int{2, 5, 29, 17}

func parseCert(crt, privateKey string) (*tls.Certificate, error) {
	var cert tls.Certificate
	//加载PEM格式证书到字节数组
	certPEMBlock, err := ioutil.ReadFile(crt)
	if err != nil {
		return nil, err
	}
	//获取下一个pem格式证书数据 -----BEGIN CERTIFICATE----- -----END CERTIFICATE-----
	certDERBlock, restPEMBlock := pem.Decode(certPEMBlock)
	if certDERBlock == nil {
		return nil, errors.New("failed to decode cert block")
	}
	//附加数字证书到返回
	cert.Certificate = append(cert.Certificate, certDERBlock.Bytes)

	//继续解析Certifacate Chan,这里要明白证书链的概念
	certDERBlockChain, _ := pem.Decode(restPEMBlock)
	if certDERBlockChain != nil {
		//追加证书链证书到返回
		cert.Certificate = append(cert.Certificate, certDERBlockChain.Bytes)
		fmt.Println("存在证书链")
	}

	//读取RSA私钥进文件到字节数组
	keyPEMBlock, err := ioutil.ReadFile(privateKey)
	if err != nil {
		return nil, err
	}

	//解码pem格式的私钥------BEGIN RSA PRIVATE KEY----- -----END RSA PRIVATE KEY-----
	keyDERBlock, _ := pem.Decode(keyPEMBlock)
	if keyDERBlock == nil {
		return nil, err
	}
	//打印出私钥类型
	fmt.Println(keyDERBlock.Type)
	fmt.Println(keyDERBlock.Headers)
	var key interface{}
	var errParsePK error
	if keyDERBlock.Type == "RSA PRIVATE KEY" {
		//RSA PKCS1
		key, errParsePK = x509.ParsePKCS1PrivateKey(keyDERBlock.Bytes)
	} else if keyDERBlock.Type == "PRIVATE KEY" {
		//pkcs8格式的私钥解析
		key, errParsePK = x509.ParsePKCS8PrivateKey(keyDERBlock.Bytes)
	}

	if errParsePK != nil {
		return nil, err
	} else {
		cert.PrivateKey = key
	}
	//第一个叶子证书就是我们https中使用的证书
	fmt.Println("Cert Type: " + certDERBlock.Type)
	fmt.Printf("Cert Headers: %v", certDERBlock.Headers)

	x509Cert, err := x509.ParseCertificate(certDERBlock.Bytes)
	if err != nil {
		fmt.Println("x509证书解析失败")
		return nil, err
	} else {
		switch x509Cert.PublicKeyAlgorithm {
		case x509.RSA:
			{
				fmt.Println("Plublic Key Algorithm:RSA")
			}
		case x509.DSA:
			{
				fmt.Println("Plublic Key Algorithm:DSA")
			}
		case x509.ECDSA:
			{
				fmt.Println("Plublic Key Algorithm:ECDSA")
			}
		case x509.UnknownPublicKeyAlgorithm:
			{
				fmt.Println("Plublic Key Algorithm:UnkNow")
			}
		}
	}

	hasSANExtension(certDERBlock.Type, x509Cert.Extensions)

	return &cert, nil
}

func oidInExtensions(oid asn1.ObjectIdentifier, extensions []pkix.Extension) bool {
	for _, e := range extensions {
		if e.Id.Equal(oid) {
			return true
		}
	}
	return false
}

func hasSANExtension(certType string, extensions []pkix.Extension) {
	if !oidInExtensions(oidExtensionSubjectAltName, extensions) {
		fmt.Println(fmt.Sprintf(WARN, certType))
	}
}

const (
	WARN   = "( WARNING: your %s certificates do not have SANs set )"
	CA     = "CA"
	SERVER = "Server"
)

func main() {
	//fmt.Println("---------pkcs8 private key ---------------")
	//parseCert("./server.crt", "pkcs8_server.key")
	fmt.Println("---------pkcs1 private key ---------------")
	_, err := parseCert("server.crt", "server.key")
	fmt.Println(err)
}
