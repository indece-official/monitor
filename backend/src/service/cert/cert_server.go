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
)

func (s *Service) GenerateServerCert(hostname string, caPEM *PEMCert) (*PEMCert, error) {
	caCrtPem, _ := pem.Decode(caPEM.Crt)
	if caCrtPem == nil {
		return nil, fmt.Errorf("error loading clients cert - no pem block found")
	}

	caCrt, err := x509.ParseCertificate(caCrtPem.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error loading clients cert: %s", err)
	}

	caKeyPem, _ := pem.Decode(caPEM.Key)
	if caKeyPem == nil {
		return nil, fmt.Errorf("error loading clients key - no pem block found")
	}

	caKey, err := x509.ParsePKCS1PrivateKey(caKeyPem.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error loading clients key: %s", err)
	}

	serialNumber, err := s.generateSerialNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %s", err)
	}

	// create our private and public key
	serverKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := x509.MarshalPKCS1PublicKey(&serverKey.PublicKey)
	pubKeyHash := sha256.Sum256(pubKeyBytes)

	// set up our server certificate
	serverCrt := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:    hostname,
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		DNSNames: []string{
			hostname,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  false,
		SubjectKeyId:          pubKeyHash[:],
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment, // Only for rsa-keys
		BasicConstraintsValid: true,
	}

	serverCrtData, err := x509.CreateCertificate(rand.Reader, serverCrt, caCrt, &serverKey.PublicKey, caKey)
	if err != nil {
		return nil, err
	}

	serverCrtPEM := new(bytes.Buffer)
	err = pem.Encode(serverCrtPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: serverCrtData,
	})
	if err != nil {
		return nil, fmt.Errorf("error encoding server crt: %s", err)
	}

	serverKeyPEM := new(bytes.Buffer)
	err = pem.Encode(serverKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(serverKey),
	})
	if err != nil {
		return nil, fmt.Errorf("error encoding server key: %s", err)
	}

	return &PEMCert{
		Crt: serverCrtPEM.Bytes(),
		Key: serverKeyPEM.Bytes(),
	}, nil
}
