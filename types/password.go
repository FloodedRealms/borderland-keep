package types

type Password struct {
	HashGenerator HashGenerator
	Hash          *HashSalt
	ProvidedKey   string
	ProvidedSalt  string
}

func NewPassword(password string) (*Password, error) {
	hash, err := DefaultArgonGenerator().GenerateHash([]byte(password), make([]byte, 0))
	if err != nil {
		return nil, err
	}
	return &Password{
		HashGenerator: *DefaultArgonGenerator(),
		Hash:          hash,
		ProvidedKey:   password,
		ProvidedSalt:  hash.Salt,
	}, nil
}

func NewPasswordFromDatabase(userProvidedPassword, hashedPassword, databaseSalt string) *Password {
	return &Password{
		HashGenerator: DefaultArgonGenerator(),
		Hash:          &HashSalt{Hash: hashedPassword, Salt: databaseSalt},
		ProvidedKey:   userProvidedPassword,
		ProvidedSalt:  databaseSalt,
	}

}

func (k *Password) HashKey() error {
	hash, err := k.HashGenerator.GenerateHash([]byte(k.ProvidedKey), []byte(k.ProvidedSalt))
	if err != nil {
		return err
	}
	k.Hash = hash
	return nil
}

func (k Password) CompareKey() (bool, error) {
	return k.HashGenerator.Compare([]byte(k.Hash.Hash), []byte(k.Hash.Salt), []byte(k.ProvidedKey))
}
