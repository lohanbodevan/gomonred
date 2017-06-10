package api

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"time"
)

func GetCars(r *Repository, c *Cache) ([]Car, error) {
	cached, err := c.Client.Get("cars_list").Result()

	if err == c.KeyNotFound {
		log.Infof("API - GetCars - Key does not exists")
	} else if err != nil {
		log.Errorf("API - GetCars - Fail to get from cache: %s", err)
		return nil, err
	}

	if cached != "" {
		var carsCached []Car
		err := json.Unmarshal([]byte(cached), &carsCached)
		if err != nil {
			log.Errorf("API - GetCars - Fail to parse cache: %s", err)
			return nil, err
		}
		log.Infof("API - GetCars - Hit from cache: %s", cached)
		return carsCached, nil
	}

	collection := r.Session.DB("example").C("cars")

	cars := []Car{}
	err = collection.Find(nil).Iter().All(&cars)
	if err != nil {
		log.Errorf("API - GetCars - Fail to select: %s", err)
		return nil, err
	}

	json_cars, err := json.Marshal(cars)
	if err != nil {
		log.Errorf("API - GetCars - Fail marshal json: %s", err)
		return nil, err
	}
	log.Infof("API - GetCars - Get from db: %s", json_cars)

	err = c.Client.Set("cars_list", json_cars, time.Minute*5).Err()
	if err != nil {
		log.Errorf("API - GetCars - Fail to set cache: %s", err)
		return nil, err
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
