package guardsman

import "github.com/google/uuid"

type Key interface {
	CompareKey() (bool, error)
}

type APIKey struct {
	HashGenerator HashGenerator
	Hash          *HashSalt
	ProvidedKey   string
	ProvidedSalt  string
}

func NewApiKey() (*APIKey, error) {
	key := uuid.New().String()
	hash, err := DefaultArgonGenerator().GenerateHash([]byte(key), make([]byte, 0))
	if err != nil {
		return nil, err
	}
	return &APIKey{
		HashGenerator: *DefaultArgonGenerator(),
		Hash:          hash,
		ProvidedKey:   key,
		ProvidedSalt:  hash.Salt,
	}, nil
}

func NewPasswordKey(password string) (*APIKey, error) {
	hash, err := DefaultArgonGenerator().GenerateHash([]byte(password), make([]byte, 0))
	if err != nil {
		return nil, err
	}
	return &APIKey{
		HashGenerator: *DefaultArgonGenerator(),
		Hash:          hash,
		ProvidedKey:   password,
		ProvidedSalt:  hash.Salt,
	}, nil
}

func NewApiKeyFromDatabase(userProvidedKey, hashedKey, databaseSalt string) *APIKey {
	return &APIKey{
		HashGenerator: DefaultArgonGenerator(),
		Hash:          &HashSalt{Hash: hashedKey, Salt: databaseSalt},
		ProvidedKey:   userProvidedKey,
		ProvidedSalt:  databaseSalt,
	}

}

type HashSalt struct {
	Hash string
	Salt string
}

func (k *APIKey) HashKey() error {
	hash, err := k.HashGenerator.GenerateHash([]byte(k.ProvidedKey), []byte(k.ProvidedSalt))
	if err != nil {
		return err
	}
	k.Hash = hash
	return nil
}

func (k APIKey) CompareKey() (bool, error) {
	return k.HashGenerator.Compare([]byte(k.Hash.Hash), []byte(k.Hash.Salt), []byte(k.ProvidedKey))
}
