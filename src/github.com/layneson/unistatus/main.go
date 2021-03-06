package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/layneson/unistatus/config"
	"github.com/layneson/unistatus/display"
	"github.com/layneson/unistatus/unicorn"
	"github.com/layneson/unistatus/weather"
)

const (
	credentialsFile = "/opt/unistatus/credentials.json"
	configFile      = "config.json"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	credentials, err := readCredentials(credentialsFile)
	if err != nil {
		panic(err)
	}

	err = config.ReadConfig(configFile)
	if err != nil {
		panic(err)
	}

	err = unicorn.InitProvider(unicorn.HATProvider{})
	if err != nil {
		panic(err)
	}

	wprovider, err := weather.NewWunderground(credentials)
	if err != nil {
		panic(err)
	}

	wstatus := display.NewWeatherStatus(wprovider)

	err = wstatus.Init()
	if err != nil {
		panic(err)
	}

	wstatus.Display(20)

	unicorn.Deinit()
}

//readCredentials reads the credentials JSON file from the specified file path, returning the map representation of its contents.
func readCredentials(file string) (map[string]string, error) {
	//Read complete file contents
	conts, err := ioutil.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("missing credentials file")
		}
		return nil, err
	}

	//Unmarshal credentials into a string->string map
	var credentials map[string]string
	err = json.Unmarshal(conts, &credentials)
	return credentials, err
}
