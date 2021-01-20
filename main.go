package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

type City struct {
	Name string `json:"city"`
}

type List struct {
	items []Object
}

type Object struct {
	Locale  string
	Celsius float64
}

type Temperature struct {
	Value float64
}

type Structure struct {
	Main interface{}
}

func ReadFile(path string) string {
	file, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Print(err)
	}

	content := string(file)

	return content
}

func WriteFile(f *os.File, path string, value string) {
	_, err := f.WriteString(value)

	if err != nil {
		fmt.Println("Error while writing txt file!")

		f.Close()
	}
}

func DeserializeJSON(content string) []City {
	jsonAsBytes := []byte(content)
	cities := make([]City, 0)

	err := json.Unmarshal(jsonAsBytes, &cities)

	if err != nil {
		panic(err)
	}

	return cities
}

func ConvertKelvinToCelsius(kelvin float64) float64 {
	temperature := kelvin - 273.15

	return temperature
}

func (l *List) Insert(item Object) {
	l.items = append(l.items, item)
}

func Request(wg *sync.WaitGroup, list *List, city string, key string) {
	wg.Add(1)

	defer wg.Done()

	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, key)

	resp, err := http.Get(url)
	if err != nil {
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	}

	var data map[string]interface{}

	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
	}

	var kelvin float64

	for index, _ := range data {
		if index == "main" {
			for nestedIndex, nestedValue := range data["main"].(map[string]interface{}) {
				if nestedIndex == "temp" {
					newValue := nestedValue.(float64)
					kelvin = newValue
				}
			}
		}
	}

	object := Object{Locale: city, Celsius: ConvertKelvinToCelsius(kelvin)}

	list.Insert(object)
}

func main() {
	// Tokens: [b5a4ac104e94da6e3288a9ee2ac02e50, ff74db949d72ebbaf303afa5fcaa49ee]

	start := time.Now()
	fmt.Println("Started!")

	var wg sync.WaitGroup

	file, err := os.Create("result.txt")

	if err != nil {
		fmt.Println("Error while writing txt file!")
	}

	list := List{}

	content := ReadFile("./cities.json")
	cities := DeserializeJSON(content)

	for _, city := range cities {
		go Request(&wg, &list, city.Name, "ff74db949d72ebbaf303afa5fcaa49ee")
	}

	wg.Wait()

	for _, result := range list.items {
		WriteFile(file, "./result.txt", fmt.Sprintf("Name: %s | Temperature: %.1f\n", result.Locale, result.Celsius))
	}

	elapsed := time.Since(start)
	fmt.Println(elapsed)
}
