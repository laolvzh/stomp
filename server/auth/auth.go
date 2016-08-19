package auth

import (
	conf "github.com/KristinaEtc/config"
	"github.com/ventu-io/slf"
)

const pwdCurr string = "github.com/go-stomp/stomp/server/auth"

var log slf.StructuredLogger

func init() {
	log = slf.WithContext(pwdCurr)
}

type AuthDB struct {
	db map[string]string // Map with login in key and password in value for authentification field
}

type AuthParams struct {
	Login    string
	Passcode string
}

// ConfFile is a file with all program options
type ConfFile struct {
	AuthData []AuthParams
}

func NewAuth() *AuthDB {
	a := AuthDB{}
	a.initAuthDB()

	return &a
}

func (a *AuthDB) Authenticate(login, passcode string) bool {
	log.Debugf("login: %s, pwd: %s ", login, passcode)
	if pwd, ok := a.db[login]; ok {
		if pwd == passcode {
			return true
		}
	}
	return false
}

// Get Login/Passcode dataBase from configure file
// Read JSON data and parsing it to AuthParams struct
func (a *AuthDB) initAuthDB() {

	var authData = ConfFile{AuthData: []AuthParams{}}
	conf.ReadGlobalConfig(&(authData), "Auth Data")

	if len(authData.AuthData) == 0 {
		log.Warn("Empty login/password database.")
	}

	dataMap := make(map[string]string)
	for _, userAuth := range authData.AuthData {
		if len(dataMap) != 0 {
			if _, userExist := dataMap[userAuth.Login]; userExist {
				log.Warn("User already exists in database; ignored")
				continue
			}
		}
		if userAuth.Login == "" || userAuth.Passcode == "" {
			log.Warnf("Empty/wrong field; igrored user=%s/password=%s.", userAuth.Login, userAuth.Passcode)
		} else {
			dataMap[userAuth.Login] = userAuth.Passcode
		}
	}
	a.db = dataMap
}
