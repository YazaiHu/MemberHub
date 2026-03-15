-- 插入测试数据

-- 1. 创建默认管理员账号
-- 用户名: admin, 密码: admin123 (bcrypt加密后的hash)
INSERT INTO `admins` (`username`, `password`, `real_name`, `phone`, `email`, `status`) VALUES
('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKXgf8xPZM19DnmN7.4PqRCnfkbi', '系统管理员', '13800138000', 'admin@example.com', 1);

-- 2. 创建测试门店
INSERT INTO `stores` (`name`, `phone`, `address`, `detailed_address`, `longitude`, `latitude`, `business_hours`, `description`, `status`) VALUES
('总店', '400-123-4567', '北京市朝阳区建国路88号', 'SOHO现代城A座1层', 116.447621, 39.910168, '09:00-21:00', '我们的旗舰店，提供全系列产品和服务', 1),
('三里屯分店', '010-12345678', '北京市朝阳区三里屯路19号', '三里屯太古里南区S6-31', 116.455464, 39.937784, '10:00-22:00', '位于时尚购物中心的精品店', 1),
('望京分店', '010-87654321', '北京市朝阳区望京街10号', '望京SOHO T3座1层', 116.481499, 39.996458, '09:00-21:00', '服务望京商圈的便利店', 1);

-- 3. 创建测试充值活动
INSERT INTO `recharge_promotions` (`name`, `description`, `recharge_amount`, `bonus_amount`, `bonus_points`, `start_time`, `end_time`, `status`) VALUES
('充200送50活动', '充值200元，赠送50元余额+100积分', 200.00, 50.00, 100, '2024-01-01 00:00:00', '2024-12-31 23:59:59', 1),
('充500送150活动', '充值500元，赠送150元余额+300积分', 500.00, 150.00, 300, '2024-01-01 00:00:00', '2024-12-31 23:59:59', 1),
('充1000送350活动', '充值1000元，赠送350元余额+600积分', 1000.00, 350.00, 600, '2024-01-01 00:00:00', '2024-12-31 23:59:59', 1);

-- 4. 创建测试积分兑换规则
INSERT INTO `points_exchange_rules` (`store_id`, `product_name`, `description`, `points_required`, `stock`, `daily_limit`, `valid_days`, `status`) VALUES
(1, '美式咖啡兑换券', '可兑换一杯中杯美式咖啡', 100, 100, 5, 30, 1),
(1, '拿铁咖啡兑换券', '可兑换一杯中杯拿铁咖啡', 150, 50, 3, 30, 1),
(1, '蛋糕券', '可兑换任意一款蛋糕', 200, 30, 2, 15, 1),
(2, '三里屯特饮券', '三里屯店专属特饮', 180, 50, 3, 30, 1),
(3, '望京早餐套餐', '咖啡+三明治套餐', 220, 40, 5, 30, 1);

-- 5. 创建测试优惠券模板
INSERT INTO `coupon_templates` (`name`, `type`, `discount_amount`, `discount_rate`, `min_amount`, `total_quantity`, `per_user_limit`, `valid_days`, `description`, `status`) VALUES
('新人专享券', 'full_discount', 20.00, NULL, 100.00, 1000, 1, 30, '新用户专享，消费满100减20', 1),
('周年庆折扣券', 'discount', NULL, 8.8, 0.00, 500, 2, 15, '全场8.8折优惠', 1),
('10元代金券', 'voucher', 10.00, NULL, 0.00, 2000, 5, 30, '无门槛10元代金券', 1),
('满200减50券', 'full_discount', 50.00, NULL, 200.00, 300, 1, 30, '消费满200减50', 1);

-- 6. 创建测试特价商品
INSERT INTO `promotion_products` (`store_id`, `product_name`, `description`, `original_price`, `promotion_price`, `stock`, `start_time`, `end_time`, `sort_order`, `status`) VALUES
(1, '特价美式咖啡', '每周特价，限量供应', 38.00, 28.00, 100, '2024-01-15 00:00:00', '2024-12-31 23:59:59', 100, 1),
(1, '精选蛋糕', '当日现做，新鲜美味', 68.00, 48.00, 50, '2024-01-15 00:00:00', '2024-12-31 23:59:59', 90, 1),
(2, '三里屯特调饮品', '网红特调，限时优惠', 58.00, 38.00, 80, '2024-01-15 00:00:00', '2024-12-31 23:59:59', 85, 1),
(3, '望京早餐套餐', '咖啡+三明治', 45.00, 35.00, 60, '2024-01-15 00:00:00', '2024-12-31 23:59:59', 80, 1);

-- 7. 创建测试用户（模拟微信用户）
INSERT INTO `users` (`openid`, `nickname`, `avatar`, `phone`, `gender`, `status`) VALUES
('wx_test_user_001', '测试用户001', 'https://thirdwx.qlogo.cn/mmopen/test001.jpg', '13800138001', 1, 1),
('wx_test_user_002', '测试用户002', 'https://thirdwx.qlogo.cn/mmopen/test002.jpg', '13800138002', 2, 1),
('wx_test_user_003', '测试用户003', 'https://thirdwx.qlogo.cn/mmopen/test003.jpg', '13800138003', 1, 1);

-- 8. 为测试用户创建积分账户
INSERT INTO `member_points` (`user_id`, `available_points`, `total_points`, `version`)
SELECT id, 500, 500, 1 FROM users WHERE openid LIKE 'wx_test_user_%';

-- 9. 为测试用户创建余额账户
INSERT INTO `member_balance` (`user_id`, `balance`, `total_recharge`, `total_consumption`, `version`)
SELECT id, 0.00, 0.00, 0.00, 1 FROM users WHERE openid LIKE 'wx_test_user_%';
