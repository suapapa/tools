package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if appKey == "" {
		return fmt.Errorf("SK_WEATHER_APIKEY is not set")
	}

	j, err := getID(ctx, 116)
	if err != nil {
		if len(j) > 0 {
			var f Fail
			if errJson := json.Unmarshal(j, &f); errJson == nil {
				return fmt.Errorf("API error: %s (%s)", f.Error.Message, f.Error.Code)
			}
		}
		return err
	}

	var w Success
	if err := json.Unmarshal(j, &w); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if len(w.Weather.Minutely) == 0 {
		return fmt.Errorf("no weather data found")
	}

	m := w.Weather.Minutely[0]
	fmt.Printf("Station: %s\n", m.Station.Name)
	fmt.Printf("Sky: %s\n", m.Sky.Name)
	fmt.Printf("Temperature: %s°C\n", m.Temperature.Tc)
	fmt.Printf("Humidity: %s%%\n", m.Humidity)

	return nil
}
