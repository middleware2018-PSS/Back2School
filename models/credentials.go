package models

import (
	"encoding/json"
)

type Credential struct {
	Email    string `json:"email" db:"-"`
	Password string `json:"password" db:"-"`
}

// String is not required by pop and may be deleted
func (c Credential) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Credentials is not required by pop and may be deleted
type Credentials []Credential

// String is not required by pop and may be deleted
func (c Credentials) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}
