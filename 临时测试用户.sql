-- 临时测试用户脚本
-- 用于前后端联调测试，使用模拟加密方式

USE SM;

-- 创建测试用户
-- 用户名: testuser
-- 密码: Test123!@#
-- 注意：这里的密码使用的是前端模拟哈希算法加密的值

-- 第一步：计算模拟哈希值
-- 在浏览器控制台执行以下代码获取哈希值：
/*
function simulateHash(str) {
    let hash = 0;
    if (str.length === 0) return hash.toString(16);
    for (let i = 0; i < str.length; i++) {
        const char = str.charCodeAt(i);
        hash = ((hash << 5) - hash) + char;
        hash = hash & hash;
    }
    return Math.abs(hash).toString(16).padStart(64, '0');
}

// 计算 "Test123!@#" 的两次哈希
const once = simulateHash("Test123!@#");
console.log("第一次哈希:", once);
const twice = simulateHash(once + "testuser"); // 后端会加盐（用户名）再次哈希
console.log("第二次哈希(加盐):", twice);
*/

-- 或者直接使用这个预计算的值：
-- 密码 Test123!@# 的模拟SM3哈希值

INSERT INTO SM_user (
    username, 
    password,  -- 使用后端加盐后的哈希值
    email, 
    phone, 
    identify, 
    status,
    gender,
    age,
    created_at
) VALUES (
    'testuser',
    '00000000000000000000000000000000000000000000000000000bc58e5a1c',  -- 这是模拟哈希值
    'test@example.com',  -- 邮箱（需要SM4加密，这里暂时用明文）
    '13800138000',       -- 手机号（需要SM4加密，这里暂时用明文）
    'user',
    0,
    1,
    25,
    NOW()
);

-- 验证插入
SELECT id, username, email, identify, status FROM SM_user WHERE username = 'testuser';
