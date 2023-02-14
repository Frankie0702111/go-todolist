CREATE TABLE IF NOT EXISTS `category_task` (
  `id`          bigint  NOT NULL AUTO_INCREMENT  PRIMARY KEY,
  `task_id`     bigint  NOT NULL,
  `category_id` bigint  NOT NULL,
  `created_at`  timestamp     NOT NULL  DEFAULT NOW()   COMMENT '新增時間',
  `updated_at`  timestamp     NOT NULL  DEFAULT NOW()   COMMENT '更新時間'
);

ALTER TABLE `category_task` ADD CONSTRAINT `tasks_task_id_foreign` FOREIGN KEY (`task_id`) REFERENCES `tasks`(`id`) ON DELETE CASCADE;
ALTER TABLE `category_task` ADD CONSTRAINT `categories_category_id_foreign` FOREIGN KEY (`category_id`) REFERENCES `categories`(`id`) ON DELETE CASCADE;

INSERT INTO `category_task` (`task_id`, `category_id`, `created_at`, `updated_at`)
VALUES
(1, 2, '2022-11-29 09:00:00', '2022-11-29 09:00:00'),
(2, 1, '2022-11-29 09:00:00', '2022-11-29 09:00:00'),
(3, 3, '2022-11-29 09:00:00', '2022-11-29 09:00:00');