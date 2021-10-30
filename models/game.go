package models

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

type Game struct {
	HomeTeam           Team
	AwayTeam           Team
	GameOver           bool
	Quarter            Quarter
	GameClockInSeconds int
	HomeTeamScore      int
	AwayTeamScore      int
	GameLog            []GameLog
	HomeTeamEndZone    EndZone
	AwayTeamEndZone    EndZone
	CurrentPossession  TeamStatus
	PossessionAtHalf   TeamStatus
	CurrentDown        Down
}

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
)

type Down string

type TeamStatus string

type EndZone string

type GameEvent string

const (
	GameStart           GameEvent  = "GameStart"
	EndOfGame           GameEvent  = "EndOfGame"
	EndofQuarter        GameEvent  = "EndOfQuater"
	EndofRegulationPlay GameEvent  = "EndofRegulationPlay"
	PlayRan             GameEvent  = "PlayRan"
	QuarterTime         int        = 900
	EastEndZone         EndZone    = "EastEndZone"
	WestEndZone         EndZone    = "WestEndZone"
	AwayTeam            TeamStatus = "AwayTeam"
	HomeTeam            TeamStatus = "HomeTeam"
	FirstDown           Down       = "FirstDown"
	SecondDown          Down       = "SecondDown"
	ThirdDown           Down       = "ThirdDown"
	FourthDown          Down       = "FourthDown"
)

func (reqBody SimulateGameRequest) InitGame() Game {
	return Game{
		HomeTeam:           reqBody.HomeTeam,
		AwayTeam:           reqBody.AwayTeam,
		GameOver:           false,
		Quarter:            First,
		GameClockInSeconds: QuarterTime,
		HomeTeamScore:      0,
		AwayTeamScore:      0,
		HomeTeamEndZone:    WestEndZone,
		AwayTeamEndZone:    EastEndZone,
		CurrentPossession:  HomeTeam,
		PossessionAtHalf:   AwayTeam,
		CurrentDown:        FirstDown,
	}
}
