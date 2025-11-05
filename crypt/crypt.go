package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	"golang.org/x/crypto/argon2"
)

// EncryptedData 加密后的信息
type EncryptedData struct {
	Salt  []byte // Argon2id加密的随机值，配合密码计算出key (base64)
	Seal  []byte // GCM加密后的内容 (base64)
	Nonce []byte // GCM加密随机数 (base64)
}

// Encryption 私钥加密配置信息
type Encryption struct {
	Memory      int
	Iterations  int // 遍历次数
	Parallelism int // 并行数 0~255
	SaltLength  int // 盐长度
	KeyLength   int
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

// EncryptArgon2id 通过password派生出的key加密私钥
func (s Encryption) EncryptArgon2id(password string, data []byte) (EncryptedData, error) {
	// 获得相应配置文件中的盐长度
	salt := s.Argon2idSalt()
	// 根据密码和随机盐派生出key
	key := s.Argon2idKey(password, salt)
	// 根据派生的密钥创建AES-GCM加密实例
	gcm, err := getGCMByKey(key)
	if err != nil {
		return EncryptedData{}, err
	}
	// 长度和gcm要求一致的随机数
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return EncryptedData{}, err
	}
	// 对私钥加密
	seal := gcm.Seal(nil, nonce, data, nil)
	// 加密后的信息，可为解密所用
	info := EncryptedData{
		Salt:  salt,
		Seal:  seal,
		Nonce: nonce,
	}
	return info, nil
}

// DecryptArgon2id 提供密码,加密信息和加密配置，返回私钥
func (s Encryption) DecryptArgon2id(password string, data EncryptedData) ([]byte, error) {
	// 根据密码和随机盐获取派生的key
	key := s.Argon2idKey(password, data.Salt)
	// 根据这个key生成gcm加密实例
	gcm, err := getGCMByKey(key)
	if err != nil {
		return nil, err
	}
	// 如果随机数或者密码不一致，会导致此gcm无法解析出seal数据
	return gcm.Open(nil, data.Nonce, data.Seal, nil)
}

// EncryptAESGCM AESGCM对称加密
func (s Encryption) EncryptAESGCM(key, data []byte) ([]byte, error) {
	gcm, err := getGCMByKey(key)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}
	// nonce塞入加密数据
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// DecryptAESGCM AESGCM对称解密
func (s Encryption) DecryptAESGCM(key, data []byte) ([]byte, error) {
	gcm, err := getGCMByKey(key)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	// 加密数据中，读取头的nonce
	nonce, ct := data[:nonceSize], data[nonceSize:]

	return gcm.Open(nil, nonce, ct, nil)
}

// Argon2idKey Argon2id密钥生成
func (s Encryption) Argon2idKey(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, uint32(s.Iterations),
		uint32(s.Memory), uint8(s.Parallelism), uint32(s.KeyLength))
}

// Argon2idSalt Argon2id随机盐生成
func (s Encryption) Argon2idSalt() []byte {
	// 根据配置信息生成salt
	salt := make([]byte, s.SaltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil
	}
	return salt
}

// 根据argon2id加密的key生成GCM实例
func getGCMByKey(key []byte) (cipher.AEAD, error) {
	// 创建 AES 密码块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}
