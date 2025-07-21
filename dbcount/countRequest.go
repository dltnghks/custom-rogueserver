package dbcount

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
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

var userNameNum int = 0
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
var classicSessionPlayedCountNum int = 0
var lastActivityNum int = 0
var unkownTableNum int = 0
var fetchRankingPageCountNum int = 0
var fetchRankingsNum int = 0
var countdailyRuns int = 0
var getDailyRunSeedNum int = 0

var userName string = "userName"
var discordId string = "discordId"
var googleId string = "googleId"
var lastLoggedIn string = "lastLoggedIn"
var key string = "key"
var salt string = "salt"
var uuid string = "uuid"
var trainerID string = "trainerID"
var secretID string = "secretID"
var dbUsername string = "dbUsername"
var data string = "data"
var slot string = "slot"
var token string = "token"
var expire string = "expire"
var hash string = "hash"
var registered string = "registered"
var banned string = "banned"
var stats string = "stats"
var timeStamp string = "timeStamp"
var id string = "id"
var clientSessionId string = "clientSessionId"
var sessionSaveData string = "sessionSaveData"
var systemDataDelete string = "systemDataDelete"
var playtime string = "playtime"
var playerCount string = "playerCount"
var battleCount string = "battleCount"
var classicSessionPlayedCount string = "classicSessionPlayedCount"
var lastActivity string = "lastActivity"
var unkownTable string = "unkownTable"
var fetchRankingPageCount string = "fetchRankingPageCount"
var fetchRankings string = "fetchRankings"

// var countdailyRuns string = ""  // 주석처리 되어 있어 유지
var getDailyRunSeed string = "getDailyRunSeed"

var initTime time.Time
var csvWriter *csv.Writer
var csvFile *os.File

func InitTimer() error {
	initTime = time.Now()
	log.Printf("init 시작 시간 : %s", initTime.Format("time.RFC3339Nano"))

	var err error

	err = os.MkdirAll("/home/user/duream/Write-back/custom-rogueserver/csv", 0777)
	if err != nil {
		log.Fatalf("Failed to create directory: %v", err)
		return fmt.Errorf("디렉토리 생성 실패: %w", err)
	}

	csvFile, err = os.Create("/app/csv/player2_db_access_count.csv")
	if err != nil {
		log.Fatalf("Failed to create CSV file: %v", err)
	}

	csvWriter = csv.NewWriter(csvFile)

	err = csvWriter.Write([]string{"Table Name", "Function Name", "Time", "R/W", "data1", "data2", "data3", "data4", "data5"})
	if err != nil {
		return fmt.Errorf("CSV 헤더 작성 실패: %w", err)
	}

	return nil
}

func LogDBAccess(tableName string, funcName string, RW string, data1 string, data2 string, data3 string, data4 string, data5 string) {
	elapsedTime := time.Since(initTime).Seconds()
	//log.Printf("DB Access +%.3fs tableName=%s funcName=%s", elapsedTime, tableName, funcName)
	record := []string{tableName, funcName, fmt.Sprintf("%.3f", elapsedTime), RW, data1, data2, data3, data4, data5}

	csvWriter.Write(record)

	if err := csvWriter.Error(); err != nil {
		log.Printf("CSV writer error: %v", err)
	}

	csvWriter.Flush()
}

func Logout() {
	totalElapsed := time.Since(initTime).Seconds()
	log.Printf("Logout game end and total time: +%.3fs", totalElapsed)
	if csvFile != nil {
		csvWriter.Flush()
		csvFile.Close()
	}
	log.Printf("csv 저장 완료.")
}

func PrintCount() {
	log.Println("Count start")
	log.Printf("------------------------------------------")
	log.Printf("Read Accounts: %d", countReadAccounts)
	log.Printf("Read AccountStats: %d", countReadAccountStats)
	log.Printf("Read Sessions: %d", countReadSessions)
	log.Printf("Read SessionSaveData: %d", countReadSessionSaveData)
	log.Printf("Read ActiveClientSessions: %d", countReadActiveClientSessions)
	log.Printf("Read SystemSaveData: %d", countReadSystemSaveData)
	log.Printf("Read DailyRuns: %d", countdailyRuns)
	log.Printf("Write Accounts: %d", countWriteAccounts)
	log.Printf("Write AccountStats: %d", countWriteAccountStats)
	log.Printf("Write Sessions: %d", countWriteSessions)
	log.Printf("Write SessionSaveData: %d", countWriteSessionSaveData)
	log.Printf("Write ActiveClientSessions: %d", countWriteActiveClientSessions)
	log.Printf("Write SystemSaveData: %d", countWriteSystemSaveData)
	log.Println(" ")
	log.Printf("userNameNum: %d", userNameNum)
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
	log.Printf("playerCountNum lastActivity: %d", lastActivityNum)
	log.Printf("battleCountNum: %d", battleCountNum)
	log.Printf("classicSessionPlayedCountNum: %d", classicSessionPlayedCountNum)
	log.Printf("lastActivityNum: %d", lastActivityNum)
	log.Printf("FetchRankingPageCountNum: %d", fetchRankingPageCountNum)
	log.Printf("FetchRankingsNum: %d", fetchRankingsNum)
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
			LogDBAccess("accounts", "AddAccountSession", "R", uuid, token, expire, "", "")
		}
		if "FetchUsernameByDiscordId" == funcName {
			userNameNum++
			log.Printf("username : %d", userNameNum)
			LogDBAccess("accounts", "FetchUsernameByDiscordId", "R", userName, "", "", "", "")
		}
		if "FetchUsernameByGoogleId" == funcName {
			userNameNum++
			log.Printf("username : %d", userNameNum)
			LogDBAccess("accounts", "FetchUsernameByGoogleId", "R", userName, "", "", "", "")
		}
		if "FetchDiscordIdByUsername" == funcName {
			discordIdNum++
			log.Printf("discordIdNum : %d", discordIdNum)
			LogDBAccess("accounts", "FetchDiscordIdByUsername", "R", discordId, "", "", "", "")
		}
		if "FetchGoogleIdByUsername" == funcName {
			googleIdNum++
			log.Printf("googleIdNum : %d", googleIdNum)
			LogDBAccess("accounts", "FetchGoogleIdByUsername", "R", googleId, "", "", "", "")
		}
		if "FetchDiscordIdByUUID" == funcName {
			discordIdNum++
			log.Printf("discordIdNum : %d", discordIdNum)
			LogDBAccess("accounts", "FetchDiscordIdByUUID", "R", discordId, "", "", "", "")
		}
		if "FetchGoogleIdByUUID" == funcName {
			googleIdNum++
			log.Printf("googleIdNum : %d", googleIdNum)
			LogDBAccess("accounts", "FetchGoogleIdByUUID", "R", googleId, "", "", "", "")
		}
		if "FetchUsernameBySessionToken" == funcName {
			userNameNum++
			log.Printf("username : %d", userNameNum)
			LogDBAccess("accounts", "FetchUsernameBySessionToken", "R", userName, "", "", "", "")
		}
		if "CheckUsernameExists" == funcName {
			dbUsernameNum++
			log.Printf("dbusernameNum : %d", dbUsernameNum)
			LogDBAccess("accounts", "CheckUsernameExists", "R", dbUsername, "", "", "", "")
		}
		if "FetchLastLoggedInDateByUsername" == funcName {
			lastLoggedInNum++
			log.Printf("lastLoggedIn : %d", lastLoggedInNum)
			LogDBAccess("accounts", "FetchLastLoggedInDateByUsername", "R", lastLoggedIn, "", "", "", "")
		}
		if "FetchAdminDetailsByUsername" == funcName {
			userNameNum++
			discordIdNum++
			googleIdNum++
			lastActivityNum++
			registeredNum++
			log.Printf("username : %d, discordIdNum : %d, googleIdNum : %d lastActivity : %d, registered : %d", userNameNum, discordIdNum, googleIdNum, lastActivityNum, registeredNum)
			LogDBAccess("accounts", "FetchAdminDetailsByUsername", "R", userName, discordId, googleId, lastActivity, registered)
		}
		if "FetchAccountKeySaltFromUsername" == funcName {
			keyNum++
			saltNum++
			log.Printf("keyNum : %d, saltNum : %d", keyNum, saltNum)
			LogDBAccess("accounts", "FetchAccountKeySaltFromUsername", "R", key, salt, "", "", "")
		}
		if "FetchTrainerIds" == funcName {
			trainerIDNum++
			secretIDNum++
			log.Printf("trainerId : %d, secretId : %d", trainerIDNum, secretIDNum)
			LogDBAccess("accounts", "FetchTrainerIds", "R", trainerID, secretID, "", "", "")
		}
		if "FetchUUIDFromToken" == funcName {
			uuidNum++
			log.Printf("uuidNum : %d", uuidNum)
			LogDBAccess("accounts", "FetchUUIDFromToken", "R", uuid, "", "", "", "")
		}
		if "FetchUsernameFromUUID" == funcName {
			userNameNum++
			log.Printf("username : %d", userNameNum)
			LogDBAccess("accounts", "FetchUsernameFromUUID", "R", userName, "", "", "", "")
		}
		if "FetchUUIDFromUsername" == funcName {
			uuidNum++
			log.Printf("uuidNum : %d", uuidNum)
			LogDBAccess("accounts", "FetchUUIDFromUsername", "R", uuid, "", "", "", "")
		}
		//account.go end.

		if "FetchPlayerCount" == funcName {
			lastActivityNum++
			log.Printf("FetchPlayerCount lastActivityNum : %d", lastActivityNum)
			LogDBAccess("accounts", "FetchPlayerCount", "R", lastActivity, "", "", "", "")
		}
		if "FetchBattleCount" == funcName {
			bannedNum++
			//battleCountNum++
			log.Printf("FetchBattleCount bannedNum: %d", bannedNum)
			LogDBAccess("accounts", "FetchBattleCount", "R", banned, "", "", "", "")
			//log.Printf("FetchBattleCount!")
		}
		if "FetchClassicSessionCount" == funcName {
			classicSessionPlayedCountNum++
			log.Printf("FetchClassicSessionCount : %d", classicSessionPlayedCountNum)
			LogDBAccess("accounts", "FetchClassicSessionCount", "R", classicSessionPlayedCount, "", "", "", "")
		}
		//game.go 항목 추가.

	case "accountStats":
		countReadAccountStats++
		if "FetchBattleCount" == funcName {
			battleCountNum++
			log.Printf("FetchBattleCount battleCountNum : %d", battleCountNum)
			LogDBAccess("accountStats", "FetchBattleCount", "R", battleCount, "", "", "", "")
		}
		if "FetchClassicSessionCount" == funcName {
			classicSessionPlayedCountNum++
			log.Printf("FetchClassicSessionPlayedCount : %d", classicSessionPlayedCountNum)
			LogDBAccess("accountStats", "FetchClassicSessionCount", "R", classicSessionPlayedCount, "", "", "", "")
		}
		//game.go 항목 추가.

		if "RetrievePlaytime" == funcName {
			playtimeNum++
			log.Printf("RetrievePlaytime playtimeNum : %d", playtimeNum)
			LogDBAccess("accountStats", "RetrievePlaytime", "R", playtime, "", "", "", "")
		}
		//savedata.go 항목 추가.

	case "sessions":
		countReadSessions++
		if "FetchUUIDFromToken" == funcName {
			uuidNum++
			log.Printf("FetchUUIDFromToken uuidNum : %d", uuidNum)
			LogDBAccess("sessions", "FetchUUIDFromToken", "R", uuid, "", "", "", "")
		}
	case "sessionSaveData":
		countReadSessionSaveData++
		if "ReadSessionSaveData" == funcName {
			dataNum++
			log.Printf("ReadSessionSaveData dataNum : %d", dataNum)
			LogDBAccess("sessionSaveData", "ReadSessionSaveData", "R", data, "", "", "", "")
			//log.Printf("ReadSessionSaveData!")
		}
		if "GetLatestSessionSaveDataSlot" == funcName {
			slotNum++
			log.Printf("GetLatestSessionSaveDataSlot slotNum : %d", slotNum)
			LogDBAccess("sessionSaveData", "GetLatestSessionSaveDataSlot", "R", slot, "", "", "", "")
			//log.Printf("GetLatestSessionSaveDataSlot!")
		}

	case "activeClientSessions":
		countReadActiveClientSessions++
		if "IsActiveSession" == funcName {
			idNum++
			log.Printf("IsActiveSession idNum : %d", idNum)
			LogDBAccess("activeClientSessions", "IsActiveSession", "R", id, "", "", "", "")
		}

	case "systemSaveData":
		countReadSystemSaveData++
		if "ReadSystemSaveData" == funcName {
			dataNum++
			log.Printf("ReadSystemSaveData dataNum : %d", dataNum)
			LogDBAccess("systemSaveData", "ReadSystemSaveData", "R", data, "", "", "", "")
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
			userNameNum++
			hashNum++ //register 과정에서 호출되는거라 해당 함수는 작동을 안 함.
			saltNum++
			registeredNum++
			log.Printf("INSERT uuid, username, hash, salt, registered : %d, %d, %d, %d, %d", uuidNum, userNameNum, hashNum, saltNum, registeredNum)
			LogDBAccess("accounts", "AddAccountRecord", "W", uuid, userName, hash, salt, registered)
		}
		if "AddAccountSession" == funcName {
			lastLoggedInNum++
			log.Printf("Update lastLoggedIn : %d", lastLoggedInNum)
			LogDBAccess("accounts", "AddAccountSession", "W", lastLoggedIn, "", "", "", "")
		}
		if "AddDiscordIdByUsername" == funcName {
			discordIdNum++
			log.Printf("Update discordId : %d", discordIdNum)
			LogDBAccess("accounts", "AddDiscordIdByUsername", "W", discordId, "", "", "", "")
		}
		if "AddGoogleIdByUsername" == funcName {
			googleIdNum++
			log.Printf("Update googleId : %d", googleIdNum)
			LogDBAccess("accounts", "AddGoogleIdByUsername", "W", googleId, "", "", "", "")
		}
		if "AddGoogleIdByUUID" == funcName {
			googleIdNum++
			log.Printf("Update googleId : %d", googleIdNum)
			LogDBAccess("accounts", "AddGoogleIdByUUID", "W", googleId, "", "", "", "")
		}
		if "AddDiscordIdByUUID" == funcName {
			discordIdNum++
			log.Printf("Update discordId : %d", discordIdNum)
			LogDBAccess("accounts", "AddDiscordIdByUUID", "W", discordId, "", "", "", "")
		}
		if "UpdateAccountPassword" == funcName {
			keyNum++
			saltNum++
			log.Printf("Update key, salt : %d, %d", keyNum, saltNum)
			LogDBAccess("accounts", "UpdateAccountPassword", "W", key, salt, "", "", "")
		}
		if "UpdateAccountLastActivity" == funcName {
			lastActivityNum++
			//lastLoggedInNum++
			log.Printf("Update lastActivity : %d", lastActivityNum)
			LogDBAccess("accounts", "UpdateAccountLastActivity", "W", lastActivity, "", "", "", "")
		}
		if "SetAccountBanned" == funcName {
			bannedNum++
			log.Printf("Update banned : %d", bannedNum)
			LogDBAccess("accounts", "SetAccountBanned", "W", banned, "", "", "", "")
		}
		if "UpdateTrainerIds" == funcName {
			trainerIDNum++
			secretIDNum++
			log.Printf("Update trainerId, secretId : %d, %d", trainerIDNum, secretIDNum)
			LogDBAccess("accounts", "UpdateTrainerIds", "W", trainerID, secretID, "", "", "")
		}
		if "RemoveDiscordIdByUUID" == funcName {
			discordIdNum++
			log.Printf("Update discordId : %d", discordIdNum)
			LogDBAccess("accounts", "RemoveDiscordIdByUUID", "W", discordId, "", "", "", "")
		}
		if "RemoveGoogleIdByUUID" == funcName {
			googleIdNum++
			log.Printf("Updat googleId : %d", googleIdNum)
			LogDBAccess("accounts", "RemoveGoogleIdByUUID", "W", googleId, "", "", "", "")
		}
		if "RemoveGoogleIdByUsername" == funcName {
			googleIdNum++
			log.Printf("Update googleId : %d", googleIdNum)
			LogDBAccess("accounts", "RemoveGoogleIdByUsername", "W", googleId, "", "", "", "")
		}
		if "RemoveDiscordIdByUsername" == funcName {
			discordIdNum++
			log.Printf("Update discordId : %d", discordIdNum)
			LogDBAccess("accounts", "RemoveDiscordIdByUsername", "W", discordId, "", "", "", "")
		}
		if "RemoveDiscordIdByDiscordId" == funcName {
			discordIdNum++
			log.Printf("Update discordId : %d", discordIdNum)
			LogDBAccess("accounts", "RemoveDiscordIdByDiscordId", "W", discordId, "", "", "", "")
		}
		if "RemoveGoogleIdByDiscordId" == funcName {
			googleIdNum++
			log.Printf("Update googleId : %d", googleIdNum)
			LogDBAccess("accounts", "RemoveGoogleIdByDiscordId", "W", googleId, "", "", "", "")
		}
		//account.go 항목 추가.

	case "accountStats":
		countWriteAccountStats++
		if "UpdateAccountStats" == funcName {
			statsNum++
			log.Printf("Insert stats playtime, battles, classicSessionPlayed 등..: %d", statsNum)
			LogDBAccess("accountStats", "UpdateAccountStats", "W", stats, "", "", "", "")
		}

	case "sessions":
		countWriteSessions++
		if "AddAccountSession" == funcName {
			uuidNum++
			tokenNum++
			expireNum++
			log.Printf("INSERT uuid, token, expire : %d, %d, %d", uuidNum, tokenNum, expireNum)
			LogDBAccess("sessions", "AddAccountSession", "W", uuid, token, expire, "", "")
		}
		if "RemoveSessionFromToken" == funcName {
			tokenNum++
			log.Printf("Remove token : %d", tokenNum)
			LogDBAccess("sessions", "RemoveSessionFromToken", "W", token, "", "", "", "")
		}

	case "sessionSaveData":
		countWriteSessionSaveData++
		if "StoreSessionSaveData" == funcName {
			uuidNum++
			slotNum++
			dataNum++
			timeStampNum++
			log.Printf("StoreSessionSaveData uuid, slot, data, timestamp : %d, %d, %d, %d", uuidNum, slotNum, dataNum, timeStampNum)
			LogDBAccess("sessionSaveData", "StoreSessionSaveData", "W", uuid, slot, data, timeStamp, "")
		}
		if "DeleteSessionSaveData" == funcName {
			sessionSaveDataNum++
			log.Printf("DeleteSessionSaveData sessionSaveDataNum : %d", sessionSaveDataNum)
			LogDBAccess("sessionSaveData", "DeleteSessionSaveData", "W", sessionSaveData, "", "", "", "")
		}

	case "activeClientSessions":
		countWriteActiveClientSessions++
		if "UpdateActiveSession" == funcName {
			uuidNum++
			clientSessionIdNum++
			log.Printf("Insert uuid, clientSessionId : %d, %d", uuidNum, clientSessionIdNum)
			LogDBAccess("activeClientSessions", "UpdateActiveSession", "W", uuid, clientSessionId, "", "", "")
		}

	case "systemSaveData":
		countWriteSystemSaveData++
		if "StoreSystemSaveData" == funcName {
			uuidNum++
			dataNum++
			timeStampNum++
			log.Printf("StoreSystemSaveData uuid, data, timestamp : %d, %d, %d", uuidNum, dataNum, timeStampNum)
			LogDBAccess("systemSaveData", "StoreSystemSaveData", "W", uuid, data, timeStamp, "", "")
		}
		if "DeleteSystemSaveData" == funcName {
			systemDataDeleteNum++
			log.Printf("DeleteSystemSaveData systemDataDeleteNum : %d", systemDataDeleteNum)
			LogDBAccess("systemSaveData", "DeleteSystemSaveData", "W", systemDataDelete, "", "", "", "")
		}

	default:
		unkownTableNum++
		log.Println("Unknown table name:", tableName)
	}
}
