package controller

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"college-football-sim/utils"

	"github.com/gin-gonic/gin"
)

func TestSimulateGame(t *testing.T) {
	r := gin.Default()
	sgc := new(SimulateGameController)

	r.POST("/", sgc.SimulateGame)

	simulateGameRequest := &SimulateGameRequest{
		Team1: Team{
			Name:    "Team1",
			Overall: 51,
		},
		Team2: Team{
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
