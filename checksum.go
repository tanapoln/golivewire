package golivewire

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

var defaultChecksum = &checksumManager{}

type checksumManager struct {
}

func (c *checksumManager) jsonstr(o interface{}) ([]byte, error) {
	return json.Marshal(o)
}

func (c *checksumManager) Generate(f fingerprint, m serverMemo) (string, error) {
	hasher := hmac.New(sha256.New, []byte(ChecksumKey))
	b, err := c.jsonstr(f)
	if err != nil {
		return "", err
	}
	hasher.Write(b)
	b2, err := c.jsonstr(m)
	if err != nil {
		return "", err
	}
	hasher.Write(b2)
	hashed := hasher.Sum(nil)
	return hex.EncodeToString(hashed), nil
}

func (c *checksumManager) Check(checksum string, f fingerprint, m serverMemo) (bool, error) {
	h, err := c.Generate(f, m)
	if err != nil {
		return false, err
	}
	return hmac.Equal([]byte(h), []byte(checksum)), nil
}
