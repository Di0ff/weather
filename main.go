package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	K "weather/KeyAPI"
)

type Data struct {
	T float64 `json:"temperature"`
	D string  `json:"description"`
}

func fetch(c string) (*Data, error) {
	key := K.KeyAPI
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", c, key)

	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, err
	}

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to fetch weather data: %s", data["message"])
	}

	w := &Data{
		T: data["main"].(map[string]interface{})["temp"].(float64),
		D: data["weather"].([]interface{})[0].(map[string]interface{})["description"].(string),
	}

	return w, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	c := r.URL.Query().Get("city")
	if c == "" {
		http.Error(w, "City parameter is required", http.StatusBadRequest)
		return
	}

	weather, err := fetch(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("The current temperature in %s is %.1f°C. Weather: %s", c, weather.T, weather.D)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, response)
}

func main() {
	http.HandleFunc("/weather", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
