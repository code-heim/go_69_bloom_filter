package main

import (
	"encoding/json"
	"fmt"
	"go_bloom_filter/models"
	"log"
	"net/http"
	"strconv"

	"github.com/bits-and-blooms/bloom/v3"
)

func init() {
	models.DBInit()

	models.BloomFilterSetup()

	// Seed the database
	models.SeedDatabase()
}

func UserFeatureCheck(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from query parameters
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w,
			"Missing 'userID' query parameter",
			http.StatusBadRequest)
		return
	}

	// Convert userID to an integer
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w,
			"Invalid 'userID'. Must be an integer.",
			http.StatusBadRequest)
		return
	}

	// Extract feature ID from query parameters
	featureIDStr := r.URL.Query().Get("feature_id")
	if featureIDStr == "" {
		http.Error(w,
			"Missing 'featureID' query parameter",
			http.StatusBadRequest)
		return
	}

	// Convert featureID to an integer
	featureID, err := strconv.Atoi(featureIDStr)
	if err != nil {
		http.Error(w,
			"Invalid 'featureID'. Must be an integer.",
			http.StatusBadRequest)
		return
	}
	access := models.UserFeatureAccess(userID, featureID)

	// Create API response
	response := map[string]bool{
		"access": access,
	}

	// Write JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// EstimateAPIResponse structure to format JSON responses
type EstimateAPIResponse struct {
	EstimatedFPRate  float64 `json:"estimated_false_positive_rate,omitempty"`
	EstimatedM       uint    `json:"estimated_m,omitempty"`
	EstimatedK       uint    `json:"estimated_k,omitempty"`
	DesiredFPRate    float64 `json:"desired_fp_rate,omitempty"`
	ExpectedAccuracy string  `json:"expected_accuracy,omitempty"`
}

// Handler function for estimating false positive rate
func EstimateFPHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	query := r.URL.Query()
	nStr := query.Get("n")
	fpStr := query.Get("desired_fp_rate")

	if nStr == "" || fpStr == "" {
		http.Error(w, "Missing 'n' (number of elements) or 'desired_fp_rate' query parameters", http.StatusBadRequest)
		return
	}

	// Convert inputs to numeric values
	n, err := strconv.Atoi(nStr)
	if err != nil || n <= 0 {
		http.Error(w, "'n' must be a positive integer", http.StatusBadRequest)
		return
	}

	fp, err := strconv.ParseFloat(fpStr, 64)
	if err != nil || fp <= 0 || fp >= 1 {
		http.Error(w, "'desired_fp_rate' must be a float between 0 and 1", http.StatusBadRequest)
		return
	}

	// Estimate Bloom filter parameters
	m, k := bloom.EstimateParameters(uint(n), fp)

	// Validate estimated parameters
	estimatedFP := bloom.EstimateFalsePositiveRate(m, k, uint(n))

	// Prepare response
	response := EstimateAPIResponse{
		EstimatedFPRate:  estimatedFP,
		EstimatedM:       m,
		EstimatedK:       k,
		DesiredFPRate:    fp,
		ExpectedAccuracy: fmt.Sprintf("Estimated false positive rate is within %.2f%% of the desired rate.", (estimatedFP-fp)/fp*100),
	}

	// Write JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Set up HTTP server
	http.HandleFunc("/feature_access", UserFeatureCheck)
	http.HandleFunc("/estimate_fp", EstimateFPHandler)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
