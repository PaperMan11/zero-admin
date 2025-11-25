-- 登录日志表
CREATE TABLE `sys_login_log` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '登录日志ID',
    `username` varchar(200) NOT NULL COMMENT '用户',
    `ip` varchar(50) NOT NULL COMMENT '登录IP',
    `location` varchar(100) NOT NULL COMMENT '登录地址',
    `browser` varchar(100) NOT NULL DEFAULT '' COMMENT '浏览器信息',
    `os` varchar(100) NOT NULL DEFAULT '' COMMENT '操作系统信息',
    `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '登录状态（0-失败，1-成功）',
    `message` varchar(255) NOT NULL DEFAULT '' COMMENT '登录信息',
    `login_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='登录日志表';

-- 系统操作日志
CREATE TABLE `sys_operate_log` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '日志ID',
    `user_id` bigint(20) NOT NULL COMMENT '操作用户',
    `operation` varchar(255) NOT NULL DEFAULT '' COMMENT '操作描述',
    `method` varchar(255) NOT NULL DEFAULT '' COMMENT '请求方法',
    `params` text COMMENT '请求参数',
    `ip` varchar(50) NOT NULL DEFAULT '' COMMENT 'IP地址',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统操作日志';