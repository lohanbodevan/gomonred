package api

import (
	"encoding/json"
	"net/http"
)

func (a *Api) GetCarsHandler(w http.ResponseWriter, r *http.Request) {
	cars, err := GetCars(a.Repository)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	response, err := json.Marshal(cars)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(response)
}

func (a *Api) CreateCarsHandler(w http.ResponseWriter, r *http.Request) {
	var payload Car
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	err = CreateCars(payload, a.Repository)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	w.WriteHeader(http.StatusCreated)
}
