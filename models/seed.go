package models

import (
	"fmt"
	"strconv"
)

func SeedDatabase() {
	var accessList []UserAccess

	// Generate 500 entries with 100 users and 10 features
	for i := 1; i <= 500; i++ {
		userID := (i-1)%100 + 1
		featureID := (i-1)%10 + 101
		accessList = append(accessList, UserAccess{
			UserID:    userID,
			FeatureID: featureID,
		})
	}

	for _, access := range accessList {
		db.FirstOrCreate(&access)
		// Add a combination of UserID and FeatureID to the Bloom filter
		key := strconv.Itoa(access.UserID) + ":" + strconv.Itoa(access.FeatureID)
		bloomFilter.Add([]byte(key))
	}

	fmt.Println("Database seeded and Bloom filter initialized.")
}
