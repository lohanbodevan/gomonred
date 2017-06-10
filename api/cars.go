package api

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"time"
)

func GetCars(r *Repository, c *Cache) ([]Car, error) {
	cached, err := getCarsFromCache(c)
	if err != nil {
		return nil, err
	}

	if len(cached) > 0 {
		return cached, nil
	}

	cars, err := getCarsFromDB(r)
	if err != nil {
		return nil, err
	}

	if len(cars) > 0 {
		err = createCarsCache(c, cars)
		if err != nil {
			return nil, err
		}
	}

	return cars, nil
}

func getCarsFromCache(c *Cache) ([]Car, error) {
	var carsCached []Car
	cached, err := c.Client.Get("cars_list").Result()

	if err == c.KeyNotFound {
		log.Infof("API - getCarsFromCache - Key does not exists")
		return carsCached, nil
	}

	if err != nil {
		log.Errorf("API - getCarsFromCache - Fail to get from cache: %s", err)
		return nil, err
	}

	err = json.Unmarshal([]byte(cached), &carsCached)
	if err != nil {
		log.Errorf("API - getCarsFromCache - Fail to parse cache: %s", err)
		return nil, err
	}

	log.Infof("API - getCarsFromCache - Hit from cache: %s", cached)
	return carsCached, nil
}

func createCarsCache(c *Cache, cars []Car) error {
	json_cars, err := json.Marshal(cars)
	if err != nil {
		log.Errorf("API - createCarsCache - Fail marshal json. Error: %s JSON: %s", err, json_cars)
		return err
	}
	err = c.Client.Set("cars_list", json_cars, time.Minute*2).Err()
	if err != nil {
		log.Errorf("API - createCarsCache - Fail to set cache: %s", err)
		return err
	}

	return nil
}

func getCarsFromDB(r *Repository) ([]Car, error) {
	collection := r.Session.DB("gomonred").C("cars")

	cars := []Car{}
	err := collection.Find(nil).Iter().All(&cars)
	if err != nil {
		log.Errorf("API - getCarsFromDB - Fail to select: %s", err)
		return nil, err
	}

	log.Infof("API - getCarsFromDB - Get from db: %s", cars)
	return cars, nil
}

func CreateCars(car Car, r *Repository) error {
	collection := r.Session.DB("gomonred").C("cars")

	err := collection.Insert(&car)
	if err != nil {
		log.Errorf("API - createCars - Fail to insert: %s", err)
		return err
	}

	log.Infof("API - CreateCars - Car created: %s", car)
	return nil
}
