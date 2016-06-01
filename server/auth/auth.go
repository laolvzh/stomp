package auth

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

type AuthDB struct {
	db map[string]string // Map with login in key and password in value for authentification field
}

var ConfigAuthDB = "../server/auth/authDB.json"

type AuthParams struct {
	Login    string `json:Login`
	Passcode string `json:Passcode`
}

func (a AuthDB) Authenticate(login, passcode string) bool {
	if a.db == nil {
		a.db = initAuthDB()
	}

	for k, v := range a.db {
		log.Println(k, " ", v)
	}

	if pwd, ok := a.db[login]; ok {
		if passcode == pwd {
			return true
		} //do something here
	}

	return false
}

// Get Login/Passcode dataBase from configure file
func initAuthDB() map[string]string {
	file, err := os.Open(ConfigAuthDB)
	if err != nil {
		log.Println("Error open ConfigAuthDB: ", err)
		os.Exit(1)
	}
	defer file.Close()

	db := make(map[string]string)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		authDataJSON := scanner.Bytes()

		log.Println("authDataJSON: ", string(authDataJSON))
		authData := AuthParams{}
		err := json.Unmarshal(authDataJSON, &authData)
		if err != nil {
			log.Println("Error get auth params from configureAuthFile: ", err)
			continue
		}

		db[authData.Login] = authData.Passcode
	}
	return db
}
