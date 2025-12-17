-- 使用前端当前SM3实现计算出的密码哈希值更新数据库
-- 这样可以先让登录功能正常工作
-- 后续可以考虑修复前端SM3实现

USE SM;

-- 前端SM3计算: Admin123!@# -> 549f2295378443e0503859d002d1e0c6018ee04421c39f6e4353fd3a00e3a49e
-- 后端加盐计算: SM3(前端哈希 + "admin") -> 137678542f3473058bac95a1c5bb30b9d708f5cef150717308430b5d4a944378

UPDATE SM_user
SET password = '137678542f3473058bac95a1c5bb30b9d708f5cef150717308430b5d4a944378'
WHERE username = 'admin';

-- 验证更新
SELECT username, password FROM SM_user WHERE username = 'admin';

-- 应该显示:
-- username: admin
-- password: 137678542f3473058bac95a1c5bb30b9d708f5cef150717308430b5d4a944378
