package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	FriendKey string    `json:"friend_key"`
	CreatedAt time.Time `json:"created_at"`
}

type WorkoutRecord struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `json:"user_id"`
	SportType   string    `json:"sport_type"`
	DurationMin int       `json:"duration_min"`
	DistanceKm  float64   `json:"distance_km"`
	Calories    float64   `json:"calories"`
	HeartRate   int       `json:"heart_rate"`
	WeightKg    float64   `json:"weight_kg"`
	OccurredAt  time.Time `json:"occurred_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Goal struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `json:"user_id"`
	WeeklySessions  int       `json:"weekly_sessions"`
	WeeklyMinutes   int       `json:"weekly_minutes"`
	WeeklyCalories  float64   `json:"weekly_calories"`
	EffectiveMonday time.Time `json:"effective_monday"`
	CreatedAt       time.Time `json:"created_at"`
}
