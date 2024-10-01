package weather

import (
    "fmt"
    "io"
    "net/http"
    "strings"
    "encoding/json"
)
// WeatherController is Controller that handles operations on weather forecasts
type WeatherController struct{
}

// [ DAILY/WEEKLY ]  FetchWeatherForecast fetches the weather forecast data for a single location
func (wc *WeatherController) FetchWeatherForecast(db Database, location Location, metrics []string) (*WeatherForecastInfo, error) { 
    weatherInfo := &WeatherForecastInfo{}

    // Construct the daily query parameter from the metrics slice
    dailyMetrics := strings.Join(metrics, ",")

    // Construct the OpenMeteo API URL with the latitude and longitude
    query := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&daily=%s&forecast_days=7&timezone=auto&format=json", 
        location.Latitude, 
        location.Longitude, 
        dailyMetrics,
    )


    // Make the HTTP request
    resp, err := http.Get(query)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch weather data: %w", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %w", err)
    }
    
    // Parse the response
    var weatherData WeatherResponse
    if err := json.Unmarshal(body, &weatherData); err != nil {
        return nil, fmt.Errorf("failed to parse weather data: %w", err)
    }
    
    // Generate LocationID based on the latitude and longitude of the current location
    var locationID = GenerateID(location.Latitude, location.Longitude)

    // Check if the location exists in the database
    var existing Location
    err = db.d.Read("locations", locationID, &existing)
    if err != nil {
        // Return a different error message if the location is not found
        return nil, fmt.Errorf("location with latitude %s and longitude %s does not exist", location.Latitude, location.Longitude)
    }

    // Map query response from OpenMeteo response
    if len(weatherData.Daily.Time) > 0 {
        weatherInfo = &WeatherForecastInfo{
            LocationName:  existing.Name, // Use the name from the existing location in the database
            Latitude:      existing.Latitude,
            Longitude:     existing.Longitude,
            Daily: weatherData.Daily,
            DailyUnits: weatherData.DailyUnits,
            Hourly: weatherData.Hourly,
        }
    }

    return weatherInfo, nil
}


// [ MAP ] FetchWeatherForLocations fetches the weather data for multiple locations at once
func (wc *WeatherController) FetchWeatherForLocations(db Database,lc LocationController) ([]*CurrentWeatherInfo, error) {
	// Create arrays of latitudes and longitudes
	var latitudes []string
	var longitudes []string
    var weatherInfos []*CurrentWeatherInfo

    locations, err := lc.GetLocations(db)
    if err != nil {
        return nil, err
    }

	for _, location := range locations {
		latitudes = append(latitudes, location.Latitude)
		longitudes = append(longitudes, location.Longitude)
	}

	// Join the latitudes and longitudes into comma-separated strings
	latStr := strings.Join(latitudes, ",")
	lonStr := strings.Join(longitudes, ",")

	// Construct the OpenMeteo API URL with the latitudes and longitudes
	query := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&current=temperature_2m,cloud_cover,wind_speed_80m,wind_direction_10m,weather_code&timezone=auto&format=json", latStr, lonStr)

	// Make the HTTP request
	resp, err := http.Get(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %w", err)
    }
    
	// Parse the response
	var weatherData []WeatherResponse
	if err := json.Unmarshal(body, &weatherData); err != nil {
        return nil, fmt.Errorf("failed to parse weather data: %w", err)
	}

        // Map the parsed data to CurrentWeatherInfo
        for index, data := range weatherData {

            requestedLocation := locations[index]

            // Generate LocationID based on the latitude and longitude of the current location
            var locationID = GenerateID(requestedLocation.Latitude, requestedLocation.Longitude)
    

            var location Location
            err := db.d.Read("locations", locationID, &location)
            if err != nil {
                return nil, fmt.Errorf("location with latitude %s and longitude %s does not exist", fmt.Sprintf("%f", data.Latitude), fmt.Sprintf("%f", data.Longitude))
            }

            // Map query response from OpenMeteo response

            weatherInfos = append(weatherInfos, &CurrentWeatherInfo{
                ID:                 location.ID, // append Location ID
                LocationName:       location.Name, // Append Location Name
                Latitude:           location.Latitude,
                Longitude:          location.Longitude,
                Temperature:        data.Current.Temperature2m,
                CloudCoverage:      data.Current.CloudCover,
                WindSpeed:          data.Current.WindSpeed80m,
                UvIndex:            data.Current.UvIndex,
                WeatherCode:        data.Current.WeatherCode,
                WindDirectionAngle: data.Current.WindDirectionAngle,
                Units:              data.CurrentUnits,
            })
        }
    
    return weatherInfos, nil
}