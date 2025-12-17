package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/sm3"
	"github.com/tjfoc/gmsm/sm4"
)

var (
	SM4Key        []byte
	SM2PrivateKey *sm2.PrivateKey
	SM2PublicKey  *sm2.PublicKey
)

// InitCrypto 初始化国密算法
func InitCrypto(sm4KeyHex string) error {
	// 初始化SM4密钥
	var err error
	SM4Key, err = hex.DecodeString(sm4KeyHex)
	if err != nil {
		return err
	}

	// 生成SM2密钥对
	SM2PrivateKey, err = sm2.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}
	SM2PublicKey = &SM2PrivateKey.PublicKey

	return nil
}

// SM3Hash SM3哈希
func SM3Hash(data string) string {
	h := sm3.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// SM3HashWithSalt SM3哈希加盐
func SM3HashWithSalt(data, salt string) string {
	return SM3Hash(data + salt)
}

// SM4Encrypt SM4加密
func SM4Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	ciphertext, err := sm4.Sm4Ecb(SM4Key, []byte(plaintext), true)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(ciphertext), nil
}

// SM4Decrypt SM4解密
func SM4Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	data, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	plaintext, err := sm4.Sm4Ecb(SM4Key, data, false)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// SM2Encrypt SM2加密
func SM2Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	ciphertext, err := sm2.Encrypt(SM2PublicKey, []byte(plaintext), rand.Reader, sm2.C1C3C2)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(ciphertext), nil
}

// SM2Decrypt SM2解密
func SM2Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	data, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	plaintext, err := sm2.Decrypt(SM2PrivateKey, data, sm2.C1C3C2)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// GetSM2PublicKeyHex 获取SM2公钥的16进制字符串
func GetSM2PublicKeyHex() string {
	if SM2PublicKey == nil {
		return ""
	}
	return hex.EncodeToString(SM2PublicKey.X.Bytes()) + hex.EncodeToString(SM2PublicKey.Y.Bytes())
}
