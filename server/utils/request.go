package utils

import (
	//"bufio"
	"encoding/json"
	l4g "github.com/alecthomas/log4go"
	"strconv"
	"strings"
)

type Req struct {
	Lat      float64 `json: Lat`
	Lon      float64 `json:Lon`
	Zoom     int     `json:Zoom`
	ClientID string  `json:ClientID`
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
