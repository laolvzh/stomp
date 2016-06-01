package auth

import (
	//	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
)

type AuthDB struct {
	configAuthDB string
	db           map[string]string // Map with login in key and password in value for authentification field
}

type AuthParams struct {
	Login    string `json:Login`
	Passcode string `json:Passcode`
}

func NewAuth(fileWithLogins string) AuthDB {
	a := AuthDB{configAuthDB: fileWithLogins}
	a.initAuthDB()

	return a
}

func (a AuthDB) Authenticate(login, passcode string) bool {

	if len(a.db) == 0 {
		log.Println("Error: no database to authorization checking")
		os.Exit(1)
	}

	if l, ok := a.db[login]; ok {
		if a.db[l] == passcode {
			return true
		}
	}
	return false
}

// Get Login/Passcode dataBase from configure file
// Read JSON data and parsing it to AuthParams struct
func (a *AuthDB) initAuthDB() {

	buf := bytes.NewBuffer(nil)

	fp, err := os.Open(a.configAuthDB)
	if err != nil {
		log.Println("Could not read data from configureAuthFile: ", err)
	}
	defer fp.Close()

	_, err = io.Copy(buf, fp)
	if err != nil {
		log.Println("Could not process data from configureAuthFile: ", err)
	}

	authDataJSON := buf.Bytes()
	//log.Println("authDataJSON: ", string(authDataJSON))

	authData := []AuthParams{}

	err = json.Unmarshal(authDataJSON, &authData)
	if err != nil {
		log.Println("Couldn't get auth params from configureAuthFile: ", err)
	}

	dataMap := make(map[string]string)
	for _, userAuth := range authData {
		if len(dataMap) != 0 {
			if _, userExist := dataMap[userAuth.Login]; userExist {
				log.Println("Warning: user already exists in database; ignored")
				continue
			}
		}
		dataMap[userAuth.Login] = userAuth.Passcode
	}
	a.db = dataMap
}
