package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

const (
	// EncryptedPrefix 加密数据前缀，用于识别是否已加密
	EncryptedPrefix = "enc:"
)

var (
	ErrInvalidCiphertext = errors.New("invalid ciphertext")
	ErrDecryptionFailed  = errors.New("decryption failed")
)

// CryptoService 加密服务
type CryptoService struct {
	key []byte
}

// NewCryptoService 创建加密服务
// masterKey 通常是 JWT secret 或专用加密密钥
func NewCryptoService(masterKey string) *CryptoService {
	// 使用 SHA-256 派生 32 字节密钥用于 AES-256
	hash := sha256.Sum256([]byte(masterKey))
	return &CryptoService{
		key: hash[:],
	}
}

// Encrypt 使用 AES-256-GCM 加密明文
// 返回 base64 编码的密文（带前缀）
func (c *CryptoService) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 生成随机 nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 加密并附加 nonce
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Base64 编码并添加前缀
	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	return EncryptedPrefix + encoded, nil
}

// Decrypt 解密 AES-256-GCM 加密的密文
// 输入应为带前缀的 base64 编码密文
func (c *CryptoService) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	// 检查是否有加密前缀
	if len(ciphertext) <= len(EncryptedPrefix) {
		return ciphertext, nil // 未加密，直接返回
	}

	if ciphertext[:len(EncryptedPrefix)] != EncryptedPrefix {
		return ciphertext, nil // 未加密，直接返回（向后兼容）
	}

	// 移除前缀并解码
	encoded := ciphertext[len(EncryptedPrefix):]
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", ErrInvalidCiphertext
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", ErrInvalidCiphertext
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", ErrDecryptionFailed
	}

	return string(plaintext), nil
}

// IsEncrypted 检查字符串是否已加密
func IsEncrypted(s string) bool {
	return len(s) > len(EncryptedPrefix) && s[:len(EncryptedPrefix)] == EncryptedPrefix
}
