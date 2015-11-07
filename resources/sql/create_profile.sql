delimiter $$

CREATE TABLE `profile` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `email` varchar(256) NOT NULL,
  `created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `county` varchar(45) NOT NULL,
  `type` varchar(45) NOT NULL,
  `profilePic` varchar(256) DEFAULT 'https://s3-eu-west-1.amazonaws.com/localsie/profile-placeholder.gif',
  `bio` varchar(1024) DEFAULT 'bio info',
  `interests` varchar(256) DEFAULT NULL,
  `phone` varchar(32) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email_UNIQUE` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=latin1$$

