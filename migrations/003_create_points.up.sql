-- 会员积分表（汇总表）
CREATE TABLE IF NOT EXISTS `member_points` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `total_points` BIGINT DEFAULT 0 COMMENT '累计获得积分',
  `available_points` BIGINT DEFAULT 0 COMMENT '可用积分',
  `frozen_points` BIGINT DEFAULT 0 COMMENT '冻结积分',
  `used_points` BIGINT DEFAULT 0 COMMENT '已使用积分',
  `expired_points` BIGINT DEFAULT 0 COMMENT '已过期积分',
  `version` INT DEFAULT 0 COMMENT '版本号（乐观锁）',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_id` (`user_id`),
  KEY `idx_available_points` (`available_points`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员积分表';

-- 积分流水表（数据源头）
CREATE TABLE IF NOT EXISTS `points_transactions` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `transaction_no` VARCHAR(64) NOT NULL COMMENT '交易流水号',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `points` BIGINT NOT NULL COMMENT '积分变动数量（正数增加，负数减少）',
  `type` TINYINT NOT NULL COMMENT '类型: 1-消费赠送, 2-充值赠送, 3-手动调整, 4-兑换扣减, 5-过期扣减',
  `source` VARCHAR(50) NOT NULL COMMENT '来源',
  `ref_id` BIGINT DEFAULT NULL COMMENT '关联ID（订单ID/兑换记录ID等）',
  `remark` VARCHAR(255) DEFAULT NULL COMMENT '备注',
  `expire_at` DATETIME DEFAULT NULL COMMENT '过期时间',
  `status` TINYINT DEFAULT 1 COMMENT '状态: 0-已撤销, 1-正常',
  `operator_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '操作人ID（手动调整时）',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_transaction_no` (`transaction_no`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_type` (`type`),
  KEY `idx_ref_id` (`ref_id`),
  KEY `idx_expire_at` (`expire_at`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='积分流水表';

-- 积分兑换规则表
CREATE TABLE IF NOT EXISTS `points_exchange_rules` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `store_id` BIGINT UNSIGNED NOT NULL COMMENT '门店ID',
  `title` VARCHAR(100) NOT NULL COMMENT '兑换标题',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '描述',
  `points_cost` INT NOT NULL COMMENT '所需积分',
  `product_type` TINYINT NOT NULL COMMENT '商品类型: 1-实物, 2-优惠券, 3-服务',
  `product_name` VARCHAR(100) NOT NULL COMMENT '商品名称',
  `product_image` VARCHAR(255) DEFAULT NULL COMMENT '商品图片',
  `total_stock` INT DEFAULT NULL COMMENT '总库存（NULL表示不限）',
  `remaining_stock` INT DEFAULT NULL COMMENT '剩余库存',
  `daily_limit` INT DEFAULT NULL COMMENT '每日限兑总数（NULL表示不限）',
  `user_daily_limit` INT DEFAULT 1 COMMENT '每人每日限兑数',
  `user_total_limit` INT DEFAULT NULL COMMENT '每人总限兑数（NULL表示不限）',
  `start_time` DATETIME DEFAULT NULL COMMENT '开始时间',
  `end_time` DATETIME DEFAULT NULL COMMENT '结束时间',
  `status` TINYINT DEFAULT 1 COMMENT '状态: 0-下架, 1-上架',
  `sort_order` INT DEFAULT 0 COMMENT '排序',
  `version` INT DEFAULT 0 COMMENT '版本号（乐观锁）',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_store_id` (`store_id`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort_order`),
  KEY `idx_time` (`start_time`, `end_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='积分兑换规则表';

-- 积分兑换记录表
CREATE TABLE IF NOT EXISTS `points_exchange_records` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `record_no` VARCHAR(64) NOT NULL COMMENT '兑换单号',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `rule_id` BIGINT UNSIGNED NOT NULL COMMENT '兑换规则ID',
  `store_id` BIGINT UNSIGNED NOT NULL COMMENT '门店ID',
  `points_cost` INT NOT NULL COMMENT '消耗积分',
  `product_name` VARCHAR(100) NOT NULL COMMENT '商品名称',
  `status` TINYINT DEFAULT 1 COMMENT '状态: 1-待核销, 2-已核销, 3-已取消',
  `verify_code` VARCHAR(20) DEFAULT NULL COMMENT '核销码',
  `verified_at` DATETIME DEFAULT NULL COMMENT '核销时间',
  `verified_by` BIGINT UNSIGNED DEFAULT NULL COMMENT '核销人ID',
  `cancelled_at` DATETIME DEFAULT NULL COMMENT '取消时间',
  `cancel_reason` VARCHAR(255) DEFAULT NULL COMMENT '取消原因',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_record_no` (`record_no`),
  UNIQUE KEY `uk_verify_code` (`verify_code`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_rule_id` (`rule_id`),
  KEY `idx_store_id` (`store_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='积分兑换记录表';
