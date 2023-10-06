package main

import (
	"encoding/json"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type Animal struct {
	gorm.Model
	Name       string   `json:"name"`
	Age        int      `json:"age"`
	CategoryID uint     `json:"-"`
	Category   Category `json:"category" gorm:"foreignKey:CategoryID"`
}
type Category struct {
	gorm.Model
	Name    string   `json:"name"`
	Animals []Animal `json:"animals"`
}

var db *gorm.DB

func initDB() {
	dsn := "user=postgres password=postgres dbname=goApiSecondDB port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Миграция (автоматическое создание таблицы "books")
	db.AutoMigrate(&Animal{}, &Category{})
}

func getMuxRoute() {
	router := mux.NewRouter()
	router.HandleFunc("/api", getAllAnimals).Methods("GET")
	router.HandleFunc("/api", addAnimal).Methods("POST")
	router.HandleFunc("/api", updAnimal).Methods("PUT")
	router.HandleFunc("/api/delete/{id}", deleteAnimal).Methods("DELETE")
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
func main() {
	initDB()
	getMuxRoute()
}
func getAllAnimals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var animals []Animal
	db.Preload("Category").Find(&animals)
	json.NewEncoder(w).Encode(&animals)
}
func addAnimal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var animal Animal
	err := json.NewDecoder(r.Body).Decode(&animal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	animal.CategoryID = animal.Category.ID
	db.Save(&animal)
	json.NewEncoder(w).Encode(&animal)
}
func updAnimal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var animal Animal
	var updAnimal Animal
	err := json.NewDecoder(r.Body).Decode(&updAnimal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	db.Find(&animal, updAnimal.ID)
	animal.Name = updAnimal.Name
	animal.Age = updAnimal.Age
	db.Save(&animal)
	json.NewEncoder(w).Encode(&animal)
}
func deleteAnimal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.ParseUint(params["id"], 10, 64)
	db.Delete(&Animal{}, id)
	json.NewEncoder(w).Encode("delete completed")
}
