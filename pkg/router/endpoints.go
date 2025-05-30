package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/roolps/lootlocker-sudoku/backend/pkg/lootlocker"
)

type logoutEndpoint struct{}

func (e *logoutEndpoint) Post(s *session, w http.ResponseWriter, raw []byte) *apiresponse {
	c := &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		Domain:   origin,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, c)
	return statusok(nil)
}

type stateEndpoint struct {
}

type cell struct {
	Value     int    `json:"value,omitempty"`
	Immutable bool   `json:"immutable,omitempty"`
	Pencil    []bool `json:"pencil,omitempty"`
}

func (e *stateEndpoint) Get(s *session, w http.ResponseWriter) *apiresponse {
	metadata, err := lootlockerClient.ListMetadata(lootlocker.MetadataSourcePlayer, s.PlayerID, s.Token)
	if err != nil {
		if err == lootlocker.ErrForbidden {
			// return a 403 error as the login token is no longer valid
			c := &http.Cookie{
				Name:     "session",
				Value:    "",
				Path:     "/",
				Domain:   origin,
				Expires:  time.Unix(0, 0),
				MaxAge:   -1,
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			}
			http.SetCookie(w, c)

			return statusForbidden(err.Error())
		}
		return statusinternalservererror(err.Error())
	}
	if gameState, ok := metadata["current_state"]; ok {
		start, ok := metadata["start_time"]
		if !ok {
			starttime := time.Now().Unix()
			start = lootlocker.Metadata{
				Access: []string{
					"game_api.read",
					"game_api.write",
				},
				Key:    "start_time",
				Tags:   []string{},
				Value:  starttime,
				Type:   lootlocker.MetadataTypeNumber,
				Action: lootlocker.MetadataActionCreate,
			}
			if err := lootlockerClient.UpdatePlayerMetadata(s.Token, []lootlocker.Metadata{start}); err != nil {
				return statusinternalservererror(err.Error())
			}
		}
		return statusok(map[string]any{"state": gameState.Value, "start_time": start.Value})
	}

	// return empty game state as there is none active
	// this will force user to use the menu
	return statusok(map[string]any{"state": [][]cell{}, "start_time": 0})
}

func (e *stateEndpoint) Post(s *session, w http.ResponseWriter, raw []byte) *apiresponse {
	var gameState [][]cell
	if err := json.Unmarshal(raw, &gameState); err != nil {
		return statusinternalservererror(err.Error())
	}

	// update current game state
	if err := lootlockerClient.UpdatePlayerMetadata(s.Token, []lootlocker.Metadata{{
		Access: []string{
			"game_api.read",
			"game_api.write",
		},
		Key:    "current_state",
		Tags:   []string{},
		Value:  gameState,
		Type:   lootlocker.MetadataTypeJSON,
		Action: lootlocker.MetadataActionUpdate,
	}}); err != nil {
		return statusinternalservererror(err.Error())
	}
	return statusok(nil)
}

func (e *stateEndpoint) Delete(s *session, w http.ResponseWriter, raw []byte) *apiresponse {
	if err := lootlockerClient.UpdatePlayerMetadata(s.Token, []lootlocker.Metadata{{
		Key:    "current_state",
		Action: lootlocker.MetadataActionDelete,
	}, {
		Key:    "start_time",
		Action: lootlocker.MetadataActionDelete,
	}}); err != nil {
		return statusinternalservererror(err.Error())
	}
	return statusok(nil)
}

type gameEndpoint struct {
	Difficulty string `json:"difficulty"`
}

func (e *gameEndpoint) Post(s *session, w http.ResponseWriter, raw []byte) *apiresponse {
	if err := json.Unmarshal(raw, e); err != nil {
		return statusnotacceptable(err.Error())
	}
	switch e.Difficulty {
	case "easy":
	default:
		return statusnotacceptable(fmt.Sprintf("invalid game difficulty %v", e.Difficulty))
	}

	gamestate := [][]cell{}
	raw, err := os.ReadFile(fmt.Sprintf("%v/example_games/%v.json", wd, e.Difficulty))
	if err != nil {
		return statusinternalservererror(err.Error())
	}
	if err := json.Unmarshal(raw, &gamestate); err != nil {
		return statusinternalservererror(err.Error())
	}

	// create current game state
	if err := lootlockerClient.UpdatePlayerMetadata(s.Token, []lootlocker.Metadata{
		{
			Access: []string{
				"game_api.read",
				"game_api.write",
			},
			Key:    "current_state",
			Tags:   []string{},
			Value:  gamestate,
			Type:   lootlocker.MetadataTypeJSON,
			Action: lootlocker.MetadataActionCreate,
		},
	}); err != nil {
		return statusinternalservererror(err.Error())
	}

	// create start time metadata
	startTime := time.Now().Unix()
	if err := lootlockerClient.UpdatePlayerMetadata(s.Token, []lootlocker.Metadata{
		{
			Access: []string{
				"game_api.read",
				"game_api.write",
			},
			Key:    "start_time",
			Tags:   []string{},
			Value:  startTime,
			Type:   lootlocker.MetadataTypeNumber,
			Action: lootlocker.MetadataActionCreate,
		},
	}); err != nil {
		return statusinternalservererror(err.Error())
	}
	return statusok(gamestate)
}

func (e *gameEndpoint) Delete(s *session, w http.ResponseWriter, raw []byte) *apiresponse {
	metadata, err := lootlockerClient.ListMetadata(lootlocker.MetadataSourcePlayer, s.PlayerID, s.Token)
	if err != nil {
		if err == lootlocker.ErrForbidden {
			// return a 403 error as the login token is no longer valid
			c := &http.Cookie{
				Name:     "session",
				Value:    "",
				Path:     "/",
				Domain:   origin,
				Expires:  time.Unix(0, 0),
				MaxAge:   -1,
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			}
			http.SetCookie(w, c)

			return statusForbidden(err.Error())
		}
		return statusinternalservererror(err.Error())
	}
	if gameState, ok := metadata["current_state"]; ok {
		raw, err := json.Marshal(gameState.Value)
		if err != nil {
			return statusinternalservererror(err.Error())
		}
		cells := gameboard{}
		if err := json.Unmarshal(raw, &cells); err != nil {
			return statusinternalservererror(err.Error())
		}

		if !cells.verifyGameComplete() {
			return statusnotacceptable("puzzle is incomplete")
		}

		wall, err := lootlockerClient.GetWalletForHolder(s.Token, s.PlayerID)
		if err != nil {
			return statusinternalservererror(err.Error())
		}
		// add reward to player wallet
		if err := lootlockerClient.CreditBalance(s.Token, &lootlocker.Credit{Amount: "50", WalletID: wall.ID, CurrencyID: LOOTLOCKER_CURRENCY_ID}); err != nil {
			return statusinternalservererror(err.Error())
		}

		// delete metadata state + start time
		if err := lootlockerClient.UpdatePlayerMetadata(s.Token, []lootlocker.Metadata{{
			Key:    "current_state",
			Action: lootlocker.MetadataActionDelete,
		}, {
			Key:    "start_time",
			Action: lootlocker.MetadataActionDelete,
		}}); err != nil {
			return statusinternalservererror(err.Error())
		}
		return statusok(nil)
	}
	return statusnotacceptable("no game state was found for the user")
}

type gameboard [][]cell

func (gb gameboard) verifyGameComplete() bool {
	isValidGroup := func(group []int) bool {
		if len(group) != 9 {
			return false
		}
		expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
		sort.Ints(group)
		for i := range group {
			if group[i] != expected[i] {
				return false
			}
		}
		return true
	}

	rows := make([][]int, 9)
	cols := make([][]int, 9)

	for i := 0; i < 9; i++ {
		rows[i] = make([]int, 0, 9)
		cols[i] = make([]int, 0, 9)
	}

	for i := range gb {
		for j := range gb[i] {
			val := gb[i][j].Value
			rows[i] = append(rows[i], val)
			cols[j] = append(cols[j], val)
		}
	}

	allRowsValid := true
	for _, row := range rows {
		if !isValidGroup(row) {
			allRowsValid = false
			break
		}
	}

	allColsValid := true
	for _, col := range cols {
		if !isValidGroup(col) {
			allColsValid = false
			break
		}
	}
	return allRowsValid && allColsValid
}

type playerEndpoint struct {
	Balance float64 `json:"balance"`
}

func (u *playerEndpoint) Get(s *session, w http.ResponseWriter) *apiresponse {
	wall, err := lootlockerClient.GetWalletForHolder(s.Token, s.PlayerID)
	if err != nil {
		logger.Error(err)
		return statusinternalservererror(err.Error())
	}
	// add reward to player wallet
	balances, err := lootlockerClient.GetBalances(s.Token, wall.ID)
	if err != nil {
		logger.Error(err)
		return statusinternalservererror(err.Error())
	}
	u.Balance = balances[LOOTLOCKER_CURRENCY_ID].Amount
	return statusok(u)
}
