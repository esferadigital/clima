package openmeteo

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Parameters for the Open-Meteo Geocoding V1 API.
// These are not exclusive. Check the docs for additional ones.
// https://open-meteo.com/en/docs/geocoding-api
type GeocodingParams struct {
	Name  string
	Count int
}

// Result item in response array.
type GeocodingResult struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Response from the Open-Meteo Geocoding V1 API.
type GeocodingResponse struct {
	Results []GeocodingResult `json:"results"`
}

const GEOCODING_API_URL = "https://geocoding-api.open-meteo.com/v1/search"

// Gets a list of location matches based on the submitted name.
// Data is provided by the Open-Meteo API.
func SearchLocation(params GeocodingParams) (GeocodingResponse, error) {
	url := fmt.Sprintf("%s?name=%s&count=%d", GEOCODING_API_URL, params.Name, params.Count)

	resp, err := http.Get(url)
	if err != nil {
		return GeocodingResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GeocodingResponse{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	decoder := json.NewDecoder(resp.Body)
	var response GeocodingResponse
	err = decoder.Decode(&response)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return GeocodingResponse{}, err
	}

	return response, nil
}

