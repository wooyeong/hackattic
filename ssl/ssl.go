package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"time"
)

const ACCESS_TOKEN = ""

const SAMPLE = `{}`
const USE_SAMPLE = true

func main() {
	problem := getProblem()

	if USE_SAMPLE {
		problem = []byte(SAMPLE)
	}

	var js map[string]interface{}
	err := json.Unmarshal(problem, &js)
	if err != nil {
		panic(err)
	}

	privateKey := js["private_key"].(string)
	requiredData := js["required_data"].(map[string]interface{})
	domain := requiredData["domain"].(string)
	serial := requiredData["serial_number"].(string)
	country := requiredData["country"].(string)

	//fmt.Println("pk", privateKey)

	fmt.Printf("what's CountryCode for '%s': ", country)
	fmt.Scanf("%s\n", &country)

	fmt.Printf("do %s, serial %s, cc %s\n", domain, serial, country)

	priv, pub := getKeyPairs(privateKey)

	cert, err := createCert(priv, pub, domain, serial, country)
	if err != nil {
		panic(err)
	}
	// debugging purpose
	//ioutil.WriteFile("created-cert.der", cert, 0644)

	submit(cert)
}

func getProblem() []byte {
	if USE_SAMPLE {
		return nil
	}
	fmt.Println("Getting problem from server")

	resp, err := http.Get("https://hackattic.com/challenges/tales_of_ssl/problem?access_token=" + ACCESS_TOKEN)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body
}

func submit(cert []byte) {
	fmt.Println("Sending cert to server")

	encoded := base64.StdEncoding.EncodeToString(cert)
	js := map[string]string{"certificate": encoded}

	jsonBuffer, err := json.Marshal(js)
	fmt.Println("json", string(jsonBuffer))

	if USE_SAMPLE {
		fmt.Println("exitting")
		os.Exit(0)
	}

	resp, err := http.Post("https://hackattic.com/challenges/tales_of_ssl/solve?access_token="+ACCESS_TOKEN, "application/json", bytes.NewReader(jsonBuffer))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("response", string(body[:]))

	return
}

func getKeyPairs(private string) (privKey *rsa.PrivateKey, pubKey crypto.PublicKey) {
	decoded, _ := base64.StdEncoding.DecodeString(private)
	privKey, _ = x509.ParsePKCS1PrivateKey([]byte(decoded))
	pubKey = privKey.Public()
	//fmt.Println("privKey", privKey)
	//fmt.Println("pubKey", pubKey)
	return
}

func createCert(privKey *rsa.PrivateKey, pubKey crypto.PublicKey, domain, serial, country string) (cert []byte, err error) {
	/* see https://golang.org/pkg/crypto/x509/#CreateCertificate
	AuthorityKeyId,
	BasicConstraintsValid
	DNSNames
	ExcludedDNSDomains
	ExtKeyUsage
	IsCA
	KeyUsage
	MaxPathLen
	MaxPathLenZero
	NotAfter
	NotBefore
	PermittedDNSDomains
	PermittedDNSDomainsCritical
	SerialNumber
	SignatureAlgorithm
	Subject
	SubjectKeyId
	UnknownExtKeyUsage
	*/
	var template x509.Certificate

	subject := pkix.Name{Country: []string{country}, CommonName: domain}

	//template.AuthorityKeyId
	template.BasicConstraintsValid = true
	template.DNSNames = []string{domain}
	//template.ExcludedDNSDomains
	//template.ExtKeyUsage
	template.IsCA = false
	template.KeyUsage = x509.KeyUsageCertSign
	template.MaxPathLen = -1
	template.MaxPathLenZero = false
	template.NotAfter = time.Now().AddDate(-1, 0, 0)
	template.NotBefore = time.Now().AddDate(1, 0, 0)
	//template.PermittedDNSDomains = []string{domain}
	//template.PermittedDNSDomainsCritical = false
	template.SerialNumber, _ = new(big.Int).SetString(serial, 0)
	template.SignatureAlgorithm = x509.SHA256WithRSA
	template.Subject = subject
	//template.SubjectKeyId
	//template.UnknownExtKeyUsage

	return x509.CreateCertificate(rand.Reader, &template, &template, pubKey, privKey)
}
