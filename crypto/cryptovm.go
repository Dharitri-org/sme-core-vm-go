package crypto

import (
	"github.com/Dharitri-org/sme-core-vm-go/crypto/hashing"
	"github.com/Dharitri-org/sme-core-vm-go/crypto/signing/bls"
	"github.com/Dharitri-org/sme-core-vm-go/crypto/signing/ed25519"
	"github.com/Dharitri-org/sme-core-vm-go/crypto/signing/secp256k1"
)

// NewVMCrypto returns a composite struct containing VMCrypto functionality implementations
func NewVMCrypto() VMCrypto {
	return struct {
		Hasher
		Ed25519
		BLS
		Secp256k1
	}{
		Hasher:    hashing.NewHasher(),
		Ed25519:   ed25519.NewEd25519Signer(),
		BLS:       bls.NewBLS(),
		Secp256k1: secp256k1.NewSecp256k1(),
	}
}
