package entity

import (
	"bytes"
	"encoding/gob"
)

type SecretData struct {
	Name     string
	Value    any
	Metadata map[string]string
}
type Secret struct {
	ServerID        int64
	ServerUpdatedAt int64
	Data            *[]byte
	// TODO: ##@@ not used?
	//IsDirty         bool
}

// TODO: ##@@ move elsewhere to make cryptography single dependency

func (s *Secret) GetData() (SecretData, error) {
	var sd SecretData
	b := bytes.NewBuffer(*s.Data)
	// TODO: ##@@ add decryption layer
	dc := gob.NewDecoder(b)
	if err := dc.Decode(&sd); err != nil {
		return SecretData{}, err
	}
	return sd, nil
}

func (s *Secret) SetData(sd SecretData) error {
	var b bytes.Buffer
	// TODO: ##@@ add encryption layer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(sd)
	if err == nil {
		d := b.Bytes()
		s.Data = &d
	}
	return err
}

func init() {
	gob.Register(Password{})
	gob.Register(SecretData{})
}
