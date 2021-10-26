package main

import (
	"college-football-sim/router"
)

func main() {
	router := router.InitRouter()
	router.Run("localhost:8080")
}
