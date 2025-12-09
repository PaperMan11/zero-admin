-- 创建系统用户表
CREATE TABLE `sys_user` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `username` varchar(50) NOT NULL COMMENT '用户名',
    `password` varchar(100) NOT NULL COMMENT '密码（加密存储）',
    `salt` varchar(50) NOT NULL DEFAULT '' COMMENT '盐值（用于密码加密）',
    `email` varchar(100) NOT NULL DEFAULT '' COMMENT '邮箱',
    `mobile` varchar(20) NOT NULL DEFAULT '' COMMENT '手机号',
    `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '头像URL',
    `real_name` varchar(50) NOT NULL DEFAULT '' COMMENT '真实姓名',
    `gender` tinyint(4) NOT NULL DEFAULT '0' COMMENT '性别（0-未知，1-男，2-女）',
    `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态（0-禁用，1-正常）',
    `last_login_time` datetime DEFAULT NULL COMMENT '最后登录时间',
    `last_login_ip` varchar(50) NOT NULL DEFAULT '' COMMENT '最后登录IP',
    `login_count` int(11) NOT NULL DEFAULT '0' COMMENT '登录次数',
    `creator` varchar(50) NOT NULL DEFAULT '' COMMENT '创建人',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updater` varchar(50) NOT NULL DEFAULT '' COMMENT '更新人',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `del_flag` tinyint(4) NOT NULL DEFAULT '0' COMMENT '删除标记（0-未删除，1-已删除）',
    `remark` varchar(500) NOT NULL DEFAULT '' COMMENT '备注',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`),
    KEY `idx_status` (`status`),
    KEY `idx_del_flag` (`del_flag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统用户表';

-- 创建角色表
CREATE TABLE `sys_role` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '角色ID',
    `role_name` varchar(50) NOT NULL COMMENT '角色名称',
    `role_code` varchar(50) NOT NULL COMMENT '角色编码（如：ADMIN, USER）',
    `description` varchar(255) NOT NULL DEFAULT '' COMMENT '角色描述',
    `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态（0-禁用，1-正常）',
    `creator` varchar(50) NOT NULL DEFAULT '' COMMENT '创建人',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updater` varchar(50) NOT NULL DEFAULT '' COMMENT '更新人',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `del_flag` tinyint(4) NOT NULL DEFAULT '0' COMMENT '删除标记（0-未删除，1-已删除）',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_role_code` (`role_code`),
    KEY `idx_status` (`status`),
    KEY `idx_del_flag` (`del_flag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统角色表';

-- 创建用户-角色关联表
CREATE TABLE `sys_user_role` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '关联ID',
    `user_id` bigint(20) NOT NULL COMMENT '用户ID',
    `role_id` bigint(20) NOT NULL COMMENT '角色ID',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_role` (`user_id`,`role_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统用户-角色关联表';


-- 创建安全范围表
CREATE TABLE `sys_scope` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '范围ID',
    `scope_name` varchar(50) NOT NULL COMMENT '范围名称',
    `scope_code` varchar(50) NOT NULL COMMENT '范围编码',
    `description` varchar(255) NOT NULL DEFAULT '' COMMENT '范围描述',
    `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
    `creator` varchar(50) NOT NULL DEFAULT '' COMMENT '创建人',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updater` varchar(50) NOT NULL DEFAULT '' COMMENT '更新人',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `del_flag` tinyint(4) NOT NULL DEFAULT '0' COMMENT '删除标记（0-未删除，1-已删除）',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_scope_code` (`scope_code`),
    KEY `idx_del_flag` (`del_flag`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='系统安全范围表';


-- 创建角色-安全范围关联表
CREATE TABLE `sys_role_scope` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '关联ID',
    `role_id` bigint(20) NOT NULL COMMENT '角色ID',
    `scope_id` bigint(20) NOT NULL COMMENT '范围ID',
    `perm` tinyint(4) NOT NULL DEFAULT '0' COMMENT '权限（0-无权限，1-读，2-写，4-创建，8-删除）',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_role_scope` (`role_id`,`scope_id`),
    KEY `idx_role_id` (`role_id`),
    KEY `idx_scope_id` (`scope_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统角色-安全范围关联表';


INSERT INTO `zeroadmin`.`sys_user`(`id`, `username`, `password`, `salt`, `email`, `mobile`, `avatar`, `real_name`, `gender`, `status`, `last_login_time`, `last_login_ip`, `login_count`, `creator`, `create_time`, `updater`, `update_time`, `del_flag`, `remark`) VALUES (1, 'admin', '$2a$10$0PeuGGGsaCocmOuKeYE7KOqBEhuW95O5j7.TyL0TkYg0kvpVbCVcS', 'pDmFnN6Q0jJwcIck', 'test@example.com', '13800138000', '', '超级管理员', 1, 1, '2025-12-09 15:44:45', '192.168.241.1:60767', 1, '', '2025-12-08 09:29:13', '', '2025-12-09 07:48:06', 0, '');
INSERT INTO `zeroadmin`.`sys_user_role`(`id`, `user_id`, `role_code`, `create_time`) VALUES (1, 1, 'SUPERUSER', '2025-12-09 07:50:43');
-- casbin rule
INSERT INTO `zeroadmin`.`casbin_rule`(`id`, `ptype`, `v0`, `v1`, `v2`, `v3`, `v4`, `v5`) VALUES (1, 'p', 'SUPERUSER', '/*', '(GET)|(POST)|(PUT)|(DELETE)', NULL, NULL, NULL);