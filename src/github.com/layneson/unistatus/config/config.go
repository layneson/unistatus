package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//Config represents the configuration for unistatus.
type Config struct {
	//Brightness is the brightness of the diplay, from 0 to 1.
	Brightness float32

	//WeatherCache contains settings for the weather provider cache.
	WeatherCache weatherCache `json:"weather_cache"`

	//Location contains settings for the current location, for weather purposes.
	Location location
}

type weatherCache struct {
	//RefreshRate is the time, in seconds, between weather cache updates.
	RefreshRate int `json:"refresh_rate"`
}

//location holds location data in the form of two-letter state code and full city name.
type location struct {
	State string
	City  string
}

//Current holds the current configuration. It is initialized to the default.
var Current = Config{
	Brightness: 0.5,

	WeatherCache: weatherCache{RefreshRate: 3600},

	Location: location{State: "NY", City: "Binghamton"},
}

//ReadConfig reads the given configuration file.
//If the file does not exist, the default configuration is used.
//If the file does exist, the default configuration is overridden with the user-supplied field values.
func ReadConfig(file string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil
	}

	conts, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(conts, &Current)

	return err
}
