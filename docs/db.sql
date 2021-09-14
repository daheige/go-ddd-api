CREATE DATABASE `test` CHARACTER SET utf8mb4;

CREATE TABLE `news` (
                        `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
                        `created_at` datetime DEFAULT NULL,
                        `updated_at` datetime DEFAULT NULL,
                        `deleted_at` datetime DEFAULT NULL,
                        `title` varchar(255) DEFAULT NULL,
                        `slug` varchar(255) DEFAULT NULL,
                        `content` text,
                        `status` varchar(255) DEFAULT NULL,
                        PRIMARY KEY (`id`),
                        KEY `idx_news_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4;

CREATE TABLE `topics` (
                          `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
                          `created_at` datetime DEFAULT NULL,
                          `updated_at` datetime DEFAULT NULL,
                          `deleted_at` datetime DEFAULT NULL,
                          `name` varchar(255) DEFAULT NULL,
                          `slug` varchar(255) DEFAULT NULL,
                          PRIMARY KEY (`id`),
                          KEY `idx_topics_deleted_at` (`deleted_at`)
) ENGINE = InnoDB CHARSET = utf8mb4;

CREATE TABLE `news_topics` (
                               `news_id` int(10) UNSIGNED NOT NULL,
                               `topic_id` int(10) UNSIGNED NOT NULL,
                               PRIMARY KEY (`news_id`, `topic_id`)
) ENGINE = InnoDB CHARSET = utf8mb4;
