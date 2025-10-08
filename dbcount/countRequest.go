package dbcount

import (
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var usernameValue string = ""
var passwordValue string = ""
var tokenValue string = ""

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

// var initTime time.Time
var initTimes = make(map[string]time.Time)
var csvWriters = make(map[string]*csv.Writer)
var csvFiles = make(map[string]*os.File)
var userNameFromUuid = make(map[string]string) //key uuid, value userName
var countConnect = make(map[string]int)
var mut sync.Mutex

// csv file data save structure.
type userFormat struct {
	Username string
	Password string
}

var userData map[userFormat]string

//var isUseUuid = make(map[string]bool)

func GetUserFormat() userFormat {
	return userFormat{}
}

func SetUserName(uName string) {
	usernameValue = uName
	//return usernameValue
}

func SetPassword(uPass string) {
	passwordValue = uPass
	//return passwordValue
}

func SetToken(uToken string) {
	tokenValue = uToken
	//return tokenValue
}

func GetUserName() string {
	return usernameValue
}

func GetPassword() string {
	return passwordValue
}

func GetToken() string {
	return tokenValue
}

func CompareUser(userF userFormat, uToken []byte) []byte {
	//map 자료구조에 저장된 비밀번호와 동일하다면 이미 계정이 로그인을 진행한 상황.
	//즉 token값 변경하면 안 됨.
	if userData[userF] != base64.StdEncoding.EncodeToString(uToken) {
		log.Printf("새로운 계정 로그인 감지, token값 변경 : %s", userF.Username)
		uToken = []byte(userData[userF])
	}
	return uToken
}

func LoadCSVFile(filePath string) (map[userFormat]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		if len(record) >= 3 {
			key := userFormat{
				Username: record[0],
				Password: record[1],
			}

			tokenInfo := record[2]
			userData[key] = tokenInfo
		}
	}

	log.Printf("load csv : %v", userData)

	return userData, nil
}

func InitTimer(uuid string, userName string) error {
	mut.Lock()
	defer mut.Unlock()

	initTimes[uuid] = time.Now()
	userNameFromUuid[uuid] = userName
	countConnect[uuid]++
	log.Printf("init 시작 시간 : %s", initTimes[uuid].Format(time.RFC3339Nano))
	log.Printf("countConnect[%s] : %d", uuid, countConnect[uuid])
	log.Printf("userName[%s]", userName)

	if err := os.MkdirAll("/app/csv", 0o755); err != nil {
		return fmt.Errorf("디렉토리 생성 실패: %w", err)
	}

	// 이미 열려 있으면 재사용
	if _, ok := csvWriters[uuid]; ok {
		return nil
	}

	fileName := fmt.Sprintf("/app/csv/%s_db_access_count%d.csv", userNameFromUuid[uuid], countConnect[uuid])
	log.Printf("/app/csv/_%s_db_access_count", userName)
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("CSV 파일 열기 실패: %w", err)
	}

	if old := csvFiles[uuid]; old != nil && old != f {
		_ = old.Close()
	}
	csvFiles[uuid] = f
	csvWriters[uuid] = csv.NewWriter(f)
	return nil
}

func LogDBAccess(uuid string, tableName string, funcName string, RW string, data1 string, data2 string, data3 string, data4 string, data5 string) {
	// writer 없으면 초기화 (덮어쓰기 없이 append 모드로 열림)
	if _, ok := csvWriters[uuid]; !ok {
		// if err := InitTimer(uuid, userNameFromUuid[uuid]); err != nil {
		// 	log.Printf("InitTimer 실패(uuid=%s): %v", uuid, err)
		// 	return
		// }
		log.Printf("LogDBAccess 실패 writer 없음")
		return
	}

	mut.Lock()
	initTime := initTimes[uuid]
	w := csvWriters[uuid]
	mut.Unlock()
	if w == nil {
		log.Printf("writer 없음(uuid=%s)", uuid)
		return
	}

	elapsed := time.Since(initTime).Seconds()
	rec := []string{tableName, funcName, fmt.Sprintf("%.3f", elapsed), RW, data1, data2, data3, data4, data5}

	mut.Lock()
	_ = w.Write(rec)
	w.Flush()
	err := w.Error()
	mut.Unlock()

	if err != nil {
		log.Printf("CSV writer error(uuid=%s): %v", uuid, err)
	}
}

func Logout(uuid string) {
	mut.Lock()
	totalElapsed := time.Since(initTimes[uuid]).Seconds()
	mut.Unlock()
	log.Printf("Logout game end and total time: +%.3fs", totalElapsed)

	mut.Lock()
	if w := csvWriters[uuid]; w != nil {
		w.Flush()
	}
	if f := csvFiles[uuid]; f != nil {
		_ = f.Close()
	}
	delete(csvWriters, uuid)
	delete(csvFiles, uuid)
	delete(initTimes, uuid)
	mut.Unlock()

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

func AddAPILog(uuidReal string, logName string) {
	if "login" == logName {
		LogDBAccess(uuidReal, logName, "", "", "", "", "", "", "")
	} else if "logout" == logName {
		LogDBAccess(uuidReal, logName, "", "", "", "", "", "", "")
	} else if "Info" == logName {
		LogDBAccess(uuidReal, logName, "", "", "", "", "", "", "")
	} else if "get system" == logName {
		LogDBAccess(uuidReal, logName, "", "", "", "", "", "", "")
	} else if "update system" == logName {
		LogDBAccess(uuidReal, logName, "", "", "", "", "", "", "")
	} else if "verify system" == logName {
		LogDBAccess(uuidReal, logName, "", "", "", "", "", "", "")
	} else if "delete system" == logName {
		LogDBAccess(uuidReal, logName, "", "", "", "", "", "", "")
	} else if "get session" == logName {
		LogDBAccess(uuidReal, logName, "", "", "", "", "", "", "")
	} else if "update session" == logName {
		LogDBAccess(uuidReal, logName, "", "", "", "", "", "", "")
	} else if "delete session" == logName {
		LogDBAccess(uuidReal, logName, "", "", "", "", "", "", "")
	} else if "clear session" == logName {
		LogDBAccess(uuidReal, logName, "", "", "", "", "", "", "")
	} else if "newclear session" == logName {
		LogDBAccess(uuidReal, logName, "", "", "", "", "", "", "")
	} else if "updateAll" == logName {
		LogDBAccess(uuidReal, logName, "", "", "", "", "", "", "")
	} else if "login start!" == logName {
		LogDBAccess(uuidReal, logName, "", "", "", "", "", "", "")
	}
}

func AddReadCount(uuidReal string, tableName string, funcName string) {
	switch tableName {
	case "accounts":
		countReadAccounts++
		//log.Printf("read account start accounts fetchuuidfromusername!!!!!!!!")
		if "AddAccountSession" == funcName {
			uuidNum++
			tokenNum++
			expireNum++
			log.Printf("AddAccountSession Read uuid, expire : %d, %d, %d", uuidNum, tokenNum, expireNum)
			LogDBAccess(uuidReal, "accounts", "AddAccountSession", "R", uuid, token, expire, "", "")
		}
		if "FetchUsernameByDiscordId" == funcName {
			userNameNum++
			log.Printf("username : %d", userNameNum)
			LogDBAccess(uuidReal, "accounts", "FetchUsernameByDiscordId", "R", userName, "", "", "", "")
		}
		if "FetchUsernameByGoogleId" == funcName {
			userNameNum++
			log.Printf("username : %d", userNameNum)
			LogDBAccess(uuidReal, "accounts", "FetchUsernameByGoogleId", "R", userName, "", "", "", "")
		}
		if "FetchDiscordIdByUsername" == funcName {
			discordIdNum++
			log.Printf("discordIdNum : %d", discordIdNum)
			LogDBAccess(uuidReal, "accounts", "FetchDiscordIdByUsername", "R", discordId, "", "", "", "")
		}
		if "FetchGoogleIdByUsername" == funcName {
			googleIdNum++
			log.Printf("googleIdNum : %d", googleIdNum)
			LogDBAccess(uuidReal, "accounts", "FetchGoogleIdByUsername", "R", googleId, "", "", "", "")
		}
		if "FetchDiscordIdByUUID" == funcName {
			discordIdNum++
			log.Printf("discordIdNum : %d", discordIdNum)
			LogDBAccess(uuidReal, "accounts", "FetchDiscordIdByUUID", "R", discordId, "", "", "", "")
		}
		if "FetchGoogleIdByUUID" == funcName {
			googleIdNum++
			log.Printf("googleIdNum : %d", googleIdNum)
			LogDBAccess(uuidReal, "accounts", "FetchGoogleIdByUUID", "R", googleId, "", "", "", "")
		}
		if "FetchUsernameBySessionToken" == funcName {
			userNameNum++
			log.Printf("username : %d", userNameNum)
			LogDBAccess(uuidReal, "accounts", "FetchUsernameBySessionToken", "R", userName, "", "", "", "")
		}
		if "CheckUsernameExists" == funcName {
			dbUsernameNum++
			log.Printf("dbusernameNum : %d", dbUsernameNum)
			LogDBAccess(uuidReal, "accounts", "CheckUsernameExists", "R", dbUsername, "", "", "", "")
		}
		if "FetchLastLoggedInDateByUsername" == funcName {
			lastLoggedInNum++
			log.Printf("lastLoggedIn : %d", lastLoggedInNum)
			LogDBAccess(uuidReal, "accounts", "FetchLastLoggedInDateByUsername", "R", lastLoggedIn, "", "", "", "")
		}
		if "FetchAdminDetailsByUsername" == funcName {
			userNameNum++
			discordIdNum++
			googleIdNum++
			lastActivityNum++
			registeredNum++
			log.Printf("username : %d, discordIdNum : %d, googleIdNum : %d lastActivity : %d, registered : %d", userNameNum, discordIdNum, googleIdNum, lastActivityNum, registeredNum)
			LogDBAccess(uuidReal, "accounts", "FetchAdminDetailsByUsername", "R", userName, discordId, googleId, lastActivity, registered)
		}
		if "FetchAccountKeySaltFromUsername" == funcName {
			keyNum++
			saltNum++
			log.Printf("keyNum : %d, saltNum : %d", keyNum, saltNum)
			LogDBAccess(uuidReal, "accounts", "FetchAccountKeySaltFromUsername", "R", key, salt, "", "", "")
		}
		if "FetchTrainerIds" == funcName {
			trainerIDNum++
			secretIDNum++
			log.Printf("trainerId : %d, secretId : %d", trainerIDNum, secretIDNum)
			LogDBAccess(uuidReal, "accounts", "FetchTrainerIds", "R", trainerID, secretID, "", "", "")
		}
		if "FetchUUIDFromToken" == funcName {
			uuidNum++
			log.Printf("uuidNum : %d", uuidNum)
			LogDBAccess(uuidReal, "accounts", "FetchUUIDFromToken", "R", uuid, "", "", "", "")
		}
		if "FetchUsernameFromUUID" == funcName {
			userNameNum++
			log.Printf("username : %d", userNameNum)
			LogDBAccess(uuidReal, "accounts", "FetchUsernameFromUUID", "R", userName, "", "", "", "")
		}
		if "FetchUUIDFromUsername" == funcName {
			uuidNum++
			log.Printf("uuidNum : %d", uuidNum)
			LogDBAccess(uuidReal, "accounts", "FetchUUIDFromUsername", "R", uuid, "", "", "", "")
		}
		//account.go end.

		if "FetchPlayerCount" == funcName {
			lastActivityNum++
			log.Printf("FetchPlayerCount lastActivityNum : %d", lastActivityNum)
			LogDBAccess(uuidReal, "accounts", "FetchPlayerCount", "R", lastActivity, "", "", "", "")
		}
		if "FetchBattleCount" == funcName {
			bannedNum++
			//battleCountNum++
			log.Printf("FetchBattleCount bannedNum: %d", bannedNum)
			LogDBAccess(uuidReal, "accounts", "FetchBattleCount", "R", banned, "", "", "", "")
			//log.Printf("FetchBattleCount!")
		}
		if "FetchClassicSessionCount" == funcName {
			classicSessionPlayedCountNum++
			log.Printf("FetchClassicSessionCount : %d", classicSessionPlayedCountNum)
			LogDBAccess(uuidReal, "accounts", "FetchClassicSessionCount", "R", classicSessionPlayedCount, "", "", "", "")
		}
		//game.go 항목 추가.

	case "accountStats":
		countReadAccountStats++
		if "FetchBattleCount" == funcName {
			battleCountNum++
			log.Printf("FetchBattleCount battleCountNum : %d", battleCountNum)
			LogDBAccess(uuidReal, "accountStats", "FetchBattleCount", "R", battleCount, "", "", "", "")
		}
		if "FetchClassicSessionCount" == funcName {
			classicSessionPlayedCountNum++
			log.Printf("FetchClassicSessionPlayedCount : %d", classicSessionPlayedCountNum)
			LogDBAccess(uuidReal, "accountStats", "FetchClassicSessionCount", "R", classicSessionPlayedCount, "", "", "", "")
		}
		//game.go 항목 추가.

		if "RetrievePlaytime" == funcName {
			playtimeNum++
			log.Printf("RetrievePlaytime playtimeNum : %d", playtimeNum)
			LogDBAccess(uuidReal, "accountStats", "RetrievePlaytime", "R", playtime, "", "", "", "")
		}
		//savedata.go 항목 추가.

	case "sessions":
		countReadSessions++
		if "FetchUUIDFromToken" == funcName {
			uuidNum++
			log.Printf("FetchUUIDFromToken uuidNum : %d", uuidNum)
			LogDBAccess(uuidReal, "sessions", "FetchUUIDFromToken", "R", uuid, "", "", "", "")
		}
	case "sessionSaveData":
		countReadSessionSaveData++
		if "ReadSessionSaveData" == funcName {
			dataNum++
			log.Printf("ReadSessionSaveData dataNum : %d", dataNum)
			LogDBAccess(uuidReal, "sessionSaveData", "ReadSessionSaveData", "R", data, "", "", "", "")
			//log.Printf("ReadSessionSaveData!")
		}
		if "GetLatestSessionSaveDataSlot" == funcName {
			slotNum++
			log.Printf("GetLatestSessionSaveDataSlot slotNum : %d", slotNum)
			LogDBAccess(uuidReal, "sessionSaveData", "GetLatestSessionSaveDataSlot", "R", slot, "", "", "", "")
			//log.Printf("GetLatestSessionSaveDataSlot!")
		}

	case "activeClientSessions":
		countReadActiveClientSessions++
		if "IsActiveSession" == funcName {
			idNum++
			log.Printf("IsActiveSession idNum : %d", idNum)
			LogDBAccess(uuidReal, "activeClientSessions", "IsActiveSession", "R", id, "", "", "", "")
		}

	case "systemSaveData":
		countReadSystemSaveData++
		if "ReadSystemSaveData" == funcName {
			dataNum++
			log.Printf("ReadSystemSaveData dataNum : %d", dataNum)
			LogDBAccess(uuidReal, "systemSaveData", "ReadSystemSaveData", "R", data, "", "", "", "")
			//log.Printf("ReadSystemSaveData!")
		}

	default:
		unkownTableNum++
		log.Println("Unknown table name:", tableName)
	}
}

func AddWriteCount(uuidReal string, tableName string, funcName string) {
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
			LogDBAccess(uuidReal, "accounts", "AddAccountRecord", "W", uuid, userName, hash, salt, registered)
		}
		if "AddAccountSession" == funcName {
			lastLoggedInNum++
			log.Printf("Update lastLoggedIn : %d", lastLoggedInNum)
			LogDBAccess(uuidReal, "accounts", "AddAccountSession", "W", lastLoggedIn, "", "", "", "")
		}
		if "AddDiscordIdByUsername" == funcName {
			discordIdNum++
			log.Printf("Update discordId : %d", discordIdNum)
			LogDBAccess(uuidReal, "accounts", "AddDiscordIdByUsername", "W", discordId, "", "", "", "")
		}
		if "AddGoogleIdByUsername" == funcName {
			googleIdNum++
			log.Printf("Update googleId : %d", googleIdNum)
			LogDBAccess(uuidReal, "accounts", "AddGoogleIdByUsername", "W", googleId, "", "", "", "")
		}
		if "AddGoogleIdByUUID" == funcName {
			googleIdNum++
			log.Printf("Update googleId : %d", googleIdNum)
			LogDBAccess(uuidReal, "accounts", "AddGoogleIdByUUID", "W", googleId, "", "", "", "")
		}
		if "AddDiscordIdByUUID" == funcName {
			discordIdNum++
			log.Printf("Update discordId : %d", discordIdNum)
			LogDBAccess(uuidReal, "accounts", "AddDiscordIdByUUID", "W", discordId, "", "", "", "")
		}
		if "UpdateAccountPassword" == funcName {
			keyNum++
			saltNum++
			log.Printf("Update key, salt : %d, %d", keyNum, saltNum)
			LogDBAccess(uuidReal, "accounts", "UpdateAccountPassword", "W", key, salt, "", "", "")
		}
		if "UpdateAccountLastActivity" == funcName {
			lastActivityNum++
			//lastLoggedInNum++
			log.Printf("Update lastActivity : %d", lastActivityNum)
			LogDBAccess(uuidReal, "accounts", "UpdateAccountLastActivity", "W", lastActivity, "", "", "", "")
		}
		if "SetAccountBanned" == funcName {
			bannedNum++
			log.Printf("Update banned : %d", bannedNum)
			LogDBAccess(uuidReal, "accounts", "SetAccountBanned", "W", banned, "", "", "", "")
		}
		if "UpdateTrainerIds" == funcName {
			trainerIDNum++
			secretIDNum++
			log.Printf("Update trainerId, secretId : %d, %d", trainerIDNum, secretIDNum)
			LogDBAccess(uuidReal, "accounts", "UpdateTrainerIds", "W", trainerID, secretID, "", "", "")
		}
		if "RemoveDiscordIdByUUID" == funcName {
			discordIdNum++
			log.Printf("Update discordId : %d", discordIdNum)
			LogDBAccess(uuidReal, "accounts", "RemoveDiscordIdByUUID", "W", discordId, "", "", "", "")
		}
		if "RemoveGoogleIdByUUID" == funcName {
			googleIdNum++
			log.Printf("Updat googleId : %d", googleIdNum)
			LogDBAccess(uuidReal, "accounts", "RemoveGoogleIdByUUID", "W", googleId, "", "", "", "")
		}
		if "RemoveGoogleIdByUsername" == funcName {
			googleIdNum++
			log.Printf("Update googleId : %d", googleIdNum)
			LogDBAccess(uuidReal, "accounts", "RemoveGoogleIdByUsername", "W", googleId, "", "", "", "")
		}
		if "RemoveDiscordIdByUsername" == funcName {
			discordIdNum++
			log.Printf("Update discordId : %d", discordIdNum)
			LogDBAccess(uuidReal, "accounts", "RemoveDiscordIdByUsername", "W", discordId, "", "", "", "")
		}
		if "RemoveDiscordIdByDiscordId" == funcName {
			discordIdNum++
			log.Printf("Update discordId : %d", discordIdNum)
			LogDBAccess(uuidReal, "accounts", "RemoveDiscordIdByDiscordId", "W", discordId, "", "", "", "")
		}
		if "RemoveGoogleIdByDiscordId" == funcName {
			googleIdNum++
			log.Printf("Update googleId : %d", googleIdNum)
			LogDBAccess(uuidReal, "accounts", "RemoveGoogleIdByDiscordId", "W", googleId, "", "", "", "")
		}
		//account.go 항목 추가.

	case "accountStats":
		countWriteAccountStats++
		if "UpdateAccountStats" == funcName {
			statsNum++
			log.Printf("Insert stats playtime, battles, classicSessionPlayed 등..: %d", statsNum)
			LogDBAccess(uuidReal, "accountStats", "UpdateAccountStats", "W", stats, "", "", "", "")
		}

	case "sessions":
		countWriteSessions++
		if "AddAccountSession" == funcName {
			uuidNum++
			tokenNum++
			expireNum++
			log.Printf("INSERT uuid, token, expire : %d, %d, %d", uuidNum, tokenNum, expireNum)
			LogDBAccess(uuidReal, "sessions", "AddAccountSession", "W", uuid, token, expire, "", "")
		}
		if "RemoveSessionFromToken" == funcName {
			tokenNum++
			log.Printf("Remove token : %d", tokenNum)
			LogDBAccess(uuidReal, "sessions", "RemoveSessionFromToken", "W", token, "", "", "", "")
		}

	case "sessionSaveData":
		countWriteSessionSaveData++
		if "StoreSessionSaveData" == funcName {
			uuidNum++
			slotNum++
			dataNum++
			timeStampNum++
			log.Printf("StoreSessionSaveData uuid, slot, data, timestamp : %d, %d, %d, %d", uuidNum, slotNum, dataNum, timeStampNum)
			LogDBAccess(uuidReal, "sessionSaveData", "StoreSessionSaveData", "W", uuid, slot, data, timeStamp, "")
		}
		if "DeleteSessionSaveData" == funcName {
			sessionSaveDataNum++
			log.Printf("DeleteSessionSaveData sessionSaveDataNum : %d", sessionSaveDataNum)
			LogDBAccess(uuidReal, "sessionSaveData", "DeleteSessionSaveData", "W", sessionSaveData, "", "", "", "")
		}

	case "activeClientSessions":
		countWriteActiveClientSessions++
		if "UpdateActiveSession" == funcName {
			uuidNum++
			clientSessionIdNum++
			log.Printf("Insert uuid, clientSessionId : %d, %d", uuidNum, clientSessionIdNum)
			LogDBAccess(uuidReal, "activeClientSessions", "UpdateActiveSession", "W", uuid, clientSessionId, "", "", "")
		}

	case "systemSaveData":
		countWriteSystemSaveData++
		if "StoreSystemSaveData" == funcName {
			uuidNum++
			dataNum++
			timeStampNum++
			log.Printf("StoreSystemSaveData uuid, data, timestamp : %d, %d, %d", uuidNum, dataNum, timeStampNum)
			LogDBAccess(uuidReal, "systemSaveData", "StoreSystemSaveData", "W", uuid, data, timeStamp, "", "")
		}
		if "DeleteSystemSaveData" == funcName {
			systemDataDeleteNum++
			log.Printf("DeleteSystemSaveData systemDataDeleteNum : %d", systemDataDeleteNum)
			LogDBAccess(uuidReal, "systemSaveData", "DeleteSystemSaveData", "W", systemDataDelete, "", "", "", "")
		}

	default:
		unkownTableNum++
		log.Println("Unknown table name:", tableName)
	}
}
