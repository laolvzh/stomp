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
	db           []AuthParams // Map with login in key and password in value for authentification field
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
	for _, param := range a.db {
		//log.Println("param: ", param)
		if param.Login == login {
			if param.Passcode == passcode {
				return true
			}
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

	a.db = authData
}
