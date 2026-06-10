package handler

import (
	"encoding/json"
	"fitnessapi/internal/model"
	"fitnessapi/internal/repository"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type mockGoalService struct {
	goals []model.Goal
	err   error
}

func (m *mockGoalService) Save(goal *model.Goal) error          { return m.err }
func (m *mockGoalService) Latest(userID uint) (model.Goal, error) {
	if len(m.goals) == 0 {
		return model.Goal{}, m.err
	}
	return m.goals[len(m.goals)-1], m.err
}
func (m *mockGoalService) List(userID uint) ([]model.Goal, error) { return m.goals, m.err }

type mockRecordService struct {
	records []model.WorkoutRecord
	err     error
}

func (m *mockRecordService) Create(record *model.WorkoutRecord) error          { return m.err }
func (m *mockRecordService) List(userID uint, start, end *time.Time) ([]model.WorkoutRecord, error) {
	return m.records, m.err
}
func (m *mockRecordService) Update(record *model.WorkoutRecord) error { return m.err }
func (m *mockRecordService) Delete(id uint) error                     { return m.err }
func (m *mockRecordService) Find(id uint) (model.WorkoutRecord, error) {
	return model.WorkoutRecord{}, m.err
}

type mockUserRepository struct {
	users []model.User
	err   error
}

func (m *mockUserRepository) List() ([]model.User, error)      { return m.users, m.err }
func (m *mockUserRepository) Find(id uint) (model.User, error) { return model.User{}, m.err }

type userDataItem struct {
	Goals   []model.Goal
	Records []model.WorkoutRecord
}

type perUserRecordSvc struct {
	data map[uint]userDataItem
}

func (s *perUserRecordSvc) Create(record *model.WorkoutRecord) error { return nil }
func (s *perUserRecordSvc) List(userID uint, start, end *time.Time) ([]model.WorkoutRecord, error) {
	if d, ok := s.data[userID]; ok {
		return d.Records, nil
	}
	return nil, nil
}
func (s *perUserRecordSvc) Update(record *model.WorkoutRecord) error  { return nil }
func (s *perUserRecordSvc) Delete(id uint) error                      { return nil }
func (s *perUserRecordSvc) Find(id uint) (model.WorkoutRecord, error) { return model.WorkoutRecord{}, nil }

type perUserGoalSvc struct {
	data map[uint]userDataItem
}

func (s *perUserGoalSvc) Save(goal *model.Goal) error { return nil }
func (s *perUserGoalSvc) Latest(userID uint) (model.Goal, error) {
	if d, ok := s.data[userID]; ok && len(d.Goals) > 0 {
		return d.Goals[len(d.Goals)-1], nil
	}
	return model.Goal{}, nil
}
func (s *perUserGoalSvc) List(userID uint) ([]model.Goal, error) {
	if d, ok := s.data[userID]; ok {
		return d.Goals, nil
	}
	return nil, nil
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestGoalProgress_ConsecutiveWeeks(t *testing.T) {
	thisWeek := repository.WeekStart(time.Now())
	prevWeek1 := thisWeek.AddDate(0, 0, -7)
	prevWeek2 := thisWeek.AddDate(0, 0, -14)
	prevWeek3 := thisWeek.AddDate(0, 0, -21)

	goals := []model.Goal{
		{ID: 1, UserID: 1, WeeklySessions: 3, WeeklyMinutes: 0, WeeklyCalories: 0, EffectiveMonday: prevWeek3},
	}

	t.Run("连续3周达标返回consecutive_weeks=3", func(t *testing.T) {
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

		h := NewGoalHandler(&mockGoalService{goals: goals}, &mockRecordService{records: records})
		r := setupRouter()
		r.GET("/goals/progress", h.Progress)

		req := httptest.NewRequest(http.MethodGet, "/goals/progress?user_id=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		var resp map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to parse response: %v", err)
		}

		streak, ok := resp["streak"].(map[string]interface{})
		if !ok {
			t.Fatal("response missing streak object")
		}

		cw, _ := streak["consecutive_weeks"].(float64)
		if int(cw) != 3 {
			t.Errorf("consecutive_weeks = %v, want 3", cw)
		}

		met, _ := streak["current_week_met"].(bool)
		if !met {
			t.Errorf("current_week_met = false, want true")
		}

		details, ok := streak["weeks_detail"].([]interface{})
		if !ok {
			t.Fatal("weeks_detail not found or not array")
		}
		if len(details) != 4 {
			t.Errorf("weeks_detail length = %d, want 4", len(details))
		}
	})

	t.Run("本周未达标返回consecutive_weeks=0", func(t *testing.T) {
		onlyPrevRecords := []model.WorkoutRecord{
			{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: prevWeek1.AddDate(0, 0, 1)},
			{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: prevWeek1.AddDate(0, 0, 2)},
			{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: prevWeek1.AddDate(0, 0, 3)},
		}

		h := NewGoalHandler(&mockGoalService{goals: goals}, &mockRecordService{records: onlyPrevRecords})
		r := setupRouter()
		r.GET("/goals/progress", h.Progress)

		req := httptest.NewRequest(http.MethodGet, "/goals/progress?user_id=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		streak := resp["streak"].(map[string]interface{})
		cw, _ := streak["consecutive_weeks"].(float64)
		if int(cw) != 0 {
			t.Errorf("consecutive_weeks = %v, want 0 (this week not met)", cw)
		}
		met, _ := streak["current_week_met"].(bool)
		if met {
			t.Errorf("current_week_met = true, want false")
		}
	})

	t.Run("只配置次数目标且连续2周达标", func(t *testing.T) {
		partialGoals := []model.Goal{
			{ID: 1, UserID: 1, WeeklySessions: 2, WeeklyMinutes: 0, WeeklyCalories: 0, EffectiveMonday: prevWeek1},
		}
		partialRecords := []model.WorkoutRecord{
			{UserID: 1, DurationMin: 30, OccurredAt: thisWeek.AddDate(0, 0, 1)},
			{UserID: 1, DurationMin: 30, OccurredAt: thisWeek.AddDate(0, 0, 2)},
			{UserID: 1, DurationMin: 30, OccurredAt: prevWeek1.AddDate(0, 0, 1)},
			{UserID: 1, DurationMin: 30, OccurredAt: prevWeek1.AddDate(0, 0, 2)},
		}

		h := NewGoalHandler(&mockGoalService{goals: partialGoals}, &mockRecordService{records: partialRecords})
		r := setupRouter()
		r.GET("/goals/progress", h.Progress)

		req := httptest.NewRequest(http.MethodGet, "/goals/progress?user_id=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		streak := resp["streak"].(map[string]interface{})
		cw, _ := streak["consecutive_weeks"].(float64)
		if int(cw) != 2 {
			t.Errorf("consecutive_weeks = %v, want 2 (only sessions target, both weeks met)", cw)
		}
	})

	t.Run("返回结构包含goal和progress字段", func(t *testing.T) {
		records := []model.WorkoutRecord{
			{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: thisWeek.AddDate(0, 0, 1)},
			{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: thisWeek.AddDate(0, 0, 2)},
			{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: thisWeek.AddDate(0, 0, 3)},
		}

		h := NewGoalHandler(&mockGoalService{goals: goals}, &mockRecordService{records: records})
		r := setupRouter()
		r.GET("/goals/progress", h.Progress)

		req := httptest.NewRequest(http.MethodGet, "/goals/progress?user_id=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		if _, ok := resp["goal"]; !ok {
			t.Error("response missing 'goal' field")
		}
		if _, ok := resp["progress"]; !ok {
			t.Error("response missing 'progress' field")
		}
		if _, ok := resp["streak"]; !ok {
			t.Error("response missing 'streak' field")
		}

		progress, _ := resp["progress"].(map[string]interface{})
		sp, _ := progress["sessions_percent"].(float64)
		if sp != 100 {
			t.Errorf("sessions_percent = %v, want 100", sp)
		}
	})

	t.Run("配置次数+卡路里但卡路里未达标", func(t *testing.T) {
		mixedGoals := []model.Goal{
			{ID: 1, UserID: 1, WeeklySessions: 3, WeeklyMinutes: 0, WeeklyCalories: 1500, EffectiveMonday: prevWeek1},
		}
		mixedRecords := []model.WorkoutRecord{
			{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: thisWeek.AddDate(0, 0, 1)},
			{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: thisWeek.AddDate(0, 0, 2)},
			{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: thisWeek.AddDate(0, 0, 3)},
			{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: prevWeek1.AddDate(0, 0, 1)},
			{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: prevWeek1.AddDate(0, 0, 2)},
			{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: prevWeek1.AddDate(0, 0, 3)},
		}

		h := NewGoalHandler(&mockGoalService{goals: mixedGoals}, &mockRecordService{records: mixedRecords})
		r := setupRouter()
		r.GET("/goals/progress", h.Progress)

		req := httptest.NewRequest(http.MethodGet, "/goals/progress?user_id=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		streak := resp["streak"].(map[string]interface{})
		cw, _ := streak["consecutive_weeks"].(float64)
		if int(cw) != 0 {
			t.Errorf("consecutive_weeks = %v, want 0 (calories target 1500 not met, each week only 600)", cw)
		}
	})
}

func TestRankings_Streak(t *testing.T) {
	thisWeek := repository.WeekStart(time.Now())
	prevWeek1 := thisWeek.AddDate(0, 0, -7)
	prevWeek2 := thisWeek.AddDate(0, 0, -14)
	prevWeek3 := thisWeek.AddDate(0, 0, -21)

	users := []model.User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Charlie"},
	}

	userData := map[uint]userDataItem{
		1: {
			Goals: []model.Goal{
				{ID: 1, UserID: 1, WeeklySessions: 3, WeeklyMinutes: 0, WeeklyCalories: 0, EffectiveMonday: prevWeek3},
			},
			Records: []model.WorkoutRecord{
				{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: thisWeek.AddDate(0, 0, 1)},
				{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: thisWeek.AddDate(0, 0, 2)},
				{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: thisWeek.AddDate(0, 0, 3)},
				{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: prevWeek1.AddDate(0, 0, 1)},
				{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: prevWeek1.AddDate(0, 0, 2)},
				{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: prevWeek1.AddDate(0, 0, 3)},
				{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: prevWeek2.AddDate(0, 0, 1)},
				{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: prevWeek2.AddDate(0, 0, 2)},
				{UserID: 1, DurationMin: 30, Calories: 200, OccurredAt: prevWeek2.AddDate(0, 0, 3)},
			},
		},
		2: {
			Goals: []model.Goal{
				{ID: 2, UserID: 2, WeeklySessions: 3, WeeklyMinutes: 0, WeeklyCalories: 0, EffectiveMonday: prevWeek3},
			},
			Records: []model.WorkoutRecord{
				{UserID: 2, DurationMin: 30, Calories: 200, OccurredAt: thisWeek.AddDate(0, 0, 1)},
				{UserID: 2, DurationMin: 30, Calories: 200, OccurredAt: thisWeek.AddDate(0, 0, 2)},
				{UserID: 2, DurationMin: 30, Calories: 200, OccurredAt: thisWeek.AddDate(0, 0, 3)},
				{UserID: 2, DurationMin: 30, Calories: 200, OccurredAt: prevWeek1.AddDate(0, 0, 1)},
				{UserID: 2, DurationMin: 30, Calories: 200, OccurredAt: prevWeek1.AddDate(0, 0, 2)},
				{UserID: 2, DurationMin: 30, Calories: 200, OccurredAt: prevWeek1.AddDate(0, 0, 3)},
			},
		},
		3: {
			Goals: []model.Goal{
				{ID: 3, UserID: 3, WeeklySessions: 3, WeeklyMinutes: 0, WeeklyCalories: 0, EffectiveMonday: prevWeek3},
			},
			Records: []model.WorkoutRecord{
				{UserID: 3, DurationMin: 30, Calories: 200, OccurredAt: thisWeek.AddDate(0, 0, 1)},
				{UserID: 3, DurationMin: 30, Calories: 200, OccurredAt: thisWeek.AddDate(0, 0, 2)},
			},
		},
	}

	newRankHandler := func() RankHandler {
		return RankHandler{
			recordSvc: &perUserRecordSvc{data: userData},
			goalSvc:   &perUserGoalSvc{data: userData},
			userRepo:  &mockUserRepository{users: users},
		}
	}

	t.Run("按连续周数降序排列", func(t *testing.T) {
		h := newRankHandler()
		r := setupRouter()
		r.GET("/rankings", h.Rankings)

		req := httptest.NewRequest(http.MethodGet, "/rankings?scope=streak&user_id=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d; body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		if scope, _ := resp["scope"].(string); scope != "streak" {
			t.Errorf("scope = %v, want streak", scope)
		}

		ranking, ok := resp["ranking"].([]interface{})
		if !ok {
			t.Fatal("ranking not found or not array")
		}
		if len(ranking) != 3 {
			t.Fatalf("ranking length = %d, want 3", len(ranking))
		}

		first := ranking[0].(map[string]interface{})
		if first["consecutive_weeks"].(float64) != 3 {
			t.Errorf("first place consecutive_weeks = %v, want 3", first["consecutive_weeks"])
		}
		if first["name"] != "Alice" {
			t.Errorf("first place name = %v, want Alice", first["name"])
		}

		second := ranking[1].(map[string]interface{})
		if second["consecutive_weeks"].(float64) != 2 {
			t.Errorf("second place consecutive_weeks = %v, want 2", second["consecutive_weeks"])
		}
		if second["name"] != "Bob" {
			t.Errorf("second place name = %v, want Bob", second["name"])
		}

		third := ranking[2].(map[string]interface{})
		if third["consecutive_weeks"].(float64) != 0 {
			t.Errorf("third place consecutive_weeks = %v, want 0", third["consecutive_weeks"])
		}
	})

	t.Run("my_rank正确返回", func(t *testing.T) {
		h := newRankHandler()
		r := setupRouter()
		r.GET("/rankings", h.Rankings)

		req := httptest.NewRequest(http.MethodGet, "/rankings?scope=streak&user_id=2", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		myRank, _ := resp["my_rank"].(float64)
		if int(myRank) != 2 {
			t.Errorf("my_rank = %v, want 2 (Bob is second with 2 consecutive weeks)", myRank)
		}
	})

	t.Run("streak排行每条记录包含consecutive_weeks字段", func(t *testing.T) {
		h := newRankHandler()
		r := setupRouter()
		r.GET("/rankings", h.Rankings)

		req := httptest.NewRequest(http.MethodGet, "/rankings?scope=streak&user_id=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		ranking := resp["ranking"].([]interface{})
		requiredFields := []string{"user_id", "name", "duration", "distance", "calories", "consecutive_weeks"}
		for i, item := range ranking {
			entry := item.(map[string]interface{})
			for _, field := range requiredFields {
				if _, ok := entry[field]; !ok {
					t.Errorf("ranking[%d] missing '%s' field", i, field)
				}
			}
		}
	})

	t.Run("非streak scope不含consecutive_weeks", func(t *testing.T) {
		h := newRankHandler()
		r := setupRouter()
		r.GET("/rankings", h.Rankings)

		req := httptest.NewRequest(http.MethodGet, "/rankings?scope=week&user_id=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		ranking := resp["ranking"].([]interface{})
		for _, item := range ranking {
			entry := item.(map[string]interface{})
			if _, ok := entry["consecutive_weeks"]; ok {
				t.Errorf("week ranking entry should NOT have 'consecutive_weeks' field: %v", entry)
			}
		}
	})

	t.Run("返回统一格式包含scope和friend_circle", func(t *testing.T) {
		h := newRankHandler()
		r := setupRouter()
		r.GET("/rankings", h.Rankings)

		req := httptest.NewRequest(http.MethodGet, "/rankings?scope=streak&friend_circle=team-a&user_id=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		if scope, _ := resp["scope"].(string); scope != "streak" {
			t.Errorf("scope = %v, want streak", scope)
		}
		if fc, _ := resp["friend_circle"].(string); fc != "team-a" {
			t.Errorf("friend_circle = %v, want team-a", fc)
		}
		if _, ok := resp["my_rank"]; !ok {
			t.Error("response missing 'my_rank' field")
		}
		if _, ok := resp["ranking"]; !ok {
			t.Error("response missing 'ranking' field")
		}
	})
}
