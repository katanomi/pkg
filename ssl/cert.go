/*
Copyright 2023 The AlaudaDevops Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ssl

import (
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/binary"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const SystemCertPath = "/etc/ssl/certs"

// CertHash compute cert hash from cert path
func CertHash(certPath string) (string, error) {
	raw, err := os.ReadFile(certPath)
	if err != nil {
		return "", fmt.Errorf("read cert failed %s", err)
	}

	return CertRawHash(raw)
}

// CertRawHash compute cert hash from raw content
func CertRawHash(raw []byte) (string, error) {
	cert, err := ParseCert(raw)
	if err != nil {
		return "", err
	}
	nameHash, err := SubjectNameHash(cert)
	if err != nil {
		return "", fmt.Errorf("generate subject name hash failed: %w", err)
	}

	// The links created are of the form HHHHHHHH.D, where each H is a hexadecimal character and D is a single decimal digit.
	// see https://www.openssl.org/docs/manmaster/man1/c_rehash.html
	return fmt.Sprintf("%08x.%d", nameHash, 0), nil
}

// ParseCert parse cert from raw content
func ParseCert(raw []byte) (*x509.Certificate, error) {
	block, rest := pem.Decode(raw)
	if block == nil {
		return nil, errors.New("failed find PEM data")
	}
	extra, _ := pem.Decode(rest)
	if extra != nil {
		return nil, errors.New("found multiple PEM blocks, expected exactly one")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse certficate failed: %w", err)
	}

	return cert, nil
}

// SubjectNameHash is a reimplementation of the X509_subject_name_hash in openssl. It computes the SHA-1
// of the canonical encoding of the certificate's subject name and returns the 32-bit integer represented by the first
// four bytes of the hash using little-endian byte order.
//
// The output should be the same as the result of the following command
// openssl x509 -hash -fingerprint -noout -in <file>
func SubjectNameHash(cert *x509.Certificate) (uint32, error) {
	name, err := CanonicalName(cert.RawSubject)
	if err != nil {
		return 0, fmt.Errorf("failed to compute canonical subject name\n%w", err)
	}
	hasher := sha1.New() // nolint: gosec // G401: used for generate hash
	_, err = hasher.Write(name)
	if err != nil {
		return 0, fmt.Errorf("failed to compute sha1sum of canonical subject name\n%w", err)
	}
	sum := hasher.Sum(nil)
	return binary.LittleEndian.Uint32(sum[:4]), nil
}

// canonicalSET holds a of canonicalATVs. Suffix SET ensures it is marshaled as a set rather than a sequence
// by asn1.Marshal.
type canonicalSET []canonicalATV

// canonicalATV is similar to pkix.AttributeTypeAndValue but includes tag to ensure all values are marshaled as
// ASN.1, UTF8String values
type canonicalATV struct {
	Type  asn1.ObjectIdentifier
	Value string `asn1:"utf8"`
}

// CanonicalName accepts a DER encoded subject name and returns a "Canonical Encoding" matching that
// returned by the x509_name_canon function in openssl. All string values are transformed with CanonicalString
// and UTF8 encoded and the leading SEQ header is removed.
//
// see https://stackoverflow.com/questions/34095440/hash-algorithm-for-certificate-crl-directory.
func CanonicalName(name []byte) ([]byte, error) {
	var origSeq pkix.RDNSequence
	_, err := asn1.Unmarshal(name, &origSeq)
	if err != nil {
		return nil, fmt.Errorf("failed to parse subject name\n%w", err)
	}
	var result []byte
	for _, origSet := range origSeq {
		var canonSet canonicalSET
		for _, origATV := range origSet {
			origVal, ok := origATV.Value.(string)
			if !ok {
				return nil, errors.New("got unexpected non-string value")
			}
			canonSet = append(canonSet, canonicalATV{
				Type:  origATV.Type,
				Value: CanonicalString(origVal),
			})
		}
		setBytes, err := asn1.Marshal(canonSet)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal canonical name\n%w", err)
		}
		result = append(result, setBytes...)
	}
	return result, nil
}

// CanonicalString transforms the given string. All leading and trailing whitespace is trimmed
// where whitespace is defined as a space, formfeed, tab, newline, carriage return, or vertical tab
// character. Any remaining sequence of one or more consecutive whitespace characters in replaced with
// a single ' '.
//
// This is a reimplementation of the asn1_string_canon in openssl
func CanonicalString(s string) string {
	s = strings.Trim(s, " \f\t\n\r\v")
	s = strings.ToLower(s)

	return string(regexp.MustCompile(`[[:space:]]+`).ReplaceAll([]byte(s), []byte(" ")))
}
