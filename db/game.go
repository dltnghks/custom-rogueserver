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

package db

import (
	"log"

	"github.com/pagefaultgames/rogueserver/dbcount"
)

func FetchPlayerCount() (int, error) {
	log.Printf("FetchPlayerCount")
	var playerCount int
	err := handle.QueryRow("SELECT COUNT(*) FROM accounts WHERE lastActivity > DATE_SUB(UTC_TIMESTAMP(), INTERVAL 5 MINUTE)").Scan(&playerCount)
	dbcount.IncrementRequestCount("accounts", false)
	dbcount.AddReadCount("accounts", "FetchPlayerCount")

	if err != nil {
		return 0, err
	}

	return playerCount, nil
}

func FetchBattleCount() (int, error) {
	log.Printf("FetchBattleCount")
	var battleCount int
	err := handle.QueryRow("SELECT COALESCE(SUM(s.battles), 0) FROM accountStats s JOIN accounts a ON a.uuid = s.uuid WHERE a.banned = 0").Scan(&battleCount)
	dbcount.IncrementRequestCount("accounts", false)
	dbcount.IncrementRequestCount("accountStats", false)
	dbcount.AddReadCount("accounts", "FetchBattleCount")
	dbcount.AddReadCount("accountStats", "FetchBattleCount")

	if err != nil {
		return 0, err
	}

	return battleCount, nil
}

func FetchClassicSessionCount() (int, error) {
	log.Printf("FetchClassicSessionCount")
	var classicSessionCount int
	err := handle.QueryRow("SELECT COALESCE(SUM(s.classicSessionsPlayed), 0) FROM accountStats s JOIN accounts a ON a.uuid = s.uuid WHERE a.banned = 0").Scan(&classicSessionCount)
	dbcount.IncrementRequestCount("accounts", false)
	dbcount.IncrementRequestCount("accountStats", false)
	dbcount.AddReadCount("accounts", "FetchClassicSessionCount")
	dbcount.AddReadCount("accountStats", "FetchClassicSessionCount")

	if err != nil {
		return 0, err
	}

	return classicSessionCount, nil
}
