-- MySQL dump 10.13  Distrib 5.5.41-37.0, for Linux (x86_64)
--
-- Host: localhost    Database: mixer
-- ------------------------------------------------------
-- Server version	5.5.41-37.0-log

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `groups`
--

DROP TABLE IF EXISTS `groups`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `groups` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `odd_person_out`
--

DROP TABLE IF EXISTS `odd_person_out`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `odd_person_out` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `person` bigint(20) NOT NULL DEFAULT '0',
  `week_id` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_person` (`person`),
  KEY `fk_week_id` (`week_id`),
  CONSTRAINT `odd_person_out_ibfk_1` FOREIGN KEY (`person`) REFERENCES `people` (`id`) ON DELETE CASCADE,
  CONSTRAINT `odd_person_out_ibfk_2` FOREIGN KEY (`week_id`) REFERENCES `week` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `pairs`
--

DROP TABLE IF EXISTS `pairs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pairs` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `person1` bigint(20) NOT NULL DEFAULT '0',
  `person2` bigint(20) NOT NULL DEFAULT '0',
  `week_id` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_week_id` (`week_id`),
  KEY `fk_person1` (`person1`),
  KEY `fk_person2` (`person2`),
  CONSTRAINT `pairs_ibfk_1` FOREIGN KEY (`week_id`) REFERENCES `week` (`id`) ON DELETE CASCADE,
  CONSTRAINT `pairs_ibfk_2` FOREIGN KEY (`person1`) REFERENCES `people` (`id`) ON DELETE CASCADE,
  CONSTRAINT `pairs_ibfk_3` FOREIGN KEY (`person2`) REFERENCES `people` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `people`
--

DROP TABLE IF EXISTS `people`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `people` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `username` text NOT NULL,
  `assistant` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`(20))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `people_groups`
--

DROP TABLE IF EXISTS `people_groups`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `people_groups` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `people_id` bigint(20) NOT NULL,
  `groups_id` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `people_groups_idx1` (`people_id`,`groups_id`),
  KEY `fk_groups_id` (`groups_id`),
  CONSTRAINT `people_groups_ibfk_1` FOREIGN KEY (`people_id`) REFERENCES `people` (`id`) ON DELETE CASCADE,
  CONSTRAINT `people_groups_ibfk_2` FOREIGN KEY (`groups_id`) REFERENCES `groups` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `staff`
--

DROP TABLE IF EXISTS `staff`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `staff` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `auth_username` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `staff_id` bigint(20) unsigned NOT NULL DEFAULT '0',
  `first_name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `last_name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `title` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `is_manager` tinyint(3) unsigned NOT NULL DEFAULT '0',
  `avatar` varchar(127) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `enabled` tinyint(3) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `staff_id` (`staff_id`),
  UNIQUE KEY `auth_username` (`auth_username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `week`
--

DROP TABLE IF EXISTS `week`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `week` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `number` bigint(20) NOT NULL,
  `date` bigint(20) DEFAULT NULL,
  `group_id` bigint(20) NOT NULL DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2015-09-15 20:30:33
