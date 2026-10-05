// Template rows for 2v2 qualification schedules: reading, writing, validating and measuring them.

package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

// A 2v2 match has four slots, in the order red1, red2, blue1, blue2. The first two are the red alliance.
const slotsPerMatch = 4

// One appearance of a team in a match.
type Appearance struct {
	Team      int // 1-based team index, as in the template files
	Surrogate bool
}

// One match of a template.
type Row [slotsPerMatch]Appearance

// Returns the number of matches needed to give each of numTeams teams matchesPerTeam appearances, and the number
// of surrogate appearances that fill out the last match.
func shape(numTeams, matchesPerTeam int) (numMatches, numSurrogates int) {
	appearances := numTeams * matchesPerTeam
	numMatches = (appearances + slotsPerMatch - 1) / slotsPerMatch
	return numMatches, numMatches*slotsPerMatch - appearances
}

// Writes rows in the 12-column format that tournament.BuildRandomSchedule reads: (team index, surrogate flag) for
// red1, red2, red3, blue1, blue2, blue3, where the 2v2 third slots are always 0, 0.
func WriteCsv(w io.Writer, rows []Row) error {
	writer := csv.NewWriter(w)
	for _, row := range rows {
		record := make([]string, 0, 12)
		for _, slot := range []int{0, 1, -1, 2, 3, -1} {
			team, flag := 0, 0
			if slot >= 0 {
				team = row[slot].Team
				if row[slot].Surrogate {
					flag = 1
				}
			}
			record = append(record, strconv.Itoa(team), strconv.Itoa(flag))
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

// Reads a template written by WriteCsv, rejecting anything that isn't a well-formed 2v2 row.
func ReadCsv(r io.Reader) ([]Row, error) {
	records, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, err
	}
	rows := make([]Row, len(records))
	for i, record := range records {
		if len(record) != 12 {
			return nil, fmt.Errorf("row %d has %d columns, expected 12", i+1, len(record))
		}
		values := make([]int, 12)
		for j, field := range record {
			if values[j], err = strconv.Atoi(field); err != nil {
				return nil, fmt.Errorf("row %d: %v", i+1, err)
			}
		}
		if values[4] != 0 || values[5] != 0 || values[10] != 0 || values[11] != 0 {
			return nil, fmt.Errorf("row %d: red3 and blue3 must be 0,0", i+1)
		}
		for slot, column := range []int{0, 2, 6, 8} {
			if values[column+1] != 0 && values[column+1] != 1 {
				return nil, fmt.Errorf("row %d: surrogate flag must be 0 or 1", i+1)
			}
			rows[i][slot] = Appearance{Team: values[column], Surrogate: values[column+1] == 1}
		}
	}
	return rows, nil
}

// Checks the hard rules: right number of matches, four distinct real teams per match, every team playing exactly
// matchesPerTeam counted matches, and surrogates filling the leftover slots with at most one per team.
func Validate(rows []Row, numTeams, matchesPerTeam int) error {
	numMatches, numSurrogates := shape(numTeams, matchesPerTeam)
	if len(rows) != numMatches {
		return fmt.Errorf("%d matches, expected %d", len(rows), numMatches)
	}
	counted := make([]int, numTeams+1)
	surrogates := make([]int, numTeams+1)
	totalSurrogates := 0
	for i, row := range rows {
		seen := make(map[int]bool)
		for _, appearance := range row {
			team := appearance.Team
			if team < 1 || team > numTeams {
				return fmt.Errorf("match %d: team index %d out of range", i+1, team)
			}
			if seen[team] {
				return fmt.Errorf("match %d: team %d appears twice", i+1, team)
			}
			seen[team] = true
			if appearance.Surrogate {
				surrogates[team]++
				totalSurrogates++
			} else {
				counted[team]++
			}
		}
	}
	for team := 1; team <= numTeams; team++ {
		if counted[team] != matchesPerTeam {
			return fmt.Errorf("team %d plays %d counted matches, expected %d", team, counted[team], matchesPerTeam)
		}
		if surrogates[team] > 1 {
			return fmt.Errorf("team %d has %d surrogate appearances", team, surrogates[team])
		}
	}
	if totalSurrogates != numSurrogates {
		return fmt.Errorf("%d surrogate appearances, expected %d", totalSurrogates, numSurrogates)
	}
	return nil
}

// Quality numbers for a template. All pair counts include surrogate appearances, since those matches are played.
type Metrics struct {
	MaxPartner     int // Most times any two teams share an alliance.
	MaxOpponent    int // Most times any two teams face each other.
	MaxRedBlueDiff int // Largest |red - blue| over teams, counted matches only.
	MaxStationDiff int // Largest |first station - second station| over teams, counted matches only.
	MinBetween     int // Fewest matches between any team's consecutive appearances.
	BackToBack     int // Number of consecutive appearances with no match in between.
	OneBetween     int // Number of consecutive appearances with exactly one match in between.
}

func Measure(rows []Row, numTeams int) Metrics {
	n := numTeams + 1
	partners := make([]int, n*n)
	opponents := make([]int, n*n)
	redBlue := make([]int, n)
	station := make([]int, n)
	last := make([]int, n) // Index of each team's previous appearance, or -1.
	for i := range last {
		last[i] = -1
	}
	metrics := Metrics{MinBetween: len(rows)}
	for i, row := range rows {
		count := func(counts []int, a, b int) {
			if a > b {
				a, b = b, a
			}
			counts[a*n+b]++
		}
		count(partners, row[0].Team, row[1].Team)
		count(partners, row[2].Team, row[3].Team)
		for _, red := range row[:2] {
			for _, blue := range row[2:] {
				count(opponents, red.Team, blue.Team)
			}
		}
		for slot, appearance := range row {
			team := appearance.Team
			if !appearance.Surrogate {
				redBlue[team] += 1 - 2*(slot/2)
				station[team] += 1 - 2*(slot%2)
			}
			if last[team] >= 0 {
				between := i - last[team] - 1
				if between < metrics.MinBetween {
					metrics.MinBetween = between
				}
				if between == 0 {
					metrics.BackToBack++
				} else if between == 1 {
					metrics.OneBetween++
				}
			}
			last[team] = i
		}
	}
	for _, c := range partners {
		metrics.MaxPartner = max(metrics.MaxPartner, c)
	}
	for _, c := range opponents {
		metrics.MaxOpponent = max(metrics.MaxOpponent, c)
	}
	for team := 1; team <= numTeams; team++ {
		metrics.MaxRedBlueDiff = max(metrics.MaxRedBlueDiff, abs(redBlue[team]))
		metrics.MaxStationDiff = max(metrics.MaxStationDiff, abs(station[team]))
	}
	return metrics
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Returns the lowest possible values of MaxPartner and MaxOpponent, if the pairings were spread perfectly. Each
// appearance has one partner and two opponents, spread over the other teams; surrogates add one appearance for
// some teams.
func Bounds(numTeams, matchesPerTeam int) (maxPartner, maxOpponent int) {
	appearances := matchesPerTeam
	if _, numSurrogates := shape(numTeams, matchesPerTeam); numSurrogates > 0 {
		appearances++
	}
	others := numTeams - 1
	return (appearances + others - 1) / others, (2*appearances + others - 1) / others
}
