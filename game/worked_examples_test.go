// Copyright 2026 Team 254. All Rights Reserved.
//
// Worked examples for the Medieval Mayhem game spec (specs/game_spec.yaml). Every expected number below was
// worked out by hand from the spec; if a change to the scoring code breaks one of these, either the change or the
// example is wrong. The treasure counters never include the crown.

package game

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

type workedExample struct {
	name     string // The purpose of the example.
	red      Score
	blue     Score
	redWant  ScoreSummary // The expected summaries; ShelfTreasureGoal is filled in by the test.
	blueWant ScoreSummary
	redRp    int // The ranking points earned in a qualification match: the result plus the bonus RPs.
	blueRp   int
	status   MatchStatus // The result of a qualification match.

	// The result with the playoff tiebreakers, checked only if there is a reason.
	playoffStatus MatchStatus
	playoffReason string
}

var workedExamples = func() []workedExample {
	type S = Score
	type W = ScoreSummary
	balance := EndgameBalance
	park := EndgamePark

	// Example D1: both alliances score 62 points, and the tie is broken by the auton points. The later D examples
	// change it.
	d1Red := S{
		LeaveStatuses: [3]bool{true, true}, AutoBalanceStatuses: [3]bool{true}, AutoFirst: 1, TeleopFloor: 1,
		TeleopFirst: 2, TeleopTop: 1, EndgameStatuses: [3]EndgameStatus{balance},
	}
	d1Blue := S{
		LeaveStatuses: [3]bool{true, true}, AutoFirst: 1, TeleopFloor: 1, TeleopFirst: 2, TeleopTop: 2, TeleopStacked: 1,
		EndgameStatuses: [3]EndgameStatus{park, park}, Toss: true,
	}
	d1RedWant := W{
		LeavePoints: 8, AutoBalancePoints: 12, AutoTreasurePoints: 8, AutonPoints: 28, TeleopTreasurePoints: 22,
		EndgamePoints: 12, MatchPoints: 62, Score: 62, TreasureCount: 5, ShelfTreasureCount: 3, AutonRankingPoint: true,
		EndgameRankingPoint: true, BonusRankingPoints: 2,
	}
	d1BlueWant := W{
		LeavePoints: 8, AutoTreasurePoints: 8, AutonPoints: 16, TeleopTreasurePoints: 40, EndgamePoints: 4,
		TossPoints: 2, MatchPoints: 62, Score: 62, TreasureCount: 7, ShelfTreasureCount: 5,
	}

	// Example D2: the same match, with one of blue's top-shelf treasures replaced by 10 points from a red major foul.
	d2Red := d1Red
	d2Red.Fouls = majors(ruleIdG415)
	d2Blue := d1Blue
	d2Blue.TeleopTop = 1
	d2BlueWant := d1BlueWant
	d2BlueWant.TeleopTreasurePoints = 30
	d2BlueWant.MatchPoints = 52
	d2BlueWant.FoulPoints = 10
	d2BlueWant.TreasureCount = 6
	d2BlueWant.ShelfTreasureCount = 4
	d2BlueWant.NumOpponentMajorFouls = 1

	// Example C: ranking points earned through opponent violations, with fouls (2v2).
	cRed := S{
		LeaveStatuses: [3]bool{true, true}, Crown: CrownAutoFirst, TeleopFirst: 2, TeleopTop: 3,
		EndgameStatuses: [3]EndgameStatus{park, park}, Toss: true,
		Fouls: majors(ruleIdMa2603, ruleIdMa2601, ruleIdMa2602),
	}
	cBlue := S{
		LeaveStatuses: [3]bool{true}, TeleopFloor: 4, TeleopFirst: 1, TeleopTop: 1,
		Fouls: append(minors(ruleIdMa2606), majors(ruleIdMa2604, ruleIdMa2604)...),
	}
	cRedWant := W{
		LeavePoints: 8, AutoTreasurePoints: 16, AutonPoints: 24, TeleopTreasurePoints: 40, EndgamePoints: 4,
		TossPoints: 2, Crown: CrownAutoFirst, CrownPoints: 16, TreasureCount: 6, ShelfTreasureCount: 5, MatchPoints: 70,
		FoulPoints: 25, Score: 95, AutonRankingPoint: true, BonusRankingPoints: 1, NumOpponentMajorFouls: 2,
	}
	// Blue earns both of its ranking points through red's violations; red earns its own on points.
	cBlueWant := W{
		LeavePoints: 4, AutonPoints: 4, TeleopTreasurePoints: 23, TreasureCount: 6, ShelfTreasureCount: 2,
		MatchPoints: 27, FoulPoints: 30, Score: 57, AutonRankingPoint: true, AutonRankingPointByFoul: true,
		EndgameRankingPoint: true, EndgameRankingPointByFoul: true, BonusRankingPoints: 2, NumOpponentMajorFouls: 3,
	}

	// Example C2: red's auto foul is entered as a major foul with no rule selected. Blue keeps its 30 foul points and
	// its score of 57, but loses the Auton RP; selecting MA2603 in Edit Match Result restores it, as in example C.
	c2Red := cRed
	c2Red.Fouls = majors(0, ruleIdMa2601, ruleIdMa2602)
	c2BlueWant := cBlueWant
	c2BlueWant.AutonRankingPoint = false
	c2BlueWant.AutonRankingPointByFoul = false
	c2BlueWant.BonusRankingPoints = 1

	return []workedExample{
		// The first three are also used for the rankings example below.
		{
			name: "A: an ordinary qualification match (2v2)",
			red: S{
				LeaveStatuses: [3]bool{true, true}, AutoFirst: 1, TeleopFloor: 3, TeleopFirst: 4, TeleopTop: 2,
				TeleopStacked: 1, EndgameStatuses: [3]EndgameStatus{balance, park}, Toss: true,
				Fouls: minors(ruleIdMa2616),
			},
			blue: S{
				LeaveStatuses: [3]bool{true}, AutoFloor: 1, TeleopFloor: 2, TeleopFirst: 2, TeleopTop: 1,
				Crown: CrownTeleopFirst, EndgameStatuses: [3]EndgameStatus{park},
			},
			redWant: W{
				LeavePoints: 8, AutoTreasurePoints: 8, AutonPoints: 16, TeleopTreasurePoints: 54, EndgamePoints: 14,
				TossPoints: 2, TreasureCount: 11, ShelfTreasureCount: 7, MatchPoints: 86, Score: 86,
				EndgameRankingPoint: true, BonusRankingPoints: 1,
			},
			blueWant: W{
				LeavePoints: 4, AutoTreasurePoints: 4, AutonPoints: 8, TeleopTreasurePoints: 34, EndgamePoints: 2,
				Crown: CrownTeleopFirst, CrownPoints: 10, TreasureCount: 7, ShelfTreasureCount: 4, MatchPoints: 44,
				FoulPoints: 5, Score: 49,
			},
			redRp: 4, blueRp: 0, status: RedWonMatch,
		},
		{
			name: "B: thresholds exactly at the boundary, plus the crown (2v2)",
			red: S{
				LeaveStatuses: [3]bool{true, true}, AutoBalanceStatuses: [3]bool{true}, TeleopFloor: 2, TeleopFirst: 6,
				TeleopTop: 4, TeleopStacked: 2, EndgameStatuses: [3]EndgameStatus{0, park},
			},
			blue: S{
				LeaveStatuses: [3]bool{true, true}, AutoFirst: 1, TeleopFloor: 1, TeleopFirst: 5, TeleopTop: 3,
				TeleopStacked: 2, Crown: CrownTeleopTop, EndgameStatuses: [3]EndgameStatus{balance, balance},
				Toss: true,
			},
			redWant: W{
				LeavePoints: 8, AutoBalancePoints: 12, AutonPoints: 20, TeleopTreasurePoints: 90, EndgamePoints: 2,
				TreasureCount: 14, ShelfTreasureCount: 12, MatchPoints: 112, Score: 112, AutonRankingPoint: true,
				ScoringRankingPoint: true, BonusRankingPoints: 2,
			},
			blueWant: W{
				LeavePoints: 8, AutoTreasurePoints: 8, AutonPoints: 16, TeleopTreasurePoints: 93, EndgamePoints: 24,
				TossPoints: 2, Crown: CrownTeleopTop, CrownPoints: 20, TreasureCount: 13, ShelfTreasureCount: 11,
				MatchPoints: 135, Score: 135, EndgameRankingPoint: true, BonusRankingPoints: 1,
			},
			redRp: 2, blueRp: 4, status: BlueWonMatch,
		},
		{
			name: "C: ranking points earned through opponent violations, with fouls (2v2)",
			red:  cRed, blue: cBlue, redWant: cRedWant, blueWant: cBlueWant, redRp: 4, blueRp: 2, status: RedWonMatch,
		},
		{
			name: "C2: a foul with no rule selected earns no ranking point",
			red:  c2Red, blue: cBlue, redWant: cRedWant, blueWant: c2BlueWant, redRp: 4, blueRp: 1, status: RedWonMatch,
		},
		{
			name: "D1: a tied playoff match broken by the auton points",
			red:  d1Red, blue: d1Blue, redWant: d1RedWant, blueWant: d1BlueWant, redRp: 3, blueRp: 1, status: TieMatch,
			playoffStatus: RedWonMatch, playoffReason: "TIEBREAK: AUTON POINTS",
		},
		{
			name: "D2: the same match decided by the major-foul tiebreaker, even though red has more auton points",
			red:  d2Red, blue: d2Blue, redWant: d1RedWant, blueWant: d2BlueWant, redRp: 3, blueRp: 1, status: TieMatch,
			playoffStatus: BlueWonMatch, playoffReason: "TIEBREAK: MAJOR FOULS",
		},
		{
			name: "D3: the same score and the same auton points, decided by the match points",
			red: S{
				LeaveStatuses: [3]bool{true, true}, AutoFirst: 1, TeleopFloor: 2, TeleopTop: 4,
			},
			blue: S{
				LeaveStatuses: [3]bool{true, true}, AutoFirst: 1, TeleopFloor: 2, TeleopFirst: 1, TeleopTop: 4,
				Fouls: minors(ruleIdMa2607),
			},
			redWant: W{
				LeavePoints: 8, AutoTreasurePoints: 8, AutonPoints: 16, TeleopTreasurePoints: 44, MatchPoints: 60,
				FoulPoints: 5, Score: 65, TreasureCount: 7, ShelfTreasureCount: 4,
			},
			blueWant: W{
				LeavePoints: 8, AutoTreasurePoints: 8, AutonPoints: 16, TeleopTreasurePoints: 49, MatchPoints: 65,
				Score: 65, TreasureCount: 8, ShelfTreasureCount: 5,
			},
			redRp: 1, blueRp: 1, status: TieMatch, playoffStatus: BlueWonMatch, playoffReason: "TIEBREAK: MATCH POINTS",
		},
		{
			name: "D4: identical inputs leave every tiebreak level equal, so the match is a true tie and is replayed",
			red:  d1Red, blue: d1Red, redWant: d1RedWant, blueWant: d1RedWant, redRp: 3, blueRp: 3, status: TieMatch,
			playoffStatus: TieMatch, playoffReason: "TRUE TIE",
		},
		{
			name: "E: the 3v3 fallback, with blue's third robot bypassed",
			red: S{
				LeaveStatuses: [3]bool{true, true, true}, AutoBalanceStatuses: [3]bool{false, false, true},
				TeleopFirst: 3, TeleopTop: 2, EndgameStatuses: [3]EndgameStatus{balance, balance, park}, Toss: true,
			},
			// The bypassed robot in station 3 leaves its statuses at their zero values.
			blue: S{
				LeaveStatuses: [3]bool{true, true}, AutoTop: 1, TeleopFirst: 8, TeleopTop: 3, TeleopStacked: 1,
				EndgameStatuses: [3]EndgameStatus{balance},
			},
			redWant: W{
				LeavePoints: 12, AutoBalancePoints: 12, AutonPoints: 24, TeleopTreasurePoints: 35, EndgamePoints: 26,
				TossPoints: 2, TreasureCount: 5, ShelfTreasureCount: 5, MatchPoints: 87, Score: 87,
				AutonRankingPoint: true, EndgameRankingPoint: true, BonusRankingPoints: 2,
			},
			blueWant: W{
				LeavePoints: 8, AutoTreasurePoints: 12, AutonPoints: 20, TeleopTreasurePoints: 78, EndgamePoints: 12,
				TreasureCount: 13, ShelfTreasureCount: 12, MatchPoints: 110, Score: 110, AutonRankingPoint: true,
				ScoringRankingPoint: true, EndgameRankingPoint: true, BonusRankingPoints: 3,
			},
			redRp: 2, blueRp: 6, status: BlueWonMatch,
		},
	}
}()

func TestWorkedExamples(t *testing.T) {
	for _, example := range workedExamples {
		t.Run(
			example.name, func(t *testing.T) {
				redSummary := example.red.Summarize(&example.blue)
				blueSummary := example.blue.Summarize(&example.red)

				redWant, blueWant := example.redWant, example.blueWant
				redWant.ShelfTreasureGoal = ScoringRpThreshold
				blueWant.ShelfTreasureGoal = ScoringRpThreshold
				assert.Equal(t, &redWant, redSummary, "red")
				assert.Equal(t, &blueWant, blueSummary, "blue")
				assert.Equal(t, example.redRp, qualificationRankingPoints(redSummary, blueSummary), "red RP")
				assert.Equal(t, example.blueRp, qualificationRankingPoints(blueSummary, redSummary), "blue RP")

				status, tiebreakReason := DetermineMatchStatus(redSummary, blueSummary, false)
				assert.Equal(t, example.status, status)
				assert.Equal(t, "", tiebreakReason)
				if example.playoffReason != "" {
					status, tiebreakReason = DetermineMatchStatus(redSummary, blueSummary, true)
					assert.Equal(t, example.playoffStatus, status)
					assert.Equal(t, example.playoffReason, tiebreakReason)
				}
			},
		)
	}
}

// Returns the ranking points an alliance earns in a qualification match: the match result plus its bonus RPs.
func qualificationRankingPoints(summary, opponentSummary *ScoreSummary) int {
	rankingPoints := summary.BonusRankingPoints
	if summary.Score > opponentSummary.Score {
		rankingPoints += GetWinRankingPoints()
	} else if summary.Score == opponentSummary.Score {
		rankingPoints++
	}
	return rankingPoints
}

// The ranking check from examples A, B and C: one qualification match each, with teams 101 to 104, 105 to 108 and 109
// to 112 (red alliance first), sorted by ranking points, then average score, then average auton points.
func TestWorkedExampleRankings(t *testing.T) {
	originalRankingRandomFloat64 := RankingRandomFloat64
	defer func() { RankingRandomFloat64 = originalRankingRandomFloat64 }()

	// Give each team a distinct random value so that the partner ordering is deterministic: the lower team number of
	// each pair sorts first.
	rankings := Rankings{}
	for i, example := range workedExamples[:3] {
		redSummary := example.red.Summarize(&example.blue)
		blueSummary := example.blue.Summarize(&example.red)
		for j := 0; j < 4; j++ {
			teamId := 101 + 4*i + j
			ownSummary, opponentSummary := redSummary, blueSummary
			if j >= 2 {
				ownSummary, opponentSummary = blueSummary, redSummary
			}
			RankingRandomFloat64 = func() float64 { return 1.0 - float64(teamId-101)/100.0 }
			ranking := Ranking{TeamId: teamId}
			ranking.AddScoreSummary(ownSummary, opponentSummary, false)
			rankings = append(rankings, ranking)
		}
	}
	sort.Sort(rankings)

	// The table from the spec: team, ranking points, score, then auton points.
	expectedRankings := [][4]int{
		{107, 4, 135, 16}, {108, 4, 135, 16}, {109, 4, 95, 24}, {110, 4, 95, 24}, {101, 4, 86, 16}, {102, 4, 86, 16},
		{105, 2, 112, 20}, {106, 2, 112, 20}, {111, 2, 57, 4}, {112, 2, 57, 4}, {103, 0, 49, 8}, {104, 0, 49, 8},
	}
	if assert.Equal(t, len(expectedRankings), len(rankings)) {
		for i, expected := range expectedRankings {
			assert.Equal(t, expected[0], rankings[i].TeamId, "rank %d", i+1)
			assert.Equal(t, expected[1], rankings[i].RankingPoints, "team %d RP", expected[0])
			assert.Equal(t, expected[2], rankings[i].ScorePoints, "team %d score", expected[0])
			assert.Equal(t, expected[3], rankings[i].AutonPoints, "team %d auton", expected[0])
		}
	}
}
