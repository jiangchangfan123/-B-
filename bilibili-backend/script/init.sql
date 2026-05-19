-- ============================================================
-- 仿B站项目 - 云服务器 MySQL 初始化脚本
-- 执行方式：将本文件复制到云服务器，通过 MySQL 客户端执行
--   mysql -h 8.210.190.129 -P 3306 -u bilibili -p < init.sql
-- 或者进入 MySQL 容器后执行：
--   docker exec -i bilibili-mysql mysql -ubilibili -pbilibili123 < init.sql
-- ============================================================

-- 1. 创建数据库（如果尚不存在）
CREATE DATABASE IF NOT EXISTS `bilibili`
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

-- 2. 切换到该数据库
USE `bilibili`;

-- 3. 创建 users 表
DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `id`              BIGINT UNSIGNED     AUTO_INCREMENT PRIMARY KEY COMMENT '用户ID',
  `username`        VARCHAR(32)         NOT NULL COMMENT '用户名，3-20位字母数字',
  `password`        VARCHAR(255)        NOT NULL COMMENT 'bcrypt 哈希后的密码，明文禁止入库',
  `email`           VARCHAR(128)                 COMMENT '邮箱',
  `avatar`          VARCHAR(255)        DEFAULT '' COMMENT '头像URL',
  `sign`            VARCHAR(200)        DEFAULT '' COMMENT '个性签名',
  `role`            TINYINT             DEFAULT 1  COMMENT '角色：1=普通用户，2=版主，3=管理员',
  `coins`           INT                 DEFAULT 0  COMMENT '硬币余额',
  `status`          TINYINT             DEFAULT 1  COMMENT '状态：1=正常，2=禁用',
  `created_at`      TIMESTAMP           DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at`      TIMESTAMP           DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',

  -- 唯一约束
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_email` (`email`)

) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='用户表';

-- 4. 显示创建结果
SELECT '数据库 `bilibili` 创建成功' AS `status`;
SELECT 'users 表创建成功，字段数: ' || COUNT(*) AS `status`
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA = 'bilibili' AND TABLE_NAME = 'users';
