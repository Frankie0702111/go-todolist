CREATE TABLE IF NOT EXISTS `category_task` (
  `id`          int NOT NULL AUTO_INCREMENT  PRIMARY KEY,
  `task_id`     int NOT NULL,
  `category_id` int NOT NULL
);

ALTER TABLE `category_task` ADD CONSTRAINT `tasks_task_id_foreign` FOREIGN KEY (`task_id`) REFERENCES `tasks`(`id`) ON DELETE CASCADE;
ALTER TABLE `category_task` ADD CONSTRAINT `categories_category_id_foreign` FOREIGN KEY (`category_id`) REFERENCES `categories`(`id`) ON DELETE CASCADE;

INSERT INTO `category_task` (`task_id`, `category_id`)
VALUES
(1, 2),
(2, 1),
(3, 3);