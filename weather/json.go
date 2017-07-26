package main

type Success struct {
	Weather struct {
		Minutely []struct {
			Station struct {
				Longitude string `json:"longitude"`
				Latitude  string `json:"latitude"`
				Name      string `json:"name"`
				ID        string `json:"id"`
				Type      string `json:"type"`
			} `json:"station"`
			Wind struct {
				Wdir string `json:"wdir"`
				Wspd string `json:"wspd"`
			} `json:"wind"`
			Precipitation struct {
				SinceOntime string `json:"sinceOntime"`
				Type        string `json:"type"`
			} `json:"precipitation"`
			Sky struct {
				Code string `json:"code"`
				Name string `json:"name"`
			} `json:"sky"`
			Rain struct {
				SinceOntime   string `json:"sinceOntime"`
				SinceMidnight string `json:"sinceMidnight"`
				Last10Min     string `json:"last10min"`
				Last15Min     string `json:"last15min"`
				Last30Min     string `json:"last30min"`
				Last1Hour     string `json:"last1hour"`
				Last6Hour     string `json:"last6hour"`
				Last12Hour    string `json:"last12hour"`
				Last24Hour    string `json:"last24hour"`
			} `json:"rain"`
			Temperature struct {
				Tc   string `json:"tc"`
				Tmax string `json:"tmax"`
				Tmin string `json:"tmin"`
			} `json:"temperature"`
			Humidity string `json:"humidity"`
			Pressure struct {
				Surface  string `json:"surface"`
				SeaLevel string `json:"seaLevel"`
			} `json:"pressure"`
			Lightning       string `json:"lightning"`
			TimeObservation string `json:"timeObservation"`
		} `json:"minutely"`
	} `json:"weather"`
	Common struct {
		AlertYn string `json:"alertYn"`
		StormYn string `json:"stormYn"`
	} `json:"common"`
	Result struct {
		Code       int    `json:"code"`
		RequestURL string `json:"requestUrl"`
		Message    string `json:"message"`
	} `json:"result"`
}

type Fail struct {
	Error struct {
		Category string `json:"category"`
		Code     string `json:"code"`
		ID       string `json:"id"`
		Link     string `json:"link"`
		Message  string `json:"message"`
	} `json:"error"`
}
