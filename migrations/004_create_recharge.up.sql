-- 会员余额表（汇总表）
CREATE TABLE IF NOT EXISTS `member_balance` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `balance` BIGINT DEFAULT 0 COMMENT '余额（分）',
  `frozen_balance` BIGINT DEFAULT 0 COMMENT '冻结金额（分）',
  `total_recharge` BIGINT DEFAULT 0 COMMENT '累计充值（分）',
  `total_consume` BIGINT DEFAULT 0 COMMENT '累计消费（分）',
  `version` INT DEFAULT 0 COMMENT '版本号（乐观锁）',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_id` (`user_id`),
  KEY `idx_balance` (`balance`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员余额表';

-- 余额流水表
CREATE TABLE IF NOT EXISTS `balance_transactions` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `transaction_no` VARCHAR(64) NOT NULL COMMENT '交易流水号',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `amount` BIGINT NOT NULL COMMENT '金额（分，正数增加，负数减少）',
  `type` TINYINT NOT NULL COMMENT '类型: 1-充值, 2-消费, 3-退款, 4-手动调整',
  `source` VARCHAR(50) NOT NULL COMMENT '来源',
  `ref_id` BIGINT DEFAULT NULL COMMENT '关联ID',
  `remark` VARCHAR(255) DEFAULT NULL COMMENT '备注',
  `operator_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '操作人ID',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_transaction_no` (`transaction_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_type` (`type`),
  KEY `idx_ref_id` (`ref_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='余额流水表';

-- 充值活动表
CREATE TABLE IF NOT EXISTS `recharge_promotions` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `title` VARCHAR(100) NOT NULL COMMENT '活动标题',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '活动描述',
  `recharge_amount` INT NOT NULL COMMENT '充值金额（元）',
  `bonus_amount` INT NOT NULL COMMENT '赠送金额（元）',
  `bonus_points` INT DEFAULT 0 COMMENT '赠送积分',
  `start_time` DATETIME DEFAULT NULL COMMENT '开始时间',
  `end_time` DATETIME DEFAULT NULL COMMENT '结束时间',
  `user_limit` INT DEFAULT NULL COMMENT '每人限购次数（NULL表示不限）',
  `total_limit` INT DEFAULT NULL COMMENT '总限购次数（NULL表示不限）',
  `status` TINYINT DEFAULT 1 COMMENT '状态: 0-停用, 1-启用',
  `sort_order` INT DEFAULT 0 COMMENT '排序',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort_order`),
  KEY `idx_time` (`start_time`, `end_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='充值活动表';

-- 充值订单表
CREATE TABLE IF NOT EXISTS `recharge_orders` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `order_no` VARCHAR(64) NOT NULL COMMENT '订单号',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `promotion_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '活动ID',
  `recharge_amount` INT NOT NULL COMMENT '充值金额（分）',
  `bonus_amount` INT DEFAULT 0 COMMENT '赠送金额（分）',
  `bonus_points` INT DEFAULT 0 COMMENT '赠送积分',
  `total_amount` INT NOT NULL COMMENT '实际到账金额（分）= recharge_amount + bonus_amount',
  `pay_amount` INT NOT NULL COMMENT '支付金额（分）= recharge_amount',
  `pay_method` TINYINT DEFAULT 1 COMMENT '支付方式: 1-微信支付',
  `transaction_id` VARCHAR(64) DEFAULT NULL COMMENT '微信支付交易号',
  `status` TINYINT DEFAULT 1 COMMENT '状态: 1-待支付, 2-已支付, 3-已取消, 4-已退款',
  `paid_at` DATETIME DEFAULT NULL COMMENT '支付时间',
  `expired_at` DATETIME DEFAULT NULL COMMENT '过期时间',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_order_no` (`order_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_promotion_id` (`promotion_id`),
  KEY `idx_transaction_id` (`transaction_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='充值订单表';

-- 插入示例充值活动
INSERT INTO `recharge_promotions` (`title`, `description`, `recharge_amount`, `bonus_amount`, `bonus_points`, `status`, `sort_order`) VALUES
('充100送10', '充值100元送10元余额', 100, 10, 0, 1, 1),
('充200送50', '充值200元送50元余额+100积分', 200, 50, 100, 1, 2),
('充500送150', '充值500元送150元余额+300积分', 500, 150, 300, 1, 3);
