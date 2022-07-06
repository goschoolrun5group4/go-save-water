CREATE DATABASE  IF NOT EXISTS `gsw_db` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `gsw_db`;
-- MySQL dump 10.13  Distrib 8.0.29, for macos12 (x86_64)
--
-- Host: gosavewater-db.ctc5pp7q4xkj.ap-southeast-1.rds.amazonaws.com    Database: gsw_db
-- ------------------------------------------------------
-- Server version	8.0.28

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
SET @MYSQLDUMP_TEMP_LOG_BIN = @@SESSION.SQL_LOG_BIN;
SET @@SESSION.SQL_LOG_BIN= 0;

--
-- GTID state at the beginning of the backup 
--

SET @@GLOBAL.GTID_PURGED=/*!80000 '+'*/ '';

--
-- Table structure for table `Address`
--

DROP TABLE IF EXISTS `Address`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Address` (
  `AccountNumber` int NOT NULL AUTO_INCREMENT,
  `PostalCode` varchar(6) NOT NULL,
  `Floor` varchar(6) NOT NULL,
  `UnitNumber` varchar(6) NOT NULL,
  `BuildingName` varchar(255) DEFAULT NULL,
  `BlockNumber` varchar(6) NOT NULL,
  `CreatedDT` datetime NOT NULL,
  `ModifiedDT` datetime DEFAULT NULL,
  `Street` varchar(255) NOT NULL,
  PRIMARY KEY (`AccountNumber`)
) ENGINE=InnoDB AUTO_INCREMENT=1423242435 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Address`
--

LOCK TABLES `Address` WRITE;
/*!40000 ALTER TABLE `Address` DISABLE KEYS */;
/*!40000 ALTER TABLE `Address` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Reward`
--

DROP TABLE IF EXISTS `Reward`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Reward` (
  `RewardID` int NOT NULL AUTO_INCREMENT,
  `Title` varchar(45) NOT NULL,
  `Description` varchar(45) NOT NULL,
  `Quantity` int NOT NULL,
  `RedeemAmt` int NOT NULL,
  `IsDeleted` tinyint NOT NULL DEFAULT '0',
  PRIMARY KEY (`RewardID`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Reward`
--

LOCK TABLES `Reward` WRITE;
/*!40000 ALTER TABLE `Reward` DISABLE KEYS */;
INSERT INTO `Reward` VALUES (1,'FairPrice S$5 Voucher','Enjoy S$5 FairPrice Voucher',87,500,0),(2,'FairPrice S$10 Voucher','Enjoy S$10 FairPrice Voucher',0,1000,0),(3,'Esso Synergy S$40 Fuel Voucher','Enjoy Esso Synergy S$40 Fuel Voucher',98,4000,0),(4,'Toys\"R\"Us S$10 Voucher','Enjoy Toys\"R\"Us S$10 Voucher',96,1000,0),(5,'Isetan S$20 Voucher','Enjoy Isetan S$20 Voucher',99,2000,0),(6,'Takashimaya S$10 Voucher','Enjoy Takashimaya Department Store S$10 Vouch',99,1000,0);
/*!40000 ALTER TABLE `Reward` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `Transaction`
--

DROP TABLE IF EXISTS `Transaction`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Transaction` (
  `TransactionID` int NOT NULL AUTO_INCREMENT,
  `UserID` int NOT NULL,
  `Type` varchar(10) NOT NULL,
  `RewardID` int DEFAULT NULL,
  `Points` int NOT NULL,
  `TransactionDT` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`TransactionID`),
  KEY `fk_Transection_UserID_idx` (`UserID`),
  KEY `fk_Transection_RewardID_idx` (`RewardID`),
  CONSTRAINT `fk_Transaction_RewardID` FOREIGN KEY (`RewardID`) REFERENCES `Reward` (`RewardID`),
  CONSTRAINT `fk_Transaction_UserID` FOREIGN KEY (`UserID`) REFERENCES `User` (`UserID`)
) ENGINE=InnoDB AUTO_INCREMENT=73 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `Transaction`
--

LOCK TABLES `Transaction` WRITE;
/*!40000 ALTER TABLE `Transaction` DISABLE KEYS */;
/*!40000 ALTER TABLE `Transaction` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `User`
--

DROP TABLE IF EXISTS `User`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `User` (
  `UserID` int NOT NULL AUTO_INCREMENT,
  `FirstName` varchar(45) NOT NULL,
  `LastName` varchar(45) NOT NULL,
  `UserName` varchar(45) NOT NULL,
  `Password` varchar(60) NOT NULL,
  `Email` varchar(100) NOT NULL,
  `Role` varchar(10) NOT NULL,
  `IsDeleted` tinyint NOT NULL,
  `CreatedDT` datetime NOT NULL,
  `ModifiedDT` datetime DEFAULT NULL,
  `Verified` varchar(45) NOT NULL DEFAULT '0',
  `PointBalance` int DEFAULT '0',
  PRIMARY KEY (`UserID`),
  UNIQUE KEY `UserName_UNIQUE` (`UserName`)
) ENGINE=InnoDB AUTO_INCREMENT=34 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `User`
--

LOCK TABLES `User` WRITE;
/*!40000 ALTER TABLE `User` DISABLE KEYS */;
/*!40000 ALTER TABLE `User` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `UserAddress`
--

DROP TABLE IF EXISTS `UserAddress`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `UserAddress` (
  `UserID` int NOT NULL,
  `AccountNumber` int NOT NULL,
  PRIMARY KEY (`UserID`),
  KEY `UserID_idx` (`UserID`),
  KEY `AddressID_idx` (`AccountNumber`),
  CONSTRAINT `AddressID` FOREIGN KEY (`AccountNumber`) REFERENCES `Address` (`AccountNumber`),
  CONSTRAINT `UserID` FOREIGN KEY (`UserID`) REFERENCES `User` (`UserID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `UserAddress`
--

LOCK TABLES `UserAddress` WRITE;
/*!40000 ALTER TABLE `UserAddress` DISABLE KEYS */;
/*!40000 ALTER TABLE `UserAddress` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `UserSession`
--

DROP TABLE IF EXISTS `UserSession`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `UserSession` (
  `UserID` int NOT NULL,
  `SessionID` varchar(36) NOT NULL,
  `ExpireDT` datetime NOT NULL,
  `LoginDT` datetime NOT NULL,
  PRIMARY KEY (`UserID`),
  CONSTRAINT `fk_UserSession_User_UserID` FOREIGN KEY (`UserID`) REFERENCES `User` (`UserID`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `UserSession`
--

LOCK TABLES `UserSession` WRITE;
/*!40000 ALTER TABLE `UserSession` DISABLE KEYS */;
/*!40000 ALTER TABLE `UserSession` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `WaterUsage`
--

DROP TABLE IF EXISTS `WaterUsage`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `WaterUsage` (
  `AccountNumber` int NOT NULL,
  `BillDate` date NOT NULL,
  `Consumption` decimal(5,1) NOT NULL,
  `ImageURL` varchar(100) DEFAULT NULL,
  `CreatedDT` datetime NOT NULL,
  `ModifiedDT` datetime DEFAULT NULL,
  PRIMARY KEY (`AccountNumber`,`BillDate`),
  KEY `AccountNumber_idx` (`AccountNumber`),
  CONSTRAINT `FK_AccountNumber` FOREIGN KEY (`AccountNumber`) REFERENCES `Address` (`AccountNumber`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `WaterUsage`
--

LOCK TABLES `WaterUsage` WRITE;
/*!40000 ALTER TABLE `WaterUsage` DISABLE KEYS */;
/*!40000 ALTER TABLE `WaterUsage` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping events for database 'gsw_db'
--
/*!50106 SET @save_time_zone= @@TIME_ZONE */ ;
/*!50106 DROP EVENT IF EXISTS `eventDeleteUserSession` */;
DELIMITER ;;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;;
/*!50003 SET character_set_client  = utf8mb4 */ ;;
/*!50003 SET character_set_results = utf8mb4 */ ;;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;;
/*!50003 SET sql_mode              = 'NO_ENGINE_SUBSTITUTION' */ ;;
/*!50003 SET @saved_time_zone      = @@time_zone */ ;;
/*!50003 SET time_zone             = 'Asia/Singapore' */ ;;
/*!50106 CREATE*/ /*!50117 DEFINER=`admin`@`%`*/ /*!50106 EVENT `eventDeleteUserSession` ON SCHEDULE EVERY 1 MINUTE STARTS '2022-06-19 17:04:08' ON COMPLETION NOT PRESERVE ENABLE DO DELETE FROM UserSession WHERE ExpireDT < NOW() */ ;;
/*!50003 SET time_zone             = @saved_time_zone */ ;;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;;
/*!50003 SET character_set_client  = @saved_cs_client */ ;;
/*!50003 SET character_set_results = @saved_cs_results */ ;;
/*!50003 SET collation_connection  = @saved_col_connection */ ;;
DELIMITER ;
/*!50106 SET TIME_ZONE= @save_time_zone */ ;

--
-- Dumping routines for database 'gsw_db'
--
/*!50003 DROP PROCEDURE IF EXISTS `spAddressCreate` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`admin`@`%` PROCEDURE `spAddressCreate`(
	pUserID INT,
	pAccountNumber INT,
    pPostalCode VARCHAR(6),
    pFloor VARCHAR(6),
    pUnitNumber VARCHAR(6),
    pBuildingName VARCHAR(255),
    pBlockNumber VARCHAR(6),
    pStreet VARCHAR(255)
)
BEGIN

DECLARE insertVal INTEGER DEFAULT 0;

INSERT INTO Address (AccountNumber, PostalCode, Floor, UnitNumber, BuildingName, BlockNumber, Street, CreatedDT, ModifiedDT)
VALUES (pAccountNumber, pPostalCode, pFloor, pUnitNumber, pBuildingName, pBlockNumber, pStreet, NOW(), null);

SELECT ROW_COUNT() INTO insertVal;

IF insertVal > 0 THEN
	INSERT INTO UserAddress (UserID, AccountNumber)
    VALUES (pUserID, pAccountNumber);
END IF;

END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spAuthenticationGet` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`admin`@`%` PROCEDURE `spAuthenticationGet`(pUsername VARCHAR(45))
BEGIN
	SELECT u.Password, u.Verified, u.Email
    FROM User u
    WHERE u.UserName = pUsername
    AND u.IsDeleted = false;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spNationalWaterUsageGetByDateRange` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`admin`@`%` PROCEDURE `spNationalWaterUsageGetByDateRange`(pStartDate date, pEndDate date)
BEGIN
	SELECT DATE_FORMAT(BillDate, "%Y-%m") AS YearMonth, AVG(Consumption) AS Consumption FROM WaterUsage
	WHERE BillDate BETWEEN pStartDate AND pEndDate;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spUserCreate` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`admin`@`%` PROCEDURE `spUserCreate`(
	pUserName VARCHAR(45),
    pFirstName VARCHAR(45),
    pLastName VARCHAR(45),
    pEmail VARCHAR(100),
    pPassword VARCHAR(60),
    pRole VARCHAR(10)
)
BEGIN

INSERT INTO User (UserName, FirstName, LastName, Email, Password, IsDeleted, Role, CreatedDT, ModifiedDT)
VALUES (pUsername, pFirstName, pLastName, pEmail, pPassword, false, pRole, NOW(), null);

END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spUserExistByUserID` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`admin`@`%` PROCEDURE `spUserExistByUserID`(pUserID INT)
BEGIN
	SELECT EXISTS (SELECT * FROM User WHERE UserID = pUserID);
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spUserExistByUserName` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`admin`@`%` PROCEDURE `spUserExistByUserName`(pUserName varchar(45))
BEGIN
	SELECT EXISTS (SELECT * FROM User WHERE UserName = pUserName);
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spUserSessionCreate` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`admin`@`%` PROCEDURE `spUserSessionCreate`(pUsername VARCHAR(45))
BEGIN
	DECLARE pUserID INT DEFAULT 0;
    
    SELECT UserID
    INTO pUserID
    FROM User
    WHERE UserName = pUsername;

	IF (SELECT COUNT(*) FROM UserSession WHERE UserID = pUserID) > 0 THEN
		DELETE FROM UserSession WHERE UserID = pUserID;
	END IF;
     
	INSERT INTO UserSession(UserID, SessionID, ExpireDT, LoginDT)
    VALUES (pUserID, UUID(), NOW() + INTERVAL 1 DAY, NOW());
    
    SELECT UserID, pUsername, SessionID, ExpireDT FROM UserSession WHERE UserID = pUserID;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spUserSessionDelete` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`admin`@`%` PROCEDURE `spUserSessionDelete`(pUserID INT, pSessionID VARCHAR(36))
BEGIN

	IF (SELECT COUNT(*) FROM UserSession WHERE UserID = pUserID AND SessionID = pSessionID) > 0 THEN
		DELETE FROM UserSession WHERE UserID = pUserID;
	END IF;

END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spUserSessionGet` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`admin`@`%` PROCEDURE `spUserSessionGet`(pSessionID VARCHAR(36))
BEGIN
	SELECT u.UserID, u.Username, u.FirstName, u.LastName, u.Email, u.Role, ua.AccountNumber, u.PointBalance
	FROM UserSession us
	INNER JOIN User u
	ON us.UserID = u.UserID
    LEFT JOIN UserAddress ua
    ON us.UserID = ua.UserID
	WHERE us.SessionID = pSessionID
	AND u.Verified = true
	AND u.IsDeleted = false;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spUserUpdate` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`admin`@`%` PROCEDURE `spUserUpdate`(
	pUserID INT,
	pFirstName VARCHAR(45),
    pLastName VARCHAR(45),
	pUserName VARCHAR(45),
    pPassword VARCHAR(60),
    pEmail VARCHAR(100),
    pRole VARCHAR(45),
    pIsDeleted BOOL,
    pVerified BOOL,
    pPointBalance INT
)
BEGIN
	UPDATE User SET
		FirstName = COALESCE(pFirstName, FirstName),
        LastName = COALESCE(pLastName, LastName),
        UserName = COALESCE(pUserName, UserName),
        Password = COALESCE(pPassword, Password),
        Email = COALESCE(pEmail, Email),
        Role = COALESCE(pRole, Role),
        IsDeleted = COALESCE(pIsDeleted, IsDeleted),
        Verified = COALESCE(pVerified, Verified),
        PointBalance = COALESCE(pPointBalance, PointBalance),
        ModifiedDT = NOW()
	WHERE UserID = pUserID;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `spWaterUsageGetByDateRange` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_0900_ai_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE DEFINER=`admin`@`%` PROCEDURE `spWaterUsageGetByDateRange`(pAccountNumber int, pStartDate date, pEndDate date)
BEGIN
	SELECT BillDate, Consumption FROM WaterUsage
	WHERE BillDate BETWEEN pStartDate AND pEndDate
	AND AccountNumber = pAccountNumber;
END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
SET @@SESSION.SQL_LOG_BIN = @MYSQLDUMP_TEMP_LOG_BIN;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-07-06 15:30:04
