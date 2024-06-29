package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Appointment struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Date  string `json:"date"`
}

/*
	var appointments = []Appointment{
		{ID: "1", Title: "Doctor Advord M. Meeting with patient John Trab", Date: "2024-06-30"},
		{ID: "2", Title: "Doctor Lamart A. Meeting with patient Sarah Crave", Date: "2024-07-01"},
	}
*/
var appointments = make(map[int]Appointment)

func getAppointments(w http.ResponseWriter, r *http.Request) {
	appointmentList := make([]Appointment, 0, len(appointments))
	for _, appointment := range appointments {
		appointmentList = append(appointmentList, appointment)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appointments)
}

func getAppointment(w http.ResponseWriter, r *http.Request) {
	paramId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		http.Error(w, "ID is not present!", http.StatusBadRequest)
		return
	}
	appointment, ok := appointments[id]
	if !ok {
		http.Error(w, "ID not found!", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appointment)
}

func createAppointment(w http.ResponseWriter, r *http.Request) {
	var newAppointment Appointment
	err := json.NewDecoder(r.Body).Decode(&newAppointment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if newAppointment.ID == 0 {
		http.Error(w, "ID is required", http.StatusBadRequest)
	}
	_, ok := appointments[newAppointment.ID]
	if ok {
		http.Error(w, "Appointment ID is already present!", http.StatusConflict)
		return
	}
	appointments[newAppointment.ID] = newAppointment

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newAppointment)
}

func main() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/appointments", getAppointments)
	router.Get("/appointments/{id}", getAppointment)
	router.Post("/appointments", createAppointment)

	fmt.Println("Starting server at port 8000")
	err := http.ListenAndServe(":8000", router)
	if err != nil {
		return
	}
}
