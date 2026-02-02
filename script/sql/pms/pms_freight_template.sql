CREATE TABLE `pms_freight_template` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `name` varchar(64) NOT NULL COMMENT '运费模版名称',
    `charge_type` tinyint(4) NOT NULL COMMENT '计费类型:0->按重量；1->按件数',
    `first_weight` bigint(20) NOT NULL COMMENT '首重kg',
    `first_fee` bigint(20) NOT NULL COMMENT '首费（元）',
    `continue_weight` bigint(20) NOT NULL COMMENT '续重kg',
    `continue_fee` bigint(20) NOT NULL COMMENT '续费（元）',
    `dest` varchar(255) NOT NULL COMMENT '目的地（省、市）',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='运费模版';