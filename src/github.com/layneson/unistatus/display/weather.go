package display

import (
	"math/rand"
	"time"

	"github.com/layneson/unistatus/config"
	"github.com/layneson/unistatus/unicorn"
	"github.com/layneson/unistatus/weather"
)

//WeatherStatus represents a Status for displaying basic info about the current day's weather.
type WeatherStatus struct {
	//provider is the weather.Provider which supplies the current weather data.
	provider weather.Provider

	//condition is the current weather condition to be rendered.
	condition weather.Condition

	//tindicator holds the temperature display information.
	tindicator tempIndicator
}

//tempIndicator represents a struct with information about how the current temperature is to be displayed.
type tempIndicator struct {
	//r, g, b are the values of the color for the current temperature.
	r, g, b int

	//width is the width, in pixels, of the temperature indicator.
	width int
}

//NewWeatherStatus creates a Weather Status with initial data and returns a pointer to it.
func NewWeatherStatus(p weather.Provider) *WeatherStatus {
	return &WeatherStatus{provider: p}
}

//Init implements the method of the Status interface.
func (w *WeatherStatus) Init() error {
	precip, err := w.provider.Precipitation() // Get precipitation percentage for next hour
	if err != nil {
		return err
	}

	condition, err := w.provider.Condition()
	if err != nil {
		return err
	}

	if precip >= 50 { // There's a significant enough chance of rain to alert the user
		w.condition = weather.LightRain
		if precip >= 80 {
			w.condition = weather.HeavyRain
		}
	} else { // Just display whatever is currently going on
		w.condition = condition
	}

	temperature, err := w.provider.Temperature()
	if err != nil {
		return err
	}

	if temperature <= 5 {
		w.tindicator = tempIndicator{255, 255, 255, 1} // White
	} else if temperature <= 20 {
		w.tindicator = tempIndicator{139, 0, 204, 2} // Purple
	} else if temperature <= 32 {
		w.tindicator = tempIndicator{0, 180, 235, 3} // Light Blue
	} else if temperature <= 48 {
		w.tindicator = tempIndicator{0, 57, 214, 4} // Blue
	} else if temperature <= 60 {
		w.tindicator = tempIndicator{0, 214, 50, 5} // Green
	} else if temperature <= 70 {
		w.tindicator = tempIndicator{255, 251, 5, 6} // Yellow
	} else if temperature <= 85 {
		w.tindicator = tempIndicator{255, 147, 5, 7} // Orange
	} else {
		w.tindicator = tempIndicator{255, 5, 5, 8} // Red
	}

	return nil
}

const tickRate = 50 // ticks/second

//Display implements the method of the Status interface
func (w WeatherStatus) Display(seconds int) error {
	unicorn.SetBrightness(config.Current.Brightness)

	if w.condition == weather.LightRain || w.condition == weather.HeavyRain {
		return w.displayRain(seconds)
	} else if w.condition == weather.Clouds {
		return w.displayClouds(seconds)
	}

	return nil
}

func (w WeatherStatus) renderTemp() error {
	for x := 0; x < w.tindicator.width; x++ {
		err := unicorn.SetPixel(x, 7, w.tindicator.r, w.tindicator.g, w.tindicator.b)
		if err != nil {
			panic(err)
		}
	}

	for x := w.tindicator.width; x < 8; x++ {
		err := unicorn.SetPixel(x, 7, 0, 0, 0)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func (w WeatherStatus) displayRain(seconds int) error {
	acounter := 0

	var skip int
	if w.condition == weather.LightRain {
		skip = 5
	} else {
		skip = 2
	}

	targets := []int{0, 0, 0, 0, 0, 0, 0, 0}
	positions := []int{0, 0, 0, 0, 0, 0, 0, 0}
	vectors := []int{0, 0, 0, 0, 0, 0, 0, 0}

	trails := make([][]int, 8*8)

	for acounter < tickRate*seconds {
		if acounter%skip == 0 {
			for x := 0; x < 8; x++ {
				if positions[x] == targets[x] {
					targets[x] = rand.Intn(7)
					if targets[x] < positions[x] {
						vectors[x] = -1
					} else {
						vectors[x] = 1
					}

					continue
				}

				positions[x] += vectors[x]

				for y := 0; y <= positions[x]; y++ {
					trails[x*8+y] = []int{0, 0, 255}
				}
			}

			for x := 0; x < 8; x++ {
				for y := 0; y < 7; y++ {
					t := trails[x*8+y]

					if len(t) == 0 {
						t = []int{0, 0, 0}
						trails[x*8+y] = t
					}

					if y > positions[x] {
						if t[2] < 10 {
							t[2] = 10
						}

						t[2] -= 10

						trails[x*8+y] = t
					}

					unicorn.SetPixel(x, y, t[0], t[1], t[2])
				}
			}
		}

		w.renderTemp()

		err := unicorn.Show()
		if err != nil {
			panic(err)
		}

		acounter++
		time.Sleep((1000 / tickRate) * time.Millisecond)
	}

	return nil
}

const (
	grayMin = 120
	grayMax = 230
)

func (w WeatherStatus) displayClouds(seconds int) error {
	cloudGrays := make([]int, 8*8)
	targetGrays := make([]int, 8*8)

	for i := range targetGrays {
		targetGrays[i] = 175
		cloudGrays[i] = 175
	}

	acounter := 0
	for acounter < tickRate*seconds {
		for x := 0; x < 8; x++ {
			for y := 0; y < 7; y++ {
				i := x*8 + y
				if cloudGrays[i] == targetGrays[i] {

					targetGrays[i] = rand.Intn(grayMax-grayMin) + grayMin
				} else if cloudGrays[i] > targetGrays[i] {
					cloudGrays[i]--
				} else {
					cloudGrays[i]++
				}

				unicorn.SetPixel(x, y, cloudGrays[i], cloudGrays[i], cloudGrays[i])
			}
		}

		w.renderTemp()

		err := unicorn.Show()
		if err != nil {
			panic(err)
		}

		acounter++
		time.Sleep((1000 / tickRate) * time.Millisecond)
	}

	return nil
}
