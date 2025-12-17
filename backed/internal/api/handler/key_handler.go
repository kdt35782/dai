package handler

import (
	"github.com/gin-gonic/gin"
	"sm-medical/internal/crypto"
	"sm-medical/pkg/utils"
)

type KeyHandler struct{}

func NewKeyHandler() *KeyHandler {
	return &KeyHandler{}
}

// Generate 生成密钥
func (h *KeyHandler) Generate(c *gin.Context) {
	publicKey := crypto.GetSM2PublicKeyHex()

	utils.Success(c, gin.H{
		"publicKey": publicKey,
		"message":   "密钥已生成",
	})
}
