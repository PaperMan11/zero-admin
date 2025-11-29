-- 创建系统菜单表
CREATE TABLE `sys_menu` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
    `scope_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '安全范围ID',
    `parent_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '父菜单ID (0表示根菜单)',
    `menu_name` varchar(50) NOT NULL COMMENT '菜单名称',
    `menu_type` char(1) NOT NULL COMMENT '菜单类型 (M-目录, C-菜单, F-按钮)',
    `path` varchar(255) NOT NULL DEFAULT '' COMMENT '路由路径',
    `component` varchar(255) NOT NULL DEFAULT '' COMMENT '组件路径',
    `redirect` varchar(255) NOT NULL DEFAULT '' COMMENT '重定向路径',
    `icon` varchar(50) NOT NULL DEFAULT '' COMMENT '菜单图标',
    `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
    `no_cache` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否缓存 (0-缓存, 1-不缓存)',
    `affix` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否固定在标签栏 (0-否, 1-是)',
    `external` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否外部链接 (0-否, 1-是)',
    `hidden` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否隐藏 (0-显示, 1-隐藏)',
    `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态 (0-禁用, 1-正常)',
    `creator` varchar(50) NOT NULL DEFAULT '' COMMENT '创建人',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updater` varchar(50) NOT NULL DEFAULT '' COMMENT '更新人',
    `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `del_flag` tinyint(4) NOT NULL DEFAULT '0' COMMENT '删除标记 (0-未删除, 1-已删除)',
    `remark` varchar(500) NOT NULL DEFAULT '' COMMENT '备注',
    PRIMARY KEY (`id`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_menu_type` (`menu_type`),
    KEY `idx_scope_id` (`scope_id`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8mb4 COMMENT='系统菜单表';


INSERT INTO `zeroadmin`.`sys_menu`(`id`, `scope_id`, `parent_id`, `menu_name`, `menu_type`, `path`, `component`, `redirect`, `icon`, `sort`, `no_cache`, `affix`, `external`, `hidden`, `status`, `creator`, `create_time`, `updater`, `update_time`, `del_flag`, `remark`) VALUES (1, 1, 0, 'Components', 'M', '/components', 'Layout', 'noRedirect', 'component', 10, 0, 0, 0, 0, 1, '', '2025-11-25 06:59:18', '', '2025-11-25 09:07:16', 0, '组件演示目录');
INSERT INTO `zeroadmin`.`sys_menu`(`id`, `scope_id`, `parent_id`, `menu_name`, `menu_type`, `path`, `component`, `redirect`, `icon`, `sort`, `no_cache`, `affix`, `external`, `hidden`, `status`, `creator`, `create_time`, `updater`, `update_time`, `del_flag`, `remark`) VALUES (2, 1, 1, 'Back To Top', 'C', 'back-to-top', 'components-demo/back-to-top', '', '', 1, 0, 0, 0, 0, 1, '', '2025-11-25 06:59:18', '', '2025-11-25 09:07:17', 0, '回到顶部组件演示');
INSERT INTO `zeroadmin`.`sys_menu`(`id`, `scope_id`, `parent_id`, `menu_name`, `menu_type`, `path`, `component`, `redirect`, `icon`, `sort`, `no_cache`, `affix`, `external`, `hidden`, `status`, `creator`, `create_time`, `updater`, `update_time`, `del_flag`, `remark`) VALUES (3, 1, 0, 'Table', 'M', '/table', 'Layout', '/table/complex-table', 'table', 10, 0, 0, 0, 0, 1, '', '2025-11-25 07:04:16', '', '2025-11-25 09:07:18', 0, '表格展示模块');
INSERT INTO `zeroadmin`.`sys_menu`(`id`, `scope_id`, `parent_id`, `menu_name`, `menu_type`, `path`, `component`, `redirect`, `icon`, `sort`, `no_cache`, `affix`, `external`, `hidden`, `status`, `creator`, `create_time`, `updater`, `update_time`, `del_flag`, `remark`) VALUES (4, 1, 3, 'Dynamic Table', 'C', 'dynamic-table', 'table/dynamic-table/index', ' ', '', 1, 0, 0, 0, 0, 1, '', '2025-11-25 07:04:16', '', '2025-11-25 09:07:19', 0, '动态表格');
INSERT INTO `zeroadmin`.`sys_menu`(`id`, `scope_id`, `parent_id`, `menu_name`, `menu_type`, `path`, `component`, `redirect`, `icon`, `sort`, `no_cache`, `affix`, `external`, `hidden`, `status`, `creator`, `create_time`, `updater`, `update_time`, `del_flag`, `remark`) VALUES (5, 1, 3, 'Drag Table', 'C', 'drag-table', 'table/drag-table', '', '', 2, 0, 0, 0, 0, 1, '', '2025-11-25 07:04:16', '', '2025-11-25 09:07:19', 0, '拖拽表格');
INSERT INTO `zeroadmin`.`sys_menu`(`id`, `scope_id`, `parent_id`, `menu_name`, `menu_type`, `path`, `component`, `redirect`, `icon`, `sort`, `no_cache`, `affix`, `external`, `hidden`, `status`, `creator`, `create_time`, `updater`, `update_time`, `del_flag`, `remark`) VALUES (6, 1, 3, 'Inline Edit', 'C', 'inline-edit-table', 'table/inline-edit-table', '', '', 3, 0, 0, 0, 0, 1, '', '2025-11-25 07:04:16', '', '2025-11-25 09:07:20', 0, '行内编辑表格');
INSERT INTO `zeroadmin`.`sys_menu`(`id`, `scope_id`, `parent_id`, `menu_name`, `menu_type`, `path`, `component`, `redirect`, `icon`, `sort`, `no_cache`, `affix`, `external`, `hidden`, `status`, `creator`, `create_time`, `updater`, `update_time`, `del_flag`, `remark`) VALUES (7, 1, 3, 'Form', 'C', 'form', 'form/index', '', '', 4, 0, 0, 0, 0, 1, '', '2025-11-25 07:04:16', '', '2025-11-25 09:07:21', 0, '表单');
INSERT INTO `zeroadmin`.`sys_menu`(`id`, `scope_id`, `parent_id`, `menu_name`, `menu_type`, `path`, `component`, `redirect`, `icon`, `sort`, `no_cache`, `affix`, `external`, `hidden`, `status`, `creator`, `create_time`, `updater`, `update_time`, `del_flag`, `remark`) VALUES (8, 1, 0, 'Charts', 'M', '/charts', 'Layout', 'noRedirect', 'chart', 20, 0, 0, 0, 0, 1, '', '2025-11-25 07:04:24', '', '2025-11-25 09:07:21', 0, '图表展示模块');
INSERT INTO `zeroadmin`.`sys_menu`(`id`, `scope_id`, `parent_id`, `menu_name`, `menu_type`, `path`, `component`, `redirect`, `icon`, `sort`, `no_cache`, `affix`, `external`, `hidden`, `status`, `creator`, `create_time`, `updater`, `update_time`, `del_flag`, `remark`) VALUES (9, 1, 8, 'Line Chart', 'C', 'line', 'charts/line', '', '', 1, 1, 0, 0, 0, 1, '', '2025-11-25 07:04:24', '', '2025-11-25 09:07:22', 0, '折线图');
INSERT INTO `zeroadmin`.`sys_menu`(`id`, `scope_id`, `parent_id`, `menu_name`, `menu_type`, `path`, `component`, `redirect`, `icon`, `sort`, `no_cache`, `affix`, `external`, `hidden`, `status`, `creator`, `create_time`, `updater`, `update_time`, `del_flag`, `remark`) VALUES (10, 1, 8, 'Mix Chart', 'C', 'mix-chart', 'charts/mix-chart', '', '', 2, 1, 0, 0, 0, 1, '', '2025-11-25 07:04:24', '', '2025-11-25 09:07:23', 0, '混合图表');
INSERT INTO `zeroadmin`.`sys_menu`(`id`, `scope_id`, `parent_id`, `menu_name`, `menu_type`, `path`, `component`, `redirect`, `icon`, `sort`, `no_cache`, `affix`, `external`, `hidden`, `status`, `creator`, `create_time`, `updater`, `update_time`, `del_flag`, `remark`) VALUES (11, 1, 0, 'Nested', 'M', '/nested', 'Layout', '/nested/menu1', 'nested', 30, 0, 0, 0, 0, 1, '', '2025-11-25 07:04:33', '', '2025-11-25 09:07:24', 0, '嵌套菜单演示');
INSERT INTO `zeroadmin`.`sys_menu`(`id`, `scope_id`, `parent_id`, `menu_name`, `menu_type`, `path`, `component`, `redirect`, `icon`, `sort`, `no_cache`, `affix`, `external`, `hidden`, `status`, `creator`, `create_time`, `updater`, `update_time`, `del_flag`, `remark`) VALUES (12, 1, 11, 'Menu1', 'M', 'menu1', 'nested/menu1/index', 'noRedirect', '', 1, 0, 0, 0, 0, 1, '', '2025-11-25 07:04:33', '', '2025-11-25 09:07:24', 0, '嵌套菜单1');
INSERT INTO `zeroadmin`.`sys_menu`(`id`, `scope_id`, `parent_id`, `menu_name`, `menu_type`, `path`, `component`, `redirect`, `icon`, `sort`, `no_cache`, `affix`, `external`, `hidden`, `status`, `creator`, `create_time`, `updater`, `update_time`, `del_flag`, `remark`) VALUES (13, 1, 12, 'Menu1-1', 'C', 'menu1-1', 'nested/menu1/menu1-1', '', '', 1, 0, 0, 0, 0, 1, '', '2025-11-25 07:04:33', '', '2025-11-25 09:07:25', 0, '嵌套菜单1-1');
INSERT INTO `zeroadmin`.`sys_menu`(`id`, `scope_id`, `parent_id`, `menu_name`, `menu_type`, `path`, `component`, `redirect`, `icon`, `sort`, `no_cache`, `affix`, `external`, `hidden`, `status`, `creator`, `create_time`, `updater`, `update_time`, `del_flag`, `remark`) VALUES (14, 1, 12, 'Menu1-2', 'M', 'menu1-2', 'nested/menu1/menu1-2', 'noRedirect', '', 2, 0, 0, 0, 0, 1, '', '2025-11-25 07:04:33', '', '2025-11-25 09:07:26', 0, '嵌套菜单1-2');
INSERT INTO `zeroadmin`.`sys_menu`(`id`, `scope_id`, `parent_id`, `menu_name`, `menu_type`, `path`, `component`, `redirect`, `icon`, `sort`, `no_cache`, `affix`, `external`, `hidden`, `status`, `creator`, `create_time`, `updater`, `update_time`, `del_flag`, `remark`) VALUES (15, 1, 14, 'Menu1-2-1', 'C', 'menu1-2-1', 'nested/menu1/menu1-2/menu1-2-1', '', '', 1, 0, 0, 0, 0, 1, '', '2025-11-25 07:04:33', '', '2025-11-25 09:07:26', 0, '嵌套菜单1-2-1');
INSERT INTO `zeroadmin`.`sys_menu`(`id`, `scope_id`, `parent_id`, `menu_name`, `menu_type`, `path`, `component`, `redirect`, `icon`, `sort`, `no_cache`, `affix`, `external`, `hidden`, `status`, `creator`, `create_time`, `updater`, `update_time`, `del_flag`, `remark`) VALUES (16, 1, 14, 'Menu1-2-2', 'C', 'menu1-2-2', 'nested/menu1/menu1-2/menu1-2-2', '', '', 2, 0, 0, 0, 0, 1, '', '2025-11-25 07:04:33', '', '2025-11-25 09:07:27', 0, '嵌套菜单1-2-2');
INSERT INTO `zeroadmin`.`sys_menu`(`id`, `scope_id`, `parent_id`, `menu_name`, `menu_type`, `path`, `component`, `redirect`, `icon`, `sort`, `no_cache`, `affix`, `external`, `hidden`, `status`, `creator`, `create_time`, `updater`, `update_time`, `del_flag`, `remark`) VALUES (17, 1, 12, 'Menu1-3', 'C', 'menu1-3', 'nested/menu1/menu1-3', '', '', 3, 0, 0, 0, 0, 1, '', '2025-11-25 07:04:33', '', '2025-11-25 09:07:28', 0, '嵌套菜单1-3');
INSERT INTO `zeroadmin`.`sys_menu`(`id`, `scope_id`, `parent_id`, `menu_name`, `menu_type`, `path`, `component`, `redirect`, `icon`, `sort`, `no_cache`, `affix`, `external`, `hidden`, `status`, `creator`, `create_time`, `updater`, `update_time`, `del_flag`, `remark`) VALUES (18, 1, 11, 'Menu2', 'C', 'menu2', 'nested/menu2/index', '', '', 2, 0, 0, 0, 0, 1, '', '2025-11-25 07:04:33', '', '2025-11-25 09:07:30', 0, '嵌套菜单2');
