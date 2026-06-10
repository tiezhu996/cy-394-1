package service

import (
	"fitnessapi/internal/model"
	"fitnessapi/internal/repository"
	"sort"
	"time"
)

type GoalService struct{ repo repository.GoalRepository }

func NewGoalService(repo repository.GoalRepository) GoalService { return GoalService{repo: repo} }
func (s GoalService) Save(goal *model.Goal) error               { return s.repo.Save(goal) }
func (s GoalService) Latest(userID uint) (model.Goal, error)    { return s.repo.Latest(userID) }
func (s GoalService) List(userID uint) ([]model.Goal, error)    { return s.repo.List(userID) }
func GoalProgress(goal model.Goal, summary Summary) map[string]float64 {
	return map[string]float64{
		"sessions_percent": percent(float64(summary.Count), float64(goal.WeeklySessions)),
		"minutes_percent":  percent(float64(summary.TotalDuration), float64(goal.WeeklyMinutes)),
		"calories_percent": percent(summary.TotalCalories, goal.WeeklyCalories),
	}
}
func percent(done, target float64) float64 {
	if target <= 0 {
		return 0
	}
	if done/target > 1 {
		return 100
	}
	return done / target * 100
}

func IsGoalMet(goal model.Goal, summary Summary) bool {
	progress := GoalProgress(goal, summary)
	return progress["sessions_percent"] >= 100 &&
		progress["minutes_percent"] >= 100 &&
		progress["calories_percent"] >= 100
}

func findGoalForWeek(goals []model.Goal, weekMonday time.Time) (model.Goal, bool) {
	var latest model.Goal
	found := false
	for _, g := range goals {
		if !g.EffectiveMonday.After(weekMonday) {
			latest = g
			found = true
		}
	}
	return latest, found
}

func groupRecordsByWeek(records []model.WorkoutRecord) map[time.Time][]model.WorkoutRecord {
	result := make(map[time.Time][]model.WorkoutRecord)
	for _, r := range records {
		weekStart := repository.WeekStart(r.OccurredAt)
		result[weekStart] = append(result[weekStart], r)
	}
	return result
}

func ConsecutiveWeeks(goals []model.Goal, records []model.WorkoutRecord, now time.Time) int {
	weekRecords := groupRecordsByWeek(records)
	currentWeek := repository.WeekStart(now)
	streak := 0
	for week := currentWeek; ; week = week.AddDate(0, 0, -7) {
		goal, hasGoal := findGoalForWeek(goals, week)
		if !hasGoal {
			break
		}
		recs := weekRecords[week]
		summary := BuildSummary(recs)
		if !IsGoalMet(goal, summary) {
			break
		}
		streak++
	}
	return streak
}

type WeeklyStreak struct {
	ConsecutiveWeeks int                `json:"consecutive_weeks"`
	CurrentWeekMet   bool               `json:"current_week_met"`
	WeeksDetail      []WeekStreakDetail `json:"weeks_detail"`
}

type WeekStreakDetail struct {
	WeekMonday   time.Time         `json:"week_monday"`
	GoalMet      bool              `json:"goal_met"`
	Progress     map[string]float64 `json:"progress"`
	Sessions     int               `json:"sessions"`
	Minutes      int               `json:"minutes"`
	Calories     float64           `json:"calories"`
	TargetSessions int             `json:"target_sessions"`
	TargetMinutes  int             `json:"target_minutes"`
	TargetCalories float64         `json:"target_calories"`
}

func BuildWeeklyStreak(goals []model.Goal, records []model.WorkoutRecord, now time.Time, recentWeeks int) WeeklyStreak {
	weekRecords := groupRecordsByWeek(records)
	currentWeek := repository.WeekStart(now)

	consecutive := 0
	for week := currentWeek; ; week = week.AddDate(0, 0, -7) {
		goal, hasGoal := findGoalForWeek(goals, week)
		if !hasGoal {
			break
		}
		recs := weekRecords[week]
		summary := BuildSummary(recs)
		if !IsGoalMet(goal, summary) {
			break
		}
		consecutive++
	}

	currentGoal, hasCurrentGoal := findGoalForWeek(goals, currentWeek)
	currentWeekMet := false
	if hasCurrentGoal {
		currentWeekMet = IsGoalMet(currentGoal, BuildSummary(weekRecords[currentWeek]))
	}

	var details []WeekStreakDetail
	weekDates := make([]time.Time, 0, recentWeeks)
	for i := 0; i < recentWeeks; i++ {
		weekDates = append(weekDates, currentWeek.AddDate(0, 0, -7*i))
	}
	sort.Slice(weekDates, func(i, j int) bool { return weekDates[i].Before(weekDates[j]) })

	for _, week := range weekDates {
		goal, hasGoal := findGoalForWeek(goals, week)
		recs := weekRecords[week]
		summary := BuildSummary(recs)
		met := false
		progress := map[string]float64{"sessions_percent": 0, "minutes_percent": 0, "calories_percent": 0}
		ts, tm, tc := 0, 0, 0.0
		if hasGoal {
			ts, tm, tc = goal.WeeklySessions, goal.WeeklyMinutes, goal.WeeklyCalories
			progress = GoalProgress(goal, summary)
			met = IsGoalMet(goal, summary)
		}
		details = append(details, WeekStreakDetail{
			WeekMonday:     week,
			GoalMet:        met,
			Progress:       progress,
			Sessions:       summary.Count,
			Minutes:        summary.TotalDuration,
			Calories:       summary.TotalCalories,
			TargetSessions: ts,
			TargetMinutes:  tm,
			TargetCalories: tc,
		})
	}

	return WeeklyStreak{
		ConsecutiveWeeks: consecutive,
		CurrentWeekMet:   currentWeekMet,
		WeeksDetail:      details,
	}
}
