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
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '编号',
    `title` varchar(50) NOT NULL COMMENT '系统模块',
    `operation_type` varchar(50) NOT NULL COMMENT '操作类型',
    `operation_name` varchar(50) NOT NULL COMMENT '操作人员',
    `request_method` varchar(200) NOT NULL COMMENT '请求方式',
    `operation_url` varchar(50) NOT NULL COMMENT '操作方法',
    `operation_params` text NOT NULL COMMENT '请求参数',
    `operation_response` text NOT NULL COMMENT '响应参数',
    `operation_status` int(11) NOT NULL COMMENT '操作状态',
    `use_time` bigint(20) NOT NULL COMMENT '执行时长(毫秒)',
    `browser` varchar(64) NOT NULL COMMENT '浏览器',
    `os` varchar(64) NOT NULL COMMENT '操作系统',
    `operation_ip` varchar(64) NOT NULL DEFAULT '' COMMENT '操作地址',
    `operation_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统操作日志表';