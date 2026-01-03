-- Sub2API 易支付功能迁移脚本
-- 添加充值订单和充值套餐功能

-- 1. 创建 recharge_products 充值套餐表（必须先创建，因为 orders 表引用它）
CREATE TABLE IF NOT EXISTS recharge_products (
    id                      BIGSERIAL PRIMARY KEY,
    name                    VARCHAR(255) NOT NULL,                -- 套餐名称
    amount                  DECIMAL(20, 2) NOT NULL,              -- 人民币金额
    balance                 DECIMAL(20, 8) NOT NULL,              -- 基础余额（USD）
    bonus_balance           DECIMAL(20, 8) DEFAULT 0,             -- 赠送余额（USD）

    description             TEXT DEFAULT '',
    sort_order              INT NOT NULL DEFAULT 0,
    is_hot                  BOOLEAN NOT NULL DEFAULT false,
    discount_label          VARCHAR(50) DEFAULT '',

    status                  VARCHAR(20) NOT NULL DEFAULT 'active',

    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- 约束
    CONSTRAINT chk_recharge_products_status CHECK (status IN ('active', 'inactive'))
);

-- recharge_products 索引
CREATE INDEX IF NOT EXISTS idx_recharge_products_status_sort ON recharge_products(status, sort_order);

-- 3. 插入默认充值套餐
INSERT INTO recharge_products (name, amount, balance, bonus_balance, sort_order, is_hot, discount_label)
VALUES
    ('10元充值', 10.00, 1.00, 0.00, 1, false, ''),
    ('30元充值', 30.00, 3.00, 0.30, 2, false, '赠送10%'),
    ('50元充值', 50.00, 5.00, 1.00, 3, true, '赠送20%'),
    ('100元充值', 100.00, 10.00, 2.50, 4, true, '赠送25%'),
    ('200元充值', 200.00, 20.00, 6.00, 5, false, '赠送30%')
ON CONFLICT DO NOTHING;

-- 4. 创建 orders 充值订单表
CREATE TABLE IF NOT EXISTS orders (
    id                      BIGSERIAL PRIMARY KEY,
    order_no                VARCHAR(32) UNIQUE NOT NULL,          -- 订单号
    user_id                 BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    product_id              BIGINT REFERENCES recharge_products(id) ON DELETE SET NULL,
    product_name            VARCHAR(255) NOT NULL,                -- 冗余存储套餐名称
    amount                  DECIMAL(20, 2) NOT NULL,              -- 支付金额（CNY）
    bonus_amount            DECIMAL(20, 8) DEFAULT 0,             -- 赠送金额（USD）
    actual_amount           DECIMAL(20, 8) NOT NULL,              -- 到账金额（USD）

    payment_method          VARCHAR(50) NOT NULL,                 -- alipay/wxpay/epusdt
    payment_gateway         VARCHAR(50) DEFAULT 'epay',
    trade_no                VARCHAR(255),                         -- 第三方订单号

    status                  VARCHAR(20) NOT NULL DEFAULT 'pending',  -- pending/paid/failed/expired

    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    paid_at                 TIMESTAMPTZ,
    expired_at              TIMESTAMPTZ,                          -- 15分钟有效期

    notes                   TEXT DEFAULT '',
    callback_data           JSONB,

    -- 索引
    CONSTRAINT chk_orders_status CHECK (status IN ('pending', 'paid', 'failed', 'expired')),
    CONSTRAINT chk_payment_method CHECK (payment_method IN ('alipay', 'wxpay', 'epusdt'))
);

-- orders 索引
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_orders_order_no ON orders(order_no);
