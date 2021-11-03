package controller

import (
	"college-football-sim/controller/mocks"
	"college-football-sim/models"
	"math/rand"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/shopspring/decimal"
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
		input      models.Game
		newYardage *int
		want       models.Game
	}{
		"EndOfFirstQuarter": {
			input: models.Game{
				GameClockInSeconds: 0,
				Quarter:            models.First,
				HomeTeamEndZone:    models.WestEndZone,
				AwayTeamEndZone:    models.EastEndZone,
			},
			want: models.Game{
				GameClockInSeconds: 900,
				Quarter:            models.Second,
				HomeTeamEndZone:    models.EastEndZone,
				AwayTeamEndZone:    models.WestEndZone,
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
				Quarter:            models.Second,
				HomeTeamEndZone:    models.EastEndZone,
				AwayTeamEndZone:    models.WestEndZone,
				CurrentPossession:  models.HomeTeam,
				PossessionAtHalf:   models.AwayTeam,
				CurrentDown:        models.SecondDown,
			},
			want: models.Game{
				GameClockInSeconds: 900,
				Quarter:            models.Third,
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
				Quarter:            models.Third,
				HomeTeamEndZone:    models.WestEndZone,
				AwayTeamEndZone:    models.EastEndZone,
			},
			want: models.Game{
				GameClockInSeconds: 900,
				Quarter:            models.Fourth,
				HomeTeamEndZone:    models.EastEndZone,
				AwayTeamEndZone:    models.WestEndZone,

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
			},
			want: models.Game{
				GameClockInSeconds: 0,
				GameOver:           true,
				Quarter:            models.Fourth,
				HomeTeamScore:      7,
				AwayTeamScore:      0,
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
				CurrentDown:        models.SecondDown,
			},
			want: models.Game{
				GameClockInSeconds: 0,
				GameOver:           false,
				Quarter:            models.Overtime1,
				HomeTeamScore:      7,
				AwayTeamScore:      7,
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
		"1stDown": {
			input: models.Game{
				GameClockInSeconds: 900,
				CurrentDown:        models.FirstDown,
			},
			want: models.Game{
				GameClockInSeconds: 900,
				CurrentDown:        models.SecondDown,
				GameLog: []models.GameLog{
					{
						Event: models.PlayRan,
					},
				},
			},
		},
		"2ndDown": {
			input: models.Game{
				GameClockInSeconds: 900,
				CurrentDown:        models.SecondDown,
			},
			want: models.Game{
				GameClockInSeconds: 900,
				CurrentDown:        models.ThirdDown,
				GameLog: []models.GameLog{
					{
						Event: models.PlayRan,
					},
				},
			},
		},
		"3rdDown": {
			input: models.Game{
				GameClockInSeconds: 900,
				CurrentDown:        models.ThirdDown,
			},
			want: models.Game{
				GameClockInSeconds: 900,
				CurrentDown:        models.FourthDown,
				GameLog: []models.GameLog{
					{
						Event: models.PlayRan,
					},
				},
			},
		},
		"4thDown": {
			input: models.Game{
				GameClockInSeconds: 900,
				CurrentPossession:  models.HomeTeam,
				CurrentDown:        models.FourthDown,
			},
			want: models.Game{
				GameClockInSeconds: 900,
				CurrentPossession:  models.AwayTeam,
				CurrentDown:        models.FirstDown,
				GameLog: []models.GameLog{
					{
						Event: models.PlayRan,
					},
				},
			},
		},
		"4thDownAwayHasPossession": {
			input: models.Game{
				GameClockInSeconds: 900,
				CurrentPossession:  models.AwayTeam,
				CurrentDown:        models.FourthDown,
			},
			want: models.Game{
				GameClockInSeconds: 900,
				CurrentPossession:  models.HomeTeam,
				CurrentDown:        models.FirstDown,
				GameLog: []models.GameLog{
					{
						Event: models.PlayRan,
					},
				},
			},
		},
		"4thDownGetsFirstDown": {
			input: models.Game{
				GameClockInSeconds: 900,
				CurrentPossession:  models.HomeTeam,
				CurrentDown:        models.FourthDown,
				FirstDownPosition:  decimal.NewFromInt(10),
				BallPosition:       decimal.NewFromInt(1),
			},
			newYardage: createInt(11),
			want: models.Game{
				GameClockInSeconds: 900,
				CurrentPossession:  models.HomeTeam,
				CurrentDown:        models.FirstDown,
				GameLog: []models.GameLog{
					{
						Event:      models.PlayRan,
						NewYardage: decimal.NewFromInt(11),
					},
				},
				BallPosition:      decimal.NewFromInt(11),
				FirstDownPosition: decimal.NewFromInt((21)),
			},
		},
		"FirstDownMax100": {
			input: models.Game{
				GameClockInSeconds: 900,
				CurrentPossession:  models.HomeTeam,
				CurrentDown:        models.ThirdDown,
				FirstDownPosition:  decimal.NewFromInt(92),
				BallPosition:       decimal.NewFromInt(90),
			},
			newYardage: createInt(93),
			want: models.Game{
				GameClockInSeconds: 900,
				CurrentPossession:  models.HomeTeam,
				CurrentDown:        models.FirstDown,
				GameLog: []models.GameLog{
					{
						Event:      models.PlayRan,
						NewYardage: decimal.NewFromInt(93),
					},
				},
				BallPosition:      decimal.NewFromInt(93),
				FirstDownPosition: decimal.NewFromInt((100)),
			},
		},
		"KickOffWhenSet": {
			input: models.Game{
				RequiresKickOff:   true,
				CurrentPossession: models.HomeTeam,
			},
			want: models.Game{
				RequiresKickOff: false,
				GameLog: []models.GameLog{
					{
						Event: models.KickOff,
					},
				},
				CurrentPossession: models.AwayTeam,
			},
		},
		"ExtraPointWhenSet": {
			input: models.Game{
				RequiresExtraPoint: true,
			},
			want: models.Game{
				RequiresExtraPoint: false,
				GameLog: []models.GameLog{
					{
						Event: models.ExtraPoint,
					},
				},
			},
		},
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for name, testData := range tests {
		t.Run(name, func(t *testing.T) {

			playRanGameLog := models.GameLog{
				Event: models.PlayRan,
			}

			if testData.newYardage != nil {
				playRanGameLog.NewYardage = decimal.NewFromInt(int64(*testData.newYardage))
			}

			mockClient := new(mocks.GameClient)
			mockClient.On("RunPlay", mock.AnythingOfType("*models.Game"), mock.AnythingOfType("*rand.Rand")).Return(playRanGameLog).Run(func(args mock.Arguments) {
				arg := args.Get(0).(*models.Game)
				arg.BallPosition = playRanGameLog.NewYardage
			})

			mockClient.On("RunKickOff", mock.AnythingOfType("*models.Game"), mock.AnythingOfType("*rand.Rand")).Return(models.GameLog{
				Event: models.KickOff,
			}).Run(func(args mock.Arguments) {
				arg := args.Get(0).(*models.Game)
				changePosession(arg)
				arg.RequiresKickOff = false
			})

			mockClient.On("RunExtraPoint", mock.AnythingOfType("*models.Game"), mock.AnythingOfType("*rand.Rand")).Return(models.GameLog{
				Event: models.ExtraPoint,
			}).Run(func(args mock.Arguments) {
				arg := args.Get(0).(*models.Game)
				arg.RequiresExtraPoint = false
			})

			runStep(mockClient, &testData.input, r1)

			diff := cmp.Diff(testData.want, testData.input, cmpopts.IgnoreFields(models.GameLog{}, "EventEndTime"))

			assert.Empty(t, diff)
		})
	}
}

func createInt(x int) *int {
	return &x
}
