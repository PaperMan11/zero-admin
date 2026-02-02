CREATE TABLE `pms_product_category_attribute_relation` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `product_category_id` bigint(20) NOT NULL COMMENT '商品分类id',
    `product_attribute_id` bigint(20) NOT NULL COMMENT '商品属性id',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='产品的分类和属性的关系表，用于设置分类筛选条件（只支持一级分类）';