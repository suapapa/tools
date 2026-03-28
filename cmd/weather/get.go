package main

import (
	"context"
	"fmt"
	"io"
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

func getGPS(ctx context.Context, lat, lon int) ([]byte, error) {
	url := fmt.Sprintf(urlFmtGPS, lat, lon)
	return getURL(ctx, url)
}

func getAddr(ctx context.Context, city, country, village string) ([]byte, error) {
	url := fmt.Sprintf(urlFmtAddr, city, country, village)
	return getURL(ctx, url)
}

func getID(ctx context.Context, id int) ([]byte, error) {
	url := fmt.Sprintf(urlFmtID, id)
	return getURL(ctx, url)
}

func getURL(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("appKey", appKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Even if it's not OK, we might want the body for the Fail struct
		body, _ := io.ReadAll(resp.Body)
		return body, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	return body, nil
}
