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

func (s *Service) GenerateClientCert(agentUID string, clientsPEM *PEMCert) (*PEMCert, error) {
	clientsCrtPem, _ := pem.Decode(clientsPEM.Crt)
	if clientsCrtPem == nil {
		return nil, fmt.Errorf("error loading clients cert - no pem block found")
	}

	clientsCrt, err := x509.ParseCertificate(clientsCrtPem.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error loading clients cert: %s", err)
	}

	clientsKeyPem, _ := pem.Decode(clientsPEM.Key)
	if clientsKeyPem == nil {
		return nil, fmt.Errorf("error loading clients key - no pem block found")
	}

	clientsKey, err := x509.ParsePKCS1PrivateKey(clientsKeyPem.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error loading clients key: %s", err)
	}

	serialNumber, err := s.generateSerialNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %s", err)
	}

	// create our private and public key
	clientKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := x509.MarshalPKCS1PublicKey(&clientKey.PublicKey)
	pubKeyHash := sha256.Sum256(pubKeyBytes)

	// set up our server certificate
	clientCrt := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:    agentUID,
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  false,
		SubjectKeyId:          pubKeyHash[:],
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment, // Only for rsa-keys
		BasicConstraintsValid: true,
	}

	clientCrtData, err := x509.CreateCertificate(rand.Reader, clientCrt, clientsCrt, &clientKey.PublicKey, clientsKey)
	if err != nil {
		return nil, err
	}

	clientCrtPEM := new(bytes.Buffer)
	err = pem.Encode(clientCrtPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: clientCrtData,
	})
	if err != nil {
		return nil, fmt.Errorf("error encoding client crt: %s", err)
	}

	clientKeyPEM := new(bytes.Buffer)
	err = pem.Encode(clientKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(clientKey),
	})
	if err != nil {
		return nil, fmt.Errorf("error encoding client key: %s", err)
	}

	return &PEMCert{
		Crt:        clientCrtPEM.Bytes(),
		Key:        clientKeyPEM.Bytes(),
		CreatedAt:  null.TimeFrom(clientCrt.NotBefore),
		ValidUntil: null.TimeFrom(clientCrt.NotAfter),
	}, nil
}
