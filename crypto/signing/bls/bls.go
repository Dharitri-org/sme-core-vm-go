package bls

import (
	"github.com/Dharitri-org/sme-dharitri/crypto/signing"
	"github.com/Dharitri-org/sme-dharitri/crypto/signing/mcl"
	"github.com/Dharitri-org/sme-dharitri/crypto/signing/mcl/singlesig"
)

type bls struct {
}

func NewBLS() *bls {
	return &bls{}
}

func (b *bls) VerifyBLS(key []byte, msg []byte, sig []byte) error {
	suite := mcl.NewSuiteBLS12()
	keyGenerator := signing.NewKeyGenerator(suite)

	publicKey, err := keyGenerator.PublicKeyFromByteArray(key)
	if err != nil {
		return err
	}

	signer := singlesig.NewBlsSigner()

	return signer.Verify(publicKey, msg, sig)
}
