package weather

// Units represents units used in the current,hourly or daily response.
type Units struct {
	Time          		string `json:"time"`
	Temperature2m 		string `json:"temperature_2m"`
	Temperature2mMax 	string `json:"temperature_2m_max"`
	Temperature2mMin 	string `json:"temperature_2m_min"`
	CloudCover    		string `json:"cloud_cover"`
	WindSpeed80m  		string `json:"wind_speed_80m"`
	WindSpeed10mMax  	string `json:"wind_speed_10m_max"`
	UvIndex       		string `json:"uv_index"`
}

// HourlyData represents the data in the hourly forecast.
type HourlyData struct {
	Time          []string  `json:"time"`
	Temperature2m []float64 `json:"temperature_2m"`
	CloudCover    []int     `json:"cloud_cover"`
	WindSpeed80m  []float64 `json:"wind_speed_80m"`
	UvIndex       []float64 `json:"uv_index"`
}

// DailyData represents the data in the daily forecast.
type DailyData struct {
	Time          		[]string   	`json:"time"`
	Temperature2mMax 	[]float64 	`json:"temperature_2m_max"`
	Temperature2mMin 	[]float64 	`json:"temperature_2m_min"`
	WindSpeed10mMax  	[]float64 	`json:"wind_speed_10m_max"`
	WeatherCode     	[]int   	`json:"weather_code"`
	WindDirectionAngle  []int   	`json:"wind_direction_10m_dominant"`
	UvIndex       		[]float64 	`json:"uv_index_max"`

}

// CurrentData represents the current weather info.
type CurrentData struct {
	Temperature2m 		float64 `json:"temperature_2m"`
	CloudCover 			int 	`json:"cloud_cover"`
	WindSpeed80m 		float64 `json:"wind_speed_80m"`
    UvIndex 			float64 `json:"uv_index"`
	WeatherCode     	int   	`json:"weather_code"`
	WindDirectionAngle  int   `json:"wind_direction_10m"`
}

// WeatherResponse represents the Open Meteo response payload.
type WeatherResponse struct {
	Latitude           		float64     `json:"latitude"`
	Longitude          		float64     `json:"longitude"`
	GenerationTimeMs   		float64     `json:"generationtime_ms"`
	UtcOffsetSeconds   		int         `json:"utc_offset_seconds"`
	Timezone           		string      `json:"timezone"`
	TimezoneAbbreviation 	string    	`json:"timezone_abbreviation"`
	Elevation          		float64     `json:"elevation"`
	HourlyUnits        		Units 		`json:"hourly_units"`
	Hourly             		HourlyData  `json:"hourly"`
	DailyUnits         		Units  		`json:"daily_units"`
	Daily              		DailyData   `json:"daily"`
	Current 		   		CurrentData `json:"current"`
	CurrentUnits 			Units 		`json:"current_units"`
}

// WeatherForecastInfo represents the  weather forecast data returned from the weatherForecast query
type WeatherForecastInfo struct {
    LocationName    string  	`json:"location_name"`
    Latitude        string  	`json:"latitude"`
    Longitude       string  	`json:"longitude"`
	HourlyUnits     Units 		`json:"hourly_units"`
	Hourly          HourlyData  `json:"hourly"`
	DailyUnits		Units  		`json:"daily_units"`
	Daily           DailyData   `json:"daily"`
}

// CurrentWeatherInfo represents the current weather data returned from the weatherForLocations query
type CurrentWeatherInfo struct {
	ID					string	`json:"id"`
    LocationName    	string  `json:"location_name"`
    Latitude        	string  `json:"latitude"`
    Longitude       	string  `json:"longitude"`
    Temperature     	float64 `json:"temperature"`
    MaxTemperature  	float64 `json:"max_temperature"`
    MinTemperature  	float64 `json:"min_temperature"`
    CloudCoverage   	int 	`json:"cloud_coverage"`
    WindSpeed       	float64 `json:"wind_speed"`
    UvIndex         	float64 `json:"uv_index"`
	WeatherCode     	int   	`json:"weather_code"`
	WindDirectionAngle 	int 	`json:"wind_direction_10m"`
	Units 				Units   `json:"units"`
}
