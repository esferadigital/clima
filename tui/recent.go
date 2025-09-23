package tui

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/esferadigital/clima/openmeteo"
)

const RECENT_LOCATIONS_FILE = "clima_recent.json"
const MAX_RECENT_LOCATIONS = 5

type RecentLocation struct {
	openmeteo.GeocodingResult
}

func getRecentLocationsPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".config", "clima")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(configDir, RECENT_LOCATIONS_FILE), nil
}

func loadRecentLocations() ([]openmeteo.GeocodingResult, error) {
	path, err := getRecentLocationsPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []openmeteo.GeocodingResult{}, nil
		}
		return nil, err
	}

	var locations []RecentLocation
	if err := json.Unmarshal(data, &locations); err != nil {
		return nil, err
	}

	results := make([]openmeteo.GeocodingResult, len(locations))
	for i, loc := range locations {
		results[i] = loc.GeocodingResult
	}

	return results, nil
}

func saveRecentLocations(locations []openmeteo.GeocodingResult) error {
	path, err := getRecentLocationsPath()
	if err != nil {
		return err
	}

	recentLocs := make([]RecentLocation, len(locations))
	for i, loc := range locations {
		recentLocs[i] = RecentLocation{loc}
	}

	data, err := json.MarshalIndent(recentLocs, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func addRecentLocation(location openmeteo.GeocodingResult) error {
	locations, err := loadRecentLocations()
	if err != nil {
		return err
	}

	// Remove if already exists
	for i, loc := range locations {
		if loc.ID == location.ID {
			locations = append(locations[:i], locations[i+1:]...)
			break
		}
	}

	// Add to front
	locations = append([]openmeteo.GeocodingResult{location}, locations...)

	// Keep only the most recent MAX_RECENT_LOCATIONS
	if len(locations) > MAX_RECENT_LOCATIONS {
		locations = locations[:MAX_RECENT_LOCATIONS]
	}

	return saveRecentLocations(locations)
}
