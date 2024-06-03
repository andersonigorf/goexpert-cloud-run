package main

import (
	"encoding/json"
	"fmt"
	"github.com/andersonigorf/goexpert-cloud-run/configs"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

type ViaCep struct {
	City  string `json:"localidade"`
	Error bool   `json:"erro"`
}
type WeatherApi struct {
	Current struct {
		TempC float64 `json:"temp_C"`
		TempF float64 `json:"temp_F"`
		TempK float64 `json:"temp_K"`
	} `json:"current"`
}

const ViaCepUrl = "http://viacep.com.br/ws/%s/json/"
const WeatherApiUrl = "http://api.weatherapi.com/v1/current.json?key=%s&aqi=no&q=%s"

var cfg *configs.Config

func main() {
	cfg, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", SearchCepAndWeatherHandler(cfg))
	fmt.Println("Starting web server on port", cfg.WebServerPort)
	http.ListenAndServe(cfg.WebServerPort, nil)
}

const (
	BadRequestMessage   = "cep parameter is required"
	UnprocessibleEntity = "invalid zipcode"
	InternalServerError = "error while searching for "
	NotFoundMessage     = "can not find "
)

func writeResponse(w http.ResponseWriter, statusCode int, errPrefix string, errMsg string) {
	w.WriteHeader(statusCode)
	w.Write([]byte(fmt.Sprintf("%s%s", errPrefix, errMsg)))
}

func searchData(url string, data interface{}) error {
	log.Printf("Requesting URL: %s", url)
	req, err := http.Get(url)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(res, data)
}

func searchTemperature(city string, weatherApiKey string) (*WeatherApi, error) {
	cityEscaped := url.QueryEscape(city)
	weatherApiURL := fmt.Sprintf(WeatherApiUrl, weatherApiKey, cityEscaped)

	var dataWeatherApi WeatherApi
	err := searchData(weatherApiURL, &dataWeatherApi)
	dataWeatherApi.Current.TempF = dataWeatherApi.Current.TempC*1.8 + 32
	dataWeatherApi.Current.TempK = dataWeatherApi.Current.TempC + 273
	return &dataWeatherApi, err
}

func searchCep(cepParam string) (*ViaCep, error) {
	requestURL := fmt.Sprintf(ViaCepUrl, cepParam)
	var dataViaCep ViaCep
	return &dataViaCep, searchData(requestURL, &dataViaCep)
}

func SearchCepAndWeatherHandler(config *configs.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		cepParam := r.URL.Query().Get("cep")
		if cepParam == "" {
			writeResponse(w, http.StatusBadRequest, "", BadRequestMessage)
			return
		}
		validate := regexp.MustCompile(`^[0-9]{5}-?[0-9]{3}$`)
		if !validate.MatchString(cepParam) {
			writeResponse(w, http.StatusUnprocessableEntity, "", UnprocessibleEntity)
			return
		}
		cep, err := searchCep(cepParam)
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, InternalServerError, "cep: "+err.Error())
			return
		}
		if cep.Error {
			writeResponse(w, http.StatusNotFound, NotFoundMessage, "zipcode")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		weather, err := searchTemperature(cep.City, config.WeatherApiKey)
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, InternalServerError, "temperature: "+err.Error())
			return
		}
		if weather != nil {
			json.NewEncoder(w).Encode(weather.Current)
		} else {
			writeResponse(w, http.StatusNotFound, NotFoundMessage, "temperature")
		}
	}
}
