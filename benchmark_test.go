package webhook

import (
	"testing"
	"time"
)

// BenchmarkSignerSign excludes construction and signs a 4 KiB body plus metadata
// with one key; default random nonce generation stays in measured SHA256/512 calls.
func BenchmarkSignerSign(b *testing.B) {
	now := time.Unix(1_700_000_000, 0)
	message := Message{
		Timestamp: now, Method: "POST", Path: "/hooks", Host: "example.com",
		Body: make([]byte, 4096), Metadata: map[string]string{"tenant": "one"},
	}
	for _, algorithm := range []Algorithm{SHA256, SHA512} {
		b.Run(string(algorithm), func(b *testing.B) {
			signer, err := NewSigner(SignerConfig{
				Algorithm: algorithm,
				Keys:      []SigningKey{{ID: "key", Secret: []byte("benchmark-key-material")}},
				Clock:     func() time.Time { return now },
			})
			if err != nil {
				b.Fatalf("NewSigner() error = %v", err)
			}
			b.ReportAllocs()
			b.SetBytes(int64(len(message.Body)))
			b.ResetTimer()
			for range b.N {
				if _, err := signer.Sign(message); err != nil {
					b.Fatalf("Sign() error = %v", err)
				}
			}
		})
	}
}

// BenchmarkVerifierVerify excludes signing/configuration and verifies one valid
// SHA256 candidate over a 4 KiB body without replay storage or HTTP decoding.
func BenchmarkVerifierVerify(b *testing.B) {
	now := time.Unix(1_700_000_000, 0)
	message := Message{Timestamp: now, Method: "POST", Path: "/hooks", Host: "example.com", Body: make([]byte, 4096)}
	key := []byte("benchmark-key-material")
	signer, err := NewSigner(SignerConfig{Algorithm: SHA256, Keys: []SigningKey{{ID: "key", Secret: key}}, Clock: func() time.Time { return now }})
	if err != nil {
		b.Fatalf("NewSigner() error = %v", err)
	}
	signatures, err := signer.Sign(message)
	if err != nil {
		b.Fatalf("Sign() error = %v", err)
	}
	verifier, err := NewVerifier(VerifierConfig{
		Algorithm: SHA256, Keys: []VerificationKey{{ID: "key", Secret: key}},
		Clock: func() time.Time { return now }, Tolerance: time.Minute,
	})
	if err != nil {
		b.Fatalf("NewVerifier() error = %v", err)
	}
	b.ReportAllocs()
	b.SetBytes(int64(len(message.Body)))
	b.ResetTimer()
	for range b.N {
		if _, err := verifier.Verify(message, signatures); err != nil {
			b.Fatalf("Verify() error = %v", err)
		}
	}
}
