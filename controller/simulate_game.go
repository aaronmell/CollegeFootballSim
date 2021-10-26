package controller

import (
	"college-football-sim/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type SimulateGameController struct{}

type Team struct {
	Name    string `json:"Name" binding:"required"`
	Overall int    `json:"Overall" binding:"required"`
}

type SimulateGameRequest struct {
	Team1 Team `json:"Team1" binding:"required"`
	Team2 Team `json:"Team2" binding:"required"`
}

type SimulateGameResponse struct {
	Winner Team `json:"Team1" binding:"required"`
	Loser  Team `json:"Team2" binding:"required"`
}

func (sgc SimulateGameController) SimulateGame(c *gin.Context) {
	reqBody := SimulateGameRequest{}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Error().Msg(err.Error())
		utils.ResponseFormat(c, utils.RequestParamError, nil)
		return
	}

	response := SimulateGameResponse{}

	if reqBody.Team1.Overall > reqBody.Team2.Overall {
		response.Winner = reqBody.Team1
		response.Loser = reqBody.Team2
	} else {
		response.Winner = reqBody.Team2
		response.Loser = reqBody.Team1
	}

	c.JSON(utils.Success.Status, response)
}
