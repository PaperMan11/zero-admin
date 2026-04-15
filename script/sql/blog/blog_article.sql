-- 文章表
CREATE TABLE `blog_article` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` DATETIME DEFAULT NULL COMMENT '删除时间',
  `title` VARCHAR(255) NOT NULL COMMENT '文章标题',
  `summary` TEXT DEFAULT NULL COMMENT '文章摘要',
  `content` LONGTEXT NOT NULL COMMENT '文章内容',
  `tags` VARCHAR(500) DEFAULT NULL COMMENT '标签',
  `category_id` BIGINT DEFAULT NULL COMMENT '分类ID',
  `author_id` BIGINT NOT NULL COMMENT '作者ID',
  `view_count` BIGINT DEFAULT 0 COMMENT '浏览量',
  `comment_count` BIGINT DEFAULT 0 COMMENT '评论数',
  `like_count` BIGINT DEFAULT 0 COMMENT '点赞数',
  `cover` VARCHAR(255) DEFAULT NULL COMMENT '封面图片',
  PRIMARY KEY (`id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_author_id` (`author_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_tags` (`tags`(255))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章表';