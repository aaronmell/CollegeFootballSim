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
		"SecondQuarterWhenTime0": {
			input: GameStatus{
				GameClockInSeconds: 0,
				GameOver:           false,
				Quarter:            First,
				HomeTeamScore:      7,
				AwayTeamScore:      0,
			},
			want: GameStatus{
				GameClockInSeconds: 0,
				GameOver:           false,
				Quarter:            Second,
				HomeTeamScore:      7,
				AwayTeamScore:      0,
			},
		},
		"ThirdQuarterWhenTime0": {
			input: GameStatus{
				GameClockInSeconds: 0,
				GameOver:           false,
				Quarter:            Second,
				HomeTeamScore:      7,
				AwayTeamScore:      0,
			},
			want: GameStatus{
				GameClockInSeconds: 0,
				GameOver:           false,
				Quarter:            Third,
				HomeTeamScore:      7,
				AwayTeamScore:      0,
			},
		},
		"FourthQuarterWhenTime0": {
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
}
