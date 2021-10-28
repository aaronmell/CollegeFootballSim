package controller

import (
	"college-football-sim/utils"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type SimulateGameController struct{}

type Team struct {
	Name    string `json:"Name" binding:"required"`
	Overall int    `json:"Overall" binding:"required"`
}

type SimulateGameRequest struct {
	HomeTeam Team `json:"HomeTeam" binding:"required"`
	AwayTeam Team `json:"AwayTeam" binding:"required"`
}

type SimulateGameResponse struct {
	Winner Team `json:"Winner" binding:"required"`
	Loser  Team `json:"Loser" binding:"required"`
}

type GameStatus struct {
	HomeTeam           Team
	AwayTeam           Team
	GameOver           bool
	Quarter            Quarter
	GameClockInSeconds int
	HomeTeamScore      int
	AwayTeamScore      int
	GameLog            []GameLog
}

type GameEvent string

const (
	GameStart           GameEvent = "GameStart"
	EndOfGame           GameEvent = "EndOfGame"
	EndofQuarter        GameEvent = "EndOfQuater"
	EndofRegulationPlay GameEvent = "EndofRegulationPlay"
	PlayRan             GameEvent = "PlayRan"
	QuarterTime         int       = 900
)

type GameLog struct {
	Event        GameEvent
	EventEndTime int
}

type Quarter int

const (
	First Quarter = iota
	Second
	Third
	Fourth
	Overtime1
	Overtime2
	Overtime3
	Overtime4
	Overtime5
	Overtime6
	Overtime7
	Overtime8
	Overtime9
	Overtime10
	Overtime11
	Overtime12
	Overtime13
	Overtime14
	Overtime15
)

func (sgc SimulateGameController) SimulateGame(c *gin.Context) {
	reqBody := SimulateGameRequest{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Error().Msg(err.Error())
		utils.ResponseFormat(c, utils.RequestParamError, nil)
		return
	}

	response := SimulateGameResponse{}
	gameStatus := initializeGameStatus(reqBody)

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for ok := true; ok; ok = gameStatus.GameOver {
		runStep(&gameStatus, r1)
	}

	c.JSON(utils.Success.Status, response)
}

func runStep(gameStatus *GameStatus, r1 *rand.Rand) {
	event := runPlay(gameStatus, r1)
	gameStatus.GameLog = append(gameStatus.GameLog, event)

	if gameOver(gameStatus) {
		return
	}

	if advanceQuarter(gameStatus) {
		return
	}
}

func gameOver(gs *GameStatus) bool {
	if gs.GameClockInSeconds != 0 {
		return false
	}

	if gs.Quarter != Fourth {
		return false
	}

	if gs.AwayTeamScore == gs.HomeTeamScore {
		return false
	}

	gs.GameOver = true
	gs.GameLog = append(gs.GameLog, GameLog{
		Event: EndOfGame,
	})

	return true
}

func advanceQuarter(gs *GameStatus) bool {
	if gs.GameClockInSeconds != 0 {
		return false
	}

	if gs.Quarter > Fourth {
		return false
	}

	gs.Quarter++

	if gs.Quarter > Fourth {
		gs.GameLog = append(gs.GameLog, GameLog{
			Event: EndofRegulationPlay,
		})

		return true
	}

	gs.GameClockInSeconds = QuarterTime
	gs.GameLog = append(gs.GameLog, GameLog{
		Event: EndofQuarter,
	})
	return true
}

func runPlay(gs *GameStatus, r *rand.Rand) GameLog {
	return GameLog{
		Event:        PlayRan,
		EventEndTime: gs.GameClockInSeconds - (r.Intn(37) + 5),
	}
}

func initializeGameStatus(reqBody SimulateGameRequest) GameStatus {
	return GameStatus{
		HomeTeam:           reqBody.HomeTeam,
		AwayTeam:           reqBody.AwayTeam,
		GameOver:           false,
		Quarter:            First,
		GameClockInSeconds: QuarterTime,
	}
}
