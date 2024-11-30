package main

import (
	"encoding/json"
	"log"
	"net/http"
	models "receipts/models"
	util "receipts/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var receiptStore = models.ReceiptStore{
	Receipts: make(map[string]models.Receipt),
}

type ErrorResponse struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

type FieldError struct {
	Field          string      `json:"field"`
	InvalidValue   interface{} `json:"invalidValue,omitempty"`
	ExpectedFormat string      `json:"expectedFormat,omitempty"`
	Constraint     string      `json:"constraint,omitempty"`
}

func sendError(w http.ResponseWriter, code string, message string, details interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   code,
		Message: message,
		Details: details,
	})
}

func createReceipt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var receipt models.Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := uuid.New().String()

	// Store the receipt with its ID
	receiptStore.Receipts[id] = receipt

	// Return the created receipt and its ID
	response := map[string]interface{}{
		"id": id,
	}
	json.NewEncoder(w).Encode(response)
}

func getReceipt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]

	receipt, exists := receiptStore.Receipts[id]
	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}
	points := util.PointsFromRetailerName(receipt)
	points = points + util.PointsFromTotalPrice(receipt)
	points = points + util.PointsFromItemCount(receipt)
	points = points + util.PointsFromItemDescription(receipt)
	points = points + util.PointsFromPurchaseDate(receipt)
	points = points + util.PointsFromPurchaseTime(receipt)
	var payload = map[string]int{}
	payload["points"] = points
	json.NewEncoder(w).Encode(payload)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/receipts/process", createReceipt).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", getReceipt).Methods("GET")
	log.Println("Server is starting...")
	log.Fatal(http.ListenAndServe(":8000", router))
}
