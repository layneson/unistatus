package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/layneson/unistatus/config"
)

//Condition represents a weather condition.
type Condition uint16

//Condition constant declarations.
const (
	LightRain Condition = iota
	HeavyRain
	Snow
	Clouds
	Sun
)

//Location represents a city within the United States.
type Location struct {
	//State is the state to which the city belongs, in its two-letter capital abbreviation.
	//Examples: Pennsylvania = "PA", California = "CA".
	State string

	//City is the city of the location.
	//It should be the official name, capitalized as a proper noun, with underscores where spaces exist.
	//Examples: "Philadelpha", "San_Francisco".
	City string
}

//A Provider provides information about the current weather status.
type Provider interface {
	//Precipitation returns the percentage chance of precipitation, using the configured location and the current local time.
	Precipitation() (int, error)

	//Temperature returns the current temperature in degrees Farenheit, using the configured location and the current local time.
	Temperature() (int, error)

	//Condition returns the current weather condition (one of the values of Condition), using the configured location and the current local time.
	Condition() (Condition, error)
}

//A weatherCache is a cache for weather data.
type weatherCache struct {
	nextUpdate time.Time

	precipitation int
	temperature   int
	condition     Condition
}

//WundergroundCredentialsKey is the key within the credentials.json file for the Wunderground API provider's key.
const WundergroundCredentialsKey = "WEATHER_PROVIDER_WUNDERGROUND_KEY"

//Wunderground represents a Provider for the Weather Underground API.
//In order to use a Wunderground Provider, one must have a Wunderground API key on the Cumulus Plan or above.
type Wunderground struct {
	//key is the Wunderground API key.
	key string

	//location is the current location.
	location *Location

	//cache is the weather value cache.
	cache *weatherCache
}

//NewWunderground takes the credentials map and returns a new Wunderground API provider.
func NewWunderground(credentials map[string]string) (Wunderground, error) {
	//Check credentials map for the key
	if key, ok := credentials[WundergroundCredentialsKey]; ok {
		return Wunderground{
			key: key,
			location: &Location{
				State: config.Current.Location.State,
				City:  config.Current.Location.City,
			},
			cache: &weatherCache{nextUpdate: time.Now()},
		}, nil
	}

	//Key does not exist within credentials map
	return Wunderground{}, fmt.Errorf("unable to find Wunderground provider API key within credentials (must have key '%s')", WundergroundCredentialsKey)
}

//callAndUnmarshalAPI calls the Wunderground API endpoint and unmarshals the response into a wundergroundAPIResponse.
func (w Wunderground) callAndUnmarshalAPI() (*wundergroundAPIResponse, error) {
	resp, err := http.Get(fmt.Sprintf("http://api.wunderground.com/api/%s/hourly/q/%s/%s.json", w.key, w.location.State, w.location.City)) // Contact API endpoint
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	conts, err := ioutil.ReadAll(resp.Body) // Read full response body
	if err != nil {
		return nil, err
	}

	var winfo wundergroundAPIResponse // Unmarshal the JSON response into a struct
	err = json.Unmarshal(conts, &winfo)
	return &winfo, err
}

//Precipitation fulfills the Provider Precipitation method.
func (w Wunderground) Precipitation() (int, error) {
	if time.Now().Before(w.cache.nextUpdate) { // Check cache first
		return w.cache.precipitation, nil
	}

	w.cache.nextUpdate = time.Now().Add(time.Duration(config.Current.WeatherCache.RefreshRate) * time.Second)

	//Contact API endpoint and get unmarshalled response
	winfo, err := w.callAndUnmarshalAPI()
	if err != nil {
		return 0, err
	}

	//Convert the string-represented precipitation percentage into an integer
	prec, err := strconv.Atoi(winfo.HourlyForecast[0].Pop)
	return prec, err
}

//Temperature fulfills the Provider Temperature method.
func (w Wunderground) Temperature() (int, error) {
	if time.Now().Before(w.cache.nextUpdate) { // Check cache first
		return w.cache.temperature, nil
	}

	w.cache.nextUpdate = time.Now().Add(time.Duration(config.Current.WeatherCache.RefreshRate) * time.Second)

	//Contact API endpoint and get unmarshalled response
	winfo, err := w.callAndUnmarshalAPI()
	if err != nil {
		return 0, err
	}

	//Convert the string-represented temperature into an integer
	prec, err := strconv.Atoi(winfo.HourlyForecast[0].Temp.English)

	fmt.Println("Temperature", prec)

	return prec, err
}

//Condition fulfills the Provider Condition method.
func (w Wunderground) Condition() (Condition, error) {
	if time.Now().Before(w.cache.nextUpdate) { // Check cache first
		return w.cache.condition, nil
	}

	w.cache.nextUpdate = time.Now().Add(time.Duration(config.Current.WeatherCache.RefreshRate) * time.Second)

	//TODO: MAKE THIS DO SOMETHING
	return Clouds, nil
}

//wundergroundAPIResponse represents an API response from the Wunderground hourly API endpoint.
type wundergroundAPIResponse struct {
	HourlyForecast []struct {
		Temp struct {
			English string `json:"english"`
		} `json:"temp"`
		Pop string `json:"pop"`
	} `json:"hourly_forecast"`
}
