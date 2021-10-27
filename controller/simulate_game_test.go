package controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"college-football-sim/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
)

func TestSimulateGame(t *testing.T) {
	r := gin.Default()
	sgc := new(SimulateGameController)

	r.POST("/", sgc.SimulateGame)

	simulateGameRequest := &SimulateGameRequest{
		HomeTeam: Team{
			Name:    "Team1",
			Overall: 51,
		},
		AwayTeam: Team{
			Name:    "Team2",
			Overall: 52,
		},
	}

	requestByte, _ := json.Marshal(simulateGameRequest)
	requestReader := bytes.NewReader(requestByte)
	req, _ := http.NewRequest("POST", "/", requestReader)

	utils.TestHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
		statusOK := w.Code == http.StatusOK

		p, err := ioutil.ReadAll(w.Body)

		if err != nil {
			t.Fail()
		}

		response := SimulateGameResponse{}
		err = json.Unmarshal(p, &response)

		if err != nil {
			t.Fail()
		}

		return statusOK && response.Winner.Name == "Team2"
	})
}

func TestRunPlay(t *testing.T) {
	tests := map[string]struct {
		input GameStatus
		want  GameStatus
	}{
		"GameOverNotTied": {
			input: GameStatus{
				GameClockInSeconds: 0,
				GameOver:           false,
				Quarter:            Fourth,
				HomeTeamScore:      7,
				AwayTeamScore:      0,
			},
			want: GameStatus{
				GameClockInSeconds: 0,
				GameOver:           true,
				Quarter:            Fourth,
				HomeTeamScore:      7,
				AwayTeamScore:      0,
			},
		},
		"NextQuarterWhenTime0": {
			input: GameStatus{
				GameClockInSeconds: 0,
				GameOver:           false,
				Quarter:            Third,
				HomeTeamScore:      7,
				AwayTeamScore:      0,
			},
			want: GameStatus{
				GameClockInSeconds: 0,
				GameOver:           false,
				Quarter:            Fourth,
				HomeTeamScore:      7,
				AwayTeamScore:      0,
			},
		},
		"OvertimeWhenScoreTied": {
			input: GameStatus{
				GameClockInSeconds: 0,
				GameOver:           false,
				Quarter:            Fourth,
				HomeTeamScore:      7,
				AwayTeamScore:      7,
			},
			want: GameStatus{
				GameClockInSeconds: 0,
				GameOver:           false,
				Quarter:            Overtime1,
				HomeTeamScore:      7,
				AwayTeamScore:      7,
			},
		},
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.input.HomeTeam = Team{
				Name: "Team1",
			}
			tc.input.AwayTeam = Team{
				Name: "Team2",
			}

			tc.want.HomeTeam = Team{
				Name: "Team1",
			}
			tc.want.AwayTeam = Team{
				Name: "Team2",
			}

			runPlay(&tc.input, r1)
			diff := cmp.Diff(tc.want, tc.input)

			if diff != "" {
				t.Fatalf(diff)
			}

		})
	}

	// t.Run("Returns Gameover when time is 0 and quarter is fourth and score is not tied", func(t *testing.T) {
	// 	gameStatus := &

	// 	s1 := rand.NewSource(1)
	// 	r1 := rand.New(s1)
	// 	runPlay(gameStatus, r1)

	// 	if gameStatus.GameOver != true {
	// 		t.Error("Expected Gameover to be true, recieved false")
	// 	}

	// 	if gameStatus.Quarter != Fourth {
	// 		t.Errorf("Expected quarter %d recieved %d", Fourth, gameStatus.Quarter)
	// 	}
	// })

	// t.Run("Increments Quarter when time is 0, and quarter is less than the fourth or when quarter is 4 or greater and score is tied", func(t *testing.T) {
	// 	cases := []struct {
	// 		Name          string
	// 		StartQuarter  Quarter
	// 		EndQuarter    Quarter
	// 		HomeTeamScore int
	// 		AwayTeamScore int
	// 	}{
	// 		{First, Second},
	// 		{Second, Third},
	// 		{Third, Fourth},
	// 		{Fourth, Overtime1},
	// 	}

	// 	for _, test := range cases {
	// 		gameStatus := &GameStatus{
	// 			GameClockInSeconds: 0,
	// 			GameOver:           false,
	// 			Quarter:            test.StartQuarter,
	// 			HomeTeam: Team{
	// 				Name: "Team1",
	// 			},
	// 			AwayTeam: Team{
	// 				Name: "Team2",
	// 			},
	// 			HomeTeamScore: 0,
	// 			AwayTeamScore: 0,
	// 		}

	// 		s1 := rand.NewSource(1)
	// 		r1 := rand.New(s1)
	// 		runPlay(gameStatus, r1)

	// 		if gameStatus.GameOver == true {
	// 			t.Errorf("Expected gameover to be false in quarter %d, was true", test)
	// 		}

	// 		if gameStatus.Quarter != test.EndQuarter {
	// 			t.Errorf("Expected Quarter to be %d was %d", test.EndQuarter, test.StartQuarter)
	// 		}

	// 	}
	// })
}
