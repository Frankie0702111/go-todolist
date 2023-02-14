CREATE TABLE IF NOT EXISTS `tasks` (
  `id`            bigint        NOT NULL  AUTO_INCREMENT  PRIMARY KEY,
  `user_id`       bigint        NOT NULL,
  `title`         varchar(255)  NOT NULL  DEFAULT ''      COMMENT '標題',
  `note`          text          NULL      COMMENT '備註',
  `url`           text          NULL      COMMENT '網址',
  `img`           varchar(100)  NULL      COMMENT '影像名稱',
  `img_link`      varchar(255)  NULL      COMMENT '影像路徑',
  `specify_date`  date          NULL  DEFAULT NULL   COMMENT '指定日期',
  `specify_time`  time          NULL  DEFAULT NULL   COMMENT '指定時間',
  `priority`      tinyint       NOT NULL  DEFAULT 0       COMMENT '0:無 1:低 2:中 3:高',
  `progress`      varchar(10)   NOT NULL  DEFAULT 'Draft' COMMENT 'Draft  Processing  Done',
  `is_complete`   bool          NOT NULL  DEFAULT false   COMMENT '是否完成',
  `created_at`    timestamp     NOT NULL  DEFAULT NOW()   COMMENT '新增時間',
  `updated_at`    timestamp     NOT NULL  DEFAULT NOW()   COMMENT '更新時間'
);

create index `idx_user_id_title` on `tasks` (`user_id`, `title`) using BTREE;
create index `idx_is_complete` on `tasks` (`is_complete`) using BTREE;
create index `idx_created_at` on `tasks` (`created_at` desc) using BTREE;
create index `idx_updated_at` on `tasks` (`updated_at` desc) using BTREE;
ALTER TABLE `tasks` ADD CONSTRAINT `users_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE;

INSERT INTO `tasks` (`user_id`, `title`, `note`, `url`, `img`, `img_link`, `specify_date`, `specify_time`, `priority`, `progress`, `is_complete`, `created_at`, `updated_at`)
VALUES
(1, 'Read book', '10 pag to 15 pag', NULL, NULL, NULL, NULL, '12:00:00', 1, 'Draft', 0, '2022-11-29 09:00:00', '2022-11-29 09:00:00'),
(1, 'Go to super market', 'apple', NULL, NULL, NULL, '2022-11-29', NULL, 2, 'Processing', 0, '2022-11-29 09:00:00', '2022-11-29 09:00:00'),
(1, 'dating with friends', 'Frankie, Daisy', 'https://goo.gl/maps/VcDMEzRLKpqPtm697', NULL, NULL, '2022-11-29', '14:00:00', 3, 'Done', 1, '2022-11-29 09:00:00', '2022-11-29 09:00:00');