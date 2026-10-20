package services

import (
	"context"
	"crypto/rand"
	"log"
	"math"
	"sync"
	"time"

	"devcode/backend/internal/models"
	"devcode/backend/internal/ws"

	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type QueuePlayer struct {
	UserID     uint      `json:"user_id"`
	Rating     int       `json:"rating"`
	Mode       string    `json:"mode"`
	EnqueuedAt time.Time `json:"enqueued_at"`
}

type MatchmakingEngine struct {
	db          *gorm.DB
	rdb         *redis.Client
	hasRedis    bool
	hub         *ws.Hub
	queue       map[uint]*QueuePlayer
	mu          sync.Mutex
	stopChan    chan struct{}
}

func NewMatchmakingEngine(db *gorm.DB, redisAddr, redisPassword string, hub *ws.Hub) *MatchmakingEngine {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
	})

	hasRedis := false
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err == nil {
		hasRedis = true
		log.Println("✅ Matchmaking connected to Redis cache/queue")
	} else {
		log.Println("⚠️ Redis not reachable, using in-memory matchmaking queue")
	}

	engine := &MatchmakingEngine{
		db:       db,
		rdb:      rdb,
		hasRedis: hasRedis,
		hub:      hub,
		queue:    make(map[uint]*QueuePlayer),
		stopChan: make(chan struct{}),
	}

	go engine.runWorker()
	return engine
}

func (m *MatchmakingEngine) Enqueue(player *QueuePlayer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	player.EnqueuedAt = time.Now()
	m.queue[player.UserID] = player
	log.Printf("Matchmaking: User %d enqueued for %s (Rating: %d)", player.UserID, player.Mode, player.Rating)
}

func (m *MatchmakingEngine) Dequeue(userID uint) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.queue, userID)
	log.Printf("Matchmaking: User %d removed from queue", userID)
}

func (m *MatchmakingEngine) IsQueued(userID uint) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, exists := m.queue[userID]
	return exists
}

func (m *MatchmakingEngine) runWorker() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopChan:
			return
		case <-ticker.C:
			m.sweepAndMatch()
		}
	}
}

func (m *MatchmakingEngine) sweepAndMatch() {
	m.mu.Lock()
	players := make([]*QueuePlayer, 0, len(m.queue))
	for _, p := range m.queue {
		players = append(players, p)
	}
	m.mu.Unlock()

	if len(players) < 2 {
		return
	}

	matchedUsers := make(map[uint]bool)

	for i := 0; i < len(players); i++ {
		p1 := players[i]
		if matchedUsers[p1.UserID] {
			continue
		}

		timeInQueue := time.Since(p1.EnqueuedAt).Seconds()
		// Expand window based on wait time: +-50 initial, +50 every 10s up to +-300
		allowedDelta := 50.0 + math.Floor(timeInQueue/10.0)*50.0
		if allowedDelta > 300.0 {
			allowedDelta = 300.0
		}

		var bestMatch *QueuePlayer
		bestDiff := 99999.0

		for j := i + 1; j < len(players); j++ {
			p2 := players[j]
			if matchedUsers[p2.UserID] || p1.Mode != p2.Mode {
				continue
			}

			diff := math.Abs(float64(p1.Rating - p2.Rating))
			if diff <= allowedDelta && diff < bestDiff {
				bestDiff = diff
				bestMatch = p2
			}
		}

		if bestMatch != nil {
			matchedUsers[p1.UserID] = true
			matchedUsers[bestMatch.UserID] = true
			m.Dequeue(p1.UserID)
			m.Dequeue(bestMatch.UserID)

			go m.createMatch(p1, bestMatch)
		}
	}
}

func (m *MatchmakingEngine) createMatch(p1, p2 *QueuePlayer) {
	// Pick random problem
	var problems []models.Problem
	if err := m.db.Find(&problems).Error; err != nil || len(problems) == 0 {
		log.Println("⚠️ No problems found to create match")
		return
	}

	// Pick random problem
	randomProblem := problems[time.Now().UnixNano()%int64(len(problems))]

	// Generate ULID
	entropy := rand.Reader
	roomULID := ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()

	now := time.Now()
	match := models.Match{
		PublicID:           roomULID,
		Mode:               p1.Mode,
		Status:             "countdown",
		ProblemID:          randomProblem.ID,
		MaxDurationSeconds: 1800,
		CountdownSeconds:   10,
		StartedAt:          nil,
		CreatedAt:          now,
	}

	if err := m.db.Create(&match).Error; err != nil {
		log.Printf("Failed to create match: %v", err)
		return
	}

	participant1 := models.MatchParticipant{
		MatchID:      match.ID,
		UserID:       p1.UserID,
		Seat:         1,
		Result:       "pending",
		RatingBefore: p1.Rating,
	}
	participant2 := models.MatchParticipant{
		MatchID:      match.ID,
		UserID:       p2.UserID,
		Seat:         2,
		Result:       "pending",
		RatingBefore: p2.Rating,
	}

	m.db.Create(&participant1)
	m.db.Create(&participant2)

	// Fetch full match with relations
	var fullMatch models.Match
	m.db.Preload("Problem").Preload("Problem.TestCases").Preload("Participants.User").First(&fullMatch, match.ID)

	// Notify both users via private channel
	payload := map[string]interface{}{
		"match_id":    fullMatch.PublicID,
		"mode":        fullMatch.Mode,
		"problem_id":  fullMatch.ProblemID,
		"problem":     fullMatch.Problem,
		"countdown":   10,
	}

	m.hub.SendToUser(p1.UserID, "match.created", payload)
	m.hub.SendToUser(p2.UserID, "match.created", payload)
	log.Printf("🎉 Match created: %s between User %d and User %d", roomULID, p1.UserID, p2.UserID)

	// Start pre-match countdown
	go m.handleMatchCountdown(fullMatch.PublicID, fullMatch.ID)
}

func (m *MatchmakingEngine) handleMatchCountdown(publicID string, matchID uint) {
	for tick := 10; tick >= 0; tick-- {
		time.Sleep(1 * time.Second)
		m.hub.BroadcastToRoom(publicID, "match.countdown_tick", map[string]int{"countdown": tick})
	}

	now := time.Now()
	m.db.Model(&models.Match{}).Where("id = ?", matchID).Updates(map[string]interface{}{
		"status":     "live",
		"started_at": &now,
	})

	m.hub.BroadcastToRoom(publicID, "match.started", map[string]interface{}{
		"status":     "live",
		"started_at": now,
	})
	log.Printf("🔥 Match %s is now LIVE!", publicID)
}

func (m *MatchmakingEngine) ForfeitMatch(matchPublicID string, forfeitUserID uint) error {
	var match models.Match
	if err := m.db.Preload("Participants.User").Where("public_id = ?", matchPublicID).First(&match).Error; err != nil {
		return err
	}

	if match.Status == "finished" || match.Status == "settled" {
		return nil
	}

	var winnerParticipant *models.MatchParticipant
	var loserParticipant *models.MatchParticipant

	for i := range match.Participants {
		if match.Participants[i].UserID == forfeitUserID {
			loserParticipant = &match.Participants[i]
		} else {
			winnerParticipant = &match.Participants[i]
		}
	}

	if winnerParticipant != nil && loserParticipant != nil {
		return m.SettleMatch(&match, winnerParticipant, loserParticipant, false)
	}

	return nil
}

func (m *MatchmakingEngine) SettleMatch(match *models.Match, winner, loser *models.MatchParticipant, isDraw bool) error {
	now := time.Now()
	scoreA := 1.0
	if isDraw {
		scoreA = 0.5
	}

	// Calculate ELO
	elo := CalculateELO(
		winner.RatingBefore,
		loser.RatingBefore,
		winner.User.MatchesPlayed,
		loser.User.MatchesPlayed,
		scoreA,
	)

	winner.Result = "won"
	if isDraw {
		winner.Result = "draw"
	}
	winner.RatingAfter = elo.PlayerANewRating
	winner.RatingDelta = elo.PlayerADelta

	loser.Result = "lost"
	if isDraw {
		loser.Result = "draw"
	}
	loser.RatingAfter = elo.PlayerBNewRating
	loser.RatingDelta = elo.PlayerBDelta

	// Update participants
	m.db.Save(winner)
	m.db.Save(loser)

	// Update user stats
	m.db.Model(&models.User{}).Where("id = ?", winner.UserID).Updates(map[string]interface{}{
		"rating":         winner.RatingAfter,
		"matches_played": gorm.Expr("matches_played + 1"),
		"wins":           gorm.Expr("wins + ?", boolToInt(!isDraw)),
	})
	m.db.Model(&models.User{}).Where("id = ?", loser.UserID).Updates(map[string]interface{}{
		"rating":         loser.RatingAfter,
		"matches_played": gorm.Expr("matches_played + 1"),
		"losses":         gorm.Expr("losses + ?", boolToInt(!isDraw)),
	})

	var winnerID *uint
	if !isDraw {
		winnerID = &winner.UserID
	}

	match.Status = "finished"
	match.EndedAt = &now
	match.SettledAt = &now
	match.WinnerID = winnerID
	m.db.Save(match)

	// Broadcast match finish to room
	m.hub.BroadcastToRoom(match.PublicID, "match.finished", map[string]interface{}{
		"match_id":  match.PublicID,
		"winner_id": winnerID,
		"is_draw":   isDraw,
		"participants": []map[string]interface{}{
			{
				"user_id":      winner.UserID,
				"name":         winner.User.Name,
				"result":       winner.Result,
				"rating_after": winner.RatingAfter,
				"rating_delta": winner.RatingDelta,
			},
			{
				"user_id":      loser.UserID,
				"name":         loser.User.Name,
				"result":       loser.Result,
				"rating_after": loser.RatingAfter,
				"rating_delta": loser.RatingDelta,
			},
		},
	})

	log.Printf("🏆 Match %s settled. Winner: %v", match.PublicID, winnerID)
	return nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
