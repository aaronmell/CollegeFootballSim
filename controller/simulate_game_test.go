package controller

import (
	"college-football-sim/controller/mocks"
	"college-football-sim/models"
	"math/rand"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		input models.Game
		want  models.Game
	}{
		"EndOfFirstQuarter": {
			input: models.Game{
				GameClockInSeconds: 0,
				GameOver:           false,
				Quarter:            models.First,
				HomeTeamScore:      7,
				AwayTeamScore:      0,
				HomeTeamEndZone:    models.WestEndZone,
				AwayTeamEndZone:    models.EastEndZone,
				CurrentPossession:  models.HomeTeam,
				PossessionAtHalf:   models.AwayTeam,
				CurrentDown:        models.SecondDown,
			},
			want: models.Game{
				GameClockInSeconds: 900,
				GameOver:           false,
				Quarter:            models.Second,
				HomeTeamScore:      7,
				AwayTeamScore:      0,
				HomeTeamEndZone:    models.EastEndZone,
				AwayTeamEndZone:    models.WestEndZone,
				CurrentPossession:  models.HomeTeam,
				PossessionAtHalf:   models.AwayTeam,
				CurrentDown:        models.SecondDown,
				GameLog: []models.GameLog{
					{
						Event: models.PlayRan,
					},
					{
						Event: models.EndofQuarter,
					},
				},
			},
		},
		"EndOfSecondQuarter": {
			input: models.Game{
				GameClockInSeconds: 0,
				GameOver:           false,
				Quarter:            models.Second,
				HomeTeamScore:      7,
				AwayTeamScore:      0,
				HomeTeamEndZone:    models.EastEndZone,
				AwayTeamEndZone:    models.WestEndZone,
				CurrentPossession:  models.HomeTeam,
				PossessionAtHalf:   models.AwayTeam,
				CurrentDown:        models.SecondDown,
			},
			want: models.Game{
				GameClockInSeconds: 900,
				GameOver:           false,
				Quarter:            models.Third,
				HomeTeamScore:      7,
				AwayTeamScore:      0,
				HomeTeamEndZone:    models.WestEndZone,
				AwayTeamEndZone:    models.EastEndZone,
				CurrentPossession:  models.AwayTeam,
				PossessionAtHalf:   models.AwayTeam,
				CurrentDown:        models.FirstDown,
				GameLog: []models.GameLog{
					{
						Event: models.PlayRan,
					},
					{
						Event: models.EndofQuarter,
					},
				},
			},
		},
		"EndOfThirdQuarter": {
			input: models.Game{
				GameClockInSeconds: 0,
				GameOver:           false,
				Quarter:            models.Third,
				HomeTeamScore:      7,
				AwayTeamScore:      0,
				HomeTeamEndZone:    models.WestEndZone,
				AwayTeamEndZone:    models.EastEndZone,
				CurrentPossession:  models.HomeTeam,
				PossessionAtHalf:   models.AwayTeam,
				CurrentDown:        models.SecondDown,
			},
			want: models.Game{
				GameClockInSeconds: 900,
				GameOver:           false,
				Quarter:            models.Fourth,
				HomeTeamScore:      7,
				AwayTeamScore:      0,
				HomeTeamEndZone:    models.EastEndZone,
				AwayTeamEndZone:    models.WestEndZone,
				CurrentPossession:  models.HomeTeam,
				PossessionAtHalf:   models.AwayTeam,
				CurrentDown:        models.SecondDown,

				GameLog: []models.GameLog{
					{
						Event: models.PlayRan,
					},
					{
						Event: models.EndofQuarter,
					},
				},
			},
		},
		"EndOfFourthQuarterNotTied": {
			input: models.Game{
				GameClockInSeconds: 0,
				GameOver:           false,
				Quarter:            models.Fourth,
				HomeTeamScore:      7,
				AwayTeamScore:      0,
				HomeTeamEndZone:    models.WestEndZone,
				AwayTeamEndZone:    models.EastEndZone,
				CurrentPossession:  models.HomeTeam,
				PossessionAtHalf:   models.AwayTeam,
				CurrentDown:        models.SecondDown,
			},
			want: models.Game{
				GameClockInSeconds: 0,
				GameOver:           true,
				Quarter:            models.Fourth,
				HomeTeamScore:      7,
				AwayTeamScore:      0,
				HomeTeamEndZone:    models.WestEndZone,
				AwayTeamEndZone:    models.EastEndZone,
				CurrentPossession:  models.HomeTeam,
				PossessionAtHalf:   models.AwayTeam,
				CurrentDown:        models.SecondDown,
				GameLog: []models.GameLog{
					{
						Event: models.PlayRan,
					},
					{
						Event: models.EndOfGame,
					},
				},
			},
		},
		"EndOfFourthQuarterTied": {
			input: models.Game{
				GameClockInSeconds: 0,
				GameOver:           false,
				Quarter:            models.Fourth,
				HomeTeamScore:      7,
				AwayTeamScore:      7,
				HomeTeamEndZone:    models.WestEndZone,
				AwayTeamEndZone:    models.EastEndZone,
				CurrentPossession:  models.HomeTeam,
				PossessionAtHalf:   models.AwayTeam,
				CurrentDown:        models.SecondDown,
			},
			want: models.Game{
				GameClockInSeconds: 0,
				GameOver:           false,
				Quarter:            models.Overtime1,
				HomeTeamScore:      7,
				AwayTeamScore:      7,
				HomeTeamEndZone:    models.WestEndZone,
				AwayTeamEndZone:    models.EastEndZone,
				CurrentPossession:  models.HomeTeam,
				PossessionAtHalf:   models.AwayTeam,
				CurrentDown:        models.FirstDown,
				GameLog: []models.GameLog{
					{
						Event: models.PlayRan,
					},
					{
						Event: models.EndofRegulationPlay,
					},
				},
			},
		},
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for name, testData := range tests {
		t.Run(name, func(t *testing.T) {

			mockClient := new(mocks.GameClient)
			mockClient.On("RunPlay", mock.AnythingOfType("*models.Game"), mock.AnythingOfType("*rand.Rand")).Return(models.GameLog{
				Event: models.PlayRan,
			}).Once()

			runStep(mockClient, &testData.input, r1)

			diff := cmp.Diff(testData.want, testData.input, cmpopts.IgnoreFields(models.GameLog{}, "EventEndTime"))

			assert.Empty(t, diff)
		})
	}
}
