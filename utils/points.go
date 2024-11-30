package points

import (
	"math"
	models "receipts/models"
	"strconv"
	"strings"
	"unicode"
)

func countAlphanumeric(s string) int {
	count := 0
	for _, char := range s {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			count++
		}
	}
	return count
}

func PointsFromTotalPrice(receipt models.Receipt) int {
	// 50 points if the total is a round dollar amount with no cents.
	// 25 points if the total is a multiple of 0.25.
	points := 0
	f, _ := strconv.ParseFloat(receipt.Total, 64)
	if math.Mod(f, 0.25) == 0 {
		points = points + 25
	}

	if math.Mod(f, 1.00) == 0 {
		points = points + 50
	}
	return points
}

func PointsFromItemCount(receipt models.Receipt) int {
	// 5 points for every two items on the receipt.
	points := 0
	itemCount := len(receipt.Items)
	points = points + (itemCount/2)*5
	return points
}

func PointsFromItemDescription(receipt models.Receipt) int {
	// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round
	// up to the nearest integer. The result is the number of points earned.
	points := 0
	for _, item := range receipt.Items {
		descriptionLength := len(strings.TrimSpace(item.ShortDescription))
		price, _ := strconv.ParseFloat(item.Price, 64)
		if descriptionLength%3 == 0 {
			price = price * 0.2
			points = points + int(math.Ceil(price))
		}
	}
	return points
}

func PointsFromPurchaseDate(receipt models.Receipt) int {
	// 6 points if the day in the purchase date is odd.
	points := 0
	day, _ := strconv.Atoi(receipt.PurchaseDate[8:])
	if day%2 != 0 {
		points = points + 6
	}
	return points
}

func PointsFromRetailerName(receipt models.Receipt) int {
	var points = 0
	var retailer = receipt.Retailer
	// One point for every alphanumeric character in the retailer name.
	points = points + countAlphanumeric(retailer)
	return points
}

func PointsFromPurchaseTime(receipt models.Receipt) int {
	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	points := 0
	purchaseTime := receipt.PurchaseTime
	hour, _ := strconv.Atoi(purchaseTime[:2])
	if hour >= 14 && hour <= 16 {
		points = points + 10
	}
	return points
}
