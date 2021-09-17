-- MySQL dump 10.13  Distrib 8.0.11, for Win64 (x86_64)
--
-- Host: localhost    Database: userschat
-- ------------------------------------------------------
-- Server version	8.0.11

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
 SET NAMES utf8mb4 ;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `allusers`
--

DROP TABLE IF EXISTS `allusers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `allusers` (
  `usersId` int(11) NOT NULL AUTO_INCREMENT,
  `chatId` int(11) NOT NULL,
  `useBotFunc` tinyint(1) NOT NULL,
  PRIMARY KEY (`usersId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `allusers`
--

LOCK TABLES `allusers` WRITE;
/*!40000 ALTER TABLE `allusers` DISABLE KEYS */;
/*!40000 ALTER TABLE `allusers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `notif_change_cost`
--

DROP TABLE IF EXISTS `notif_change_cost`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `notif_change_cost` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `ChatId` int(11) NOT NULL,
  `Change_cost` decimal(7,2) DEFAULT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `notif_change_cost`
--

LOCK TABLES `notif_change_cost` WRITE;
/*!40000 ALTER TABLE `notif_change_cost` DISABLE KEYS */;
INSERT INTO `notif_change_cost` VALUES (19,1561775851,100.00);
/*!40000 ALTER TABLE `notif_change_cost` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `notifcost`
--

DROP TABLE IF EXISTS `notifcost`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `notifcost` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `ChatId` int(11) NOT NULL,
  `Cost` decimal(7,2) DEFAULT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `notifcost`
--

LOCK TABLES `notifcost` WRITE;
/*!40000 ALTER TABLE `notifcost` DISABLE KEYS */;
/*!40000 ALTER TABLE `notifcost` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `userschat`
--

DROP TABLE IF EXISTS `userschat`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
 SET character_set_client = utf8mb4 ;
CREATE TABLE `userschat` (
  `usersId` int(11) NOT NULL AUTO_INCREMENT,
  `chatId` int(11) NOT NULL,
  `useBotFunc` varchar(1) NOT NULL,
  PRIMARY KEY (`usersId`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `userschat`
--

LOCK TABLES `userschat` WRITE;
/*!40000 ALTER TABLE `userschat` DISABLE KEYS */;
INSERT INTO `userschat` VALUES (6,613870799,'0'),(11,1561775851,'1'),(12,417406411,'1'),(13,860762049,'0');
/*!40000 ALTER TABLE `userschat` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-09-17  5:40:13
