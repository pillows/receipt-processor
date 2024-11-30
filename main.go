package main

import (
	"encoding/json"
	"log"
	"net/http"
	models "receipts/models"
	util "receipts/utils"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var receiptStore = models.ReceiptStore{
	Receipts: make(map[string]models.Receipt),
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

	// error checking
	if len(receipt.Retailer) == 0 {
		http.Error(w, "The receipt is invalid", http.StatusBadRequest)
		return
	}

	_, timeParseErr := time.Parse("2006-01-02", receipt.PurchaseDate)
	if timeParseErr != nil {
		// response := map[string]string{
		// 	"error": "Invalid purchase date. Please provide a valid date in YYYY-MM-DD format",
		// }
		// w.WriteHeader(http.StatusBadRequest)
		// json.NewEncoder(w).Encode(response)
		http.Error(w, "The receipt is invalid", http.StatusBadRequest)
		return
	}

	// Check if purchaseDate is not in the future
	purchaseDate, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if purchaseDate.After(time.Now()) {
		// response := map[string]string{
		// 	"error": "Purchase date cannot be in the future",
		// }
		// w.WriteHeader(http.StatusBadRequest)
		// json.NewEncoder(w).Encode(response)
		http.Error(w, "The receipt is invalid", http.StatusBadRequest)
		return
	}

	if len(receipt.PurchaseTime) == 0 {
		// response := map[string]string{
		// 	"error": "Purchase time is required",
		// }
		// w.WriteHeader(http.StatusBadRequest)
		// json.NewEncoder(w).Encode(response)
		http.Error(w, "The receipt is invalid", http.StatusBadRequest)
		return
	}

	// Check if purchaseTime is in HH:MM format
	_, purchaseTimeErr := time.Parse("15:04", receipt.PurchaseTime)
	if purchaseTimeErr != nil {
		// response := map[string]string{
		// 	"error": "Invalid purchase time. Please provide a valid purchaseTime in HH:MM format",
		// }
		// w.WriteHeader(http.StatusBadRequest)
		// json.NewEncoder(w).Encode(response)
		http.Error(w, "The receipt is invalid", http.StatusBadRequest)
		return
	}
	var totalCost float64
	if len(receipt.Items) == 0 {
		// response := map[string]string{
		// 	"error": "At least one item is required",
		// }
		// w.WriteHeader(http.StatusBadRequest)
		// json.NewEncoder(w).Encode(response)
		http.Error(w, "The receipt is invalid", http.StatusBadRequest)
		return
	}

	for _, item := range receipt.Items {
		if len(item.ShortDescription) == 0 {
			// response := map[string]string{
			// 	"error": "Short description is required for all items",
			// }
			// w.WriteHeader(http.StatusBadRequest)
			// json.NewEncoder(w).Encode(response)
			http.Error(w, "The receipt is invalid", http.StatusBadRequest)
			return
		}

		// Check if price is a valid float
		_, priceErr := strconv.ParseFloat(item.Price, 64)
		if priceErr != nil {
			// response := map[string]string{
			// 	"error": "Price for item #" + (strconv.Itoa(index + 1)) + " " + "with description: " + item.ShortDescription + " must be a string in decimal format without currency symbols (e.g., '4.99' not '$4.99')",
			// }
			// w.WriteHeader(http.StatusBadRequest)
			// json.NewEncoder(w).Encode(response)
			http.Error(w, "The receipt is invalid", http.StatusBadRequest)
			return
		}
		price, _ := strconv.ParseFloat(item.Price, 64)
		totalCost = totalCost + price
	}

	_, totalCostErr := strconv.ParseFloat(receipt.Total, 64)
	if totalCostErr != nil {
		// response := map[string]string{
		// 	"error": "Invalid total cost provided",
		// }
		// w.WriteHeader(http.StatusBadRequest)
		// json.NewEncoder(w).Encode(response)
		http.Error(w, "The receipt is invalid", http.StatusBadRequest)
		return
	}

	// Check if the total cost matches the provided total
	providedTotal, _ := strconv.ParseFloat(receipt.Total, 64)
	if totalCost != providedTotal {
		// response := map[string]string{
		// 	"error": "Total cost does not match the sum of item prices",
		// }
		// w.WriteHeader(http.StatusBadRequest)
		// json.NewEncoder(w).Encode(response)
		http.Error(w, "The receipt is invalid", http.StatusBadRequest)
		return
	}

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
		http.Error(w, "No receipt found for that id", http.StatusNotFound)
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
