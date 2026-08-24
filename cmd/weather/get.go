package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const (
	weatherAPIBase = "https://apis.skplanetx.com/weather/current/minutely"
	maxWeatherBody = 1 << 20 // 1 MiB
)

var appKey = os.Getenv("SK_WEATHER_APIKEY")

func getGPS(ctx context.Context, lat, lon int) ([]byte, error) {
	q := url.Values{}
	q.Set("version", "1")
	q.Set("lat", strconv.Itoa(lat))
	q.Set("lon", strconv.Itoa(lon))
	return getURL(ctx, weatherAPIBase+"?"+q.Encode())
}

func getAddr(ctx context.Context, city, country, village string) ([]byte, error) {
	q := url.Values{}
	q.Set("version", "1")
	q.Set("city", city)
	q.Set("county", country)
	q.Set("village", village)
	return getURL(ctx, weatherAPIBase+"?"+q.Encode())
}

func getID(ctx context.Context, id int) ([]byte, error) {
	q := url.Values{}
	q.Set("version", "1")
	q.Set("stnid", strconv.Itoa(id))
	return getURL(ctx, weatherAPIBase+"?"+q.Encode())
}

func getURL(ctx context.Context, rawURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("appKey", appKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxWeatherBody+1))
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	if len(body) > maxWeatherBody {
		return nil, fmt.Errorf("response body exceeds %d bytes", maxWeatherBody)
	}

	if resp.StatusCode != http.StatusOK {
		return body, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return body, nil
}
