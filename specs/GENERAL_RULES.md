# General rules

The rules that carry over from one FRC season to the next. With no game applied, the FMS's foul list
(`game/rule.go`) holds only these; a game replaces them with the fouls from its own manual.

Rule numbers don't say whether a rule carries over. FRC renumbers its rules every year, its `G` rules include season
rules (G405 names that year's OUTPOST, G412 its RUNGS), and a game's own manual may renumber them again. A game's own
rules (`MA26xx` for M-Ayhem 2026) are always season rules. The numbers below are from the 2026 FRC manual.

| Rule | Foul | 2026 FRC |
| --- | --- | --- |
| No colluding to influence Ranking Points | minor, affects ranking points | G206 |
| No strategy aimed at forcing the opponent into a violation | major | G210 |
| No delaying the start of a match | major | G301 |
| Drive team stays in its staged area in auto | minor | G401 |
| Drive team doesn't touch the robot or console in auto | minor | G402 |
| No contacting an opponent from across the center line in auto | major | G403 |
| No using a game piece to ease or amplify a field element's challenge | major | G404 |
| No ejecting game pieces from the field | minor or major | G405 |
| No damaging game pieces | major | G406 |
| No grabbing or attaching to field elements | major | G412 |
| Stay within the extension limits | major | G413 |
| No damaging or impairing an opponent by contact inside its frame perimeter | major | G415 |
| No deliberately damaging or impairing an opponent | major | G416 |
| No deliberately attaching to, tipping or entangling an opponent | major | G417 |
| No pinning for more than 3 seconds | minor or major | G418 |
| No robots working together to close off part of the game | major | G419 |
| Drive team stays in its designated area | minor | G421 |
| Only the team's drivers and human players operate the robot | major | G422 |
| Coaches don't touch game pieces | minor | G426 |

## Is a rule general?

- It's about conduct, safety, the drive team, or robots interacting, not a season's field element or game piece.
- The same idea appears in most recent seasons, even when the wording, limits and zone names change.
- If it only makes sense once you name that season's field element, game piece or mechanism, it's a season rule.

## When FRC starts a new season

`sync-upstream` checks this table against the new season's manual: match rules by meaning, not number, update the
numbers, add any rule that has become a recurring idea, and change `game/rule.go` to match.
