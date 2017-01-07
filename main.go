package main

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
)

type Repository struct {
	Session *mgo.Session
}

type Car struct {
	Name  string `json:"name"`
	Brand string `json:"brand"`
}

func main() {
	repository := DatabseInit()
	defer repository.Session.Close()

	http.HandleFunc("/cars", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: Endpoint /cars called", r.Method)

		if r.Method == http.MethodGet {
			cars := getCars(repository)
			response, err := json.Marshal(cars)
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
			}
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.Write(response)
		} else if r.Method == http.MethodPost {
			var payload Car
			err := json.NewDecoder(r.Body).Decode(&payload)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}

			err = createCars(repository, payload)
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
			}

			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	log.Print("UP and Running at 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func DatabseInit() Repository {
	session, err := mgo.Dial("mongo_db")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	repo := Repository{
		Session: session,
	}
	return repo
}

func getCars(repo Repository) []Car {
	collection := repo.Session.DB("example").C("cars")

	cars := []Car{}
	err := collection.Find(nil).Iter().All(&cars)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return cars
}

func createCars(repo Repository, car Car) error {
	collection := repo.Session.DB("example").C("cars")

	err := collection.Insert(&car)
	if err != nil {
		log.Fatal(err)
	}

	return err
}
