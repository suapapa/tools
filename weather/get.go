package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var (
	urlFmtBase = "http://apis.skplanetx.com/weather/current/minutely?" +
		"version=1"
	urlFmtGps  = urlFmtBase + "&lat=%d&lon=%d"
	urlFmtAddr = urlFmtBase + "&city=%s&county=%s&village=%s"
	urlFmtID   = urlFmtBase + "&stnid=%d"

	appKey = os.Getenv("SK_WEATHER_APIKEY")
)

func getID(id int) (string, error) {
	url := fmt.Sprintf(urlFmtID, id)
	req, err := http.NewRequest("GET", url, nil)

	req.Header.Set("appKey", appKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(s), nil
}
