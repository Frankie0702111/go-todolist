CREATE TABLE IF NOT EXISTS `categories` (
  `id`    int           NOT NULL  AUTO_INCREMENT  PRIMARY KEY,
  `name`  varchar(100)  NOT NULL  DEFAULT ''      COMMENT '類別名稱'
);

INSERT INTO `categories` (`name`)
VALUES
("提醒"),
("工作"),
("活動");