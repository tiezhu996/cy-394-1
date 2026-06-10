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
