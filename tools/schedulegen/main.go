// Generates the 2v2 qualification schedule templates in schedules/2p_<teams>_<matchesPerTeam>.csv.
//
// Usage, from the repository root:
//
//	go run ./tools/schedulegen                       # the whole grid
//	go run ./tools/schedulegen -teams 14 -matches 8  # one template
//
// The output is deterministic, so rerunning it reproduces the checked-in files.

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

func main() {
	minTeams := flag.Int("min-teams", 6, "smallest team count to generate")
	maxTeams := flag.Int("max-teams", 14, "largest team count to generate")
	minMatches := flag.Int("min-matches", 4, "smallest matches per team to generate")
	maxMatches := flag.Int("max-matches", 10, "largest matches per team to generate")
	teams := flag.Int("teams", 0, "generate only this team count")
	matches := flag.Int("matches", 0, "generate only this many matches per team")
	outDir := flag.String("out", "schedules", "directory to write the templates to")
	report := flag.Bool("report", false, "print each template's quality numbers instead of writing files")
	flag.Parse()
	if *teams > 0 {
		*minTeams, *maxTeams = *teams, *teams
	}
	if *matches > 0 {
		*minMatches, *maxMatches = *matches, *matches
	}

	type job struct{ numTeams, matchesPerTeam int }
	jobs := make(chan job)
	var wg sync.WaitGroup
	var mu sync.Mutex
	failed := false
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				if err := generateFile(*outDir, j.numTeams, j.matchesPerTeam, *report); err != nil {
					mu.Lock()
					failed = true
					fmt.Fprintln(os.Stderr, err)
					mu.Unlock()
				}
			}
		}()
	}
	for numTeams := *minTeams; numTeams <= *maxTeams; numTeams++ {
		for matchesPerTeam := *minMatches; matchesPerTeam <= *maxMatches; matchesPerTeam++ {
			jobs <- job{numTeams, matchesPerTeam}
		}
	}
	close(jobs)
	wg.Wait()
	if failed {
		os.Exit(1)
	}
}

// The seed depends only on the shape, so a template doesn't change when others are added or removed.
func seedFor(numTeams, matchesPerTeam int) int64 {
	return int64(numTeams)*100 + int64(matchesPerTeam)
}

func generateFile(outDir string, numTeams, matchesPerTeam int, report bool) error {
	rows := Generate(numTeams, matchesPerTeam, seedFor(numTeams, matchesPerTeam))
	if err := Validate(rows, numTeams, matchesPerTeam); err != nil {
		return fmt.Errorf("%d teams, %d matches per team: %v", numTeams, matchesPerTeam, err)
	}
	if report {
		partnerBound, opponentBound := Bounds(numTeams, matchesPerTeam)
		fmt.Printf(
			"%3d teams %2d matches: %+v bounds partner %d opponent %d\n", numTeams, matchesPerTeam,
			Measure(rows, numTeams), partnerBound, opponentBound,
		)
		return nil
	}
	file, err := os.Create(filepath.Join(outDir, fmt.Sprintf("2p_%d_%d.csv", numTeams, matchesPerTeam)))
	if err != nil {
		return err
	}
	defer file.Close()
	return WriteCsv(file, rows)
}
