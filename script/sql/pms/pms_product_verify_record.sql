CREATE TABLE `pms_product_verify_record` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `product_id` bigint(20) NOT NULL COMMENT '商品id',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `review_man` varchar(64) NOT NULL COMMENT '审核人',
    `status` tinyint(4) NOT NULL COMMENT '审核状态：0->未通过；1->通过',
    `detail` varchar(255) NOT NULL COMMENT '反馈详情',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品审核记录';