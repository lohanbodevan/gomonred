package api

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"os"
)

type Repository struct {
	Session *mgo.Session
}

type Cache struct {
	Client      *redis.Client
	KeyNotFound interface{}
}

type Api struct {
	Repository *Repository
	Cache      *Cache
}

func (app *Api) ConfigureRoutes(router *mux.Router) {
	router.HandleFunc("/cars", app.GetCarsHandler).Methods("get")
	router.HandleFunc("/cars", app.CreateCarsHandler).Methods("post")
}

func InitServer() {
	repository := DatabseInit()
	defer repository.Session.Close()

	cache := CacheInit()
	defer cache.Client.Close()

	app := Api{
		Repository: &repository,
		Cache:      &cache,
	}

	mux := mux.NewRouter()
	app.ConfigureRoutes(mux)

	server := negroni.New(negroni.NewRecovery())
	server.UseHandler(mux)

	serverAddr := ":" + os.Getenv("PORT")
	server.Run(serverAddr)
}

func DatabseInit() Repository {
	session, err := mgo.Dial(os.Getenv("DB_HOST"))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	repo := Repository{
		Session: session,
	}
	return repo
}

func CacheInit() Cache {
	serverAddr := os.Getenv("CACHE_HOST") + ":" + os.Getenv("CACHE_PORT")
	client := redis.NewClient(&redis.Options{
		Addr:     serverAddr,
		Password: "",
		DB:       0,
	})

	cache := Cache{
		Client:      client,
		KeyNotFound: redis.Nil,
	}
	return cache
}
