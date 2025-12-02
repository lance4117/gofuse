package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"

	"github.com/lance4117/gofuse/errs"
	"golang.org/x/crypto/argon2"
)

// EncryptedData 加密后的信息
type EncryptedData struct {
	Salt  []byte // Argon2id使用的随机盐（生成 key）
	Seal  []byte // GCM密文
	Nonce []byte // GCM随机nonce
}

// Encryption 密钥派生和对称加解密配置
type Encryption struct {
	Memory      int
	Iterations  int // 迭代次数
	Parallelism int // 并行度 0~255
	SaltLength  int // 盐长度
	KeyLength   int // Argon2生成的 key 长度，典型为16/24/32
}

// New 获取默认配置
func New() *Encryption {
	return &Encryption{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 4,
		SaltLength:  32,
		KeyLength:   32,
	}
}

// EncryptArgon2id 通过password派生key并加密
func (s Encryption) EncryptArgon2id(password string, data []byte) (EncryptedData, error) {
	salt := s.Argon2idSalt()
	key := s.Argon2idKey(password, salt)
	gcm, err := getGCMByKey(key)
	if err != nil {
		return EncryptedData{}, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return EncryptedData{}, err
	}
	seal := gcm.Seal(nil, nonce, data, nil)
	return EncryptedData{Salt: salt, Seal: seal, Nonce: nonce}, nil
}

// DecryptArgon2id 提供密码、加密信息和配置，解密数据
func (s Encryption) DecryptArgon2id(password string, data EncryptedData) ([]byte, error) {
	key := s.Argon2idKey(password, data.Salt)
	gcm, err := getGCMByKey(key)
	if err != nil {
		return nil, err
	}
	return gcm.Open(nil, data.Nonce, data.Seal, nil)
}

// EncryptAESGCM AESGCM对称加密
func (s Encryption) EncryptAESGCM(key, data []byte) ([]byte, error) {
	normKey := normalizeAESKey(key)
	gcm, err := getGCMByKey(normKey)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	// nonce放在前面，方便解密时提取
	return gcm.Seal(nonce, nonce, data, nil), nil
}

// DecryptAESGCM AESGCM对称解密
func (s Encryption) DecryptAESGCM(key, data []byte) ([]byte, error) {
	normKey := normalizeAESKey(key)
	gcm, err := getGCMByKey(normKey)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errs.ErrAESKeyLength
	}
	nonce, ct := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ct, nil)
}

// Argon2idKey Argon2id密钥派生
func (s Encryption) Argon2idKey(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, uint32(s.Iterations),
		uint32(s.Memory), uint8(s.Parallelism), uint32(s.KeyLength))
}

// Argon2idSalt Argon2id随机盐
func (s Encryption) Argon2idSalt() []byte {
	salt := make([]byte, s.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil
	}
	return salt
}

// getGCMByKey 创建 GCM 对象
func getGCMByKey(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

// normalizeAESKey 将任意长度 key 规范到 AES 支持的 16/24/32 字节，使用 sha256 扩展到 32 字节。
func normalizeAESKey(key []byte) []byte {
	switch len(key) {
	case 16, 24, 32:
		return key
	default:
		sum := sha256.Sum256(key)
		return sum[:]
	}
}
