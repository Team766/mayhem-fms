// Copyright 2017 Team 254. All Rights Reserved.
// Author: pat@patfairbank.com (Patrick Fairbank)

package game

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Returns a list of major fouls, one for each rule ID given.
func majors(ruleIds ...int) []Foul {
	fouls := make([]Foul, len(ruleIds))
	for i, ruleId := range ruleIds {
		fouls[i] = Foul{i + 1, true, 0, ruleId}
	}
	return fouls
}

// Returns a list of minor fouls, one for each rule ID given.
func minors(ruleIds ...int) []Foul {
	fouls := majors(ruleIds...)
	for i := range fouls {
		fouls[i].IsMajor = false
	}
	return fouls
}

// Each row lists the expected primary fields of the summary; the sums, the foul-adjusted score and the bonus ranking
// point count are derived below. Every field of the summary is compared, so each row also pins down everything that
// must stay zero. The treasure counters never include the crown.
func TestScoreSummarize(t *testing.T) {
	type S = Score
	type W = ScoreSummary
	const (
		ma2601 = ruleIdMa2601
		ma2602 = ruleIdMa2602
		ma2603 = ruleIdMa2603
		ma2604 = ruleIdMa2604
		ma2606 = ruleIdMa2606
		ma2616 = ruleIdMa2616
	)
	all3 := [3]bool{true, true, true}
	park := EndgamePark
	balance := EndgameBalance
	testCases := []struct {
		name          string
		score         Score
		opponentFouls []Foul
		want          W
	}{
		{"nothing", S{}, nil, W{}},

		// Auto treasures count toward the auton points and the treasure count, but never the shelf count.
		{"auto floor x3", S{AutoFloor: 3}, nil, W{AutoTreasurePoints: 12, TreasureCount: 3}},
		{"auto first x2", S{AutoFirst: 2}, nil, W{AutoTreasurePoints: 16, TreasureCount: 2}},
		{"auto top x2", S{AutoTop: 2}, nil, W{AutoTreasurePoints: 24, TreasureCount: 2, AutonRankingPoint: true}},

		// Teleop treasures never count toward the auton points; shelf and stacked ones count toward the shelf count.
		{"teleop floor x5", S{TeleopFloor: 5}, nil, W{TeleopTreasurePoints: 10, TreasureCount: 5}},
		{
			"teleop first x3", S{TeleopFirst: 3}, nil,
			W{TeleopTreasurePoints: 15, TreasureCount: 3, ShelfTreasureCount: 3},
		},
		{
			"teleop top x4", S{TeleopTop: 4}, nil,
			W{TeleopTreasurePoints: 40, TreasureCount: 4, ShelfTreasureCount: 4},
		},
		{
			"teleop stacked x2", S{TeleopStacked: 2}, nil,
			W{TeleopTreasurePoints: 16, TreasureCount: 2, ShelfTreasureCount: 2},
		},
		{
			// Every counter counts once and the crown counts as one more treasure; the toss cube is not a treasure.
			"every counter and the crown",
			S{
				AutoFloor: 1, AutoFirst: 2, AutoTop: 3, TeleopFloor: 4, TeleopFirst: 5, TeleopTop: 6, TeleopStacked: 7,
				Crown: CrownTeleopTop, Toss: true,
			},
			nil,
			W{
				AutoTreasurePoints: 56, TeleopTreasurePoints: 149 + 20, TossPoints: 2, TreasureCount: 1 + 2 + 3 + 4 + 5 + 6 + 7 + 1, ShelfTreasureCount: 5 + 6 + 7 + 1,
				AutonRankingPoint: true, ScoringRankingPoint: true,
			},
		},

		// The crown is worth twice a treasure in the same place; see also TestScoreCrown.
		{
			// 3 x 10 + 20: the other treasures are not doubled.
			"crown with treasures in the same place", S{TeleopTop: 3, Crown: CrownTeleopTop}, nil,
			W{
				TeleopTreasurePoints: 50, TreasureCount: 4,
				ShelfTreasureCount: 4,
			},
		},
		{
			// Nothing else is doubled: not the robot points, not the toss, not the foul points.
			"crown with robots, toss and a foul",
			S{
				Crown: CrownTeleopTop, LeaveStatuses: [3]bool{true, true}, EndgameStatuses: [3]EndgameStatus{balance},
				Toss: true,
			},
			minors(ma2616),
			W{
				LeavePoints: 8, TeleopTreasurePoints: 20, EndgamePoints: 12, TossPoints: 2, TreasureCount: 1, ShelfTreasureCount: 1, FoulPoints: 5, EndgameRankingPoint: true,
			},
		},

		// Leave and auto balance are independent per-robot awards and stack: 4 + 12 = 16 for one robot.
		{"one robot left", S{LeaveStatuses: [3]bool{true}}, nil, W{LeavePoints: 4}},
		{"three robots left", S{LeaveStatuses: all3}, nil, W{LeavePoints: 12}},
		{"3v3 with the third robot bypassed", S{LeaveStatuses: [3]bool{true, true}}, nil, W{LeavePoints: 8}},
		{"balanced without a leave", S{AutoBalanceStatuses: [3]bool{false, true}}, nil, W{AutoBalancePoints: 12}},
		{
			"three robots balanced", S{AutoBalanceStatuses: all3}, nil,
			W{AutoBalancePoints: 36, AutonRankingPoint: true},
		},
		{
			"one robot left and balanced", S{LeaveStatuses: [3]bool{true}, AutoBalanceStatuses: [3]bool{true}}, nil,
			W{LeavePoints: 4, AutoBalancePoints: 12},
		},
		{
			// The third station is unused in a 2v2 match, so its statuses simply stay false.
			"2v2 with both robots doing both",
			S{LeaveStatuses: [3]bool{true, true}, AutoBalanceStatuses: [3]bool{true, true}}, nil,
			W{LeavePoints: 8, AutoBalancePoints: 24, AutonRankingPoint: true},
		},

		// Each robot has one endgame status, so park and balance are exclusive. Only balancing earns the Endgame RP, and
		// an auto balance is not an endgame balance (see the auto balance rows above).
		{"one parked", S{EndgameStatuses: [3]EndgameStatus{park}}, nil, W{EndgamePoints: 2}},
		{"two parked", S{EndgameStatuses: [3]EndgameStatus{park, park}}, nil, W{EndgamePoints: 4}},
		{"three parked", S{EndgameStatuses: [3]EndgameStatus{park, park, park}}, nil, W{EndgamePoints: 6}},
		{
			"balanced in the second station", S{EndgameStatuses: [3]EndgameStatus{0, balance}}, nil,
			W{EndgamePoints: 12, EndgameRankingPoint: true},
		},
		{
			"balanced in the third station", S{EndgameStatuses: [3]EndgameStatus{0, 0, balance}}, nil,
			W{EndgamePoints: 12, EndgameRankingPoint: true},
		},
		{
			"two balanced", S{EndgameStatuses: [3]EndgameStatus{balance, balance}}, nil,
			W{EndgamePoints: 24, EndgameRankingPoint: true},
		},
		{
			"three balanced", S{EndgameStatuses: [3]EndgameStatus{balance, balance, balance}}, nil,
			W{EndgamePoints: 36, EndgameRankingPoint: true},
		},

		// The toss is worth 2 points once per alliance and counts toward nothing else.
		{"toss", S{Toss: true}, nil, W{TossPoints: 2}},

		// Foul points come from the opponent's fouls, and a foul with no rule selected scores normally.
		{"opponent minor", S{}, minors(ma2616), W{FoulPoints: 5}},
		{"opponent major", S{}, majors(ma2604), W{FoulPoints: 10, NumOpponentMajorFouls: 1}},
		{
			"opponent minor and two majors", S{}, append(minors(ma2606), majors(ma2604, ma2604)...),
			W{FoulPoints: 25, NumOpponentMajorFouls: 2},
		},
		{"opponent minor with no rule", S{}, minors(0), W{FoulPoints: 5}},
		{"opponent major with no rule", S{}, majors(0), W{FoulPoints: 10, NumOpponentMajorFouls: 1}},
		{"an alliance's own fouls", S{Fouls: majors(ma2604)}, nil, W{}},

		// Auton RP: at least 20 auton points, or the opponent committed MA2603. Every auto value is a multiple of 4, so
		// 16 and 20 are the boundary cases. Teleop points and foul points never count.
		{
			"auton one step below the threshold", S{LeaveStatuses: [3]bool{true}, AutoTop: 1}, nil,
			W{LeavePoints: 4, AutoTreasurePoints: 12, TreasureCount: 1},
		},
		{
			"auton exactly at the threshold", S{LeaveStatuses: [3]bool{true, true}, AutoBalanceStatuses: [3]bool{true}},
			nil, W{LeavePoints: 8, AutoBalancePoints: 12, AutonRankingPoint: true},
		},
		{
			"auton not reached with foul points or other fouls", S{LeaveStatuses: [3]bool{true}},
			append(majors(ma2604, ma2604), minors(ma2606)...), W{LeavePoints: 4, FoulPoints: 25, NumOpponentMajorFouls: 2},
		},
		{
			"auton from the opponent's MA2603", S{LeaveStatuses: [3]bool{true}}, majors(ma2603),
			W{LeavePoints: 4, FoulPoints: 10, NumOpponentMajorFouls: 1, AutonRankingPoint: true},
		},
		{
			"auton earned both ways is one ranking point", S{AutoTop: 2}, majors(ma2603),
			W{
				AutoTreasurePoints: 24, TreasureCount: 2, FoulPoints: 10, NumOpponentMajorFouls: 1,
				AutonRankingPoint: true},
		},

		// Scoring RP: at least 12 teleop treasures on the first shelf, on the top shelf or stacked, counting the crown
		// as one (see the crown rows above). Floor treasures, auto placements and the toss do not count. The loop below
		// checks that the opponent's fouls make no difference.
		{
			"scoring one below the threshold", S{TeleopFirst: 6, TeleopTop: 4, TeleopStacked: 1}, nil,
			W{TeleopTreasurePoints: 78, TreasureCount: 11, ShelfTreasureCount: 11},
		},
		{
			"scoring exactly at the threshold", S{TeleopFirst: 6, TeleopTop: 4, TeleopStacked: 2}, nil,
			W{TeleopTreasurePoints: 86, TreasureCount: 12, ShelfTreasureCount: 12, ScoringRankingPoint: true},
		},
		{
			"scoring above the threshold", S{TeleopFirst: 13}, nil,
			W{TeleopTreasurePoints: 65, TreasureCount: 13, ShelfTreasureCount: 13, ScoringRankingPoint: true},
		},
		{
			"scoring not helped by floor treasures", S{TeleopFloor: 20, TeleopFirst: 11}, nil,
			W{TeleopTreasurePoints: 95, TreasureCount: 31, ShelfTreasureCount: 11},
		},
		{
			"scoring not helped by auto placements", S{AutoFloor: 5, AutoFirst: 5, AutoTop: 5, TeleopFirst: 11}, nil,
			W{
				AutoTreasurePoints: 120, TeleopTreasurePoints: 55, TreasureCount: 26, ShelfTreasureCount: 11,
				AutonRankingPoint: true,
			},
		},
		{
			"scoring not helped by the toss", S{TeleopTop: 11, Toss: true}, nil,
			W{TeleopTreasurePoints: 110, TossPoints: 2, TreasureCount: 11, ShelfTreasureCount: 11},
		},

		// Endgame RP: at least one robot balanced, or the opponent committed MA2601 or MA2602.
		{
			"endgame from the opponent's MA2601", S{}, majors(ma2601),
			W{FoulPoints: 10, NumOpponentMajorFouls: 1, EndgameRankingPoint: true},
		},
		{
			"endgame from the opponent's MA2602", S{}, majors(ma2602),
			W{FoulPoints: 10, NumOpponentMajorFouls: 1, EndgameRankingPoint: true},
		},
		{
			"endgame earned both ways is one ranking point", S{EndgameStatuses: [3]EndgameStatus{balance}}, majors(ma2602),
			W{
				EndgamePoints: 12, FoulPoints: 10, NumOpponentMajorFouls: 1, EndgameRankingPoint: true,
			},
		},
		{
			"no ranking points from other fouls or fouls with no rule", S{}, append(majors(ma2604, 0), minors(ma2606)...),
			W{FoulPoints: 25, NumOpponentMajorFouls: 2},
		},
		{"no ranking points from the alliance's own violations", S{Fouls: majors(ma2603, ma2601)}, nil, W{}},
		{
			"ranking points from both opponent violations", S{}, majors(ma2603, ma2601),
			W{
				FoulPoints: 20, NumOpponentMajorFouls: 2, AutonRankingPoint: true, EndgameRankingPoint: true},
		},
		{
			"all three bonus ranking points",
			S{AutoTop: 2, TeleopFirst: 12, EndgameStatuses: [3]EndgameStatus{balance}}, nil,
			W{
				AutoTreasurePoints: 24, TeleopTreasurePoints: 60, EndgamePoints: 12, TreasureCount: 14,
				ShelfTreasureCount: 12, AutonRankingPoint: true, ScoringRankingPoint: true, EndgameRankingPoint: true,
			},
		},

		// A well-populated pair of scores: red committed 7 fouls, 5 of them major.
		{
			"test score 1", *TestScore1(), TestScore2().Fouls,
			W{
				LeavePoints: 8, AutoBalancePoints: 12, AutoTreasurePoints: 32, TeleopTreasurePoints: 64,
				EndgamePoints: 14, TossPoints: 2, TreasureCount: 14,
				ShelfTreasureCount: 7, AutonRankingPoint: true, EndgameRankingPoint: true,
			},
		},
		{
			"test score 2", *TestScore2(), TestScore1().Fouls,
			W{
				LeavePoints: 4, AutoTreasurePoints: 8, TeleopTreasurePoints: 29, EndgamePoints: 2, TreasureCount: 7,
				ShelfTreasureCount: 4, FoulPoints: 60, NumOpponentMajorFouls: 5,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name, func(t *testing.T) {
				// Derive the sums from the primary fields.
				want := testCase.want
				want.AutonPoints = want.LeavePoints + want.AutoBalancePoints + want.AutoTreasurePoints
				want.MatchPoints = want.AutonPoints + want.TeleopTreasurePoints + want.EndgamePoints + want.TossPoints
				want.Score = want.MatchPoints + want.FoulPoints
				want.ShelfTreasureGoal = ScoringRpThreshold
				for _, earned := range []bool{want.AutonRankingPoint, want.ScoringRankingPoint, want.EndgameRankingPoint} {
					if earned {
						want.BonusRankingPoints++
					}
				}
				assert.Equal(t, &want, testCase.score.Summarize(&Score{Fouls: testCase.opponentFouls}))

				// There is no opponent-violation alternative for the Scoring RP.
				summary := testCase.score.Summarize(&Score{Fouls: majors(ma2601, ma2603)})
				assert.Equal(t, want.ScoringRankingPoint, summary.ScoringRankingPoint)
				assert.Equal(t, want.ShelfTreasureCount, summary.ShelfTreasureCount)
			},
		)
	}
}

// The crown is one treasure and is worth twice an ordinary treasure in the same place. It adds to the points of the
// period in which it was placed, and counts toward the Scoring RP only on a shelf or stacked in teleop. Each spot is
// added to 11 first-shelf treasures (55 points), one below the Scoring RP threshold.
func TestScoreCrown(t *testing.T) {
	testCases := []struct {
		crown   CrownPlacement
		points  int
		inAuto  bool
		onShelf bool
	}{
		{CrownNone, 0, false, false},
		{CrownAutoFloor, 8, true, false},
		{CrownAutoFirst, 16, true, false},
		{CrownAutoTop, 24, true, false},
		{CrownTeleopFloor, 4, false, false},
		{CrownTeleopFirst, 10, false, true},
		{CrownTeleopTop, 20, false, true},
		{CrownTeleopStacked, 16, false, true},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.crown.String(), func(t *testing.T) {
				points := testCase.points
				assert.Equal(t, points, testCase.crown.PointValue())
				want := ScoreSummary{
					TeleopTreasurePoints: 55, TreasureCount: 12,
					ShelfTreasureCount: 11, ShelfTreasureGoal: 12,
				}
				if testCase.crown == CrownNone {
					want.TreasureCount = 11
				}
				if testCase.inAuto {
					want.AutoTreasurePoints = points
					want.AutonPoints = points
					if points >= 20 {
						want.AutonRankingPoint = true
						want.BonusRankingPoints = 1
					}
				} else {
					want.TeleopTreasurePoints += points
				}
				if testCase.onShelf {
					want.ShelfTreasureCount = 12
					want.ScoringRankingPoint = true
					want.BonusRankingPoints = 1
				}
				want.MatchPoints = 55 + points
				want.Score = 55 + points

				score := Score{TeleopFirst: 11, Crown: testCase.crown}
				assert.Equal(t, &want, score.Summarize(&Score{}))
			},
		)
	}
}

func TestScorePlayoffDq(t *testing.T) {
	redScore := TestScore1()
	blueScore := TestScore2()

	// Unsetting the team and rule ID don't invalidate the foul.
	redScore.Fouls[0].TeamId = 0
	redScore.Fouls[0].RuleId = 0
	assert.Equal(t, 60, blueScore.Summarize(redScore).FoulPoints)

	// A disqualified alliance's summary is entirely zero, even if the opponent committed violations.
	redScore.PlayoffDq = true
	assert.Equal(t, &ScoreSummary{PlayoffDq: true}, redScore.Summarize(blueScore))
	assert.Equal(
		t, &ScoreSummary{PlayoffDq: true}, redScore.Summarize(&Score{Fouls: majors(ruleIdMa2603, ruleIdMa2601)}),
	)

	// Red's disqualification does not affect blue's own summary, including the foul points red committed.
	blueSummary := blueScore.Summarize(redScore)
	assert.Equal(t, 103, blueSummary.Score)
	assert.False(t, blueSummary.PlayoffDq)
	blueScore.PlayoffDq = true
	assert.Equal(t, &ScoreSummary{PlayoffDq: true}, blueScore.Summarize(redScore))
}

func TestScoreEquals(t *testing.T) {
	score1 := TestScore1()
	assert.True(t, score1.Equals(TestScore1()))
	assert.True(t, TestScore1().Equals(score1))
	assert.False(t, score1.Equals(TestScore2()))
	assert.False(t, TestScore2().Equals(score1))

	// Each mutation must be detected by Equals, so that every field of the score is covered.
	mutators := map[string]func(score *Score){
		"AutoFloor":              func(score *Score) { score.AutoFloor++ },
		"AutoFirst":              func(score *Score) { score.AutoFirst++ },
		"AutoTop":                func(score *Score) { score.AutoTop++ },
		"TeleopFloor":            func(score *Score) { score.TeleopFloor++ },
		"TeleopFirst":            func(score *Score) { score.TeleopFirst++ },
		"TeleopTop":              func(score *Score) { score.TeleopTop++ },
		"TeleopStacked":          func(score *Score) { score.TeleopStacked++ },
		"Crown":                  func(score *Score) { score.Crown = CrownAutoFloor },
		"LeaveStatuses[0]":       func(score *Score) { score.LeaveStatuses[0] = !score.LeaveStatuses[0] },
		"LeaveStatuses[1]":       func(score *Score) { score.LeaveStatuses[1] = !score.LeaveStatuses[1] },
		"LeaveStatuses[2]":       func(score *Score) { score.LeaveStatuses[2] = !score.LeaveStatuses[2] },
		"AutoBalanceStatuses[0]": func(score *Score) { score.AutoBalanceStatuses[0] = !score.AutoBalanceStatuses[0] },
		"AutoBalanceStatuses[1]": func(score *Score) { score.AutoBalanceStatuses[1] = !score.AutoBalanceStatuses[1] },
		"AutoBalanceStatuses[2]": func(score *Score) { score.AutoBalanceStatuses[2] = !score.AutoBalanceStatuses[2] },
		"EndgameStatuses[0]":     func(score *Score) { score.EndgameStatuses[0] = EndgameNone },
		"EndgameStatuses[1]":     func(score *Score) { score.EndgameStatuses[1] = EndgameBalance },
		"EndgameStatuses[2]":     func(score *Score) { score.EndgameStatuses[2] = EndgamePark },
		"Toss":                   func(score *Score) { score.Toss = !score.Toss },
		"Fouls length":           func(score *Score) { score.Fouls = []Foul{} },
		"Foul.FoulId":            func(score *Score) { score.Fouls[0].FoulId++ },
		"Foul.IsMajor":           func(score *Score) { score.Fouls[0].IsMajor = !score.Fouls[0].IsMajor },
		"Foul.TeamId":            func(score *Score) { score.Fouls[0].TeamId++ },
		"Foul.RuleId":            func(score *Score) { score.Fouls[0].RuleId++ },
		"PlayoffDq":              func(score *Score) { score.PlayoffDq = !score.PlayoffDq },
	}
	for name, mutate := range mutators {
		mutated := TestScore1()
		mutate(mutated)
		assert.False(t, score1.Equals(mutated), "Equals did not detect a change to %s", name)
		assert.False(t, mutated.Equals(score1), "Equals did not detect a change to %s", name)
	}
}

func TestScoreValidate(t *testing.T) {
	assert.Nil(t, TestScore1().Validate())
	for _, testCase := range []struct {
		name    string
		score   Score
		message string
	}{
		{"negative auto counter", Score{AutoFirst: -1}, "AutoFirst can't be negative"},
		{"negative teleop counter", Score{TeleopStacked: -2}, "TeleopStacked can't be negative"},
		{"crown below range", Score{Crown: -1}, "invalid crown placement -1"},
		{"crown above range", Score{Crown: CrownTeleopStacked + 1}, "invalid crown placement 8"},
		{"endgame above range", Score{EndgameStatuses: [3]EndgameStatus{0, 0, EndgameBalance + 1}}, "robot 3"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.score.Validate()
			if assert.NotNil(t, err) {
				assert.Contains(t, err.Error(), testCase.message)
			}
		})
	}
}

func TestCrownSpots(t *testing.T) {
	// Every placement except none appears exactly once, in enum order, with a unique panel id.
	ids := map[string]bool{}
	for i, spot := range CrownSpots {
		assert.Equal(t, CrownPlacement(i+1), spot.Placement)
		assert.False(t, ids[spot.Id], spot.Id)
		ids[spot.Id] = true
	}
	assert.Equal(t, int(CrownTeleopStacked), len(CrownSpots))

	assert.Equal(t, "None", CrownNone.String())
	assert.Equal(t, "Auto Floor", CrownAutoFloor.String())
	assert.Equal(t, "Teleop Top Shelf", CrownTeleopTop.String())
	labels := CrownLabels()
	assert.Equal(t, "-", labels[CrownNone])
	assert.Equal(t, "Teleop Stacked", labels[CrownTeleopStacked])
	assert.Equal(t, len(CrownSpots)+1, len(labels))
}
