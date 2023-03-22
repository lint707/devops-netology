package acctest

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"strings"
	"time"
)

const (
	pemBlockTypeCertificate        = `CERTIFICATE`
	pemBlockTypeRsaPrivateKey      = `RSA PRIVATE KEY`
	pemBlockTypePublicKey          = `PUBLIC KEY`
	pemBlockTypeCertificateRequest = `CERTIFICATE REQUEST`
)

var tlsX509CertificateSerialNumberLimit = new(big.Int).Lsh(big.NewInt(1), 128) //nolint:gomnd

// TLSRSAPrivateKeyPEM generates a RSA private key PEM string.
// Wrap with TLSPEMEscapeNewlines() to allow simple fmt.Sprintf()
// configurations such as: private_key_pem = "%[1]s"
func TLSRSAPrivateKeyPEM(bits int) string {
	key, err := rsa.GenerateKey(rand.Reader, bits)

	if err != nil {
		//lintignore:R009
		panic(err)
	}

	block := &pem.Block{
		Bytes: x509.MarshalPKCS1PrivateKey(key),
		Type:  pemBlockTypeRsaPrivateKey,
	}

	return string(pem.EncodeToMemory(block))
}

// TLSRSAPublicKeyPEM generates a RSA public key PEM string.
// Wrap with TLSPEMEscapeNewlines() to allow simple fmt.Sprintf()
// configurations such as: public_key_pem = "%[1]s"
func TLSRSAPublicKeyPEM(keyPem string) string {
	keyBlock, _ := pem.Decode([]byte(keyPem))

	key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)

	if err != nil {
		//lintignore:R009
		panic(err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)

	if err != nil {
		//lintignore:R009
		panic(err)
	}

	block := &pem.Block{
		Bytes: publicKeyBytes,
		Type:  pemBlockTypePublicKey,
	}

	return string(pem.EncodeToMemory(block))
}

// TLSRSAX509LocallySignedCertificatePEM generates a local CA x509 certificate PEM string.
// Wrap with TLSPEMEscapeNewlines() to allow simple fmt.Sprintf()
// configurations such as: certificate_pem = "%[1]s"
func TLSRSAX509LocallySignedCertificatePEM(caKeyPem, caCertificatePem, keyPem, commonName string) string {
	caCertificateBlock, _ := pem.Decode([]byte(caCertificatePem))

	caCertificate, err := x509.ParseCertificate(caCertificateBlock.Bytes)

	if err != nil {
		//lintignore:R009
		panic(err)
	}

	caKeyBlock, _ := pem.Decode([]byte(caKeyPem))

	caKey, err := x509.ParsePKCS1PrivateKey(caKeyBlock.Bytes)

	if err != nil {
		//lintignore:R009
		panic(err)
	}

	keyBlock, _ := pem.Decode([]byte(keyPem))

	key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)

	if err != nil {
		//lintignore:R009
		panic(err)
	}

	serialNumber, err := rand.Int(rand.Reader, tlsX509CertificateSerialNumberLimit)

	if err != nil {
		//lintignore:R009
		panic(err)
	}

	certificate := &x509.Certificate{
		BasicConstraintsValid: true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		NotAfter:              time.Now().Add(24 * time.Hour), //nolint:gomnd
		NotBefore:             time.Now(),
		SerialNumber:          serialNumber,
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"ACME Examples, Inc"},
		},
	}

	certificateBytes, err := x509.CreateCertificate(rand.Reader, certificate, caCertificate, &key.PublicKey, caKey)

	if err != nil {
		//lintignore:R009
		panic(err)
	}

	certificateBlock := &pem.Block{
		Bytes: certificateBytes,
		Type:  pemBlockTypeCertificate,
	}

	return string(pem.EncodeToMemory(certificateBlock))
}

// TLSRSAX509SelfSignedCACertificatePEM generates a x509 CA certificate PEM string.
// Wrap with TLSPEMEscapeNewlines() to allow simple fmt.Sprintf()
// configurations such as: root_certificate_pem = "%[1]s"
func TLSRSAX509SelfSignedCACertificatePEM(keyPem string) string {
	keyBlock, _ := pem.Decode([]byte(keyPem))

	key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)

	if err != nil {
		//lintignore:R009
		panic(err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)

	if err != nil {
		//lintignore:R009
		panic(err)
	}

	publicKeyBytesSha1 := sha1.Sum(publicKeyBytes)

	serialNumber, err := rand.Int(rand.Reader, tlsX509CertificateSerialNumberLimit)

	if err != nil {
		//lintignore:R009
		panic(err)
	}

	certificate := &x509.Certificate{
		BasicConstraintsValid: true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		NotAfter:              time.Now().Add(24 * time.Hour), //nolint:gomnd
		NotBefore:             time.Now(),
		SerialNumber:          serialNumber,
		Subject: pkix.Name{
			CommonName:   "ACME Root CA",
			Organization: []string{"ACME Examples, Inc"},
		},
		SubjectKeyId: publicKeyBytesSha1[:],
	}

	certificateBytes, err := x509.CreateCertificate(rand.Reader, certificate, certificate, &key.PublicKey, key)

	if err != nil {
		//lintignore:R009
		panic(err)
	}

	certificateBlock := &pem.Block{
		Bytes: certificateBytes,
		Type:  pemBlockTypeCertificate,
	}

	return string(pem.EncodeToMemory(certificateBlock))
}

// TLSRSAX509SelfSignedCACertificateForRolesAnywhereTrustAnchorPEM generates a x509 CA certificate PEM string.
// The CA certificate is suitable for use as an IAM RolesAnywhere Trust Anchor.
// See https://docs.aws.amazon.com/rolesanywhere/latest/userguide/trust-model.html#signature-verification.
// Wrap with TLSPEMEscapeNewlines() to allow simple fmt.Sprintf()
// configurations such as: root_certificate_pem = "%[1]s"
func TLSRSAX509SelfSignedCACertificateForRolesAnywhereTrustAnchorPEM(keyPem string) string {
	keyBlock, _ := pem.Decode([]byte(keyPem))

	key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)

	if err != nil {
		//lintignore:R009
		panic(err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)

	if err != nil {
		//lintignore:R009
		panic(err)
	}

	publicKeyBytesSha1 := sha1.Sum(publicKeyBytes)

	serialNumber, err := rand.Int(rand.Reader, tlsX509CertificateSerialNumberLimit)

	if err != nil {
		//lintignore:R009
		panic(err)
	}

	certificate := &x509.Certificate{
		BasicConstraintsValid: true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		NotAfter:              time.Now().Add(24 * time.Hour), //nolint:gomnd
		NotBefore:             time.Now(),
		SerialNumber:          serialNumber,
		SignatureAlgorithm:    x509.SHA256WithRSA,
		Subject: pkix.Name{
			CommonName:   "ACME Root CA",
			Organization: []string{"ACME Examples, Inc"},
		},
		SubjectKeyId: publicKeyBytesSha1[:],
	}

	certificateBytes, err := x509.CreateCertificate(rand.Reader, certificate, certificate, &key.PublicKey, key)

	if err != nil {
		//lintignore:R009
		panic(err)
	}

	certificateBlock := &pem.Block{
		Bytes: certificateBytes,
		Type:  pemBlockTypeCertificate,
	}

	return string(pem.EncodeToMemory(certificateBlock))
}

// TLSRSAX509SelfSignedCertificatePEM generates a x509 certificate PEM string.
// Wrap with TLSPEMEscapeNewlines() to allow simple fmt.Sprintf()
// configurations such as: private_key_pem = "%[1]s"
func TLSRSAX509SelfSignedCertificatePEM(keyPem, commonName string) string {
	keyBlock, _ := pem.Decode([]byte(keyPem))

	key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)

	if err != nil {
		//lintignore:R009
		panic(err)
	}

	serialNumber, err := rand.Int(rand.Reader, tlsX509CertificateSerialNumberLimit)

	if err != nil {
		//lintignore:R009
		panic(err)
	}

	certificate := &x509.Certificate{
		BasicConstraintsValid: true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		NotAfter:              time.Now().Add(24 * time.Hour), //nolint:gomnd
		NotBefore:             time.Now(),
		SerialNumber:          serialNumber,
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"ACME Examples, Inc"},
		},
	}

	certificateBytes, err := x509.CreateCertificate(rand.Reader, certificate, certificate, &key.PublicKey, key)

	if err != nil {
		//lintignore:R009
		panic(err)
	}

	certificateBlock := &pem.Block{
		Bytes: certificateBytes,
		Type:  pemBlockTypeCertificate,
	}

	return string(pem.EncodeToMemory(certificateBlock))
}

// TLSRSAX509CertificateRequestPEM generates a x509 certificate request PEM string
// and a RSA private key PEM string.
// Wrap with TLSPEMEscapeNewlines() to allow simple fmt.Sprintf()
// configurations such as: certificate_signing_request_pem = "%[1]s" private_key_pem = "%[2]s"
func TLSRSAX509CertificateRequestPEM(keyBits int, commonName string) (string, string) {
	keyBytes, err := rsa.GenerateKey(rand.Reader, keyBits)
	if err != nil {
		//lintignore:R009
		panic(err)
	}

	csr := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"ACME Examples, Inc"},
		},
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &csr, keyBytes)
	if err != nil {
		//lintignore:R009
		panic(err)
	}

	csrBlock := &pem.Block{
		Bytes: csrBytes,
		Type:  pemBlockTypeCertificateRequest,
	}

	keyBlock := &pem.Block{
		Bytes: x509.MarshalPKCS1PrivateKey(keyBytes),
		Type:  pemBlockTypeRsaPrivateKey,
	}

	return string(pem.EncodeToMemory(csrBlock)), string(pem.EncodeToMemory(keyBlock))
}

func TLSPEMEscapeNewlines(pem string) string {
	return strings.ReplaceAll(pem, "\n", "\\n")
}
