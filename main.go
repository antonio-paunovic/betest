package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DB_NAME = "roycebetest.db"

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func main() {
	a := &App{}
	a.InitializeDB()

	fmt.Println("DB main", a.DB)
	a.HandleHTTP()
}

func (a *App) InitializeDB() {
	db, err := gorm.Open(sqlite.Open(DB_NAME), &gorm.Config{})
	if err != nil {
		panic("Starting database failed")
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("Failed to initialize the DB")
	}
	a.DB = db
}

func (a *App) HandleHTTP() {
	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/user/{id}", a.get).Methods(http.MethodGet)
	a.Router.HandleFunc("/users", a.getUsers).Methods(http.MethodGet)
	a.Router.HandleFunc("/user", a.post).Methods(http.MethodPost)
	a.Router.HandleFunc("/user/{id}", a.put).Methods(http.MethodPut)
	a.Router.HandleFunc("/user/{id}", a.delete).Methods(http.MethodDelete)
	// r.HandleFunc("/", notFound)
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, err string) {
	fmt.Println("Responding with error")
	respondWithJSON(w, code, map[string]string{"error": err})
}

func (a *App) get(w http.ResponseWriter, r *http.Request) {
	// vars := r.URL.Query()
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64) //vars["id"]
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid ID")
		return
	}

	u := &User{ID: id}
	if err := u.getUser(a.DB); err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {
	users := []User{}
	err := dbGetUsers(&users, a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, users)
}

type DateOfBirth time.Time

func (dob *DateOfBirth) UnmarshalJSON(b []byte) error {
	fmt.Println("UNMARSHALLING DOB")
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*dob = DateOfBirth(t)
	return nil
}

func (dob DateOfBirth) MarshalJSON() ([]byte, error) {
	return json.Marshal(dob)
}

func (a *App) post(w http.ResponseWriter, r *http.Request) {
	u := &User{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(u); err != nil {
		fmt.Println(err.Error()) // TODO
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	defer r.Body.Close()

	// Check if time format is valid according to ISO 8601
	if _, err := time.Parse("2006-01-02", u.Dob); err != nil {
		fmt.Println(err.Error()) // TODO
		respondWithError(w, http.StatusBadRequest, "invalid date of birth format")
		return

	}

	if err := u.createUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) put(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64) //vars["id"]
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid ID")
		return
	}
	u := &User{ID: id}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(u); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := u.updateUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid ID")
		return
	}
	u := &User{ID: id}
	if err := u.deleteUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, u)
}
