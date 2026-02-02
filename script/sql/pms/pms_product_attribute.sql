CREATE TABLE `pms_product_attribute` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `product_attribute_category_id` bigint(20) NOT NULL COMMENT '商品属性分类id',
    `name` varchar(64) NOT NULL COMMENT '商品属性分类id',
    `select_type` tinyint(4) NOT NULL COMMENT '属性选择类型：0->唯一；1->单选；2->多选',
    `input_type` tinyint(4) NOT NULL COMMENT '属性录入方式：0->手工录入；1->从列表中选取',
    `input_list` varchar(255) NOT NULL COMMENT '可选值列表，以逗号隔开',
    `sort` int(11) NOT NULL COMMENT '排序字段：最高的可以单独上传图片',
    `filter_type` tinyint(4) NOT NULL COMMENT '分类筛选样式：1->普通；1->颜色',
    `search_type` tinyint(4) NOT NULL COMMENT '检索类型；0->不需要进行检索；1->关键字检索；2->范围检索',
    `related_status` tinyint(4) NOT NULL COMMENT '相同属性产品是否关联；0->不关联；1->关联',
    `hand_add_status` tinyint(4) NOT NULL COMMENT '是否支持手动新增；0->不支持；1->支持',
    `type` tinyint(4) NOT NULL COMMENT '属性的类型；0->规格；1->参数',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品属性参数表';