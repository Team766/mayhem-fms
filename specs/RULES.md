# Foul Rules

This is the foul list the FMS offers referees, taken from section 4.5 of the 2026 M-Ayhem manual, and `game/rule.go` is written from this table. Descriptions are the manual's own words, with struck-out text left out; wording we would like changed is in [RULE_SUGGESTIONS.md](RULE_SUGGESTIONS.md). A rule the referee may call at either severity has one row per severity. Two fouls have no rule number in the manual, so the Rule column holds a short label instead. Yellow and red cards and no-score calls are not FMS fouls and are not listed.

| Id | Rule | Foul | Description |
| --- | --- | --- | --- |
| 1 | G210 | Major | Do not force another robot to commit a foul. |
| 2 | G401 | Major | DRIVE TEAM members must be behind the lines during AUTO and |
| 3 | G402 | Major | DRIVE TEAM members must not control their robots during auton, directly or through operator console, unless it’s for a-stop or e-stop |
| 4 | G403 | Minor | Do not cross the midfield line into the opponents half during AUTO. |
| 5 | G408 | Major | Neither a Human or Human Player may damage a Scoring Element |
| 6 | G415 | Major | Contacting another robot inside the robot’s perimeter [Stackable with taking a piece within a robots bounding box (re: yellow card)] |
| 7 | G423 | Minor | This isn’t combat robotics. A ROBOT may not damage or functionally impair an opponent ROBOT: A. deliberately. B. regardless of intent, by initiating contact, either directly or transitively via a SCORING ELEMENT CONTROLLED by the ROBOT, inside the vertical projection of an opponent’s ROBOT PERIMETER. Damage or functional impairment because of contact with a tipped-over opponent ROBOT, which is not perceived by a referee to be deliberate, is not a violation of this rule. |
| 8 | G423 | Major | This isn’t combat robotics. A ROBOT may not damage or functionally impair an opponent ROBOT: A. deliberately. B. regardless of intent, by initiating contact, either directly or transitively via a SCORING ELEMENT CONTROLLED by the ROBOT, inside the vertical projection of an opponent’s ROBOT PERIMETER. Damage or functional impairment because of contact with a tipped-over opponent ROBOT, which is not perceived by a referee to be deliberate, is not a violation of this rule. |
| 9 | G424 | Major | Attaching, entangling or tipping over another robot intentionally |
| 10 | G425 | Minor | Pinning a robot for longer than 3 seconds |
| 11 | G429 | Major | Drive team members must stay designated areas |
| 12 | G430 | Major | Robot should only be operated by Drivers |
| 13 | G434 | Major | Coaches many not touch Scoring Elements, unless for safety purposes |
| 14 | MA2601 | Major | Contacting the opposing team's balance beam at during the endgame [Ref Discretion, Can stack with MA2602]. This results in an automatic endgame RP for the other team, no matter the scoring |
| 15 | MA2602 | Major | Contacting the opposing team AT ALL while they’re in their own Balance Beam during the endgame [Can stack with MA2601]. This results in an automatic endgame RP for the other team, no matter the scoring. |
| 16 | MA2603 | Major | Going into another team's safe house or safe zone within auton. This also results in an automatic auto RP for the other team, no matter the scoring. |
| 17 | MA2604 | Major | No Hoarding: More than 6 non-scoring treasures in your safe house or safe zone. Major Foul is given at the beginning of infraction, given every 10 seconds 6 non-scoring treasures are in the safe house or safe zone. |
| 18 | MA2606 | Minor | Entering another alliance’s Human Player Loading zone, Safe House, Safe Zone, or touching the balance beam before endgame. This can be reapplied every 10 seconds if contact does not end. Exception: Ref may allow an exception due to momentary contact with the balance beam or contact due to contact with another robot (e.g. being pushed into the balance beam). Exception: The Apron facing the middle of the field is okay to touch during anytime before the endgame. Note that MA2601 applies during endgame. |
| 19 | MA2607 | Minor | Placing/Throwing a treasure into the field not through the Human Player Loading holes [doesn’t apply during the “toss”, may be escalated by the head ref if flagrant] |
| 20 | MA2608 | Minor | Placing/Throwing a treasure that bounces not first on the same alliance’s robot, or the Human Player Loading zone Floor. |
| 21 | MA2609 | Minor | Contacting scoring own shelf repeatedly or in an unsafe manner [Ref Discretion] |
| 22 | MA2610 | Minor | Starting with over 1 piece during auton. The robot is also ineligible to gain points within auton by scoring in the shelf |
| 23 | MA2611 | Minor | Descoring pieces from other teams [can stack with MA2620] |
| 24 | MA2612 | Minor | Robot deliberately permanently destroying a treasure |
| 25 | MA2613 | Minor | Adding treasures within endgame [treasure does not count for scoring] |
| 26 | MA2614 | Minor | Placing/Shooting treasures within the opposing alliance zones |
| 27 | MA2615 | Minor | Failure to obey ref area signs |
| 28 | MA2615 | Major | Failure to obey ref area signs |
| 29 | MA2616 | Minor | Controlling more than 2 gamepieces at a non-momentarily time [additional treasures do not count for scoring][ref discretion] |
| 30 | MA2617 | Minor | Throwing a treasure into/on the shelf not inside own Safe House[treasure does not count for scoring] |
| 31 | Damaging contact | Major | Initiating damaging contact [Ref discretion] |
| 32 | Reaching onto field | Major | Having a part of your body reaching on to the field [Ref Discretion] |

Left out on purpose: the card-only violations (G101, G102, G206, G209, G246, MA2605, MA2618, MA2619, MA2620, MA2621, MA2623, MA2624) and the no-score violation MA2622. Rule ids must not change within an event, because stored fouls refer to them.
