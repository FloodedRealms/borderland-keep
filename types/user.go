package types

import "github.com/google/uuid"

type User interface {
	DisplayUUID() string
	DisplayAPIKey() string
	DisplayUserName() string
	RetreiveSalt() string
	RetreiveHash() string
	Validate() (bool, error)
}
type APIRequest struct {
	Auth APIAuthSection `json:"auth"`
}
type APIAuthSection struct {
	ProvidedClientId string `json:"client_id"`
	ProvidedAPIKey   string `json:"api_key"`
}
type APIUser struct {
	id            string
	friendly_name string
	key           APIKey
}

func GenerateNewUser(name string) (*APIUser, error) {
	key, err := NewApiKey()
	if err != nil {
		return nil, err
	}
	return &APIUser{
		id:            uuid.New().String(),
		friendly_name: name,
		key:           *key,
	}, nil
}
func LoadUser(clientId, providedKey, name, hashedKey, salt string) (*APIUser, error) {
	key := NewApiKeyFromDatabase(providedKey, hashedKey, salt)
	return &APIUser{
		id:            clientId,
		friendly_name: name,
		key:           *key,
	}, nil
}

func (au APIUser) DisplayUUID() string {
	return au.id
}

func (au APIUser) DisplayAPIKey() string {
	return au.key.ProvidedKey
}
func (au APIUser) RetreiveSalt() string {
	return au.key.Hash.Salt
}
func (au APIUser) RetreiveHash() string {
	return au.key.Hash.Hash
}

func (au APIUser) DisplayUserName() string {
	return au.friendly_name
}

func (au APIUser) Validate() (bool, error) {
	return au.key.CompareKey()
}
