package controller

import (
	"math/rand"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// func TestSimulateGame(t *testing.T) {
// 	r := gin.Default()
// 	sgc := new(SimulateGameController)

// 	r.POST("/", sgc.SimulateGame)

// 	simulateGameRequest := &SimulateGameRequest{
// 		HomeTeam: Team{
// 			Name:    "Team1",
// 			Overall: 51,
// 		},
// 		AwayTeam: Team{
// 			Name:    "Team2",
// 			Overall: 52,
// 		},
// 	}

// 	requestByte, _ := json.Marshal(simulateGameRequest)
// 	requestReader := bytes.NewReader(requestByte)
// 	req, _ := http.NewRequest("POST", "/", requestReader)

// 	utils.TestHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
// 		statusOK := w.Code == http.StatusOK

// 		p, err := ioutil.ReadAll(w.Body)

// 		if err != nil {
// 			t.Fail()
// 		}

// 		response := SimulateGameResponse{}
// 		err = json.Unmarshal(p, &response)

// 		if err != nil {
// 			t.Fail()
// 		}

// 		return statusOK && response.Winner.Name == "Team2"
// 	})
// }

func TestRunStep(t *testing.T) {
	tests := map[string]struct {
		input GameStatus
		want  GameStatus
	}{
		"EndOfFirstQuarter": {
			input: GameStatus{
				GameClockInSeconds: 0,
				GameOver:           false,
				Quarter:            First,
				HomeTeamScore:      7,
				AwayTeamScore:      0,
			},
			want: GameStatus{
				GameClockInSeconds: 900,
				GameOver:           false,
				Quarter:            Second,
				HomeTeamScore:      7,
				AwayTeamScore:      0,
				GameLog: []GameLog{
					{
						Event: PlayRan,
					},
					{
						Event: EndofQuarter,
					},
				},
			},
		},
		"EndOfSecondQuarter": {
			input: GameStatus{
				GameClockInSeconds: 0,
				GameOver:           false,
				Quarter:            Second,
				HomeTeamScore:      7,
				AwayTeamScore:      0,
			},
			want: GameStatus{
				GameClockInSeconds: 900,
				GameOver:           false,
				Quarter:            Third,
				HomeTeamScore:      7,
				AwayTeamScore:      0,
				GameLog: []GameLog{
					{
						Event: PlayRan,
					},
					{
						Event: EndofQuarter,
					},
				},
			},
		},
		"EndOfThirdQuarter": {
			input: GameStatus{
				GameClockInSeconds: 0,
				GameOver:           false,
				Quarter:            Third,
				HomeTeamScore:      7,
				AwayTeamScore:      0,
			},
			want: GameStatus{
				GameClockInSeconds: 900,
				GameOver:           false,
				Quarter:            Fourth,
				HomeTeamScore:      7,
				AwayTeamScore:      0,
				GameLog: []GameLog{
					{
						Event: PlayRan,
					},
					{
						Event: EndofQuarter,
					},
				},
			},
		},
		"EndOfFourthQuarterNotTied": {
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
				GameLog: []GameLog{
					{
						Event: PlayRan,
					},
					{
						Event: EndOfGame,
					},
				},
			},
		},
		"EndOfFourthQuarterTied": {
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
				GameLog: []GameLog{
					{
						Event: PlayRan,
					},
					{
						Event: EndofRegulationPlay,
					},
				},
			},
		},
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for name, testData := range tests {
		t.Run(name, func(t *testing.T) {
			runStep(&testData.input, r1)

			diff := cmp.Diff(testData.want, testData.input, cmpopts.IgnoreFields(GameLog{}, "EventEndTime"))

			if diff != "" {
				t.Fatalf(diff)
			}
		})
	}
}
