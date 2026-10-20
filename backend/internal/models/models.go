package models

import (
	"time"

	"gorm.io/datatypes"
)

type User struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"size:255;not null" json:"name"`
	Email         string    `gorm:"size:255;uniqueIndex;not null" json:"email"`
	PasswordHash  string    `gorm:"size:255;not null" json:"-"`
	Rating        int       `gorm:"default:1200;not null" json:"rating"`
	MatchesPlayed int       `gorm:"default:0;not null" json:"matches_played"`
	Wins          int       `gorm:"default:0;not null" json:"wins"`
	Losses        int       `gorm:"default:0;not null" json:"losses"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Problem struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Slug          string         `gorm:"size:128;uniqueIndex;not null" json:"slug"`
	Title         string         `gorm:"size:255;not null" json:"title"`
	Difficulty    string         `gorm:"size:32;not null" json:"difficulty"` // easy, medium, hard
	Description   string         `gorm:"type:text;not null" json:"description"`
	StarterCode   datatypes.JSON `gorm:"type:text;not null" json:"starter_code"`
	TimeLimitMs   int            `gorm:"default:2000;not null" json:"time_limit_ms"`
	MemoryLimitKb int            `gorm:"default:131072;not null" json:"memory_limit_kb"`
	TestCases     []TestCase     `gorm:"foreignKey:ProblemID" json:"test_cases,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type TestCase struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ProblemID      uint      `gorm:"not null;index" json:"problem_id"`
	Input          string    `gorm:"type:text;not null" json:"input"`
	ExpectedOutput string    `gorm:"type:text;not null" json:"expected_output"`
	IsHidden       bool      `gorm:"default:true;not null" json:"is_hidden"`
	SequenceOrder  int       `gorm:"default:1;not null" json:"sequence_order"`
	CreatedAt      time.Time `json:"created_at"`
}

type Match struct {
	ID                 uint               `gorm:"primaryKey" json:"id"`
	PublicID           string             `gorm:"size:32;uniqueIndex;not null" json:"public_id"` // ULID
	Mode               string             `gorm:"size:32;not null" json:"mode"`                  // ranked, casual, tournament
	Status             string             `gorm:"size:32;default:'created';not null" json:"status"`
	ProblemID          uint               `gorm:"not null" json:"problem_id"`
	Problem            Problem            `gorm:"foreignKey:ProblemID" json:"problem"`
	MaxDurationSeconds int                `gorm:"default:1800;not null" json:"max_duration_seconds"`
	CountdownSeconds   int                `gorm:"default:10;not null" json:"countdown_seconds"`
	StartedAt          *time.Time         `json:"started_at"`
	EndedAt            *time.Time         `json:"ended_at"`
	SettledAt          *time.Time         `json:"settled_at"`
	WinnerID           *uint              `json:"winner_id"`
	Version            int                `gorm:"default:1;not null" json:"version"`
	Metadata           datatypes.JSON     `gorm:"type:text" json:"metadata"`
	Participants       []MatchParticipant `gorm:"foreignKey:MatchID" json:"participants"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
}

type MatchParticipant struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	MatchID         uint      `gorm:"not null;index" json:"match_id"`
	UserID          uint      `gorm:"not null;index" json:"user_id"`
	User            User      `gorm:"foreignKey:UserID" json:"user"`
	Seat            int       `gorm:"not null" json:"seat"` // 1 or 2
	Result          string    `gorm:"size:32;default:'pending';not null" json:"result"`
	RatingBefore    int       `json:"rating_before"`
	RatingAfter     int       `json:"rating_after"`
	RatingDelta     int       `json:"rating_delta"`
	SubmissionCount int       `gorm:"default:0;not null" json:"submission_count"`
	PassedTests     int       `gorm:"default:0;not null" json:"passed_tests"`
	TotalTests      int       `gorm:"default:0;not null" json:"total_tests"`
	IsConnected     bool      `gorm:"default:false;not null" json:"is_connected"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Submission struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	MatchID         *uint      `gorm:"index" json:"match_id,omitempty"`
	UserID          uint       `gorm:"not null;index" json:"user_id"`
	User            User       `gorm:"foreignKey:UserID" json:"user"`
	ProblemID       uint       `gorm:"not null;index" json:"problem_id"`
	Language        string     `gorm:"size:32;not null" json:"language"`
	Code            string     `gorm:"type:text;not null" json:"code"`
	Verdict         string     `gorm:"size:32;not null" json:"verdict"` // AC, WA, TLE, MLE, CE, RE, PENDING
	PassedTestCases int        `gorm:"default:0;not null" json:"passed_test_cases"`
	TotalTestCases  int        `gorm:"default:0;not null" json:"total_test_cases"`
	RuntimeMs       int        `json:"runtime_ms"`
	MemoryKb        int        `json:"memory_kb"`
	CompileOutput   string     `gorm:"type:text" json:"compile_output,omitempty"`
	SubmittedAt     time.Time  `json:"submitted_at"`
}

type MatchEvent struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	MatchID       uint           `gorm:"not null;index" json:"match_id"`
	ParticipantID *uint          `json:"participant_id"`
	EventType     string         `gorm:"size:64;not null" json:"event_type"`
	EventSequence int64          `gorm:"not null" json:"event_sequence"`
	Payload       datatypes.JSON `gorm:"type:text;not null" json:"payload"`
	OccurredAt    time.Time      `json:"occurred_at"`
}
