package api

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

type Repository struct {
	Session *mgo.Session
}

type Api struct {
	Repository *Repository
}

func (app *Api) ConfigureRoutes(router *mux.Router) {
	router.HandleFunc("/cars", app.GetCarsHandler).Methods("get")
	router.HandleFunc("/cars", app.CreateCarsHandler).Methods("post")
}

func InitServer() {
	repository := DatabseInit()
	defer repository.Session.Close()

	app := Api{Repository: &repository}

	mux := mux.NewRouter()
	app.ConfigureRoutes(mux)

	server := negroni.New(negroni.NewRecovery())
	server.UseHandler(mux)

	server.Run(":8080")
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
