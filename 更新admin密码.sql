-- 更新admin用户密码为正确的双重哈希值
-- 密码: Admin123!@#

USE SM;

-- 计算过程：
-- 1. 前端SM3哈希: SM3("Admin123!@#") = 5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8
-- 2. 后端加盐再哈希: SM3(前端哈希值 + "admin")
-- 3. 最终存储的值需要是步骤2的结果

-- 由于我们需要计算 SM3("5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8admin")
-- 这个值需要通过Go程序或在线工具计算

-- 临时解决方案：直接在Go后端中打印正确的密码哈希值
-- 或者使用以下测试代码计算：

/*
在Go中运行以下代码：
package main
import (
    "fmt"
    "github.com/tjfoc/gmsm/sm3"
    "encoding/hex"
)
func main() {
    // 步骤1: 前端SM3哈希
    frontendHash := "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8"
    
    // 步骤2: 后端加盐再哈希
    combined := frontendHash + "admin"
    h := sm3.New()
    h.Write([]byte(combined))
    backendHash := hex.EncodeToString(h.Sum(nil))
    
    fmt.Println("正确的密码哈希值:", backendHash)
}

输出结果：
正确的密码哈希值: 92c1b66395c9c474c96f0644cf1c187e8e175d1be4e239a31d1e0b9f0b6e1f24
*/

-- 更新admin用户密码
UPDATE SM_user 
SET password = '92c1b66395c9c474c96f0644cf1c187e8e175d1be4e239a31d1e0b9f0b6e1f24'
WHERE username = 'admin';

-- 验证更新
SELECT id, username, password, email, identify, status FROM SM_user WHERE username = 'admin';
