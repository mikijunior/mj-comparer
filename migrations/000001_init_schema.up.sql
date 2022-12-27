CREATE TABLE `auth` (
    `id` bigint(20) NOT NULL AUTO_INCREMENT,
    `api-key` varchar(32) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;