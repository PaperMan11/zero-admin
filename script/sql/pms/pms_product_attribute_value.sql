CREATE TABLE `pms_product_attribute_value` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `product_id` bigint(20) NOT NULL COMMENT '商品id',
    `product_attribute_id` bigint(20) NOT NULL COMMENT '商品属性id',
    `value` varchar(64) NOT NULL COMMENT '手动添加规格或参数的值，参数单值，规格有多个时以逗号隔开',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='存储产品参数信息的表';