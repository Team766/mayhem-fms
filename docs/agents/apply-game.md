# Playbook: apply-game

Replace the game in the tree with the game described by `specs/game_spec.yaml`.

For a coding agent (any LLM) or a careful person. It is phase 2 of two; phase 1 is
[sync-upstream.md](sync-upstream.md). The spec format is in [game-spec-format.md](game-spec-format.md). The game
being replaced is whatever `specs/game_spec.yaml` held before (`git show <base>:specs/game_spec.yaml`), if any.

## The idea

The output is ordinary Cheesy Arena code: a typed score struct, a `Summarize()` a volunteer could check against
the manual's point table, and screens that name the game's elements. No runtime config, no generator, no engine.
The spec is the durable artifact; when the tree moves to a new upstream, the same spec is applied again.

## Ground rules

1. **Stay inside the game seam** as [../DEVELOPMENT.md](../DEVELOPMENT.md) defines it. If the spec needs a base change (arena, PLC, playoffs, a new display), stop and ask, or do it as a separate pull request.
2. **The spec decides; last year's code does not.** Where the spec is silent or contradicts itself, ask. If nobody can answer, choose, record the question in the questions file (kept out of the pull requests), and say so.
3. **Find every consumer.** This playbook names no files on purpose. Anything that enters, reads or displays the score, the summary or the rankings must be updated. Find them with `git grep`; step 6 checks you did.
4. **Robot count.** Rules about "all robots" use the robots actually playing (see [../TwoVTwo.md](../TwoVTwo.md)), never a literal 3.
5. **One place for each fact.** Point values and ranking point thresholds are named constants next to the score model, not event settings. A list the screens share (such as the places a special piece can end up) is defined once in Go and read by every template and script.
6. Upstream's style: typed fields, table-driven tests, imports alphabetical and ungrouped.
7. New branches. Do not push or open pull requests unless asked. If you think you need a destructive command, stop and ask.

## Three pull requests, in this order

1. **Spec:** `specs/game_spec.yaml`. The game's authors review it.
2. **Rules:** `specs/RULES.md`, the foul list in the manual's own words (struck-out text dropped), one row per severity a referee can call; `game/rule.go` written from it, with a test that checks they match and that the fouls the spec's `opponent_fouls` name give a ranking point. Wording we would like changed goes in `specs/RULE_SUGGESTIONS.md`, never into the list. The authors review this too.
3. **Game:** everything else in the seam.

## Procedure

1. **Read the whole spec and settle open questions first.** They are cheap now and expensive after forty files. The game being replaced shows how a game plugs into every screen.
2. **Name what goes and what comes.** Outgoing: the ids and labels in the previous spec and the Go, JSON and CSS names derived from them. Incoming: the names you will use, from the new spec's ids. Put both in the pull request.
3. **Model and math, with table-driven tests in the same commit:** every scoring element in every phase; every ranking point at its threshold and one below; every tiebreak level with each side winning; every rule that mentions robots, in each alliance size; a table of whole matches with hand-worked totals. List the tests that encode a ruling on an open question, so a person can check those numbers.
4. **Entry, then display and reporting.** Each spec construct becomes code the same way every year:

   | Spec construct | Score model | Entered on | Summary and screens |
   |---|---|---|---|
   | Counter, points per phase | An int per phase | +/- counter on the scorer's panel, under its phase | Its group's points and count |
   | Yes/no status, per robot | A bool per robot | A toggle per robot | Its phase's points |
   | One-of status, per robot | A typed enum per robot | A button group per robot | Its phase's points |
   | One-of status, per alliance | One typed enum | A button group under each phase, one shared value | Its phase's points; counts as a piece in its group |
   | Scoring group | Nothing | Nothing | Summary points and count; live count on the audience overlay |
   | Bonus RP | A summary bool; threshold as a constant | Nothing | A mark on the final score; counted in rankings; live progress on the overlay as "‹name› RP n/threshold" |
   | Foul that awards an RP | The rule's RP flag | Referee panel, unchanged | Same as a bonus RP |
   | Ranking tiebreakers | Ranking fields and their sort | Nothing | Rankings display columns |
   | Playoff tiebreakers | The tied-match cascade | Nothing | The reason on the final score |

   Conventions: panels follow phase order (auto, then teleop, then endgame), items in spec order, per-robot items in station order; every element is editable in match review, and the server rejects values no scorer could enter; the final score shows each alliance's total ranking points. If the spec uses something this table doesn't cover, stop and ask, then add a row here.
5. **Install the game's assets**, replacing upstream's files in place (the logos keep their names); a new sound cue gets its own file.
6. **Completeness checks.**
   - No outgoing name is left: `git grep -n -i -E '<outgoing words>'` over Go, templates, scripts, styles and CSV, excluding `docs/` and `specs/`, prints nothing you cannot explain.
   - Every spec id appears in the score model, on its scorer's panel, in match review and in the summary.
   - Every score, summary or ranking field that a script or template reads exists in the Go structs.
   - `git diff --stat <branch you started from>` touches only the game seam. Explain any exception.
7. **Verify.**
   - `go generate ./... && go fmt ./... && go vet ./... && go build ./... && go test ./...`, with no new vet warnings.
   - Run a whole event through the real screens, scripted or by hand: settings, teams, schedule, qualification matches covering every element, fouls and cards, an edited result, rankings, alliance selection and playoffs. Check every score against a hand calculation and that no page shows an error, `NaN` or `undefined`.
   - Keep the screenshots out of the branch; a short Markdown write-up of them goes in a pull request comment so reviewers can see the game.
8. **Review before reviewers.** Build a local review page of the game pull request (the `review-local` skill) for the author, and work through their comments before the pull request goes out.

## Report

The three pull requests and their commits; assumptions made where the spec was silent; anything touched outside
the seam and why; the tests that encode rulings; screens checked and not checked; proposed edits to the spec or
to this playbook. Do not claim a check you did not run.
