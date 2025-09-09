package main

import (
	"fmt"

	"github.com/esferadigital/clima/openmeteo"
)

const DEFAULT_LATITUDE = -2.1885029262472884
const DEFAULT_LONGITUDE = -81.00949365444902

func main() {
	params := openmeteo.ForecastParams{
		Latitude:  DEFAULT_LATITUDE,
		Longitude: DEFAULT_LONGITUDE,
		Current: []openmeteo.CurrentWeatherVariables{
			openmeteo.Temperature2m,
			openmeteo.RelativeHumidity2m,
			openmeteo.ApparentTemperature,
			openmeteo.IsDay,
			openmeteo.WeatherCode,
			openmeteo.WindSpeed10m,
			openmeteo.WindDirection10m,
			openmeteo.WindGusts10m,
			openmeteo.Rain,
			openmeteo.Showers,
			openmeteo.Snowfall,
			openmeteo.CloudCover,
			openmeteo.SeaLevelPressure,
			openmeteo.SurfacePressure,
			openmeteo.PressureAtGround,
		},
	}

	forecast, err := openmeteo.GetForecast(params)
	if err != nil {
		panic(err)
	}

	temperature := fmt.Sprintf("Temperature: %.1f %s", forecast.Current[string(openmeteo.Temperature2m)], forecast.CurrentUnits[string(openmeteo.Temperature2m)])
	fmt.Println(temperature)

	windSpeed := fmt.Sprintf("Wind Speed: %.1f %s", forecast.Current[string(openmeteo.WindSpeed10m)], forecast.CurrentUnits[string(openmeteo.WindSpeed10m)])
	fmt.Println(windSpeed)

	windDirection := fmt.Sprintf("Wind Direction: %.1f %s", forecast.Current[string(openmeteo.WindDirection10m)], forecast.CurrentUnits[string(openmeteo.WindDirection10m)])
	fmt.Println(windDirection)

	windGusts := fmt.Sprintf("Wind Gusts: %.1f %s", forecast.Current[string(openmeteo.WindGusts10m)], forecast.CurrentUnits[string(openmeteo.WindGusts10m)])
	fmt.Println(windGusts)

	rain := fmt.Sprintf("Rain: %.1f %s", forecast.Current[string(openmeteo.Rain)], forecast.CurrentUnits[string(openmeteo.Rain)])
	fmt.Println(rain)

	showers := fmt.Sprintf("Showers: %.1f %s", forecast.Current[string(openmeteo.Showers)], forecast.CurrentUnits[string(openmeteo.Showers)])
	fmt.Println(showers)

	snowfall := fmt.Sprintf("Snowfall: %.1f %s", forecast.Current[string(openmeteo.Snowfall)], forecast.CurrentUnits[string(openmeteo.Snowfall)])
	fmt.Println(snowfall)

	cloudCover := fmt.Sprintf("Cloud Cover: %.1f %s", forecast.Current[string(openmeteo.CloudCover)], forecast.CurrentUnits[string(openmeteo.CloudCover)])
	fmt.Println(cloudCover)

	seaLevelPressure := fmt.Sprintf("Sea Level Pressure: %.1f %s", forecast.Current[string(openmeteo.SeaLevelPressure)], forecast.CurrentUnits[string(openmeteo.SeaLevelPressure)])
	fmt.Println(seaLevelPressure)

	surfacePressure := fmt.Sprintf("Surface Pressure: %.1f %s", forecast.Current[string(openmeteo.SurfacePressure)], forecast.CurrentUnits[string(openmeteo.SurfacePressure)])
	fmt.Println(surfacePressure)

	pressureAtGround := fmt.Sprintf("Pressure at Ground: %.1f %s", forecast.Current[string(openmeteo.PressureAtGround)], forecast.CurrentUnits[string(openmeteo.PressureAtGround)])
	fmt.Println(pressureAtGround)
}
