package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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
	//Precipitation returns the percentage chance of precipitation, using the given location and the current local time.
	Precipitation(*Location) (int, error)

	//Temperature returns the current temperature in degrees Farenheit, using the given location and the current local time.
	Temperature(*Location) (int, error)

	//Condition returns the current weather condition (one of the values of Condition), using the given location and the current local time.
	Condition(*Location) (Condition, error)
}

//WundergroundCredentialsKey is the key within the credentials.json file for the Wunderground API provider's key.
const WundergroundCredentialsKey = "WEATHER_PROVIDER_WUNDERGROUND_KEY"

//Wunderground represents a Provider for the Weather Underground API.
//In order to use a Wunderground Provider, one must have a Wunderground API key on the Cumulus Plan or above.
type Wunderground struct {
	//key is the Wunderground API key.
	key string
}

//NewWunderground takes the credentials map and returns a new Wunderground API provider.
func NewWunderground(credentials map[string]string) (Wunderground, error) {
	//Check credentials map for the key
	if key, ok := credentials[WundergroundCredentialsKey]; ok {
		return Wunderground{key: key}, nil
	}

	//Key does not exist within credentials map
	return Wunderground{}, fmt.Errorf("unable to find Wunderground provider API key within credentials (must have key '%s')", WundergroundCredentialsKey)
}

//callAndUnmarshalAPI calls the Wunderground API endpoint and unmarshals the response into a wundergroundAPIResponse.
func (w Wunderground) callAndUnmarshalAPI(loc *Location) (*wundergroundAPIResponse, error) {
	resp, err := http.Get(fmt.Sprintf("http://api.wunderground.com/api/%s/hourly/q/%s/%s.json", w.key, loc.State, loc.City)) // Contact API endpoint
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
func (w Wunderground) Precipitation(loc *Location) (int, error) {
	//Contact API endpoint and get unmarshalled response
	winfo, err := w.callAndUnmarshalAPI(loc)
	if err != nil {
		return 0, err
	}

	//Convert the string-represented precipitation percentage into an integer
	prec, err := strconv.Atoi(winfo.HourlyForecast[0].Pop)
	return prec, err
}

//Temperature fulfills the Provider Temperature method.
func (w Wunderground) Temperature(loc *Location) (int, error) {
	//Contact API endpoint and get unmarshalled response
	winfo, err := w.callAndUnmarshalAPI(loc)
	if err != nil {
		return 0, err
	}

	//Convert the string-represented temperature into an integer
	prec, err := strconv.Atoi(winfo.HourlyForecast[0].Temp.English)

	fmt.Println("Temperature", prec)

	return prec, err
}

//Condition fulfills the Provider Condition method.
func (w Wunderground) Condition(loc *Location) (Condition, error) {
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
