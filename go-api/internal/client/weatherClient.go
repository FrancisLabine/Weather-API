package client

import "weather-api/internal/util"
import "weather-api/internal/models"
import (
    "encoding/json"
    "io"
    "net/http"
	"time"
	"fmt"
    // "log"
)

//DAILY
func dailyProcessing(url string) ([]models.DailyWeather, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    var parsed map[string]interface{}
    json.Unmarshal(body, &parsed)

    daily := parsed["daily"].(map[string]interface{})
    dates := daily["time"].([]interface{})
    maxTemps := daily["temperature_2m_max"].([]interface{})
    minTemps := daily["temperature_2m_min"].([]interface{})
    precip := daily["precipitation_sum"].([]interface{})
    wind := daily["wind_speed_10m_max"].([]interface{})
    uv := daily["uv_index_max"].([]interface{})
    humidity := daily["relative_humidity_2m_mean"].([]interface{})
    weatherCode := daily["weather_code"].([]interface{})

    result := []models.DailyWeather{}
    for i := 0; i < len(dates); i++ {
		result = append(result, models.DailyWeather{
			Date:             dates[i].(string),
			MaxTemp:          util.ToFloat(maxTemps[i]),
			MinTemp:          util.ToFloat(minTemps[i]),
			Precip:           util.ToFloat(precip[i]),
			Wind:             util.ToFloat(wind[i]),
			UV:               util.ToFloat(uv[i]),
            Humidity:         util.ToFloat(humidity[i]),
            WeatherCode:      util.ToFloat(weatherCode[i]),
		})
    }

    return result, nil
}

func FetchWeatherHistory(lat, lon string, days int) ([]models.DailyWeather, error){
    end := time.Now().AddDate(0, 0, -1)
    start := end.AddDate(0, 0, -days)

	url := fmt.Sprintf(
        "https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&daily=weather_code,temperature_2m_max,temperature_2m_min,precipitation_sum,uv_index_max,wind_speed_10m_max,relative_humidity_2m_mean&timezone=America%%2FToronto&start_date=%s&end_date=%s",
        lat, lon,
        start.Format("2006-01-02"),
        end.Format("2006-01-02"),
    )

    result, err := dailyProcessing(url)
    return result, err
}

func FetchDailyForecast(lat, lon string, days int) ([]models.DailyWeather, error) {
    now := time.Now()
    end := now.AddDate(0, 0, days-1)

    url := fmt.Sprintf(
        "https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&daily=weather_code,temperature_2m_max,temperature_2m_min,precipitation_sum,uv_index_max,wind_speed_10m_max,relative_humidity_2m_mean&timezone=America/Toronto&start_date=%s&end_date=%s",
        lat, lon,
        now.Format("2006-01-02"),
        end.Format("2006-01-02"),
    )

    result, err := dailyProcessing(url)
    return result, err
}

//HOURLY
//TODO : REFACTOR THIS LIKE DAILY CALLS
func hourlyProcessing(url string) ([]models.HourlyWeather, error) {
    return nil, fmt.Errorf("hourlyProcessing not implemented yet")
}

func FetchHourlyForecast(lat, lon string, hours int) ([]map[string]interface{}, error) {
    now := time.Now()
    end := now.Add(time.Duration(hours) * time.Hour)

    url := fmt.Sprintf(
        "https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&hourly=weather_code,temperature_2m,relative_humidity_2m,precipitation_probability,uv_index,wind_speed_10m&timezone=America%%2FToronto&start_date=%s&end_date=%s",
        lat, lon,
        now.Format("2006-01-02"),
        end.Format("2006-01-02"),
    )

    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)
    var parsed map[string]interface{}
    json.Unmarshal(body, &parsed)

    hourly := parsed["hourly"].(map[string]interface{})
    times := hourly["time"].([]interface{})
    temps := hourly["temperature_2m"].([]interface{})
    humidity := hourly["relative_humidity_2m"].([]interface{})
    rain := hourly["precipitation_probability"].([]interface{})
    wind := hourly["wind_speed_10m"].([]interface{})
    uv := hourly["uv_index"].([]interface{})
    weatherCode := hourly["weather_code"].([]interface{})

    result := []map[string]interface{}{}
    for i := 0; i < len(times) && i < hours; i++ {
        result = append(result, map[string]interface{}{
			"timestamp":     times[i],                         // ISO 8601 string
			"temp_c":        util.ToFloat(temps[i]),                // temperature in Celsius
			"humidity_pct":  util.ToFloat(humidity[i]),             // percentage
			"precip_pct":    util.ToFloat(rain[i]),                 // chance of rain (percentage)
			"wind_kph":      util.ToFloat(wind[i]),                 // wind speed in km/h
			"uv_index":      util.ToFloat(uv[i]),
            "weather_code": util.ToFloat(weatherCode[i]),
		})
    }

    return result, nil
}

func FetchCurrentWeather(lat, lon string) (map[string]interface{}, error) {

    loc, err := time.LoadLocation("America/Toronto")
    if err != nil {
        panic(err)
    }

    // Get the current time in the specified timezone
    nowUTC := time.Now().In(loc)

    // Truncate to the nearest hour
    now := time.Date(
        nowUTC.Year(), nowUTC.Month(), nowUTC.Day(),
        nowUTC.Hour(), 0, 0, 0, loc,
    )

	dateStr := now.Format("2006-01-02")
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&current_weather=true&hourly=relative_humidity_2m,uv_index&timezone=America/Toronto&start_date=%s&end_date=%s",
		lat, lon, dateStr, dateStr,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
    var parsed map[string]interface{}
    json.Unmarshal(body, &parsed)

	// Extract current weather block
	cw := parsed["current_weather"].(map[string]interface{})

	// Extract matching hourly humidity and UV index
	hourly := parsed["hourly"].(map[string]interface{})
	times := hourly["time"].([]interface{})
	humidityList := hourly["relative_humidity_2m"].([]interface{})
	uvList := hourly["uv_index"].([]interface{})

	// Match current timestamp with hourly
	var humidity, uv float64

    cwTimeStr := cw["time"].(string)
    cwTime, err := time.Parse("2006-01-02T15:04", cwTimeStr)
    if err != nil {
    } else {
        rounded := cwTime.Truncate(time.Hour).Format("2006-01-02T15:00")  
        for i, t := range times {
            if t == rounded {
                humidity = util.ToFloat(humidityList[i])
                uv = util.ToFloat(uvList[i])
                break
            }
        }
    }
	return map[string]interface{}{
		"timestamp":    cw["time"],
		"temp_c":       util.ToFloat(cw["temperature"]),
		"wind_kph":     util.ToFloat(cw["windspeed"]),
		"wind_dir_deg": util.ToFloat(cw["winddirection"]),
		"weather_code": util.ToFloat(cw["weathercode"]),
		"humidity_pct": humidity,
		"uv_index":     uv,
	}, nil
}


