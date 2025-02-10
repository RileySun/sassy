--Create
CREATE TABLE `Images` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `model_id` int(11) DEFAULT NULL,
  `path` varchar(256) DEFAULT NULL,
  `desc` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `Models` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) DEFAULT NULL,
  `desc` varchar(2560) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `Videos` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `model_id` int(11) DEFAULT NULL,
  `path` varchar(256) DEFAULT NULL,
  `desc` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;

--Add Base Data