package models

import (
	"fmt"
	"strconv"
)

// UserAccess represents a user's access rights in the database
type UserAccess struct {
	UserID    int `gorm:"primaryKey"`
	FeatureID int
}

func UserFeatureAccess(userID, featureID int) bool {

	// Check Bloom filter
	key := strconv.Itoa(userID) + ":" + strconv.Itoa(featureID)
	if !bloomFilter.Test([]byte(key)) {
		return false
	}

	// Query the database for confirmation
	var access UserAccess
	result := db.Where("user_id = ? AND feature_id = ?",
		userID, featureID).First(&access)
	if result.Error != nil {
		fmt.Printf("Database: User %d does not have access to Feature %d.\n",
			userID, featureID)
		return false
	}

	fmt.Printf("Database: User %d has access to Feature %d.\n", userID, featureID)
	return true
}
