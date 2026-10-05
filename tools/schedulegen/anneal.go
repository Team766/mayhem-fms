// Builds a 2v2 template by simulated annealing.
//
// The schedule is a stream of slots, four per match. The first numTeams*matchesPerTeam slots are split into
// matchesPerTeam rounds of numTeams slots, and each round holds every team exactly once. If the appearances don't
// fill the last match, the leftover slots form a short tail of surrogate appearances (distinct teams, so at most
// one each). Because the rounds are fixed, every team plays once per round and so is spread evenly across the
// schedule. The search only ever moves teams around inside a round, so it can never break that rule.

package main

import (
	"math"
	"math/rand"
)

// Cost weights. Partner repeats matter most, then keeping a team's matches apart, then facing the same opponent
// again, then red/blue and station balance. A team twice in one match is effectively forbidden.
const (
	weightDuplicate = 1_000_000
	weightShortGap  = 1_000 // Per (wanted gap - actual gap)^2, see wantedGap.
	weightTieBreak  = 25    // A small push toward a bigger gap at every size.
	weightPartner   = 1_000 // Per extra repeat: the cost grows as C(times together, 2).
	weightOpponent  = 250   // The same, for facing each other.
	weightRedBlue   = 300   // Per (red - blue)^2 for a team.
	weightStation   = 20    // Per (first station - second station)^2 for a team.

	// Annealing: the temperature falls geometrically over movesPerSlot moves for each slot in the schedule. Then
	// a zero-temperature polish keeps only improvements.
	startTemperature   = 300.0
	endTemperature     = 4.0
	movesPerSlot       = 3_000
	finishingMovesSlot = 300

	// Runs restart from a new random start until one is good enough (see goodEnough) or this many have been tried;
	// the best is kept. Easy shapes stop after one or two, so the budget goes to the hard ones.
	maxRestarts = 30
)

// The chance of each kind of move, in percent.
const (
	chanceRoundSwap = 70
	chanceSlotSwap  = 25 // The remaining percent replaces a surrogate, when there are any.
)

type annealer struct {
	numTeams       int
	matchesPerTeam int
	numMatches     int
	numSurrogates  int
	wantedGap      int
	rng            *rand.Rand

	// The stream of teams (0-based), and for each match which stream slot goes in which alliance slot.
	stream []int
	order  [][slotsPerMatch]int

	partners  []int   // Times each pair of teams (smaller*n+larger) shares an alliance.
	opponents []int   // Times each pair of teams faces each other.
	redBlue   []int   // Per team: counted red appearances minus blue.
	station   []int   // Per team: counted first-station appearances minus second-station.
	times     [][]int // Per team: stream slot of its appearance in each round, plus the tail (-1 if absent).
	inTail    []bool
	cost      int64
}

// Generates the best template found for the given shape, deterministically for a given seed.
func Generate(numTeams, matchesPerTeam int, seed int64) []Row {
	var best *annealer
	for restart := 0; restart < maxRestarts; restart++ {
		a := newAnnealer(numTeams, matchesPerTeam, seed*1000+int64(restart))
		a.anneal()
		if best == nil || a.cost < best.cost {
			best = a
		}
		if goodEnough(best.rows(), numTeams, matchesPerTeam) {
			break
		}
	}
	return best.rows()
}

// Whether a template reaches the ideal on every soft goal: pairings spread as evenly as the counts would allow
// with no other constraints (Bounds), red/blue within one, and the wanted gap between a team's matches. Some
// shapes can't, and then every restart is used.
func goodEnough(rows []Row, numTeams, matchesPerTeam int) bool {
	metrics := Measure(rows, numTeams)
	partnerBound, opponentBound := Bounds(numTeams, matchesPerTeam)
	return metrics.MaxPartner <= partnerBound && metrics.MaxOpponent <= opponentBound &&
		metrics.MaxRedBlueDiff <= 1 && metrics.MinBetween >= wantedGap(numTeams)
}

func newAnnealer(numTeams, matchesPerTeam int, seed int64) *annealer {
	numMatches, numSurrogates := shape(numTeams, matchesPerTeam)
	a := &annealer{
		numTeams:       numTeams,
		matchesPerTeam: matchesPerTeam,
		numMatches:     numMatches,
		numSurrogates:  numSurrogates,
		wantedGap:      wantedGap(numTeams),
		rng:            rand.New(rand.NewSource(seed)),
		stream:         make([]int, numMatches*slotsPerMatch),
		order:          make([][slotsPerMatch]int, numMatches),
		partners:       make([]int, numTeams*numTeams),
		opponents:      make([]int, numTeams*numTeams),
		redBlue:        make([]int, numTeams),
		station:        make([]int, numTeams),
		times:          make([][]int, numTeams),
		inTail:         make([]bool, numTeams),
	}
	// Each round starts as a random permutation of the teams.
	for round := 0; round < matchesPerTeam; round++ {
		copy(a.stream[round*numTeams:], a.rng.Perm(numTeams))
	}
	// The surrogate tail is a random set of distinct teams.
	copy(a.stream[numTeams*matchesPerTeam:], a.rng.Perm(numTeams)[:numSurrogates])
	for m := range a.order {
		a.order[m] = [slotsPerMatch]int{0, 1, 2, 3}
	}
	for team := range a.times {
		a.times[team] = make([]int, matchesPerTeam+1)
		for i := range a.times[team] {
			a.times[team][i] = -1
		}
	}
	for slot, team := range a.stream {
		a.times[team][a.roundOf(slot)] = slot
	}
	for slot := numTeams * matchesPerTeam; slot < len(a.stream); slot++ {
		a.inTail[a.stream[slot]] = true
	}
	a.cost = a.fullCost()
	return a
}

// The round a stream slot belongs to; the surrogate tail counts as one extra round.
func (a *annealer) roundOf(slot int) int {
	return min(slot/a.numTeams, a.matchesPerTeam)
}

func (a *annealer) isCounted(slot int) bool {
	return slot < a.numTeams*a.matchesPerTeam
}

// The team in each alliance slot of a match: red1, red2, blue1, blue2.
func (a *annealer) lineup(match int) [slotsPerMatch]int {
	var lineup [slotsPerMatch]int
	for i, streamOffset := range a.order[match] {
		lineup[i] = a.stream[match*slotsPerMatch+streamOffset]
	}
	return lineup
}

func (a *annealer) anneal() {
	slots := len(a.stream)
	total := movesPerSlot * slots
	ratio := math.Pow(endTemperature/startTemperature, 1/float64(total))
	temperature := startTemperature
	for i := 0; i < total; i++ {
		a.tryMove(temperature)
		temperature *= ratio
	}
	for i := 0; i < finishingMovesSlot*slots; i++ {
		a.tryMove(0)
	}
}

// Makes one random move and keeps it if it lowers the cost (or, when warm, sometimes if it raises it).
func (a *annealer) tryMove(temperature float64) {
	before := a.cost
	undo := a.randomMove()
	delta := a.cost - before
	if delta <= 0 {
		return
	}
	if temperature > 0 && a.rng.Float64() < math.Exp(-float64(delta)/temperature) {
		return
	}
	undo()
}

// Applies a random move and returns a function that reverses it.
func (a *annealer) randomMove() func() {
	roll := a.rng.Intn(100)
	switch {
	case roll < chanceRoundSwap:
		first := a.rng.Intn(a.numTeams * a.matchesPerTeam)
		round := first / a.numTeams
		second := round*a.numTeams + a.rng.Intn(a.numTeams)
		a.swapStream(first, second)
		return func() { a.swapStream(first, second) }
	case roll < chanceRoundSwap+chanceSlotSwap || a.numSurrogates == 0:
		match := a.rng.Intn(a.numMatches)
		first := a.rng.Intn(slotsPerMatch)
		second := (first + 1 + a.rng.Intn(slotsPerMatch-1)) % slotsPerMatch
		a.swapOrder(match, first, second)
		return func() { a.swapOrder(match, first, second) }
	default:
		slot := a.numTeams*a.matchesPerTeam + a.rng.Intn(a.numSurrogates)
		replacement := a.rng.Intn(a.numTeams)
		old := a.stream[slot]
		if a.inTail[replacement] {
			return func() {}
		}
		a.replaceSurrogate(slot, replacement)
		return func() { a.replaceSurrogate(slot, old) }
	}
}

// Swaps the teams in two stream slots of the same round.
func (a *annealer) swapStream(first, second int) {
	if first == second {
		return
	}
	matchA, matchB := first/slotsPerMatch, second/slotsPerMatch
	teamA, teamB := a.stream[first], a.stream[second]
	a.cost -= a.gapCost(teamA) + a.gapCost(teamB)
	a.removeMatch(matchA)
	if matchB != matchA {
		a.removeMatch(matchB)
	}
	a.stream[first], a.stream[second] = teamB, teamA
	round := a.roundOf(first)
	a.times[teamA][round], a.times[teamB][round] = second, first
	a.addMatch(matchA)
	if matchB != matchA {
		a.addMatch(matchB)
	}
	a.cost += a.gapCost(teamA) + a.gapCost(teamB)
}

// Swaps which alliance slots two of a match's teams take. Their times don't change.
func (a *annealer) swapOrder(match, first, second int) {
	a.removeMatch(match)
	a.order[match][first], a.order[match][second] = a.order[match][second], a.order[match][first]
	a.addMatch(match)
}

// Puts a different team in a surrogate slot.
func (a *annealer) replaceSurrogate(slot, team int) {
	match := slot / slotsPerMatch
	old := a.stream[slot]
	a.cost -= a.gapCost(old) + a.gapCost(team)
	a.removeMatch(match)
	a.stream[slot] = team
	a.times[old][a.matchesPerTeam] = -1
	a.times[team][a.matchesPerTeam] = slot
	a.inTail[old], a.inTail[team] = false, true
	a.addMatch(match)
	a.cost += a.gapCost(old) + a.gapCost(team)
}

func (a *annealer) removeMatch(match int) { a.updateMatch(match, -1) }
func (a *annealer) addMatch(match int)    { a.updateMatch(match, 1) }

// Adds (direction 1) or removes (-1) a match's contribution to the pair counts and the cost.
func (a *annealer) updateMatch(match, direction int) {
	lineup := a.lineup(match)
	a.bumpPair(a.partners, lineup[0], lineup[1], direction, weightPartner)
	a.bumpPair(a.partners, lineup[2], lineup[3], direction, weightPartner)
	for _, red := range lineup[:2] {
		for _, blue := range lineup[2:] {
			a.bumpPair(a.opponents, red, blue, direction, weightOpponent)
		}
	}
	for i, team := range lineup {
		if !a.isCounted(match*slotsPerMatch + a.order[match][i]) {
			continue
		}
		a.bumpBalance(a.redBlue, team, direction*(1-2*(i/2)), weightRedBlue)
		a.bumpBalance(a.station, team, direction*(1-2*(i%2)), weightStation)
	}
}

// Changes a pair's count by direction and updates the cost, which grows as C(count, 2) times the weight.
func (a *annealer) bumpPair(counts []int, x, y, direction int, weight int64) {
	if x > y {
		x, y = y, x
	}
	index := x*a.numTeams + y
	if direction > 0 {
		a.cost += int64(counts[index]) * weight
		counts[index]++
	} else {
		counts[index]--
		a.cost -= int64(counts[index]) * weight
	}
}

func (a *annealer) bumpBalance(balance []int, team, change int, weight int64) {
	a.cost -= int64(balance[team]*balance[team]) * weight
	balance[team] += change
	a.cost += int64(balance[team]*balance[team]) * weight
}

// The cost of how close together a team's consecutive appearances are.
func (a *annealer) gapCost(team int) int64 {
	var cost int64
	previous := -1
	for _, slot := range a.times[team] {
		if slot < 0 {
			continue
		}
		if previous >= 0 {
			cost += a.gapPenalty(slot/slotsPerMatch - previous/slotsPerMatch - 1)
		}
		previous = slot
	}
	return cost
}

// The penalty for a team's next appearance coming this many matches after its previous one.
func (a *annealer) gapPenalty(between int) int64 {
	if between < 0 {
		return weightDuplicate
	}
	penalty := int64(weightTieBreak * max(0, 2-between))
	if shortfall := a.wantedGap - between; shortfall > 0 {
		penalty += int64(weightShortGap * shortfall * shortfall)
	}
	return penalty
}

// How many matches we want between a team's own matches, at least. A team averages about numTeams/4 - 1 matches
// between its own, so demanding close to that average pins every team to nearly the same place in each round. The
// same teams then meet match after match. So small events are asked for less: a gap of g needs numTeams to be at
// least 4*(g+1) plus some room, which is about 2 at 14 teams, 1 at 9 and nothing below that.
func wantedGap(numTeams int) int {
	switch {
	case numTeams < 9:
		return 0
	case numTeams < 14:
		return 1
	}
	return 2
}

// Computes the cost from scratch; used once at the start, after which moves update it incrementally.
func (a *annealer) fullCost() int64 {
	a.cost = 0
	for match := 0; match < a.numMatches; match++ {
		a.addMatch(match)
	}
	for team := 0; team < a.numTeams; team++ {
		a.cost += a.gapCost(team)
	}
	return a.cost
}

// Converts the stream into template rows with 1-based team indexes.
func (a *annealer) rows() []Row {
	rows := make([]Row, a.numMatches)
	for match := range rows {
		for i, streamOffset := range a.order[match] {
			slot := match*slotsPerMatch + streamOffset
			rows[match][i] = Appearance{Team: a.stream[slot] + 1, Surrogate: !a.isCounted(slot)}
		}
	}
	return rows
}
