package services

import (
	"math"
)

type ELOResult struct {
	PlayerANewRating int
	PlayerADelta     int
	PlayerBNewRating int
	PlayerBDelta     int
}

// CalculateELO computes rating changes for Player A and Player B.
// scoreA: 1.0 = Player A wins, 0.5 = Draw, 0.0 = Player B wins.
func CalculateELO(ratingA, ratingB, matchesA, matchesB int, scoreA float64) ELOResult {
	scoreB := 1.0 - scoreA

	// Determine K factor
	kA := 32.0
	if matchesA >= 30 {
		kA = 16.0
	}

	kB := 32.0
	if matchesB >= 30 {
		kB = 16.0
	}

	// Expected scores
	expectedA := 1.0 / (1.0 + math.Pow(10.0, float64(ratingB-ratingA)/400.0))
	expectedB := 1.0 / (1.0 + math.Pow(10.0, float64(ratingA-ratingB)/400.0))

	// Rating deltas
	deltaA := int(math.Round(kA * (scoreA - expectedA)))
	deltaB := int(math.Round(kB * (scoreB - expectedB)))

	newRatingA := ratingA + deltaA
	if newRatingA < 100 {
		newRatingA = 100
	}

	newRatingB := ratingB + deltaB
	if newRatingB < 100 {
		newRatingB = 100
	}

	return ELOResult{
		PlayerANewRating: newRatingA,
		PlayerADelta:     deltaA,
		PlayerBNewRating: newRatingB,
		PlayerBDelta:     deltaB,
	}
}
