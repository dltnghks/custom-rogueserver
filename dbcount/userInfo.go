package dbcount

import (
	"encoding/csv"
	"log"
	"os"
	"sync"
)

// 1. 전역 뮤텍스(Mutex)를 생성하여 파일 접근을 동기화합니다.
// 여러 사용자가 동시에 이 함수를 호출해도 파일이 깨지지 않도록 보장합니다.
var csvMutex = &sync.Mutex{}

// 이 함수를 로그인 성공 및 토큰 생성 직후에 호출합니다.
func WriteCredentialsToCSV(username, password, token string) {
	// 2. 뮤텍스를 잠가 한 번에 하나의 작업만 파일에 쓸 수 있도록 합니다.
	csvMutex.Lock()
	defer csvMutex.Unlock() // 함수가 끝나면 (에러 발생 여부와 상관없이) 자동으로 잠금을 해제합니다.

	filePath := "/app/csv/credentials5.csv"

	// 3. 파일을 추가 쓰기 모드(Append)로 엽니다. 파일이 없으면 새로 생성합니다.
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("CSV 파일 열기 실패: %s", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	// 파일 크기를 확인하여 파일이 비어있는지 확인합니다.
	fileInfo, _ := file.Stat()
	isNewFile := fileInfo.Size() == 0

	// 4. 파일이 새로 생성되었다면, 맨 위에 헤더(header)를 추가합니다.
	if isNewFile {
		header := []string{"username", "password", "token"}
		if err := writer.Write(header); err != nil {
			log.Printf("CSV 헤더 쓰기 실패: %s", err)
			return
		}
	}

	// 5. 실제 사용자 데이터를 CSV 파일에 한 줄 씁니다.
	record := []string{username, password, token}
	if err := writer.Write(record); err != nil {
		log.Printf("CSV 데이터 쓰기 실패: %s", err)
		return
	}

	// 6. 버퍼에 남아있는 데이터를 파일에 완전히 쓰도록 합니다. (매우 중요)
	writer.Flush()
}

// // --- API 서버의 로그인 핸들러에 적용하는 방법 (예시) ---
// func LoginHandler(w http.ResponseWriter, r *http.Request) {
//     // ... (사용자 이름과 비밀번호를 파싱하는 로직) ...
//     username := r.FormValue("username")
//     password := r.FormValue("password")

//     // ... (데이터베이스에서 사용자 인증하는 로직) ...

//     // 인증 성공 후 토큰 생성
//     tokenString, err := generateNewToken(username)
//     if err != nil {
//         // ... 에러 처리 ...
//         return
//     }

//     // *** CSV 파일에 자격 증명 기록 ***
//     // 'go' 키워드를 사용해 비동기적으로 실행하면 응답 지연을 최소화할 수 있습니다.
//     go writeCredentialsToCSV(username, password, tokenString)

//     // ... (클라이언트에게 토큰을 포함한 응답을 보내는 로직) ...
// }
