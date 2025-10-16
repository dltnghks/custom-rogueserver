/*
	Copyright (C) 2024  Pagefault Games

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package account

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"

	//"log"

	"github.com/pagefaultgames/rogueserver/db"
	"github.com/pagefaultgames/rogueserver/dbcount"
	//"github.com/pagefaultgames/rogueserver/dbcount/userInfo"
)

type LoginResponse GenericAuthResponse

// /account/login - log into account
func Login(username, password string) (LoginResponse, error) {
	var response LoginResponse

	log.Printf("username : %s", username)
	if !isValidUsername(username) {
		return response, fmt.Errorf("invalid username")
	}

	if len(password) < 6 {
		return response, fmt.Errorf("invalid password")
	}

	//log.Printf("login start !!!")

	key, salt, err := db.FetchAccountKeySaltFromUsername(username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return response, fmt.Errorf("account doesn't exist")
		}

		return response, err
	}

	if !bytes.Equal(key, deriveArgon2IDKey([]byte(password), salt)) {
		return response, fmt.Errorf("password doesn't match")
	}

	dbcount.SetUserName(username)
	dbcount.SetPassword(password)

	response.Token, err = GenerateTokenForUsername(username)

	//dbcount.SetToken(response.Token)

	dbcount.WriteCredentialsToCSV(username, password, response.Token)

	uuid, err := db.UUIDFromUsername(username)
	if err != nil {
		return response, fmt.Errorf("failed to fetch UUID: %s", err)
	}

	fmt.Println("uuid length:", len(uuid)) // 16이 나와야 정상

	encodeUuid := base64.RawURLEncoding.EncodeToString(uuid)
	//encodeUserName := base64.RawURLEncoding.EncodeToString([]byte(username))
	//fmt.Printf("uuid: %v\n", uuid)
	//fmt.Printf("encodedUuid: %s\n", encodeUuid)

	//userNameFromUuid[encodeUuid] = username
	dbcount.InitTimer(encodeUuid, username)
	dbcount.AddAPILog(encodeUuid, "login start!")
	dbcount.AddAPILog(encodeUuid, "login")

	if err != nil {
		return response, fmt.Errorf("failed to generate token: %s", err)
	}

	return response, nil
}

func GenerateTokenForUsername(username string) (string, error) {
	token := make([]byte, TokenSize)
	//log.Printf("token2 : %s", token)
	_, err := rand.Read(token)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %s", err)
	}

	//이미 login을 진행했던 계정이면 token값 수정 안 되게 설정하기.
	//어떻게 구현하지? username이랑 password를 받아서 동일한 username과 password가 있으면 token값 수정 안 되게 if문으로 처리.
	key := dbcount.GetUserFormat()
	key.Username = username
	key.Password = dbcount.GetPassword()

	log.Printf("token before compare : %s", token)
	token = dbcount.CompareUser(key, token)
	log.Printf("token after compare : %s", token)

	err = db.AddAccountSession(username, token)
	if err != nil {
		return "", fmt.Errorf("failed to add account session")
	}

	return base64.StdEncoding.EncodeToString(token), nil
}
