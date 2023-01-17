ALTER TABLE `category_task` DROP FOREIGN KEY `tasks_task_id_foreign`;
ALTER TABLE `category_task` DROP FOREIGN KEY `categories_category_id_foreign`;
DROP TABLE IF EXISTS `category_task`;
