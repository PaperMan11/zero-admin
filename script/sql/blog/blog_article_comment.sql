-- 文章评论表
CREATE TABLE `blog_article_comment` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
  `article_id` BIGINT NOT NULL COMMENT '文章ID',
  `parent_id` BIGINT DEFAULT 0 COMMENT '父评论ID，0表示一级评论',
  `reply_id` bigint(20) DEFAULT '0' COMMENT '回复评论ID',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `reply_to` BIGINT DEFAULT 0 COMMENT '回复评论ID，0表示不是回复评论',
  `content` TEXT NOT NULL COMMENT '评论内容',
  `images` TEXT DEFAULT NULL COMMENT '图片',
  `reply_count` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '回复数',
  `read` tinyint(1) NOT NULL DEFAULT '0' COMMENT '0-未读 1-已读',
  PRIMARY KEY (`id`),
  KEY `idx_article_id` (`article_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_parent_id` (`parent_id`),
  KEY `idx_reply_to` (`reply_to`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_article_id_created_at` (`article_id`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章评论表';