package keys

import (
	"crypto/ecdsa"
	"os"

	"github.com/Cedro-Finance/CedroAVS/verifierOperator/config"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

type KeyStore struct {
	PublicKey  *ecdsa.PublicKey
	PrivateKey *ecdsa.PrivateKey
}

var kstore *KeyStore;

func LoadKeys(cfg *config.KeyStore) *KeyStore {
	file, err := os.ReadFile(cfg.Path)
	if err != nil {
		panic(err)
	}
	key, err := keystore.DecryptKey(file, cfg.Password)
	if err != nil {
		panic(err)
	}

	

	return &KeyStore{
		PrivateKey: key.PrivateKey,
		PublicKey:  &key.PrivateKey.PublicKey,
	}
}
