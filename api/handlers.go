package api

import (
	"encoding/json"
	"net/http"
)

func (a *Api) GetCarsHandler(w http.ResponseWriter, r *http.Request) {
	cars, err := GetCars(a.Repository, a.Cache)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	var response []byte

	if len(cars) > 0 {
		response, err = json.Marshal(cars)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
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
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}
