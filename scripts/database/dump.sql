-- MySQL dump 10.13  Distrib 8.0.31, for macos12 (x86_64)
--
-- Host: 127.0.0.1    Database: dating
-- ------------------------------------------------------
-- Server version	5.5.5-10.4.27-MariaDB-1:10.4.27+maria~ubu2004

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

CREATE OR REPLACE DATABASE dating;
use dating;

--
-- Table structure for table `swipes`
--

DROP TABLE IF EXISTS `swipes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `swipes` (
  `first_user_id` int(11) NOT NULL,
  `second_user_id` int(11) NOT NULL,
  `first_user_swiped` tinyint(1) NOT NULL,
  `second_user_swiped` tinyint(1) NOT NULL,
  `id` int(11) NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`),
  UNIQUE KEY `swipes_UN` (`first_user_id`,`second_user_id`),
  KEY `swipes_firstUserId_IDX` (`first_user_id`,`second_user_id`) USING BTREE,
  KEY `swipes_first_user_id_IDX` (`first_user_id`,`first_user_swiped`) USING BTREE,
  KEY `swipes_second_user_id_IDX` (`second_user_id`,`second_user_swiped`) USING BTREE,
  CONSTRAINT `swipes_FK` FOREIGN KEY (`first_user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `swipes_FK_1` FOREIGN KEY (`second_user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `swipes`
--

LOCK TABLES `swipes` WRITE;
/*!40000 ALTER TABLE `swipes` DISABLE KEYS */;
/*!40000 ALTER TABLE `swipes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `email` varchar(100) NOT NULL,
  `password_hash` binary(200) NOT NULL,
  `name` varchar(100) DEFAULT NULL,
  `gender` smallint(6) NOT NULL,
  `age` smallint(6) NOT NULL,
  `location` point NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `users_UN` (`email`),
  KEY `users_age_IDX` (`age`,`gender`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=43 DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (39,'Andrew.Martin@gmail.com',_binary 'W��PasswordHash��\0Hash\n\0Salt\n\0Time\0Memory\0Threads\0KeyLen\0\0\0A�� ,^e佅�;�Yz{*\�(`\��d$q\�D\�/�ģR{��`q\�$;*Oe��<\0 \0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0','Andrew Martin',1,49,_binary '\0\0\0\0\0\0\0�fd���4@\�4�\\C@'),(40,'Alexander.White@gmail.com',_binary 'W��PasswordHash��\0Hash\n\0Salt\n\0Time\0Memory\0Threads\0KeyLen\0\0\0A�� �{�ˁDz�$�g\�Bl8X��FTKiZY����Y\�\'\�\�\�b2��ZQ��\'��<\0 \0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0','Alexander White',1,34,_binary '\0\0\0\0\0\0\0m\�\�~�T�/3l�\��C@'),(41,'Anthony.Brown@gmail.com',_binary 'W��PasswordHash��\0Hash\n\0Salt\n\0Time\0Memory\0Threads\0KeyLen\0\0\0A�� �Ҍ���7G��p\0m[�\�XĶ\�Po\�\��L닩\�\�\�(��:`�r�E2i�<\0 \0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0','Anthony Brown',1,14,_binary '\0\0\0\0\0\0\0\�\�\�\�6�R@lv��NZf�'),(42,'Avery.Martinez@gmail.com',_binary 'W��PasswordHash��\0Hash\n\0Salt\n\0Time\0Memory\0Threads\0KeyLen\0\0\0A�� �D\�乀¢\�0\�)O�\�\�<d�\�\�Ą\�\�݊\�y��g\��:����<\0 \0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0','Avery Martinez',2,34,_binary '\0\0\0\0\0\0\0�(&o��U@�}s\��X@');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping routines for database 'dating'
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-11-19 13:22:33
