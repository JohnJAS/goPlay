package main

import (
	"crypto/x509"
	"encoding/pem"
)

func main() {
	const provided_CA = `
-----BEGIN CERTIFICATE-----
MIIFMjCCAxoCCQC3JYP/A1NRkjANBgkqhkiG9w0BAQsFADBbMQswCQYDVQQGEwJD
TjEVMBMGA1UEBwwMRGVmYXVsdCBDaXR5MRwwGgYDVQQKDBNEZWZhdWx0IENvbXBh
bnkgTHRkMRcwFQYDVQQDDA4xNi4xNTUuMTk4LjI1MjAeFw0yMTExMDgwNjAyNTha
Fw0yMTEyMDgwNjAyNThaMFsxCzAJBgNVBAYTAkNOMRUwEwYDVQQHDAxEZWZhdWx0
IENpdHkxHDAaBgNVBAoME0RlZmF1bHQgQ29tcGFueSBMdGQxFzAVBgNVBAMMDjE2
LjE1NS4xOTguMjUyMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEA0y0q
2Ij4fKqnVJkL0J4qfl95IcBsaWPIgD7Bxjk2s3P8NceDV0Lq3Mexmc6HWCcxG3yn
+dv5ZELYni9/USpKn+PTHV6S3/gNCIqk8/0sa6aPW0gXv6uezb5MRG1ktilOTD2r
sXufRUrJghjbs3iYThZxpbGhezDeHYX3ms5buHMpJUERkrOeYzxJ3xZXQ9HN9rF2
wuR+TaaPENABlhW716it5ZV7l95sC3v+LzJ3YMQ8ZNbPzxhYYDPidiRXdCMHnPMz
8htBwih5h5khZDE5/ppTgFUhVyG7XPGQHQLtKwIYsKbYh3j/m0XcBD+5F6wyQcx+
Ucw1PNu4mJ+WRzRyrJBEj5iVejg1mrH7XdiBhqP7z0CJb6vCOxsXBmNtC9vRG95X
dJ1XQ1JuoNnaywJ5Z5ktNNEbrTf6BsYEmQTddmnTopEPgBpo+KfY153IYzAZVONS
zK6HbheXuAUfbJq5zALywnpgy6T+W+YsLQNJ66AahJY7FxdKYK/fNUKkBNSDADjI
XsyeRI5fGJBp8zM8mpGXL+OiUAoNEhYgrJPem3DW+IrIDzLcsKwrgFqya3JiABX2
eO4plxxdro7tqs/5sMvQId/kAkJYdXQq4E9jrJVK4PO0MiEP9wrZg8kttrrWNqP1
naQeVag21l4EeUmv26GSobM9UBgx+21EVr4HiXkCAwEAATANBgkqhkiG9w0BAQsF
AAOCAgEAJ29BBdF7ldUClaA+oHmtKmduGDJZp/qrxsdwaA/MqlWpKCgz8sZSE6LH
TYZWcgeJU/R4gQ+7FEXJVMkt6WjqITUN8K9SYeN6oJ2mm4kkgLGEhfbIg7Bid9UD
GZFKUEb5gewX2eRrT/PUkZD0036OlgJE5gbHo4sTA4i+Dt2RXmvehjyK4kYQRerd
Mqf30wBs0JXuW+Ih4Aml3px2mRtkLa65kodOjTt9GqGbAgSlAeZUtDmzuXnhxZVI
D3l9a6vWltuQgSOmHdXwdoBIjFNyyftKF8QhuAGdK7AyiUMVPbcLqB2SqsPiRjcL
dEPwG6J2ZLv5Hcn1bbAHm1EzJ15MPuPUsQ+uJm4kEwXfJ6/r1yXDrvGhWRGZj7Ml
7tBtdDN2s9Sp/ImKZC06PcTjSVFJ06IFMNsF7LRHBIw0vcQN6PnSs2ogRjrq0yFK
4eZ7WnMQ/Y0gAJ3ZhCdcM0qdyqVQQHioZdhDBzI3uwKBLCX+kJM+8KUUT5EAPeSk
XQKpCf0BR46ZV0U6RQmg8lIdLFExEFd3nuBv+dbgN5pVmsLQh1HUth+OMUox0YLo
CmXKzmHlAxCdVqr85cN/98l7OgYH/AwGAWedhC23Gzni+KsvSTRGdMJnPzI9uqSe
OdtubtNbNegCnl6JyT4Ycvcbx7FD4/Ec3adKLM0ntSdiQ0Pe8r0=
-----END CERTIFICATE-----`
	// First, create the set of root certificates. For this example we only
	// have one. It's also possible to omit this in order to use the
	// default root set of the current operating system.
	roots := x509.NewCertPool()

	block, _ := pem.Decode([]byte(provided_CA))
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	opts := x509.VerifyOptions{
		Roots:   roots,
	}

	if _, err := cert.Verify(opts); err != nil {
		panic("failed to verify certificate: " + err.Error())
	}
}
