CREATE TABLE `pms_comment` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `product_id` bigint(20) NOT NULL COMMENT '商品id',
    `parent_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '父评论 ID（用于回复）0 表示顶级评论',
    `user_id` bigint(20) NOT NULL COMMENT '评价者id',
    `reply_user_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '回复用户',
    `star` int(3) NOT NULL COMMENT '评价星数：0->5',
    `user_ip` varchar(64) NOT NULL COMMENT '评价的ip',
    `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '评价时间',
    `show_status` tinyint(4) NOT NULL COMMENT '是否显示，0-不显示，1-显示',
    `collect_count` int(11) NOT NULL COMMENT '点赞数',
    `content` text NOT NULL COMMENT '内容',
    `pics` varchar(1000) NOT NULL COMMENT '上传图片地址，以逗号隔开',
    `replay_count` int(11) NOT NULL COMMENT '回复数量',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品评价表';