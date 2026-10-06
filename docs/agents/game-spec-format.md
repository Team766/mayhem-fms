# Game spec format

The game spec is one YAML file, `specs/game_spec.yaml`, that tells [apply-game.md](apply-game.md) what this year's
game is. It always holds the current game; earlier games live in git history. People and coding agents read it; the
FMS never parses it.

It holds only the game's scoring: game pieces, the elements scorers enter and who enters them, the totals the
screens show, the bonus ranking points and the tiebreakers. These live elsewhere:

- **Match timing:** event settings, with the game's defaults in code.
- **Foul values and the rules list:** the rules pull request (`specs/RULES.md`, in the manual's own words).
- **Screen layouts:** they follow from the structure below (see apply-game's table).
- **Worked examples:** `apply-game` writes them as tests.
- **Open questions:** a questions file kept out of the pull requests until they are answered.

## Sections

| Section | Each entry | Notes |
|---------|------------|-------|
| `game_pieces` | `id`, `display_name` | The physical pieces. |
| `scoring_groups` | `id`, `display_name` | Totals shown on the audience screens. A group adds up its elements' points and counts their pieces; a status in a group counts as one piece when it isn't none. An element with no group is shown on its own. |
| `scoring_counts` | `id`, `display_name`, `game_piece`, `scoring_group`, `scorer`, `phases` | A +/- counter per alliance. `scorer` is `near` or `far`. `phases` lists `{phase, points}`; teleop covers the endgame period too. |
| `statuses` | `id`, `display_name`, `per`, `scorer`, `phases`, optional `game_piece`, `scoring_group` | `per` is `robot` or `alliance`. A phase with `points` is yes/no; with `values` it is one of the values, or none (0 points, never listed). One status across several phases is still one choice. |
| `bonus_ranking_points` | `id`, `display_name`, `earned_when`, optional `threshold`, `opponent_fouls` | 1 RP each in qualification matches. `earned_when` is one plain sentence. `opponent_fouls`: the other alliance committing any of these rule numbers also earns it. Thresholds are fixed game constants. |
| `ranking_tiebreakers` | `metric` | After ranking points; each is an average per match. |
| `playoff_tiebreakers` | `metric` | For a tied playoff match, in order; say which side wins. A true tie plays the next game of the series. |

A short example, from 2026:

```yaml
scoring_counts:
  - id: top_shelf
    display_name: "Top Shelf"
    game_piece: treasure
    scoring_group: treasure
    scorer: near
    phases:
      - { phase: auto, points: 12 }
      - { phase: teleop, points: 10 }

statuses:
  - id: endgame
    display_name: "Endgame"
    per: robot
    scorer: far
    phases:
      - { phase: endgame, values: { park: 2, balance: 12 } }

bonus_ranking_points:
  - id: endgame
    display_name: "Endgame"
    earned_when: At least one of the alliance's robots is balanced at the end.
    opponent_fouls: [MA2601, MA2602]
```

## Writing rules that save a round trip

- Give every element a short, unique `id`; ids become Go fields, JSON keys and CSS ids.
- State comparisons and directions outright: "at least the threshold", "fewer major fouls wins".
- Say what does *not* count ("floor and auto treasures don't count toward the Scoring RP").
- If the spec needs something the sections above can't say, stop and add a section here first.

## Sample prompts

Apply a new game:
> Follow `docs/agents/apply-game.md` with `specs/game_spec.yaml`. Ask me the open questions before you change code. Work on new branches off the base and do not push.

Try a rule change during design:
> Using `docs/agents/apply-game.md`, update the game in this branch for the edits I just made to `specs/game_spec.yaml`. Show me the tests that changed.

Bring the base up to date:
> Follow `docs/agents/sync-upstream.md` in `update` mode up to upstream tag `v2027.0.1`. Give me the classification table before porting anything.

Re-apply this year's game after a base update:
> The base moved to a new upstream. Re-apply `specs/game_spec.yaml` with `docs/agents/apply-game.md` on fresh branches off the base and tell me what differed from the previous application.
