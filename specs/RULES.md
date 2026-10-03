# Foul Rules

This is the foul list the FMS offers referees, taken from section 4.5 of the 2026 M-Ayhem manual, and `game/rule.go` is written from this table. A rule the referee may call at either severity has one row per severity. Two fouls have no rule number in the manual, so the Rule column holds a short label instead. Yellow and red cards and no-score calls are not FMS fouls and are not listed.

| Id | Rule | Foul | Description |
| --- | --- | --- | --- |
| 1 | G210 | Major | Do not force an opponent robot to commit a foul. |
| 2 | G401 | Major | Drive team members must stay behind the lines during auto. |
| 3 | G402 | Major | Drive team members must not control their robot during auto, except to press A-stop or E-stop. |
| 4 | G403 | Minor | A robot may not cross the midfield line into the opponent's half during auto. |
| 5 | G408 | Major | No human or human player may damage a treasure. |
| 6 | G415 | Major | Contacting another robot inside its perimeter. Can stack with MA2618. |
| 7 | G423 | Minor | Damaging or impairing an opponent robot (referee's choice of severity). |
| 8 | G423 | Major | Damaging or impairing an opponent robot (referee's choice of severity). |
| 9 | G424 | Major | Intentionally attaching to, entangling or tipping another robot. |
| 10 | G425 | Minor | Pinning a robot for longer than 3 seconds. |
| 11 | G429 | Major | Drive team members must stay in their designated areas. |
| 12 | G430 | Major | A robot may be operated only by its drivers. |
| 13 | G434 | Major | Coaches may not touch treasures, unless for safety. |
| 14 | MA2601 | Major | Contacting the opposing alliance's balance beam during endgame. Opponent gets the Endgame RP. Can stack with MA2602. |
| 15 | MA2602 | Major | Contacting an opposing robot at all while it is on its own balance beam during endgame. Opponent gets the Endgame RP. Can stack with MA2601. |
| 16 | MA2603 | Major | Entering the opposing alliance's safe house or safe zone during auto. Opponent gets the Auton RP. |
| 17 | MA2604 | Major | Hoarding: more than 6 non-scoring treasures in your safe house or safe zone. Repeats every 10 seconds. |
| 18 | MA2606 | Minor | Entering the opposing human loading zone, safe house or safe zone, or touching their balance beam before endgame. Repeats every 10 seconds. The apron facing midfield is allowed. |
| 19 | MA2607 | Minor | Putting a treasure into the field other than through the human loading holes (not during the Toss). |
| 20 | MA2608 | Minor | Loaded treasure does not first touch an own-alliance robot or the human loading zone floor. |
| 21 | MA2609 | Minor | Contacting own shelf repeatedly or in an unsafe manner. |
| 22 | MA2610 | Minor | Starting auto with more than 1 treasure. Treasures that robot places in auto do not count. |
| 23 | MA2611 | Minor | Descoring the other alliance's treasures. Can stack with MA2620. |
| 24 | MA2612 | Minor | Robot deliberately destroying a treasure. |
| 25 | MA2613 | Minor | Adding treasures to the field during endgame. The treasure does not count. |
| 26 | MA2614 | Minor | Placing or shooting treasures into the opposing alliance's zones. |
| 27 | MA2615 | Minor | Failure to obey referee-area signs (referee's choice of severity). |
| 28 | MA2615 | Major | Failure to obey referee-area signs (referee's choice of severity). |
| 29 | MA2616 | Minor | Controlling more than 2 treasures for more than a moment. Extra treasures do not count. |
| 30 | MA2617 | Minor | Throwing a treasure onto the shelf from outside own safe house. The treasure does not count. |
| 31 | Damaging contact | Major | Initiating damaging contact (referee's discretion). |
| 32 | Reaching onto field | Major | Having a part of your body reach onto the field (referee's discretion). |

Left out on purpose: the card-only violations (G101, G102, G206, G209, G246, MA2605, MA2618, MA2619, MA2620, MA2621, MA2623, MA2624) and the no-score violation MA2622. Rule ids must not change within an event, because stored fouls refer to them.
