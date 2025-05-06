package pkg

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type HashConfig struct {
	Time    uint32
	Memory  uint32
	KeyLen  uint32
	SaltLen uint32
	Threads uint8
}

func InitHashConfig() *HashConfig {
	return &HashConfig{}
}

func (h *HashConfig) UseConfig(time, memory, keyLen, saltLen uint32, threads uint8) {
	h.Time = time
	h.Memory = memory
	h.KeyLen = keyLen
	h.SaltLen = saltLen
	h.Threads = threads
}

func (h *HashConfig) UseDefaultConfig() {
	h.Threads = 2
	h.Time = 3
	h.Memory = 64 * 1024
	h.KeyLen = 32
	h.SaltLen = 16
}

func (h *HashConfig) genSalt() ([]byte, error) {
	salt := make([]byte, h.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

func (h *HashConfig) GenHashedPassword(password string) (string, error) {
	salt, err := h.genSalt()
	if err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, h.Time, h.Memory, h.Threads, h.KeyLen)
	// $jenisKey$versiKey$konfigurasi(memory, time, thread)$salt$hash
	version := argon2.Version
	bash64Salt := base64.RawStdEncoding.EncodeToString(salt)
	base64Hash := base64.RawStdEncoding.EncodeToString(hash)
	hashedPwd := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", version, h.Memory, h.Time, h.Threads, bash64Salt, base64Hash)
	return hashedPwd, nil
}

func (h *HashConfig) CompareHashAndPassword(hadhedPass string, password string) (bool, error) {
	salt, hash, err := h.decodeHash(hadhedPass)
	if err != nil {
		return false, err
	}
	newHash := argon2.IDKey([]byte(password), salt, h.Time, h.Memory, h.Threads, h.KeyLen)
	if subtle.ConstantTimeCompare(hash, newHash) == 0 {
		return false, err
	}
	return true, nil
}

func (h *HashConfig) decodeHash(hashedPass string) (salt []byte, hash []byte, err error) {
	// $jenisKey$versiKey$konfigurasi(memory, time, thread)$salt$hash
	values := strings.Split(hashedPass, "$")
	if len(values) != 6 {
		return nil, nil, fmt.Errorf("invalid length format")
	}
	if values[1] != "argon2id" {
		return nil, nil, fmt.Errorf("invalid hash type")
	}
	var version int
	if _, err := fmt.Sscanf(values[2], "v=%d", &version); err != nil {
		return nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, fmt.Errorf("invalid hash version")
	}
	if _, err := fmt.Sscanf(values[3], "m=%d,t=%d,p=%d", &h.Memory, &h.Time, &h.Threads); err != nil {
		return nil, nil, err
	}
	salt, err = base64.RawStdEncoding.DecodeString(values[4])
	if err != nil {
		return nil, nil, err
	}
	hash, err = base64.RawStdEncoding.DecodeString(values[5])
	if err != nil {
		return nil, nil, err
	}
	return salt, hash, nil
}
