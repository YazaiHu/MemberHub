-- 特价商品表
CREATE TABLE IF NOT EXISTS `promotion_products` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `store_id` BIGINT UNSIGNED NOT NULL COMMENT '门店ID',
  `product_name` VARCHAR(100) NOT NULL COMMENT '商品名称',
  `product_image` VARCHAR(255) DEFAULT NULL COMMENT '商品图片',
  `original_price` INT NOT NULL COMMENT '原价（分）',
  `promotion_price` INT NOT NULL COMMENT '特价（分）',
  `stock` INT DEFAULT NULL COMMENT '库存（NULL表示不限）',
  `description` TEXT DEFAULT NULL COMMENT '商品描述',
  `start_time` DATETIME NOT NULL COMMENT '开始时间',
  `end_time` DATETIME NOT NULL COMMENT '结束时间',
  `status` TINYINT DEFAULT 1 COMMENT '状态: 0-下架, 1-上架',
  `sort_order` INT DEFAULT 0 COMMENT '排序',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_store_id` (`store_id`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort_order`),
  KEY `idx_time` (`start_time`, `end_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='特价商品表';

-- 操作日志表
CREATE TABLE IF NOT EXISTS `operation_logs` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `operator_id` BIGINT UNSIGNED NOT NULL COMMENT '操作人ID',
  `operator_name` VARCHAR(50) DEFAULT NULL COMMENT '操作人姓名',
  `operator_type` TINYINT NOT NULL COMMENT '操作人类型: 1-管理员, 2-会员, 3-系统',
  `module` VARCHAR(50) NOT NULL COMMENT '模块',
  `action` VARCHAR(50) NOT NULL COMMENT '操作',
  `resource` VARCHAR(100) DEFAULT NULL COMMENT '资源',
  `resource_id` BIGINT UNSIGNED DEFAULT NULL COMMENT '资源ID',
  `description` VARCHAR(255) DEFAULT NULL COMMENT '操作描述',
  `ip` VARCHAR(50) DEFAULT NULL COMMENT 'IP地址',
  `user_agent` VARCHAR(255) DEFAULT NULL COMMENT 'User-Agent',
  `request_data` JSON DEFAULT NULL COMMENT '请求数据',
  `response_data` JSON DEFAULT NULL COMMENT '响应数据',
  `status` TINYINT DEFAULT 1 COMMENT '状态: 0-失败, 1-成功',
  `error_msg` TEXT DEFAULT NULL COMMENT '错误信息',
  `duration` INT DEFAULT NULL COMMENT '执行时长（毫秒）',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_operator` (`operator_id`, `operator_type`),
  KEY `idx_module_action` (`module`, `action`),
  KEY `idx_resource` (`resource`, `resource_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';

-- 插入示例特价商品
INSERT INTO `promotion_products` (`store_id`, `product_name`, `original_price`, `promotion_price`, `stock`, `description`, `start_time`, `end_time`, `status`, `sort_order`) VALUES
(1, '本周特价-咖啡', 3000, 1999, 100, '精选咖啡豆，限时特惠', '2026-03-15 00:00:00', '2026-03-21 23:59:59', 1, 1),
(1, '本周特价-蛋糕', 5800, 3999, 50, '新鲜现做，当日优惠', '2026-03-15 00:00:00', '2026-03-21 23:59:59', 1, 2);
