// Copyright 2014 Team 254. All Rights Reserved.
// Author: pat@patfairbank.com (Patrick Fairbank)

package web

import (
	"github.com/Team254/cheesy-arena/field"
	"github.com/Team254/cheesy-arena/game"
	"github.com/Team254/cheesy-arena/websocket"
	gorillawebsocket "github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestScoringPanel(t *testing.T) {
	web := setupTestWeb(t)

	recorder := web.getHttpResponse("/panels/scoring/invalidposition")
	assert.Equal(t, 500, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Invalid position")

	// All six positions should render, each with the controls appropriate to it.
	for _, position := range []string{"red", "blue", "red_near", "red_far", "blue_near", "blue_far"} {
		recorder = web.getHttpResponse("/panels/scoring/" + position)
		assert.Equal(t, 200, recorder.Code)
		assert.Contains(t, recorder.Body.String(), "Scoring Panel - Untitled Event - Cheesy Arena")
		body := recorder.Body.String()

		parameters := positionParameters[position]
		if parameters.ShowsNear() {
			assert.Contains(t, body, "Auto treasure")
			assert.Contains(t, body, "Stacked")
		} else {
			assert.NotContains(t, body, "Auto treasure")
		}
		if parameters.ShowsFar() {
			assert.Contains(t, body, "Teleop crown")
			assert.Contains(t, body, "crown-auto_top")
			assert.Contains(t, body, "crown-teleop_stacked")
			assert.Contains(t, body, "The Toss")
			assert.Contains(t, body, "Balance")
		} else {
			assert.NotContains(t, body, "The Toss")
			assert.NotContains(t, body, "Teleop crown")
		}
	}
}

func TestScoringPanelWebsocket(t *testing.T) {
	web := setupTestWeb(t)
	web.arena.EventSettings.TwoVsTwoMode = false

	server, wsUrl := web.startTestServer()
	defer server.Close()
	_, _, err := gorillawebsocket.DefaultDialer.Dial(wsUrl+"/panels/scoring/blorpy/websocket", nil)
	assert.NotNil(t, err)

	// Connect the near and far panels for each alliance; both register under the alliance, not the position.
	registry := &web.arena.ScoringPanelRegistry
	redNearWs := dialScoringPanel(t, wsUrl, "red_near")
	assert.Equal(t, 1, registry.GetNumPanels("red"))
	redFarWs := dialScoringPanel(t, wsUrl, "red_far")
	assert.Equal(t, 2, registry.GetNumPanels("red"))
	assert.Equal(t, 0, registry.GetNumPanels("blue"))
	blueWs := dialScoringPanel(t, wsUrl, "blue")
	assert.Equal(t, 1, registry.GetNumPanels("blue"))

	// Should get a few status updates right after connection.
	for _, ws := range []*websocket.Websocket{redNearWs, redFarWs, blueWs} {
		readInitialScoringPanelMessages(t, ws)
	}

	// Exercise each websocket command in turn. Each accepted command notifies all three panels and satisfies its check;
	// a command with no check is rejected, and silently dropped without a notification.
	type m = map[string]any
	type check func(s *game.Score) bool
	steps := []struct {
		name    string
		command string
		data    any
		check   check
	}{
		{
			"treasure", "treasure", m{"Counter": "auto_top", "Adjustment": 1},
			func(s *game.Score) bool { return s.AutoTop == 1 },
		},
		{
			"treasure clamps at zero", "treasure", m{"Counter": "auto_top", "Adjustment": -5},
			func(s *game.Score) bool { return s.AutoTop == 0 },
		},
		{"unknown counter", "treasure", m{"Counter": "nonexistent", "Adjustment": 1}, nil},
		{"crown", "crown", m{"Value": "teleop_top"}, func(s *game.Score) bool { return s.Crown == game.CrownTeleopTop }},
		{"leave on", "leave", m{"TeamPosition": 1}, func(s *game.Score) bool { return s.LeaveStatuses[0] }},
		{"leave toggles off", "leave", m{"TeamPosition": 1}, func(s *game.Score) bool { return !s.LeaveStatuses[0] }},
		// Balancing in auto also sets Leave, since the robot had to leave its safe house to reach the beam.
		{
			"auto balance sets leave", "auto_balance", m{"TeamPosition": 2},
			func(s *game.Score) bool { return s.AutoBalanceStatuses[1] && s.LeaveStatuses[1] },
		},
		{
			"leave cleared by hand", "leave", m{"TeamPosition": 2},
			func(s *game.Score) bool { return !s.LeaveStatuses[1] && s.AutoBalanceStatuses[1] },
		},
		{"leave set again", "leave", m{"TeamPosition": 2}, func(s *game.Score) bool { return s.LeaveStatuses[1] }},
		{
			"clearing auto balance keeps leave", "auto_balance", m{"TeamPosition": 2},
			func(s *game.Score) bool { return !s.AutoBalanceStatuses[1] && s.LeaveStatuses[1] },
		},
		{
			"endgame", "endgame", m{"TeamPosition": 3, "Value": int(game.EndgameBalance)},
			func(s *game.Score) bool { return s.EndgameStatuses[2] == game.EndgameBalance },
		},
		{"no station 4", "leave", m{"TeamPosition": 4}, nil},
		{"endgame value out of range", "endgame", m{"TeamPosition": 1, "Value": 3}, nil},
		{"toss", "toss", nil, func(s *game.Score) bool { return s.Toss }},
	}
	for _, step := range steps {
		sender := redFarWs
		if step.command == "treasure" {
			sender = redNearWs
		}
		sender.Write(step.command, step.data)
		if step.check != nil {
			for _, ws := range []*websocket.Websocket{redNearWs, redFarWs, blueWs} {
				readWebsocketType(t, ws, "realtimeScore")
			}
			assert.True(t, step.check(&web.arena.RedRealtimeScore.CurrentScore), step.name)
		}
	}

	// Test that some invalid commands do nothing and don't result in score change notifications.
	redNearWs.Write("invalid", nil)

	// Test committing logic; the alliance isn't ready until both its near and far panels have committed.
	redNearWs.Write("commitMatch", nil)
	readWebsocketType(t, redNearWs, "error")
	blueWs.Write("commitMatch", nil)
	readWebsocketType(t, blueWs, "error")
	assert.Equal(t, 0, registry.GetNumScoreCommitted("red"))
	assert.Equal(t, 0, registry.GetNumScoreCommitted("blue"))
	web.arena.MatchState = field.PostMatch

	redNearWs.Write("commitMatch", nil)
	time.Sleep(time.Millisecond * 10)
	assert.Equal(t, 1, registry.GetNumScoreCommitted("red"))
	assert.Less(t, registry.GetNumScoreCommitted("red"), registry.GetNumPanels("red"))

	redFarWs.Write("commitMatch", nil)
	blueWs.Write("commitMatch", nil)
	time.Sleep(time.Millisecond * 10) // Allow some time for the commands to be processed.
	assert.Equal(t, 2, registry.GetNumScoreCommitted("red"))
	assert.Equal(t, 1, registry.GetNumScoreCommitted("blue"))

	// Load another match to reset the results.
	web.arena.ResetMatch()
	web.arena.LoadTestMatch()
	readWebsocketType(t, redNearWs, "matchLoad")
	readWebsocketType(t, redNearWs, "realtimeScore")
	readWebsocketType(t, redFarWs, "matchLoad")
	readWebsocketType(t, redFarWs, "realtimeScore")
	readWebsocketType(t, blueWs, "matchLoad")
	readWebsocketType(t, blueWs, "realtimeScore")
	assert.Equal(t, field.NewRealtimeScore(), web.arena.RedRealtimeScore)
	assert.Equal(t, field.NewRealtimeScore(), web.arena.BlueRealtimeScore)
	assert.Equal(t, 0, registry.GetNumScoreCommitted("red"))
	assert.Equal(t, 0, registry.GetNumScoreCommitted("blue"))
}

func TestScoringPanelWebsocketRejectsThirdRobotInTwoVsTwoMode(t *testing.T) {
	web := setupTestWeb(t)
	web.arena.EventSettings.TwoVsTwoMode = true

	server, wsUrl := web.startTestServer()
	defer server.Close()
	redFarWs := dialScoringPanel(t, wsUrl, "red_far")
	readInitialScoringPanelMessages(t, redFarWs)

	// Station 3 is not in play in 2v2 mode, so the command should be silently rejected.
	redFarWs.Write("leave", struct{ TeamPosition int }{3})
	redFarWs.Write("leave", struct{ TeamPosition int }{1})
	readWebsocketType(t, redFarWs, "realtimeScore")
	assert.False(t, web.arena.RedRealtimeScore.CurrentScore.LeaveStatuses[2])
	assert.True(t, web.arena.RedRealtimeScore.CurrentScore.LeaveStatuses[0])
}

func dialScoringPanel(t *testing.T, wsUrl, position string) *websocket.Websocket {
	conn, _, err := gorillawebsocket.DefaultDialer.Dial(wsUrl+"/panels/scoring/"+position+"/websocket", nil)
	assert.Nil(t, err)
	t.Cleanup(func() { conn.Close() })
	return websocket.NewTestWebsocket(conn)
}

func readInitialScoringPanelMessages(t *testing.T, ws *websocket.Websocket) {
	for _, messageType := range []string{"resetLocalState", "matchLoad", "matchTime", "realtimeScore"} {
		readWebsocketType(t, ws, messageType)
	}
}
