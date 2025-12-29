-- 为 users 表添加 SSO 相关字段

-- 添加用户头像字段
ALTER TABLE users
ADD COLUMN IF NOT EXISTS avatar VARCHAR(500) DEFAULT '';

COMMENT ON COLUMN users.avatar IS '用户头像URL（来自SSO提供商）';

-- 添加完整SSO回调数据字段
ALTER TABLE users
ADD COLUMN IF NOT EXISTS sso_data TEXT DEFAULT '';

COMMENT ON COLUMN users.sso_data IS 'SSO登录回调的完整数据（JSON格式）';
