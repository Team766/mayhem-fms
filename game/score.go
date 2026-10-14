// Copyright 2023 Team 254. All Rights Reserved.
// Author: pat@patfairbank.com (Patrick Fairbank)
//
// Model representing the instantaneous score of a match.

package game

import "fmt"

type Score struct {
	AutoFloor           int
	AutoFirst           int
	AutoTop             int
	TeleopFloor         int
	TeleopFirst         int
	TeleopTop           int
	TeleopStacked       int
	Crown               CrownPlacement
	LeaveStatuses       [3]bool
	AutoBalanceStatuses [3]bool
	EndgameStatuses     [3]EndgameStatus
	Toss                bool
	Fouls               []Foul
	PlayoffDq           bool
}

// Point values and ranking point thresholds, from the game spec.
const (
	leavePoints            = 4
	autoBalancePoints      = 12
	autoFloorPoints        = 4
	autoFirstShelfPoints   = 8
	autoTopShelfPoints     = 12
	teleopFloorPoints      = 2
	teleopFirstShelfPoints = 5
	teleopTopShelfPoints   = 10
	teleopStackedPoints    = 8
	crownMultiplier        = 2 // The crown is worth twice an ordinary treasure in the same place.
	parkPoints             = 2
	balancePoints          = 12
	tossPoints             = 2

	// AutonRpThreshold is the auto points an alliance needs for the Auton ranking point.
	AutonRpThreshold = 20
	// ScoringRpThreshold is the teleop treasures on the first shelf, top shelf or stacked needed for the Scoring
	// ranking point.
	ScoringRpThreshold = 12
)

// Represents where on the field the dragon's crown ended the match, if the alliance scored it. The treasure counters
// never include the crown.
type CrownPlacement int

const (
	CrownNone CrownPlacement = iota
	CrownAutoFloor
	CrownAutoFirst
	CrownAutoTop
	CrownTeleopFloor
	CrownTeleopFirst
	CrownTeleopTop
	CrownTeleopStacked
)

// Returns the crown's points, which are twice the value of an ordinary treasure in the same place.
func (placement CrownPlacement) PointValue() int {
	switch placement {
	case CrownAutoFloor:
		return crownMultiplier * autoFloorPoints
	case CrownAutoFirst:
		return crownMultiplier * autoFirstShelfPoints
	case CrownAutoTop:
		return crownMultiplier * autoTopShelfPoints
	case CrownTeleopFloor:
		return crownMultiplier * teleopFloorPoints
	case CrownTeleopFirst:
		return crownMultiplier * teleopFirstShelfPoints
	case CrownTeleopTop:
		return crownMultiplier * teleopTopShelfPoints
	case CrownTeleopStacked:
		return crownMultiplier * teleopStackedPoints
	default:
		return 0
	}
}

// Returns true if the crown was placed during the autonomous period, meaning that its points belong to the auto
// treasure points and count toward the Auton ranking point.
func (placement CrownPlacement) IsAuto() bool {
	return placement == CrownAutoFloor || placement == CrownAutoFirst || placement == CrownAutoTop
}

// CrownSpot is one place the crown can end up, as the screens show it.
type CrownSpot struct {
	Placement CrownPlacement
	Id        string // The id the scoring panel sends, e.g. "teleop_top".
	Phase     string // "Auto" or "Teleop".
	Spot      string // "Floor", "First Shelf", "Top Shelf" or "Stacked".
}

// CrownSpots lists every place the crown can end up, in display order. The scoring panel, the edit-result form and
// the referee panel read it instead of keeping their own copies.
var CrownSpots = []CrownSpot{
	{CrownAutoFloor, "auto_floor", "Auto", "Floor"},
	{CrownAutoFirst, "auto_first", "Auto", "First Shelf"},
	{CrownAutoTop, "auto_top", "Auto", "Top Shelf"},
	{CrownTeleopFloor, "teleop_floor", "Teleop", "Floor"},
	{CrownTeleopFirst, "teleop_first", "Teleop", "First Shelf"},
	{CrownTeleopTop, "teleop_top", "Teleop", "Top Shelf"},
	{CrownTeleopStacked, "teleop_stacked", "Teleop", "Stacked"},
}

// Returns the name of the placement, for display, e.g. "Teleop Top Shelf", or "None".
func (placement CrownPlacement) String() string {
	for _, spot := range CrownSpots {
		if spot.Placement == placement {
			return spot.Phase + " " + spot.Spot
		}
	}
	return "None"
}

// CrownLabels returns the display name of every placement, indexed by its value, with "-" for no crown. The referee
// panel uses it to label the live crown placement.
func CrownLabels() []string {
	labels := []string{"-"}
	for _, spot := range CrownSpots {
		labels = append(labels, spot.Placement.String())
	}
	return labels
}

// Represents where a robot ended the match: nowhere in particular, parked in its own safe house, or balanced on its own
// mountain top. Park and balance are mutually exclusive.
type EndgameStatus int

const (
	EndgameNone EndgameStatus = iota
	EndgamePark
	EndgameBalance
)

// Summarize calculates and returns the summary fields used for ranking and display.
func (score *Score) Summarize(opponentScore *Score) *ScoreSummary {
	summary := new(ScoreSummary)
	summary.PlayoffDq = score.PlayoffDq

	// Leave the score at zero if the alliance was disqualified.
	if score.PlayoffDq {
		return summary
	}

	// Calculate autonomous period points. Leaving the safe house and balancing on the mountain top are independent
	// per-robot awards; a robot that does both earns both.
	for _, leave := range score.LeaveStatuses {
		if leave {
			summary.LeavePoints += leavePoints
		}
	}
	for _, balanced := range score.AutoBalanceStatuses {
		if balanced {
			summary.AutoBalancePoints += autoBalancePoints
		}
	}
	summary.AutoTreasurePoints = autoFloorPoints*score.AutoFloor + autoFirstShelfPoints*score.AutoFirst +
		autoTopShelfPoints*score.AutoTop

	// Calculate teleoperated period points. The teleop counters also cover the endgame period, which has no placement
	// values of its own.
	summary.TeleopTreasurePoints =
		teleopFloorPoints*score.TeleopFloor + teleopFirstShelfPoints*score.TeleopFirst +
			teleopTopShelfPoints*score.TeleopTop + teleopStackedPoints*score.TeleopStacked

	// The dragon's crown is worth twice an ordinary treasure in the same place. The treasure counters never include it,
	// so its points are added to the period in which it was placed.
	crownPoints := score.Crown.PointValue()
	if score.Crown.IsAuto() {
		summary.AutoTreasurePoints += crownPoints
	} else {
		// CrownNone is worth zero, so adding it here is harmless.
		summary.TeleopTreasurePoints += crownPoints
	}

	summary.AutonPoints = summary.LeavePoints + summary.AutoBalancePoints + summary.AutoTreasurePoints

	// Calculate endgame points, which are also per robot.
	numRobotsBalanced := 0
	for _, status := range score.EndgameStatuses {
		switch status {
		case EndgamePark:
			summary.EndgamePoints += parkPoints
		case EndgameBalance:
			summary.EndgamePoints += balancePoints
			numRobotsBalanced++
		default:
		}
	}

	// The toss scores once per alliance, whatever the alliance size.
	if score.Toss {
		summary.TossPoints = tossPoints
	}

	summary.MatchPoints =
		summary.AutonPoints + summary.TeleopTreasurePoints + summary.EndgamePoints + summary.TossPoints

	// Count treasures for the live displays and for the Scoring ranking point. The crown counts as one treasure.
	summary.TreasureCount = score.AutoFloor + score.AutoFirst + score.AutoTop + score.TeleopFloor +
		score.TeleopFirst + score.TeleopTop + score.TeleopStacked
	summary.ShelfTreasureCount = score.TeleopFirst + score.TeleopTop + score.TeleopStacked
	if score.Crown != CrownNone {
		summary.TreasureCount++
	}
	if score.Crown == CrownTeleopFirst || score.Crown == CrownTeleopTop || score.Crown == CrownTeleopStacked {
		summary.ShelfTreasureCount++
	}
	summary.ShelfTreasureGoal = ScoringRpThreshold

	// Calculate penalty points.
	for _, foul := range opponentScore.Fouls {
		summary.FoulPoints += foul.PointValue()
		// Store the number of major fouls since it is used to break ties in playoffs.
		if foul.IsMajor {
			summary.NumOpponentMajorFouls++
		}
	}

	summary.Score = summary.MatchPoints + summary.FoulPoints

	// Auton ranking point: at least the threshold in autonomous points, or the opponent entered the alliance's safe
	// house or safe zone during auto (MA2603).
	summary.AutonRankingPoint = summary.AutonPoints >= AutonRpThreshold || opponentScore.hasFoulForRule("MA2603")

	// Scoring ranking point: at least the threshold in treasures placed during teleop on the first shelf, on the top
	// shelf, or stacked, counting the crown as one. There is no opponent-violation alternative.
	summary.ScoringRankingPoint = summary.ShelfTreasureCount >= ScoringRpThreshold

	// Endgame ranking point: at least one robot balanced on its mountain top, or the opponent contacted the alliance's
	// balance beam (MA2601) or a robot on its own beam (MA2602) during endgame. Parking does not count.
	summary.EndgameRankingPoint = numRobotsBalanced > 0 || opponentScore.hasFoulForRule("MA2601", "MA2602")

	// Add up the bonus ranking points.
	if summary.AutonRankingPoint {
		summary.BonusRankingPoints++
	}
	if summary.ScoringRankingPoint {
		summary.BonusRankingPoints++
	}
	if summary.EndgameRankingPoint {
		summary.BonusRankingPoints++
	}

	return summary
}

// Returns true if the alliance committed at least one foul for any of the given rule numbers. A foul entered without a
// rule never satisfies this.
func (score *Score) hasFoulForRule(ruleNumbers ...string) bool {
	for _, foul := range score.Fouls {
		rule := foul.Rule()
		if rule == nil {
			continue
		}
		for _, ruleNumber := range ruleNumbers {
			if rule.RuleNumber == ruleNumber {
				return true
			}
		}
	}
	return false
}

// Validate returns an error if the score holds a value no scorer could enter: a negative treasure count, an unknown
// crown placement or an unknown endgame status. The match review form uses it, since it can submit any value.
func (score *Score) Validate() error {
	counters := map[string]int{
		"AutoFloor": score.AutoFloor, "AutoFirst": score.AutoFirst, "AutoTop": score.AutoTop,
		"TeleopFloor": score.TeleopFloor, "TeleopFirst": score.TeleopFirst, "TeleopTop": score.TeleopTop,
		"TeleopStacked": score.TeleopStacked,
	}
	for name, value := range counters {
		if value < 0 {
			return fmt.Errorf("%s can't be negative (got %d)", name, value)
		}
	}
	if score.Crown < CrownNone || score.Crown > CrownTeleopStacked {
		return fmt.Errorf("invalid crown placement %d", score.Crown)
	}
	for i, status := range score.EndgameStatuses {
		if status < EndgameNone || status > EndgameBalance {
			return fmt.Errorf("invalid endgame status %d for robot %d", status, i+1)
		}
	}
	return nil
}

// Equals returns true if and only if all fields of the two scores are equal.
func (score *Score) Equals(other *Score) bool {
	if score.AutoFloor != other.AutoFloor ||
		score.AutoFirst != other.AutoFirst ||
		score.AutoTop != other.AutoTop ||
		score.TeleopFloor != other.TeleopFloor ||
		score.TeleopFirst != other.TeleopFirst ||
		score.TeleopTop != other.TeleopTop ||
		score.TeleopStacked != other.TeleopStacked ||
		score.Crown != other.Crown ||
		score.LeaveStatuses != other.LeaveStatuses ||
		score.AutoBalanceStatuses != other.AutoBalanceStatuses ||
		score.EndgameStatuses != other.EndgameStatuses ||
		score.Toss != other.Toss ||
		score.PlayoffDq != other.PlayoffDq ||
		len(score.Fouls) != len(other.Fouls) {
		return false
	}

	for i, foul := range score.Fouls {
		if foul != other.Fouls[i] {
			return false
		}
	}

	return true
}
