package model

// User represents an application user authenticated via Firebase.
type User struct {
	BaseModel
	FirebaseUID  string `gorm:"type:text;uniqueIndex;not null"`
	Email        string `gorm:"type:text;not null"`
	DisplayName  string `gorm:"type:text;not null"`
	CurrencyCode string `gorm:"type:char(3);not null"`
}
