// Copyright 2020 Team 254. All Rights Reserved.
// Author: pat@patfairbank.com (Patrick Fairbank)

package game

import (
	"github.com/stretchr/testify/assert"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// Rule IDs used by the tests in this package, from the rules list in rule.go.
const (
	ruleIdG415   = 6  // Major, no ranking point.
	ruleIdMa2601 = 14 // Major, grants the opponent the Endgame RP.
	ruleIdMa2602 = 15 // Major, grants the opponent the Endgame RP.
	ruleIdMa2603 = 16 // Major, grants the opponent the Auton RP.
	ruleIdMa2604 = 17 // Major, no ranking point.
	ruleIdMa2606 = 18 // Minor, no ranking point.
	ruleIdMa2607 = 19 // Minor, no ranking point.
	ruleIdMa2616 = 29 // Minor, no ranking point.
)

func TestGetRuleById(t *testing.T) {
	assert.Nil(t, GetRuleById(0))
	assert.Equal(t, rules[0], GetRuleById(1))
	assert.Equal(t, rules[20], GetRuleById(21))
	assert.Nil(t, GetRuleById(1000))
}

func TestGetAllRules(t *testing.T) {
	allRules := GetAllRules()
	assert.Equal(t, len(rules), len(allRules))
	for _, rule := range rules {
		assert.Equal(t, rule, allRules[rule.Id])
	}
}

func TestRuleIdsAreSequentialAndUnique(t *testing.T) {
	for i, rule := range rules {
		assert.Equal(t, i+1, rule.Id, "rule %s is out of order", rule.RuleNumber)
		assert.NotEmpty(t, rule.RuleNumber)
		assert.NotEmpty(t, rule.Description)
	}

	// A duplicate ID would silently hide a rule, since the lookup is a map keyed by ID.
	assert.Equal(t, len(rules), len(GetAllRules()))
}

func TestRuleCounts(t *testing.T) {
	numMajor, numMinor, numRankingPoint := 0, 0, 0
	for _, rule := range rules {
		if rule.IsMajor {
			numMajor++
		} else {
			numMinor++
		}
		if rule.IsRankingPoint {
			numRankingPoint++
			// Only a major foul can grant the opposing alliance a ranking point.
			assert.True(t, rule.IsMajor, "%s grants a ranking point but is not major", rule.RuleNumber)
		}
	}

	assert.Equal(t, 32, len(rules))
	assert.Equal(t, 17, numMajor)
	assert.Equal(t, 15, numMinor)
	assert.Equal(t, 3, numRankingPoint)
}

func TestRankingPointRules(t *testing.T) {
	// The ranking-point conditions look these up by rule number, so both the numbers and the IDs matter.
	for ruleId, ruleNumber := range map[int]string{
		ruleIdMa2601: "MA2601",
		ruleIdMa2602: "MA2602",
		ruleIdMa2603: "MA2603",
	} {
		rule := GetRuleById(ruleId)
		if assert.NotNil(t, rule) {
			assert.Equal(t, ruleNumber, rule.RuleNumber)
			assert.True(t, rule.IsMajor)
			assert.True(t, rule.IsRankingPoint)
		}
	}

	// Rules that may be assessed at either severity appear twice, once per severity.
	for _, ruleNumber := range []string{"G423", "MA2615"} {
		severities := map[bool]int{}
		for _, rule := range rules {
			if rule.RuleNumber == ruleNumber {
				severities[rule.IsMajor]++
			}
		}
		assert.Equal(t, map[bool]int{false: 1, true: 1}, severities, "%s", ruleNumber)
	}
}

// Checks that the rules list matches the table in specs/RULES.md, row for row.
func TestRulesMatchSpecTable(t *testing.T) {
	specText, err := os.ReadFile("../specs/RULES.md")
	assert.Nil(t, err)

	var tableRows [][]string
	for _, line := range strings.Split(string(specText), "\n") {
		if !strings.HasPrefix(line, "| ") {
			continue
		}
		cells := strings.Split(strings.Trim(line, "| "), " | ")
		if cells[0] == "Id" || strings.HasPrefix(cells[0], "---") {
			continue // Header and separator rows.
		}
		tableRows = append(tableRows, cells)
	}

	if assert.Equal(t, len(rules), len(tableRows)) {
		for i, rule := range rules {
			cells := tableRows[i]
			if assert.Equal(t, 4, len(cells), "row %d", i+1) {
				assert.Equal(t, strconv.Itoa(rule.Id), cells[0])
				assert.Equal(t, rule.RuleNumber, cells[1])
				assert.Equal(t, rule.IsMajor, cells[2] == "Major", "%s", rule.RuleNumber)
				assert.Equal(t, rule.Description, cells[3])
			}
		}
	}
}

// Checks that every rule the game spec says gives the opponent a ranking point is in the list and flagged as such, and
// that no other rule is flagged.
func TestRankingPointRulesMatchGameSpec(t *testing.T) {
	specText, err := os.ReadFile("../specs/2026_medieval_mayhem.yaml")
	assert.Nil(t, err)

	specRuleNumbers := map[string]bool{}
	for _, match := range regexp.MustCompile(`opponent_fouls: \[([^\]]*)\]`).FindAllStringSubmatch(string(specText), -1) {
		for _, ruleNumber := range strings.Split(match[1], ",") {
			specRuleNumbers[strings.TrimSpace(ruleNumber)] = true
		}
	}
	assert.NotEmpty(t, specRuleNumbers)

	for _, rule := range rules {
		assert.Equal(t, specRuleNumbers[rule.RuleNumber], rule.IsRankingPoint, "%s", rule.RuleNumber)
	}
	for ruleNumber := range specRuleNumbers {
		found := false
		for _, rule := range rules {
			found = found || rule.RuleNumber == ruleNumber
		}
		assert.True(t, found, "%s is in the game spec but not in the rules list", ruleNumber)
	}
}
