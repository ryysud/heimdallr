-- This file generated by protoc-ddl.
-- Generated by MySQL Generator (v0.1)

SET foreign_key_checks=0;

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
	`id` INTEGER NOT NULL AUTO_INCREMENT,
	`identity` VARCHAR(255) NOT NULL,
	`login_name` VARCHAR(255) NOT NULL,
	`admin` TINYINT(1) NOT NULL,
	`type` VARCHAR(255) NOT NULL,
	`comment` VARCHAR(255) NOT NULL,
	`created_at` DATETIME NOT NULL,
	`updated_at` DATETIME NULL,
	UNIQUE `idx_identity` (`identity`),
	PRIMARY KEY(`id`)
) Engine=InnoDB;

DROP TABLE IF EXISTS `user_state`;
CREATE TABLE `user_state` (
	`id` INTEGER NOT NULL AUTO_INCREMENT,
	`state` VARCHAR(255) NOT NULL,
	`unique` VARCHAR(255) NOT NULL,
	`created_at` DATETIME NOT NULL,
	`updated_at` DATETIME NULL,
	UNIQUE `idx_state` (`state`),
	PRIMARY KEY(`id`)
) Engine=InnoDB;

DROP TABLE IF EXISTS `ssh_key`;
CREATE TABLE `ssh_key` (
	`user_id` INTEGER NOT NULL,
	`key` VARCHAR(255) NOT NULL,
	`created_at` DATETIME NOT NULL,
	`updated_at` DATETIME NULL,
	PRIMARY KEY(`user_id`)
) Engine=InnoDB;

DROP TABLE IF EXISTS `gpg_key`;
CREATE TABLE `gpg_key` (
	`user_id` INTEGER NOT NULL,
	`key` VARCHAR(255) NOT NULL,
	`created_at` DATETIME NOT NULL,
	`updated_at` DATETIME NULL,
	PRIMARY KEY(`user_id`)
) Engine=InnoDB;

DROP TABLE IF EXISTS `role_binding`;
CREATE TABLE `role_binding` (
	`id` INTEGER NOT NULL AUTO_INCREMENT,
	`user_id` INTEGER NOT NULL,
	`role` VARCHAR(255) NOT NULL,
	`maintainer` TINYINT(1) NOT NULL,
	`created_at` DATETIME NOT NULL,
	`updated_at` DATETIME NULL,
	UNIQUE `idx_user_id_role` (`user_id`, `role`),
	PRIMARY KEY(`id`)
) Engine=InnoDB;

DROP TABLE IF EXISTS `access_token`;
CREATE TABLE `access_token` (
	`id` INTEGER NOT NULL AUTO_INCREMENT,
	`name` VARCHAR(255) NOT NULL,
	`value` VARCHAR(255) NOT NULL,
	`user_id` INTEGER NOT NULL,
	`issuer_id` INTEGER NOT NULL,
	`created_at` DATETIME NOT NULL,
	`updated_at` DATETIME NULL,
	UNIQUE `idx_value` (`value`),
	PRIMARY KEY(`id`)
) Engine=InnoDB;

DROP TABLE IF EXISTS `token`;
CREATE TABLE `token` (
	`id` INTEGER NOT NULL AUTO_INCREMENT,
	`token` VARCHAR(255) NOT NULL,
	`user_id` INTEGER NOT NULL,
	`issued_at` DATETIME NOT NULL,
	PRIMARY KEY(`id`)
) Engine=InnoDB;

DROP TABLE IF EXISTS `code`;
CREATE TABLE `code` (
	`id` INTEGER NOT NULL AUTO_INCREMENT,
	`code` VARCHAR(255) NOT NULL,
	`challenge` VARCHAR(255) NOT NULL,
	`challenge_method` VARCHAR(255) NOT NULL,
	`user_id` INTEGER NOT NULL,
	`issued_at` DATETIME NOT NULL,
	UNIQUE `idx_code` (`code`),
	PRIMARY KEY(`id`)
) Engine=InnoDB;

DROP TABLE IF EXISTS `relay`;
CREATE TABLE `relay` (
	`id` INTEGER NOT NULL AUTO_INCREMENT,
	`name` VARCHAR(255) NOT NULL,
	`addr` VARCHAR(255) NOT NULL,
	`from_addr` VARCHAR(255) NOT NULL,
	`connected_at` DATETIME NOT NULL,
	`created_at` DATETIME NOT NULL,
	`updated_at` DATETIME NULL,
	UNIQUE `idx_name_addr` (`name`, `addr`),
	PRIMARY KEY(`id`)
) Engine=InnoDB;

DROP TABLE IF EXISTS `serial_number`;
CREATE TABLE `serial_number` (
	`id` BIGINT NOT NULL AUTO_INCREMENT,
	`serial_number` VARBINARY(20) NOT NULL,
	UNIQUE `idx_serial_number` (`serial_number`),
	PRIMARY KEY(`id`)
) Engine=InnoDB;

DROP TABLE IF EXISTS `signed_certificate`;
CREATE TABLE `signed_certificate` (
	`id` INTEGER NOT NULL AUTO_INCREMENT,
	`certificate` BLOB NOT NULL,
	`serial_number_id` BIGINT NOT NULL,
	`p12` BLOB NOT NULL,
	`agent` TINYINT(1) NOT NULL,
	`comment` VARCHAR(255) NOT NULL,
	`issued_at` DATETIME NOT NULL,
	PRIMARY KEY(`id`)
) Engine=InnoDB;

DROP TABLE IF EXISTS `revoked_certificate`;
CREATE TABLE `revoked_certificate` (
	`id` INTEGER NOT NULL AUTO_INCREMENT,
	`common_name` VARCHAR(255) NOT NULL,
	`serial_number` BLOB NOT NULL,
	`agent` TINYINT(1) NOT NULL,
	`comment` VARCHAR(255) NOT NULL,
	`revoked_at` DATETIME NOT NULL,
	`issued_at` DATETIME NOT NULL,
	`created_at` DATETIME NOT NULL,
	`updated_at` DATETIME NULL,
	PRIMARY KEY(`id`)
) Engine=InnoDB;

DROP TABLE IF EXISTS `node`;
CREATE TABLE `node` (
	`id` INTEGER NOT NULL AUTO_INCREMENT,
	`hostname` VARCHAR(255) NOT NULL,
	`created_at` DATETIME NOT NULL,
	`updated_at` DATETIME NULL,
	PRIMARY KEY(`id`)
) Engine=InnoDB;

SET foreign_key_checks=1;
