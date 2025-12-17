package main

import (
	"encoding/hex"
	"fmt"
	"github.com/tjfoc/gmsm/sm3"
)

func main() {
	// 原始密码
	password := "Admin123!@#"
	username := "admin"
	
	// 步骤1: 前端SM3哈希
	h1 := sm3.New()
	h1.Write([]byte(password))
	frontendHash := hex.EncodeToString(h1.Sum(nil))
	
	fmt.Println("=== 密码加密过程 ===")
	fmt.Println("原始密码:", password)
	fmt.Println("用户名:", username)
	fmt.Println()
	fmt.Println("步骤1 - 前端SM3哈希:")
	fmt.Println(frontendHash)
	fmt.Println()
	
	// 步骤2: 后端加盐再哈希（盐是用户名）
	combined := frontendHash + username
	h2 := sm3.New()
	h2.Write([]byte(combined))
	backendHash := hex.EncodeToString(h2.Sum(nil))
	
	fmt.Println("步骤2 - 后端加盐哈希 (盐=用户名):")
	fmt.Println("输入: 前端哈希值 + 用户名")
	fmt.Println("计算:", combined)
	fmt.Println("结果:", backendHash)
	fmt.Println()
	
	fmt.Println("=== 数据库更新SQL ===")
	fmt.Println("UPDATE SM_user")
	fmt.Println("SET password = '" + backendHash + "'")
	fmt.Println("WHERE username = '" + username + "';")
	fmt.Println()
	
	// 测试其他常用密码
	fmt.Println("=== 其他测试密码 ===")
	testPasswords := []string{"Test123!@#", "Password123!@#", "User123!@#"}
	testUsernames := []string{"testuser", "testuser", "user"}
	
	for i, pwd := range testPasswords {
		h1 := sm3.New()
		h1.Write([]byte(pwd))
		hash1 := hex.EncodeToString(h1.Sum(nil))
		
		combined := hash1 + testUsernames[i]
		h2 := sm3.New()
		h2.Write([]byte(combined))
		hash2 := hex.EncodeToString(h2.Sum(nil))
		
		fmt.Printf("密码: %s, 用户名: %s\n", pwd, testUsernames[i])
		fmt.Printf("数据库密码: %s\n\n", hash2)
	}
}
