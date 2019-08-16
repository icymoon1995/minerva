DROP TABLE IF EXISTS `auth`
CREATE TABLE `auth` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `email` varchar(60) NOT NULL DEFAULT '' COMMENT '邮箱',
  `status` int(2) NOT NULL DEFAULT '0' COMMENT '状态',
  `password` varchar(100) NOT NULL DEFAULT '' COMMENT '密码(加密后)',
  `telephone` varchar(15) NOT NULL DEFAULT '' COMMENT '电话',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `index_email` (`email`),
  KEY `index_telephone` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS `users`
CREATE TABLE `users` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL DEFAULT '',
  `avatar` varchar(250) NOT NULL DEFAULT '' COMMENT '头像url',
  `gender` int(1) NOT NULL DEFAULT '3' COMMENT '性别 1男2女3未知',
  `birthday` date NOT NULL COMMENT '出生年月',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `index_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;