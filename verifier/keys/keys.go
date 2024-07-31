package keys

import (
	"crypto/ecdsa"
	"os"

	"github.com/Cedro-Finance/CedroAVS/core/config"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

type ECSDAKeyStore struct {
	PublicKey  *ecdsa.PublicKey
	PrivateKey *ecdsa.PrivateKey
}

// type BLSKeyStore struct {
// 	PublicKey  *bls.PublicKey
// 	PrivateKey *bls.SecretKey
// }

type KeyStore struct {
	// BlsKeyStore   *BLSKeyStore
	ECSDAKeyStore *ECSDAKeyStore
}

var Kstore *KeyStore

func LoadKeys(cfg *config.KeyStore) *KeyStore {
	file, err := os.ReadFile(cfg.Path)
	if err != nil {
		panic(err)
	}
	key, err := keystore.DecryptKey(file, cfg.Password)
	if err != nil {
		panic(err)
	}

	ecdsa := &ECSDAKeyStore{
		PublicKey:  &key.PrivateKey.PublicKey,
		PrivateKey: key.PrivateKey,
	}

	return &KeyStore{
		ECSDAKeyStore: ecdsa,
	}
}
