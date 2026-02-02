CREATE TABLE `pms_product_attribute_category` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `name` varchar(64) NOT NULL COMMENT '商品属性分类名称',
    `attribute_count` int(11) NOT NULL DEFAULT '0' COMMENT '属性数量',
    `param_count` int(11) NOT NULL DEFAULT '0' COMMENT '参数数量',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='产品属性分类表';