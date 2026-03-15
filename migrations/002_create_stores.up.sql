-- 门店表
CREATE TABLE IF NOT EXISTS `stores` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '门店ID',
  `name` VARCHAR(100) NOT NULL COMMENT '门店名称',
  `code` VARCHAR(50) NOT NULL COMMENT '门店编码',
  `address` VARCHAR(255) DEFAULT NULL COMMENT '地址',
  `province` VARCHAR(50) DEFAULT NULL COMMENT '省',
  `city` VARCHAR(50) DEFAULT NULL COMMENT '市',
  `district` VARCHAR(50) DEFAULT NULL COMMENT '区',
  `longitude` DECIMAL(10, 7) DEFAULT NULL COMMENT '经度',
  `latitude` DECIMAL(10, 7) DEFAULT NULL COMMENT '纬度',
  `phone` VARCHAR(20) DEFAULT NULL COMMENT '联系电话',
  `business_hours` VARCHAR(100) DEFAULT NULL COMMENT '营业时间',
  `images` JSON DEFAULT NULL COMMENT '门店图片',
  `description` TEXT DEFAULT NULL COMMENT '门店介绍',
  `status` TINYINT DEFAULT 1 COMMENT '状态: 0-停业, 1-营业',
  `sort_order` INT DEFAULT 0 COMMENT '排序',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`),
  KEY `idx_status` (`status`),
  KEY `idx_sort` (`sort_order`),
  KEY `idx_location` (`longitude`, `latitude`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='门店表';

-- 插入示例门店
INSERT INTO `stores` (`name`, `code`, `address`, `province`, `city`, `district`, `phone`, `business_hours`, `status`, `sort_order`) VALUES
('总店', 'STORE001', '北京市朝阳区xxx路xxx号', '北京市', '北京市', '朝阳区', '010-12345678', '09:00-21:00', 1, 1),
('分店A', 'STORE002', '上海市浦东新区xxx路xxx号', '上海市', '上海市', '浦东新区', '021-12345678', '09:00-21:00', 1, 2);
