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
	urlFmtGPS  = urlFmtBase + "&lat=%d&lon=%d"
	urlFmtAddr = urlFmtBase + "&city=%s&county=%s&village=%s"
	urlFmtID   = urlFmtBase + "&stnid=%d"

	appKey = os.Getenv("SK_WEATHER_APIKEY")
)

func getGPS(lat, lon int) ([]byte, error) {
	url := fmt.Sprintf(urlFmtGPS, lat, lon)
	return getURL(url)
}

func getAddr(city, country, village string) ([]byte, error) {
	url := fmt.Sprintf(urlFmtAddr, city, country, village)
	return getURL(url)
}

func getID(id int) ([]byte, error) {
	url := fmt.Sprintf(urlFmtID, id)
	return getURL(url)
}

func getURL(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)

	req.Header.Set("appKey", appKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return s, nil
}
