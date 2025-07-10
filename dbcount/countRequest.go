package dbcount

import (
	"log"
)

var countReadAccounts int = 0
var countReadAccountStats int = 0
var countReadSessions int = 0
var countReadSessionSaveData int = 0
var countReadActiveClientSessions int = 0
var countReadSystemSaveData int = 0

var countWriteAccounts int = 0
var countWriteAccountStats int = 0
var countWriteSessions int = 0
var countWriteSessionSaveData int = 0
var countWriteActiveClientSessions int = 0
var countWriteSystemSaveData int = 0

var userName int = 0
var discordIdNum int = 0
var googleIdNum int = 0
var lastLoggedInNum int = 0
var keyNum int = 0
var saltNum int = 0
var uuidNum int = 0
var trainerIDNum int = 0
var secretIDNum int = 0
var dbUsernameNum int = 0
var dataNum int = 0
var slotNum int = 0
var tokenNum int = 0
var expireNum int = 0
var hashNum int = 0
var registeredNum int = 0
var bannedNum int = 0
var statsNum int = 0
var timeStampNum int = 0
var idNum int = 0
var clientSessionIdNum int = 0
var sessionSaveDataNum int = 0
var systemDataDeleteNum int = 0
var playtimeNum int = 0
var playerCountNum int = 0
var battleCountNum int = 0
var classicSessionCountNum int = 0
var lastActivityNum int = 0
var unkownTableNum int = 0

func PrintCount() {
	log.Println("Count start test")
	log.Printf("------------------------------------------")
	log.Printf("Read Accounts: %d", countReadAccounts)
	log.Printf("Read AccountStats: %d", countReadAccountStats)
	log.Printf("Read Sessions: %d", countReadSessions)
	log.Printf("Read SessionSaveData: %d", countReadSessionSaveData)
	log.Printf("Read ActiveClientSessions: %d", countReadActiveClientSessions)
	log.Printf("Read SystemSaveData: %d", countReadSystemSaveData)
	log.Printf("Write Accounts: %d", countWriteAccounts)
	log.Printf("Write AccountStats: %d", countWriteAccountStats)
	log.Printf("Write Sessions: %d", countWriteSessions)
	log.Printf("Write SessionSaveData: %d", countWriteSessionSaveData)
	log.Printf("Write ActiveClientSessions: %d", countWriteActiveClientSessions)
	log.Printf("Write SystemSaveData: %d", countWriteSystemSaveData)
	log.Println(" ")
	log.Printf("userName: %d", userName)
	log.Printf("discordIdNum: %d", discordIdNum)
	log.Printf("googleIdNum: %d", googleIdNum)
	log.Printf("lastLoggedInNum: %d", lastLoggedInNum)
	log.Printf("keyNum: %d", keyNum)
	log.Printf("saltNum: %d", saltNum)
	log.Printf("uuidNum: %d", uuidNum)
	log.Printf("trainerIDNum: %d", trainerIDNum)
	log.Printf("secretIDNum: %d", secretIDNum)
	log.Printf("dbUsernameNum: %d", dbUsernameNum)
	log.Printf("dataNum: %d", dataNum)
	log.Printf("slotNum: %d", slotNum)
	log.Printf("tokenNum: %d", tokenNum)
	log.Printf("expireNum: %d", expireNum)
	log.Printf("hashNum: %d", hashNum)
	log.Printf("registeredNum: %d", registeredNum)
	log.Printf("bannedNum: %d", bannedNum)
	log.Printf("statsNum: %d", statsNum)
	log.Printf("timeStampNum: %d", timeStampNum)
	log.Printf("idNum: %d", idNum)
	log.Printf("clientSessionIdNum: %d", clientSessionIdNum)
	log.Printf("sessionSaveDataNum: %d", sessionSaveDataNum)
	log.Printf("systemDataDeleteNum: %d", systemDataDeleteNum)
	log.Printf("playtimeNum: %d", playtimeNum)
	log.Printf("playerCountNum: %d", playerCountNum)
	log.Printf("battleCountNum: %d", battleCountNum)
	log.Printf("classicSessionCountNum: %d", classicSessionCountNum)
	log.Printf("lastActivityNum: %d", lastActivityNum)
	log.Printf("Unknown table name: %d", unkownTableNum)
}

func AddReadCount(tableName string, funcName string) {
	switch tableName {
	case "accounts":
		countReadAccounts++
		if "AddAccountSession" == funcName {
			uuidNum++
			tokenNum++
			expireNum++
			log.Printf("AddAccountSession Read uuid, expire : %d, %d, %d", uuidNum, tokenNum, expireNum)
		}
		if "FetchUsernameByDiscordId" == funcName {
			userName++
			log.Printf("username : %d", userName)
		}
		if "FetchUsernameByGoogleId" == funcName {
			userName++
			log.Printf("username : %d", userName)
		}
		if "FetchDiscordIdByUsername" == funcName {
			discordIdNum++
			log.Printf("discordIdNum : %d", discordIdNum)
		}
		if "FetchGoogleIdByUsername" == funcName {
			googleIdNum++
			log.Printf("discordIdNum : %d", googleIdNum)
		}
		if "FetchDiscordIdByUUID" == funcName {
			discordIdNum++
			log.Printf("discordIdNum : %d", discordIdNum)
		}
		if "FetchGoogleIdByUUID" == funcName {
			googleIdNum++
			log.Printf("googleIdNum : %d", googleIdNum)
		}
		if "FetchUsernameBySessionToken" == funcName {
			userName++
			log.Printf("username : %d", userName)
		}
		if "CheckUsernameExists" == funcName {
			dbUsernameNum++
			log.Printf("dbusernameNum : %d", dbUsernameNum)
		}
		if "FetchLastLoggedInDateByUsername" == funcName {
			lastLoggedInNum++
			log.Printf("lastLoggedIn : %d", lastLoggedInNum)
		}
		if "FetchAdminDetailsByUsername" == funcName {
			userName++
			discordIdNum++
			googleIdNum++
			lastActivityNum++
			registeredNum++
			log.Printf("username : %d, discordIdNum : %d, googleIdNum : %d lastActivity : %d, registered : %d", userName, discordIdNum, googleIdNum, lastActivityNum, registeredNum)
		}
		if "FetchAccountKeySaltFromUsername" == funcName {
			keyNum++
			saltNum++
			log.Printf("keyNum : %d, saltNum : %d", keyNum, saltNum)
		}
		if "FetchTrainerIds" == funcName {
			trainerIDNum++
			secretIDNum++
			log.Printf("trainerId : %d, secretId : %d", trainerIDNum, secretIDNum)
		}
		if "FetchUUIDFromToken" == funcName {
			uuidNum++
			log.Printf("uuidNum : %d", uuidNum)
		}
		if "FetchUsernameFromUUID" == funcName {
			userName++
			log.Printf("username : %d", userName)
		}
		if "FetchUUIDFromUsername" == funcName {
			uuidNum++
			log.Printf("uuidNum : %d", uuidNum)
		}
		//account.go end.

		if "FetchRankingPageCount" == funcName {
			log.Printf("FetchRankingPageCount!")
		}
		if "FetchRankings" == funcName {
			log.Printf("FetchRankings!")
		}
		if "FetchPlayerCount" == funcName {
			playerCountNum++
			log.Printf("FetchPlayerCount : %d", playerCountNum)
		}
		if "FetchBattleCount" == funcName {
			battleCountNum++
			log.Printf("FetchBattleCount : %d", battleCountNum)
			//log.Printf("FetchBattleCount!")
		}
		if "FetchClassicSessionCount" == funcName {
			classicSessionCountNum++
			log.Printf("FetchClassicSessionCount : %d", classicSessionCountNum)
			//log.Printf("FetchClassicSessionCount!")
		}

	case "accountStats":
		countReadAccountStats++
		if "FetchRankings" == funcName {
			log.Printf("FetchRankings!")
		}
		if "RetrievePlaytime" == funcName {
			playtimeNum++
			log.Printf("RetrievePlaytime playtimeNum : %d", playtimeNum)
		}
	case "sessions":
		countReadSessions++
		if "FetchUUIDFromToken" == funcName {
			uuidNum++
			log.Printf("FetchUUIDFromToken uuidNum : %d", uuidNum)
		}
	case "sessionSaveData":
		countReadSessionSaveData++
		if "ReadSessionSaveData" == funcName {
			dataNum++
			log.Printf("ReadSessionSaveData dataNum : %d", dataNum)
			//log.Printf("ReadSessionSaveData!")
		}
		if "GetLatestSessionSaveDataSlot" == funcName {
			slotNum++
			log.Printf("GetLatestSessionSaveDataSlot slotNum : %d", slotNum)
			//log.Printf("GetLatestSessionSaveDataSlot!")
		}

	case "activeClientSessions":
		countReadActiveClientSessions++
		if "IsActiveClientSession" == funcName {
			idNum++
			log.Printf("IsActiveClientSession idNum : %d", idNum)
		}

	case "systemSaveData":
		countReadSystemSaveData++
		if "ReadSystemSaveData" == funcName {
			dataNum++
			log.Printf("ReadSystemSaveData dataNum : %d", dataNum)
			//log.Printf("ReadSystemSaveData!")
		}

	default:
		unkownTableNum++
		log.Println("Unknown table name:", tableName)
	}
}

func AddWriteCount(tableName string, funcName string) {
	switch tableName {
	case "accounts":
		countWriteAccounts++
		if "AddAccountRecord" == funcName {
			uuidNum++
			userName++
			hashNum++
			saltNum++
			registeredNum++
			log.Printf("INSERT uuid, username, hash, salt, registered : %d, %d, %d, %d, %d", uuidNum, userName, hashNum, saltNum, registeredNum)
		}
		if "AddAccountSession" == funcName {
			lastLoggedInNum++
			log.Printf("Update lastLoggedIn : %d", lastLoggedInNum)
		}
		if "AddDiscordIdByUsername" == funcName {
			discordIdNum++
			log.Printf("Update discordId : %d", discordIdNum)
		}
		if "AddGoogleIdByUsername" == funcName {
			googleIdNum++
			log.Printf("Update googleId : %d", googleIdNum)
		}
		if "AddGoogleIdByUUID" == funcName {
			googleIdNum++
			log.Printf("Update googleId : %d", googleIdNum)
		}
		if "AddDiscordIdByUUID" == funcName {
			discordIdNum++
			log.Printf("Update discordId : %d", discordIdNum)
		}
		if "UpdateAccountPassword" == funcName {
			keyNum++
			saltNum++
			log.Printf("Update key, salt : %d, %d", keyNum, saltNum)
		}
		if "UpdateAccountLastActivity" == funcName {
			lastLoggedInNum++
			log.Printf("Update lastActivity : %d", lastLoggedInNum)
		}
		if "SetAccountBanned" == funcName {
			bannedNum++
			log.Printf("Update banned : %d", bannedNum)
		}
		if "UpdateTrainerIds" == funcName {
			trainerIDNum++
			secretIDNum++
			log.Printf("Update trainerId, secretId : %d, %d", trainerIDNum, secretIDNum)
		}
		if "RemoveDiscordIdByUUID" == funcName {
			discordIdNum++
			log.Printf("Update discordId : %d", discordIdNum)
		}
		if "RemoveGoogleIdByUUID" == funcName {
			googleIdNum++
			log.Printf("Updat googleId : %d", googleIdNum)
		}
		if "RemoveGoogleIdByUsername" == funcName {
			googleIdNum++
			log.Printf("Update googleId : %d", googleIdNum)
		}
		if "RemoveDiscordIdByUsername" == funcName {
			discordIdNum++
			log.Printf("Update discordId : %d", discordIdNum)
		}
		if "RemoveDiscordIdByDiscordId" == funcName {
			discordIdNum++
			log.Printf("Update discordId : %d", discordIdNum)
		}
		if "RemoveGoogleIdByDiscordId" == funcName {
			googleIdNum++
			log.Printf("Update googleId : %d", googleIdNum)
		}

	case "accountStats":
		countWriteAccountStats++
		if "UpdateAccountStats" == funcName {
			statsNum++
			log.Printf("Insert stats : %d", statsNum)
		}

	case "sessions":
		countWriteSessions++
		if "AddAccountSession" == funcName {
			uuidNum++
			tokenNum++
			expireNum++
			log.Printf("INSERT uuid, token, expire : %d, %d, %d", uuidNum, tokenNum, expireNum)
		}
		if "RemoveSessionFromToken" == funcName {
			tokenNum++
			log.Printf("Remove token : %d", tokenNum)
		}

	case "sessionSaveData":
		countWriteSessionSaveData++
		if "StoreSessionSaveData" == funcName {
			uuidNum++
			slotNum++
			dataNum++
			timeStampNum++
			log.Printf("StoreSessionSaveData uuid, slot, data, timestamp : %d, %d, %d, %d", uuidNum, slotNum, dataNum, timeStampNum)
		}
		if "DeleteSessionSaveData" == funcName {
			sessionSaveDataNum++
			log.Printf("DeleteSessionSaveData sessionSaveDataNum : %d", sessionSaveDataNum)
		}

	case "activeClientSessions":
		countWriteActiveClientSessions++
		if "UpdateActiveSession" == funcName {
			uuidNum++
			clientSessionIdNum++
			log.Printf("Insert uuid, clientSessionId : %d, %d", uuidNum, clientSessionIdNum)
		}

	case "systemSaveData":
		countWriteSystemSaveData++
		if "StoreSystemSaveData" == funcName {
			uuidNum++
			dataNum++
			timeStampNum++
			log.Printf("StoreSystemSaveData uuid, data, timestamp : %d, %d, %d", uuidNum, dataNum, timeStampNum)
		}
		if "DeleteSystemSaveData" == funcName {
			systemDataDeleteNum++
			log.Printf("DeleteSystemSaveData systemDataDeleteNum : %d", systemDataDeleteNum)
		}

	default:
		unkownTableNum++
		log.Println("Unknown table name:", tableName)
	}
}
