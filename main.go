package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/haisrig/dynamo/db"
)

func defaultEmpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "You are On!!!!")
}

func employeeListHandler(w http.ResponseWriter, r *http.Request) {
	result, _ := db.ListEmployees()
	t, _ := template.ParseFiles("template.html")
	t.Execute(w, result)
}

func employeeAddHandler(w http.ResponseWriter, r *http.Request) {
	result, _ := db.InsertEmployee(r.Body)
	json.NewEncoder(w).Encode("{status :" + result + "}")
}

func main() {
	r := mux.NewRouter()
	empRouter := r.PathPrefix("/employees").Subrouter()
	empRouter.HandleFunc("/", defaultEmpHandler).Methods(http.MethodGet)
	empRouter.HandleFunc("/list", employeeListHandler).Methods(http.MethodGet)
	empRouter.HandleFunc("/add", employeeAddHandler).Methods(http.MethodPut)
	server := http.Server{
		Addr:    ":8000",
		Handler: r,
	}
	log.Println("Server is listening on port 8000....")
	log.Fatal(server.ListenAndServe())
}
