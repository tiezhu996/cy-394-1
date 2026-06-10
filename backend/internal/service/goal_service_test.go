package service

import (
	"fitnessapi/internal/model"
	"fitnessapi/internal/repository"
	"testing"
	"time"
)

func TestIsGoalMet_PartialTargets(t *testing.T) {
	tests := []struct {
		name     string
		goal     model.Goal
		summary  Summary
		expected bool
	}{
		{
			name: "仅配置次数，达标",
			goal: model.Goal{WeeklySessions: 3, WeeklyMinutes: 0, WeeklyCalories: 0},
			summary: Summary{Count: 3},
			expected: true,
		},
		{
			name: "仅配置次数，未达标",
			goal: model.Goal{WeeklySessions: 3, WeeklyMinutes: 0, WeeklyCalories: 0},
			summary: Summary{Count: 2},
			expected: false,
		},
		{
			name: "仅配置时长，达标",
			goal: model.Goal{WeeklySessions: 0, WeeklyMinutes: 150, WeeklyCalories: 0},
			summary: Summary{TotalDuration: 150},
			expected: true,
		},
		{
			name: "仅配置卡路里，达标",
			goal: model.Goal{WeeklySessions: 0, WeeklyMinutes: 0, WeeklyCalories: 1500},
			summary: Summary{TotalCalories: 1500},
			expected: true,
		},
		{
			name: "配置次数+时长，全部达标",
			goal: model.Goal{WeeklySessions: 3, WeeklyMinutes: 150, WeeklyCalories: 0},
			summary: Summary{Count: 3, TotalDuration: 150},
			expected: true,
		},
		{
			name: "配置次数+时长，次数未达标",
			goal: model.Goal{WeeklySessions: 3, WeeklyMinutes: 150, WeeklyCalories: 0},
			summary: Summary{Count: 2, TotalDuration: 150},
			expected: false,
		},
		{
			name: "全部配置，全部达标",
			goal: model.Goal{WeeklySessions: 3, WeeklyMinutes: 150, WeeklyCalories: 1500},
			summary: Summary{Count: 3, TotalDuration: 150, TotalCalories: 1500},
			expected: true,
		},
		{
			name: "全部配置，卡路里未达标",
			goal: model.Goal{WeeklySessions: 3, WeeklyMinutes: 150, WeeklyCalories: 1500},
			summary: Summary{Count: 3, TotalDuration: 150, TotalCalories: 1000},
			expected: false,
		},
		{
			name: "全部为0，视为未达标",
			goal: model.Goal{WeeklySessions: 0, WeeklyMinutes: 0, WeeklyCalories: 0},
			summary: Summary{Count: 10, TotalDuration: 1000, TotalCalories: 10000},
			expected: false,
		},
		{
			name: "超额完成也算达标",
			goal: model.Goal{WeeklySessions: 3, WeeklyMinutes: 150, WeeklyCalories: 0},
			summary: Summary{Count: 5, TotalDuration: 200},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsGoalMet(tt.goal, tt.summary)
			if result != tt.expected {
				t.Errorf("IsGoalMet() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConsecutiveWeeks_SingleTarget(t *testing.T) {
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	thisWeek := repository.WeekStart(now)
	prevWeek1 := thisWeek.AddDate(0, 0, -7)
	prevWeek2 := thisWeek.AddDate(0, 0, -14)
	prevWeek3 := thisWeek.AddDate(0, 0, -21)

	goals := []model.Goal{
		{ID: 1, UserID: 1, WeeklySessions: 3, WeeklyMinutes: 0, WeeklyCalories: 0, EffectiveMonday: prevWeek3},
	}

	records := []model.WorkoutRecord{
		{UserID: 1, OccurredAt: thisWeek.AddDate(0, 0, 1)},
		{UserID: 1, OccurredAt: thisWeek.AddDate(0, 0, 2)},
		{UserID: 1, OccurredAt: thisWeek.AddDate(0, 0, 3)},
		{UserID: 1, OccurredAt: prevWeek1.AddDate(0, 0, 1)},
		{UserID: 1, OccurredAt: prevWeek1.AddDate(0, 0, 2)},
		{UserID: 1, OccurredAt: prevWeek1.AddDate(0, 0, 3)},
		{UserID: 1, OccurredAt: prevWeek2.AddDate(0, 0, 1)},
		{UserID: 1, OccurredAt: prevWeek2.AddDate(0, 0, 2)},
		{UserID: 1, OccurredAt: prevWeek2.AddDate(0, 0, 3)},
	}

	result := ConsecutiveWeeks(goals, records, now)
	if result != 3 {
		t.Errorf("ConsecutiveWeeks() = %v, want 3", result)
	}
}

func TestConsecutiveWeeks_WithGap(t *testing.T) {
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	thisWeek := repository.WeekStart(now)
	prevWeek1 := thisWeek.AddDate(0, 0, -7)
	prevWeek2 := thisWeek.AddDate(0, 0, -14)
	prevWeek3 := thisWeek.AddDate(0, 0, -21)

	goals := []model.Goal{
		{ID: 1, UserID: 1, WeeklySessions: 3, WeeklyMinutes: 0, WeeklyCalories: 0, EffectiveMonday: prevWeek3},
	}

	records := []model.WorkoutRecord{
		{UserID: 1, OccurredAt: thisWeek.AddDate(0, 0, 1)},
		{UserID: 1, OccurredAt: thisWeek.AddDate(0, 0, 2)},
		{UserID: 1, OccurredAt: thisWeek.AddDate(0, 0, 3)},
		{UserID: 1, OccurredAt: prevWeek1.AddDate(0, 0, 1)},
		{UserID: 1, OccurredAt: prevWeek1.AddDate(0, 0, 2)},
		{UserID: 1, OccurredAt: prevWeek2.AddDate(0, 0, 1)},
		{UserID: 1, OccurredAt: prevWeek2.AddDate(0, 0, 2)},
		{UserID: 1, OccurredAt: prevWeek2.AddDate(0, 0, 3)},
	}

	result := ConsecutiveWeeks(goals, records, now)
	if result != 1 {
		t.Errorf("ConsecutiveWeeks() = %v, want 1 (prevWeek1 only had 2 sessions)", result)
	}
}

func TestConsecutiveWeeks_MultipleTargets(t *testing.T) {
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	thisWeek := repository.WeekStart(now)
	prevWeek1 := thisWeek.AddDate(0, 0, -7)
	prevWeek2 := thisWeek.AddDate(0, 0, -14)

	goals := []model.Goal{
		{ID: 1, UserID: 1, WeeklySessions: 3, WeeklyMinutes: 150, WeeklyCalories: 0, EffectiveMonday: prevWeek2},
	}

	records := []model.WorkoutRecord{
		{UserID: 1, DurationMin: 60, OccurredAt: thisWeek.AddDate(0, 0, 1)},
		{UserID: 1, DurationMin: 50, OccurredAt: thisWeek.AddDate(0, 0, 2)},
		{UserID: 1, DurationMin: 40, OccurredAt: thisWeek.AddDate(0, 0, 3)},
		{UserID: 1, DurationMin: 60, OccurredAt: prevWeek1.AddDate(0, 0, 1)},
		{UserID: 1, DurationMin: 60, OccurredAt: prevWeek1.AddDate(0, 0, 2)},
		{UserID: 1, DurationMin: 40, OccurredAt: prevWeek1.AddDate(0, 0, 3)},
	}

	result := ConsecutiveWeeks(goals, records, now)
	if result != 2 {
		t.Errorf("ConsecutiveWeeks() = %v, want 2 (both weeks: 3 sessions + 150 min)", result)
	}
}

func TestConsecutiveWeeks_GoalChanges(t *testing.T) {
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	thisWeek := repository.WeekStart(now)
	prevWeek1 := thisWeek.AddDate(0, 0, -7)
	prevWeek2 := thisWeek.AddDate(0, 0, -14)
	prevWeek3 := thisWeek.AddDate(0, 0, -21)

	goals := []model.Goal{
		{ID: 1, UserID: 1, WeeklySessions: 5, WeeklyMinutes: 0, WeeklyCalories: 0, EffectiveMonday: prevWeek3},
		{ID: 2, UserID: 1, WeeklySessions: 3, WeeklyMinutes: 0, WeeklyCalories: 0, EffectiveMonday: prevWeek1},
	}

	records := []model.WorkoutRecord{
		{UserID: 1, OccurredAt: thisWeek.AddDate(0, 0, 1)},
		{UserID: 1, OccurredAt: thisWeek.AddDate(0, 0, 2)},
		{UserID: 1, OccurredAt: thisWeek.AddDate(0, 0, 3)},
		{UserID: 1, OccurredAt: prevWeek1.AddDate(0, 0, 1)},
		{UserID: 1, OccurredAt: prevWeek1.AddDate(0, 0, 2)},
		{UserID: 1, OccurredAt: prevWeek1.AddDate(0, 0, 3)},
		{UserID: 1, OccurredAt: prevWeek2.AddDate(0, 0, 1)},
		{UserID: 1, OccurredAt: prevWeek2.AddDate(0, 0, 2)},
		{UserID: 1, OccurredAt: prevWeek2.AddDate(0, 0, 3)},
		{UserID: 1, OccurredAt: prevWeek2.AddDate(0, 0, 4)},
		{UserID: 1, OccurredAt: prevWeek2.AddDate(0, 0, 5)},
	}

	result := ConsecutiveWeeks(goals, records, now)
	if result != 3 {
		t.Errorf("ConsecutiveWeeks() = %v, want 3 (goal changes from 5 to 3, prevWeek2 had 5 sessions, prevWeek1+thisWeek had 3)", result)
	}
}

func TestBuildWeeklyStreak(t *testing.T) {
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	thisWeek := repository.WeekStart(now)
	prevWeek1 := thisWeek.AddDate(0, 0, -7)
	prevWeek2 := thisWeek.AddDate(0, 0, -14)
	prevWeek3 := thisWeek.AddDate(0, 0, -21)

	goals := []model.Goal{
		{ID: 1, UserID: 1, WeeklySessions: 3, WeeklyMinutes: 0, WeeklyCalories: 0, EffectiveMonday: prevWeek3},
	}

	records := []model.WorkoutRecord{
		{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: thisWeek.AddDate(0, 0, 1)},
		{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: thisWeek.AddDate(0, 0, 2)},
		{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: thisWeek.AddDate(0, 0, 3)},
		{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: prevWeek1.AddDate(0, 0, 1)},
		{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: prevWeek1.AddDate(0, 0, 2)},
		{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: prevWeek1.AddDate(0, 0, 3)},
		{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: prevWeek2.AddDate(0, 0, 1)},
		{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: prevWeek2.AddDate(0, 0, 2)},
		{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: prevWeek2.AddDate(0, 0, 3)},
	}

	streak := BuildWeeklyStreak(goals, records, now, 4)

	if streak.ConsecutiveWeeks != 3 {
		t.Errorf("ConsecutiveWeeks = %v, want 3", streak.ConsecutiveWeeks)
	}
	if !streak.CurrentWeekMet {
		t.Errorf("CurrentWeekMet = false, want true")
	}
	if len(streak.WeeksDetail) != 4 {
		t.Errorf("WeeksDetail length = %v, want 4", len(streak.WeeksDetail))
	}
	if streak.WeeksDetail[3].WeekMonday != thisWeek {
		t.Errorf("Last week detail should be current week")
	}
	for i, d := range streak.WeeksDetail {
		if i == 0 {
			if d.GoalMet {
				t.Errorf("Week %d (%s) GoalMet = true, want false (no records that week)", i, d.WeekMonday)
			}
		} else {
			if !d.GoalMet {
				t.Errorf("Week %d (%s) GoalMet = false, want true", i, d.WeekMonday)
			}
		}
		if d.TargetSessions != 3 {
			t.Errorf("Week %d TargetSessions = %v, want 3", i, d.TargetSessions)
		}
		if d.TargetMinutes != 0 {
			t.Errorf("Week %d TargetMinutes = %v, want 0", i, d.TargetMinutes)
		}
	}
}

func TestConsecutiveWeeks_NoGoals(t *testing.T) {
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	thisWeek := repository.WeekStart(now)

	records := []model.WorkoutRecord{
		{UserID: 1, OccurredAt: thisWeek.AddDate(0, 0, 1)},
		{UserID: 1, OccurredAt: thisWeek.AddDate(0, 0, 2)},
		{UserID: 1, OccurredAt: thisWeek.AddDate(0, 0, 3)},
	}

	result := ConsecutiveWeeks(nil, records, now)
	if result != 0 {
		t.Errorf("ConsecutiveWeeks() = %v, want 0 (no goals configured)", result)
	}
}

func TestConsecutiveWeeks_NoRecords(t *testing.T) {
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	prevWeek1 := repository.WeekStart(now).AddDate(0, 0, -7)

	goals := []model.Goal{
		{ID: 1, UserID: 1, WeeklySessions: 3, WeeklyMinutes: 0, WeeklyCalories: 0, EffectiveMonday: prevWeek1},
	}

	result := ConsecutiveWeeks(goals, nil, now)
	if result != 0 {
		t.Errorf("ConsecutiveWeeks() = %v, want 0 (no records)", result)
	}
}

func TestConsecutiveWeeks_AllZeroTargets(t *testing.T) {
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	thisWeek := repository.WeekStart(now)
	prevWeek1 := thisWeek.AddDate(0, 0, -7)

	goals := []model.Goal{
		{ID: 1, UserID: 1, WeeklySessions: 0, WeeklyMinutes: 0, WeeklyCalories: 0, EffectiveMonday: prevWeek1},
	}

	records := []model.WorkoutRecord{
		{UserID: 1, OccurredAt: thisWeek.AddDate(0, 0, 1)},
		{UserID: 1, OccurredAt: prevWeek1.AddDate(0, 0, 1)},
	}

	result := ConsecutiveWeeks(goals, records, now)
	if result != 0 {
		t.Errorf("ConsecutiveWeeks() = %v, want 0 (all targets are zero, not met)", result)
	}
}

func TestConsecutiveWeeks_PartialGoalMet(t *testing.T) {
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	thisWeek := repository.WeekStart(now)
	prevWeek1 := thisWeek.AddDate(0, 0, -7)
	prevWeek2 := thisWeek.AddDate(0, 0, -14)

	goals := []model.Goal{
		{ID: 1, UserID: 1, WeeklySessions: 3, WeeklyMinutes: 100, WeeklyCalories: 0, EffectiveMonday: prevWeek2},
	}

	records := []model.WorkoutRecord{
		{UserID: 1, DurationMin: 40, OccurredAt: thisWeek.AddDate(0, 0, 1)},
		{UserID: 1, DurationMin: 40, OccurredAt: thisWeek.AddDate(0, 0, 2)},
		{UserID: 1, DurationMin: 40, OccurredAt: thisWeek.AddDate(0, 0, 3)},
		{UserID: 1, DurationMin: 40, OccurredAt: prevWeek1.AddDate(0, 0, 1)},
		{UserID: 1, DurationMin: 40, OccurredAt: prevWeek1.AddDate(0, 0, 2)},
	}

	result := ConsecutiveWeeks(goals, records, now)
	if result != 1 {
		t.Errorf("ConsecutiveWeeks() = %v, want 1 (this week: 3 sessions + 120 min >= 100; prevWeek1: only 2 sessions < 3)", result)
	}
}

func TestGoalProgress_PartialConfig(t *testing.T) {
	goal := model.Goal{WeeklySessions: 3, WeeklyMinutes: 0, WeeklyCalories: 1500}
	summary := Summary{Count: 3, TotalDuration: 200, TotalCalories: 1500}

	progress := GoalProgress(goal, summary)

	if progress["sessions_percent"] != 100 {
		t.Errorf("sessions_percent = %v, want 100", progress["sessions_percent"])
	}
	if progress["minutes_percent"] != 0 {
		t.Errorf("minutes_percent = %v, want 0 (target is 0, not configured)", progress["minutes_percent"])
	}
	if progress["calories_percent"] != 100 {
		t.Errorf("calories_percent = %v, want 100", progress["calories_percent"])
	}
}

func TestBuildWeeklyStreak_EmptyData(t *testing.T) {
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)

	streak := BuildWeeklyStreak(nil, nil, now, 4)

	if streak.ConsecutiveWeeks != 0 {
		t.Errorf("ConsecutiveWeeks = %v, want 0", streak.ConsecutiveWeeks)
	}
	if streak.CurrentWeekMet {
		t.Errorf("CurrentWeekMet = true, want false")
	}
	if len(streak.WeeksDetail) != 4 {
		t.Errorf("WeeksDetail length = %v, want 4", len(streak.WeeksDetail))
	}
	for _, d := range streak.WeeksDetail {
		if d.GoalMet {
			t.Errorf("GoalMet should be false when no goals exist, got true for week %s", d.WeekMonday)
		}
	}
}

func TestFindGoalForWeek(t *testing.T) {
	now := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)
	thisWeek := repository.WeekStart(now)
	prevWeek1 := thisWeek.AddDate(0, 0, -7)
	prevWeek3 := thisWeek.AddDate(0, 0, -21)

	goals := []model.Goal{
		{ID: 1, UserID: 1, WeeklySessions: 5, EffectiveMonday: prevWeek3},
		{ID: 2, UserID: 1, WeeklySessions: 3, EffectiveMonday: prevWeek1},
	}

	t.Run("returns latest goal for given week", func(t *testing.T) {
		goal, found := findGoalForWeek(goals, prevWeek1)
		if !found {
			t.Fatal("expected goal to be found")
		}
		if goal.WeeklySessions != 3 {
			t.Errorf("WeeklySessions = %v, want 3 (should use newer goal)", goal.WeeklySessions)
		}
	})

	t.Run("falls back to older goal", func(t *testing.T) {
		goal, found := findGoalForWeek(goals, thisWeek.AddDate(0, 0, -14))
		if !found {
			t.Fatal("expected goal to be found")
		}
		if goal.WeeklySessions != 5 {
			t.Errorf("WeeklySessions = %v, want 5 (should use older goal)", goal.WeeklySessions)
		}
	})

	t.Run("returns false when no goal exists", func(t *testing.T) {
		_, found := findGoalForWeek(goals, thisWeek.AddDate(0, 0, -28))
		if found {
			t.Error("expected no goal to be found for week before any goal")
		}
	})

	t.Run("current week uses most recent goal", func(t *testing.T) {
		goal, found := findGoalForWeek(goals, thisWeek)
		if !found {
			t.Fatal("expected goal to be found")
		}
		if goal.WeeklySessions != 3 {
			t.Errorf("WeeklySessions = %v, want 3 (should use most recent goal)", goal.WeeklySessions)
		}
	})
}

func TestPercent(t *testing.T) {
	tests := []struct {
		name     string
		done     float64
		target   float64
		expected float64
	}{
		{"exact match", 100, 100, 100},
		{"over achieved", 200, 100, 100},
		{"half done", 50, 100, 50},
		{"zero target", 50, 0, 0},
		{"negative target", 50, -10, 0},
		{"zero done positive target", 0, 100, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := percent(tt.done, tt.target)
			if result != tt.expected {
				t.Errorf("percent(%v, %v) = %v, want %v", tt.done, tt.target, result, tt.expected)
			}
		})
	}
}

func TestGroupRecordsByWeek(t *testing.T) {
	monday := time.Date(2026, 6, 8, 0, 0, 0, 0, time.UTC)
	tuesday := time.Date(2026, 6, 9, 10, 0, 0, 0, time.UTC)
	sunday := time.Date(2026, 6, 14, 23, 0, 0, 0, time.UTC)
	nextMonday := time.Date(2026, 6, 15, 0, 0, 0, 0, time.UTC)

	records := []model.WorkoutRecord{
		{OccurredAt: monday},
		{OccurredAt: tuesday},
		{OccurredAt: sunday},
		{OccurredAt: nextMonday},
	}

	result := groupRecordsByWeek(records)

	if len(result) != 2 {
		t.Errorf("expected 2 weeks, got %d", len(result))
	}

	week1 := repository.WeekStart(monday)
	week2 := repository.WeekStart(nextMonday)

	if len(result[week1]) != 3 {
		t.Errorf("week1 should have 3 records, got %d", len(result[week1]))
	}
	if len(result[week2]) != 1 {
		t.Errorf("week2 should have 1 record, got %d", len(result[week2]))
	}
}
