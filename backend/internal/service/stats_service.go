package service

import "fitnessapi/internal/model"

type TypeShare struct {
	SportType string  `json:"sport_type"`
	Calories  float64 `json:"calories"`
	Ratio     float64 `json:"ratio"`
}

type Summary struct {
	TotalDuration int         `json:"total_duration"`
	TotalDistance float64     `json:"total_distance"`
	TotalCalories float64     `json:"total_calories"`
	Count         int         `json:"count"`
	TypeShares    []TypeShare `json:"type_shares"`
}

func BuildSummary(records []model.WorkoutRecord) Summary {
	result := Summary{Count: len(records)}
	byType := map[string]float64{}
	for _, item := range records {
		result.TotalDuration += item.DurationMin
		result.TotalDistance += item.DistanceKm
		result.TotalCalories += item.Calories
		byType[item.SportType] += item.Calories
	}
	for sport, calories := range byType {
		ratio := 0.0
		if result.TotalCalories > 0 {
			ratio = calories / result.TotalCalories
		}
		result.TypeShares = append(result.TypeShares, TypeShare{SportType: sport, Calories: calories, Ratio: ratio})
	}
	return result
}

func PersonalRecords(records []model.WorkoutRecord) map[string]float64 {
	pr := map[string]float64{"longest_distance": 0, "fastest_pace": 0, "max_weight": 0}
	for _, r := range records {
		if r.DistanceKm > pr["longest_distance"] {
			pr["longest_distance"] = r.DistanceKm
		}
		if r.WeightKg > pr["max_weight"] {
			pr["max_weight"] = r.WeightKg
		}
		if r.DistanceKm > 0 {
			pace := float64(r.DurationMin) / r.DistanceKm
			if pr["fastest_pace"] == 0 || pace < pr["fastest_pace"] {
				pr["fastest_pace"] = pace
			}
		}
	}
	return pr
}
