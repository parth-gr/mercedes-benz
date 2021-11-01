package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

type User struct {
	Vin         string `json:"vin"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

func getCurrentChargeLevel(vin string) (string, error) {
	reqBody, err := json.Marshal(map[string]string{
		"vin": vin,
	})
	if err != nil {
		return "0", err
	}
	resp, err := http.Post("https://restmock.techgig.com/merc/charge_level",
		"application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "0", err
	}
	defer resp.Body.Close()
	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "0", err
	}
	charge := fmt.Sprintf("%v", data["currentChargeLevel"])
	fmt.Print(data["error"])
	if data["error"] != nil {
		return charge, fmt.Errorf("error", data["error"])
	}
	return charge, nil
}

// func getChargingStations(source string, destination string) (string, error) {

// }

func index(w http.ResponseWriter, r *http.Request) {
	// if the method is GET so redirect it back to the Home page
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	var t User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	errors := map[string]interface{}{
		"transactionId": strconv.Itoa(rand.Intn(100)),
		"errors":        map[string]string{"id": "9999", "description": "Technical1 Exception"},
	}
	charge, err := getCurrentChargeLevel(t.Vin)
	if err != nil {
		json.NewEncoder(w).Encode(errors)
	} else {
		// chargingStations, err := getChargingStations(t.Source, t.Destination)
		// if err != nil {
		// 	json.NewEncoder(w).Encode(errors)
		// } else {
		// 	calculate()
		json.NewEncoder(w).Encode(charge)
		//}
	}
}

func main() {
	http.HandleFunc("/charging", index) // router and function
	fmt.Println("server starting")
	log.Fatal(http.ListenAndServe(":8080", nil)) // Start Server at https://localhost:5000
}
