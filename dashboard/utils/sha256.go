package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

func digest(h hash.Hash, datas ...[]byte) ([]byte, error) {
	for _, data := range datas {
		_, err := h.Write(data)
		if err != nil {
			return nil, err
		}
	}

	return h.Sum(nil), nil
}

// Sha256Hex returns the Sha256 checksum of the data with HEX format
func Sha256Hex(datas ...[]byte) (string, error) {
	b, err := digest(sha256.New(), datas...)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
