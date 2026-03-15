-- 优惠券模板表
CREATE TABLE IF NOT EXISTS `coupon_templates` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `name` VARCHAR(100) NOT NULL COMMENT '优惠券名称',
  `type` TINYINT NOT NULL COMMENT '类型: 1-满减券, 2-折扣券, 3-代金券',
  `discount_type` TINYINT NOT NULL COMMENT '优惠类型: 1-金额, 2-折扣',
  `discount_value` INT NOT NULL COMMENT '优惠值（金额单位分，折扣单位0.1%，如95表示9.5折）',
  `min_amount` INT DEFAULT 0 COMMENT '最低消费金额（分，0表示无门槛）',
  `total_quantity` INT NOT NULL COMMENT '总发行量',
  `remaining_quantity` INT NOT NULL COMMENT '剩余数量',
  `per_user_limit` INT DEFAULT 1 COMMENT '每人限领数量',
  `valid_days` INT NOT NULL COMMENT '有效天数（从领取日算起）',
  `applicable_stores` JSON DEFAULT NULL COMMENT '适用门店ID列表（NULL表示全部门店）',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '使用说明',
  `start_time` DATETIME DEFAULT NULL COMMENT '领取开始时间',
  `end_time` DATETIME DEFAULT NULL COMMENT '领取结束时间',
  `status` TINYINT DEFAULT 1 COMMENT '状态: 0-停用, 1-启用',
  `sort_order` INT DEFAULT 0 COMMENT '排序',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_type` (`type`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort_order`),
  KEY `idx_time` (`start_time`, `end_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='优惠券模板表';

-- 用户优惠券表
CREATE TABLE IF NOT EXISTS `user_coupons` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `coupon_code` VARCHAR(32) NOT NULL COMMENT '优惠券码',
  `user_id` BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
  `template_id` BIGINT UNSIGNED NOT NULL COMMENT '模板ID',
  `name` VARCHAR(100) NOT NULL COMMENT '优惠券名称',
  `type` TINYINT NOT NULL COMMENT '类型',
  `discount_type` TINYINT NOT NULL COMMENT '优惠类型',
  `discount_value` INT NOT NULL COMMENT '优惠值',
  `min_amount` INT DEFAULT 0 COMMENT '最低消费金额（分）',
  `applicable_stores` JSON DEFAULT NULL COMMENT '适用门店',
  `status` TINYINT DEFAULT 1 COMMENT '状态: 1-未使用, 2-已使用, 3-已过期',
  `used_at` DATETIME DEFAULT NULL COMMENT '使用时间',
  `used_order_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '使用订单ID',
  `expired_at` DATETIME NOT NULL COMMENT '过期时间',
  `received_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '领取时间',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_coupon_code` (`coupon_code`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_template_id` (`template_id`),
  KEY `idx_status` (`status`),
  KEY `idx_expired_at` (`expired_at`),
  KEY `idx_received_at` (`received_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户优惠券表';

-- 优惠券推送任务表
CREATE TABLE IF NOT EXISTS `coupon_push_tasks` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `task_no` VARCHAR(64) NOT NULL COMMENT '任务号',
  `template_id` BIGINT UNSIGNED NOT NULL COMMENT '优惠券模板ID',
  `title` VARCHAR(100) NOT NULL COMMENT '推送标题',
  `target_type` TINYINT NOT NULL COMMENT '目标类型: 1-全部用户, 2-指定用户, 3-条件筛选',
  `target_users` JSON DEFAULT NULL COMMENT '目标用户ID列表（target_type=2时）',
  `target_conditions` JSON DEFAULT NULL COMMENT '筛选条件（target_type=3时）',
  `total_count` INT DEFAULT 0 COMMENT '总推送数',
  `success_count` INT DEFAULT 0 COMMENT '成功数',
  `fail_count` INT DEFAULT 0 COMMENT '失败数',
  `status` TINYINT DEFAULT 1 COMMENT '状态: 1-待执行, 2-执行中, 3-已完成, 4-已取消',
  `send_wechat_msg` TINYINT DEFAULT 1 COMMENT '是否发送微信通知: 0-否, 1-是',
  `started_at` DATETIME DEFAULT NULL COMMENT '开始时间',
  `completed_at` DATETIME DEFAULT NULL COMMENT '完成时间',
  `created_by` BIGINT UNSIGNED NOT NULL COMMENT '创建人ID',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_task_no` (`task_no`),
  KEY `idx_template_id` (`template_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='优惠券推送任务表';

-- 插入示例优惠券模板
INSERT INTO `coupon_templates` (`name`, `type`, `discount_type`, `discount_value`, `min_amount`, `total_quantity`, `remaining_quantity`, `per_user_limit`, `valid_days`, `description`, `status`, `sort_order`) VALUES
('新人礼包', 3, 1, 1000, 0, 10000, 10000, 1, 7, '新用户专享10元代金券，无门槛使用', 1, 1),
('满50减5', 1, 1, 500, 5000, 5000, 5000, 3, 30, '消费满50元可用', 1, 2),
('9折优惠券', 2, 2, 90, 10000, 3000, 3000, 2, 30, '全场9折，满100元可用', 1, 3);
