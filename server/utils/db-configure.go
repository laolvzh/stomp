package nominatim-utils

package main

import (
	"bufio"
	//"fmt"
	//"io"
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"os"
	//"time"
)

type ConfigDB struct {
	DBname, Host, User, Password string
}

type Req struct {
	Lat      float64 `json: Lat`
	Lon      float64 `json:Lon`
	Zoom     int     `json:Zoom`
	ClientID string  `json:ClientID`
}

type Params struct {
	locParams      LocationParams
	format         string
	addressDetails bool
	sqlOpenStr     string
	config         ConfigDB
	db             *sql.DB
	speedTest      bool
}

func ConfigurateDB(configFile string) string {

	file, err := os.Open(configFile)
	if err != nil {
		log.Printf("No configurate file")
	} else {
		defer file.Close()
		decoder := json.NewDecoder(file)
		config := ConfigDB{}
		err := decoder.Decode(&config)
		if err != nil {
			log.Println("error: ", err)
		}

		sqlOpenStr := "dbname=" + config.DBname +
			" host=" + config.Host +
			" user=" + config.User +
			" password=" + config.Password
		return sqlOpenStr
	}
}

func GetLocationFromNominatim(coordinates, sqlOpenStr string) (*Nominatim.DataWithoutDetails, error) {


	reverseGeocode, err := Nominatim.NewReverseGeocode(sqlOpenStr)
	if err != nil {
		log.Critical(err)
		os.Exit(1)
	}
	defer reverseGeocode.Close()

	//oReverseGeocode.SetLanguagePreference()
	reverseGeocode.SetIncludeAddressDetails(p.addressDetails)
	reverseGeocode.SetZoom(p.clientReq.Zoom)
	reverseGeocode.SetLocation(p.clientReq.Lat, p.clientReq.Lon)
	place, err := reverseGeocode.Lookup()
	if err != nil {
		return nil, err
	}

	return place, nil
}


func (r *Req) getLocationJSON() (string, error) {

	dataJSON, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(dataJSON), nil
}

type error interface {
	Error() string
}

func MakeReq(parameters, clientID string, log l4g.Logger) (reqInJSON *string, err error) {

	locSlice := strings.Split(parameters, ",")
	r := Req{}
	r.Lat, err = strconv.ParseFloat(locSlice[0], 32)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	r.Lon, err = strconv.ParseFloat(locSlice[1], 32)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	r.Zoom, err = strconv.Atoi(locSlice[2])
	if err != nil {
		log.Error(err)
		return nil, err
	}
	r.ClientID = clientID

	jsonReq, err := r.getLocationJSON()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &jsonReq, nil
}

func GetLocationJSON(data Nominatim.DataWithoutDetails, log l4g.Logger) (string, error) {

	dataJSON, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
		return "", err
	}
	return string(dataJSON), nil
}

