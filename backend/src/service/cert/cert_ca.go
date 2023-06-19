// indece Monitor
// Copyright (C) 2023 indece UG (haftungsbeschränkt)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License or any
// later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package cert

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"time"

	"gopkg.in/guregu/null.v4"
)

func (s *Service) GenerateCACert() (*PEMCert, error) {
	serialNumber, err := s.generateSerialNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %s", err)
	}

	// create our private and public key
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := x509.MarshalPKCS1PublicKey(&caPrivKey.PublicKey)
	pubKeyHash := sha256.Sum256(pubKeyBytes)

	ca := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:    "indece Monitor auto generated CA",
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		SubjectKeyId:          pubKeyHash[:],
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// create the CA
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, err
	}

	// pem encode
	caCrtPEM := new(bytes.Buffer)
	err = pem.Encode(caCrtPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("error encoding ca crt: %s", err)
	}

	caPrivKeyPEM := new(bytes.Buffer)
	err = pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})
	if err != nil {
		return nil, fmt.Errorf("error encoding ca key: %s", err)
	}

	return &PEMCert{
		Crt:        caCrtPEM.Bytes(),
		Key:        caPrivKeyPEM.Bytes(),
		CreatedAt:  null.TimeFrom(ca.NotBefore),
		ValidUntil: null.TimeFrom(ca.NotAfter),
	}, nil
}
