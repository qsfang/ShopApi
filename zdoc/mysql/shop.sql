# Dump of table carts
# ------------------------------------------------------------

DROP TABLE IF EXISTS `carts`;

CREATE TABLE `carts` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `productid` int(11) unsigned NOT NULL,
  `name` varchar(200) NOT NULL,
  `count` int(11) unsigned NOT NULL,
  `size`    varchar(50) DEFAULT '0',
  `color`   varchar(50) DEFAULT '0',
  `imageid` int(11) unsigned NOT NULL,
  `userid` int(11) NOT NULL,
  `created` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# Dump of table categories
# ------------------------------------------------------------

DROP TABLE IF EXISTS `categories`;

CREATE TABLE `categories` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(200) NOT NULL DEFAULT '',
  `pid` int(11) NOT NULL DEFAULT '0',
  `status` int(11) NOT NULL,
  `remark` varchar(1000) DEFAULT NULL,
  `created` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# Dump of table contact
# ------------------------------------------------------------
# 收货地址
DROP TABLE IF EXISTS `contact`;

CREATE TABLE `contact` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL DEFAULT '',
  `phone` varchar(20) DEFAULT '',
  `province` int(11) NOT NULL,
  `city` int(11) NOT NULL,
  `street` int(11) NOT NULL,
  `address` varchar(200) NOT NULL DEFAULT '',
  `createdat` datetime NOT NULL,
  `isdefault` TINYINT(1) DEFAULT NULL COMMENT '0:非默认;1:默认',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# Dump of table images
# ------------------------------------------------------------

DROP TABLE IF EXISTS `images`;

CREATE TABLE `images` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `url` varchar(200) NOT NULL DEFAULT '',
  `image` varchar(200) NOT NULL,
  `type` int(11) NOT NULL,
  `title` varchar(100) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# Dump of table orders
# ------------------------------------------------------------

DROP TABLE IF EXISTS `orders`;

CREATE TABLE `orders` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `userid` int(11) unsigned NOT NULL,
  `totalprice` double NOT NULL,
  `payment` double NOT NULL,
  `freight` double DEFAULT '0' COMMENT '运费',
  `remark` text COMMENT '备注',
  `discount` int(11) DEFAULT '0',
  `size`    varchar(50) DEFAULT '0',
  `color`   varchar(50) DEFAULT '0',
  `status` int(11) NOT NULL,
  `created` datetime NOT NULL,
  `payway` INT  NOT NULL ,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# Dump of table products
# ------------------------------------------------------------

DROP TABLE IF EXISTS `products`;

CREATE TABLE `products` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(200) NOT NULL DEFAULT '',
  `totalsale` double NOT NULL DEFAULT '0',
  `categories` int(11) NOT NULL,
  `price` double NOT NULL,
  `originalprice` double NOT NULL,
  `status` int(11) NOT NULL,
  `size`  varchar(200),
  `color` varchar(200),
  `imageid` int(11) unsigned NOT NULL COMMENT '商品封面图片',
  `imageids` varchar(200) NOT NULL DEFAULT '' COMMENT '商品图片集',
  `remark` varchar(1000) DEFAULT '',
  `detail` longtext NOT NULL,
  `created` datetime NOT NULL,
  `inventory` int(11) unsigned NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# Dump of table users
# ------------------------------------------------------------

DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `openid` text,
  `name` varchar(100) DEFAULT NULL,
  `password` varchar(20) NOT NULL DEFAULT '',
  `status` int(11) DEFAULT NULL,
  `type` INT(11)  NOT NULL,
  `created` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

# Dump of table userinfo
# ------------------------------------------------------------

DROP TABLE IF EXISTS `userinfo`;

CREATE TABLE `userinfo` (
  `userid`   INT(11),
  `avatar` text,
  `nickname` VARCHAR(100)         DEFAULT NULL,
  `email`    VARCHAR(100)         DEFAULT NULL,
  `phone`    VARCHAR(20) NOT NULL DEFAULT '',
  `sex`      TINYINT(1)           DEFAULT NULL COMMENT '0:男;1:女',
  PRIMARY KEY (`userid`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;
