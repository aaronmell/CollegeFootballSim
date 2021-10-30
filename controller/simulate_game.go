package controller

import (
	"college-football-sim/models"
	"college-football-sim/utils"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type SimulateGameController struct{}

type GameClient interface {
	RunPlay(g *models.Game, r *rand.Rand) models.GameLog
}

type client struct {
}

func (sgc SimulateGameController) SimulateGame(c *gin.Context) {
	reqBody := models.SimulateGameRequest{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Error().Msg(err.Error())
		utils.ResponseFormat(c, utils.RequestParamError, nil)
		return
	}

	response := models.SimulateGameResponse{}
	game := reqBody.InitGame()

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	client := client{}

	for ok := true; ok; ok = game.GameOver {
		runStep(client, &game, r1)
	}

	c.JSON(utils.Success.Status, response)
}

func runStep(c GameClient, g *models.Game, r1 *rand.Rand) {
	event := c.RunPlay(g, r1)
	g.GameLog = append(g.GameLog, event)

	if GameOver(g) {
		return
	}

	if AdvanceQuarter(g) {
		return
	}

	if event.Event == models.PlayRan {
		IncrementDown(g)
	}
}

func (c client) RunPlay(g *models.Game, r *rand.Rand) models.GameLog {

	return models.GameLog{
		Event:        models.PlayRan,
		EventEndTime: g.GameClockInSeconds - (r.Intn(37) + 5),
	}
}

func GameOver(g *models.Game) bool {
	if g.GameClockInSeconds != 0 {
		return false
	}

	if g.Quarter != models.Fourth {
		return false
	}

	if g.AwayTeamScore == g.HomeTeamScore {
		return false
	}

	g.GameOver = true
	g.GameLog = append(g.GameLog, models.GameLog{
		Event: models.EndOfGame,
	})

	return true
}

func AdvanceQuarter(g *models.Game) bool {
	if g.GameClockInSeconds != 0 {
		return false
	}

	switch g.Quarter {
	case models.Second:
		g.CurrentPossession = g.PossessionAtHalf
		g.CurrentDown = models.FirstDown
	case models.Overtime1:
		return false
	default:
	}

	g.Quarter++

	if g.Quarter == models.Overtime1 {
		g.CurrentDown = models.FirstDown
		g.GameLog = append(g.GameLog, models.GameLog{
			Event: models.EndofRegulationPlay,
		})

		return true
	}

	updateEndZone(g)
	g.GameClockInSeconds = models.QuarterTime
	g.GameLog = append(g.GameLog, models.GameLog{
		Event: models.EndofQuarter,
	})
	return true
}

func IncrementDown(g *models.Game) {
	switch g.CurrentDown {
	case models.FirstDown:
		g.CurrentDown = models.SecondDown
	case models.SecondDown:
		g.CurrentDown = models.ThirdDown
	case models.ThirdDown:
		g.CurrentDown = models.FourthDown
	case models.FourthDown:
		g.CurrentDown = models.FirstDown
		changePosession(g)
	}
}

func updateEndZone(g *models.Game) {
	switch g.HomeTeamEndZone {
	case models.WestEndZone:
		g.HomeTeamEndZone = models.EastEndZone
		g.AwayTeamEndZone = models.WestEndZone
	case models.EastEndZone:
		g.HomeTeamEndZone = models.WestEndZone
		g.AwayTeamEndZone = models.EastEndZone
	}
}

func changePosession(g *models.Game) {
	if g.CurrentPossession == models.HomeTeam {
		g.CurrentPossession = models.AwayTeam
	} else {
		g.CurrentPossession = models.HomeTeam
	}
}
