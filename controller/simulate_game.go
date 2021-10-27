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
		runPlay(&gameStatus, r1)
	}

	if reqBody.HomeTeam.Overall > reqBody.AwayTeam.Overall {
		response.Winner = reqBody.HomeTeam
		response.Loser = reqBody.AwayTeam
	} else {
		response.Winner = reqBody.AwayTeam
		response.Loser = reqBody.HomeTeam
	}

	c.JSON(utils.Success.Status, response)
}

func runPlay(gs *GameStatus, r *rand.Rand) {
	if gs.Quarter <= Fourth {
		runRegulationPlay(gs, r)
		return
	}

}

func runRegulationPlay(gs *GameStatus, r *rand.Rand) {
	if gs.GameClockInSeconds == 0 {
		if gs.Quarter == Fourth {
			if gs.HomeTeamScore != gs.AwayTeamScore {
				gs.GameOver = true
				return
			}
		}
		gs.Quarter++
		return
	}
}

func initializeGameStatus(reqBody SimulateGameRequest) GameStatus {
	return GameStatus{
		HomeTeam:           reqBody.HomeTeam,
		AwayTeam:           reqBody.AwayTeam,
		GameOver:           false,
		Quarter:            First,
		GameClockInSeconds: 900,
	}
}
