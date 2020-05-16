package database

import (
	"context"
	"crypto"
	"crypto/x509"
	"math/big"
	"time"
)

const (
	DefaultPrivateKeyType = "ecdsa"
	DefaultPrivateKeyBits = 256
)

type CertificateAuthority interface {
	GetSignedCertificates(ctx context.Context) ([]*SignedCertificate, error)
	GetSignedCertificate(ctx context.Context, serial *big.Int) (*SignedCertificate, error)
	GetRevokedCertificates() []*RevokedCertificate
	NewClientCertificate(ctx context.Context, name, keyType string, keyBits int, password, comment string) (*SignedCertificate, error)
	NewAgentCertificate(ctx context.Context, name, comment string) (*SignedCertificate, error)
	NewServerCertificate(commonName string) (*x509.Certificate, crypto.PrivateKey, error)
	SignCertificateRequest(ctx context.Context, csr *x509.CertificateRequest, comment string, agent bool) (*SignedCertificate, error)
	Revoke(ctx context.Context, certificate *SignedCertificate) error
	WatchRevokeCertificate() chan *RevokedCertificate
}

type SignedCertificate struct {
	Certificate *x509.Certificate
	P12         []byte
	IssuedAt    time.Time
	Agent       bool
	Comment     string
}

type RevokedCertificate struct {
	CommonName   string
	SerialNumber *big.Int
	IssuedAt     time.Time
	RevokedAt    time.Time
	Agent        bool
	Comment      string
}
