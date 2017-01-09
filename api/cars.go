package api

import (
	log "github.com/Sirupsen/logrus"
)

func GetCars(r *Repository) ([]Car, error) {
	collection := r.Session.DB("example").C("cars")

	cars := []Car{}
	err := collection.Find(nil).Iter().All(&cars)
	if err != nil {
		log.Errorf("API - getCars - Fail to select: %s", err)
		return cars, err
	}

	return cars, nil
}

func CreateCars(car Car, r *Repository) error {
	collection := r.Session.DB("example").C("cars")

	err := collection.Insert(&car)
	if err != nil {
		log.Errorf("API - createCars - Fail to insert: %s", err)
		return err
	}

	return nil
}
