package main

import (
	"context"
	"fmt"
	"os"

	"github.com/machinebox/graphql"
)

const graphqlEndpoint = "https://api.stats.premierlacrosseleague.com/graphql"

// Team
type Team struct {
	FullName     string `json:"fullName"`
	Location     any    `json:"location"`
	LocationCode any    `json:"locationCode"`
	OfficialID   string `json:"officialId"`
	URLLogo      string `json:"urlLogo"` // URL to image
}

// Standing
type Standing struct {
	Conference              any  `json:"conference"`
	ConferenceLosses        int  `json:"conferenceLosses"`
	ConferenceScores        int  `json:"conferenceScores"`
	ConferenceScoresAgainst int  `json:"conferenceScoresAgainst"`
	ConferenceSeed          any  `json:"conferenceSeed"`
	ConferenceTies          int  `json:"conferenceTies"`
	ConferenceWins          int  `json:"conferenceWins"`
	Losses                  int  `json:"losses"`
	ScoreDiff               int  `json:"scoreDiff"`
	Scores                  int  `json:"scores"`
	ScoresAgainst           int  `json:"scoresAgainst"`
	Seed                    int  `json:"seed"`
	Team                    Team `json:"team"`
	Ties                    int  `json:"ties"`
	Wins                    int  `json:"wins"`
}

// StandingsResponse
type StandingsResponse struct {
	Standings []Standing `json:"standings"`
}

// standingsQuery contains the GraphQL query to get all standings
// by year and championship series.
const standingsQuery = `query($year: Int!, $champSeries: Boolean!) {
	standings(season: $year, champSeries: $champSeries){
	team {
		officialId
		location
		locationCode
		urlLogo
		fullName
	}
	seed
	wins @skip(if: $champSeries)
	losses @skip(if: $champSeries)
	ties @skip(if: $champSeries)
	scores @skip(if: $champSeries)
	scoresAgainst @skip(if: $champSeries)
	scoreDiff @skip(if: $champSeries)
	csWins @include(if: $champSeries)
	csLosses @include(if: $champSeries)
	csTies @include(if: $champSeries)
	csScores @include(if: $champSeries)
	csScoresAgainst @include(if: $champSeries)
	csScoreDiff @include(if: $champSeries)
	conferenceWins
	conferenceLosses
	conferenceTies
	conferenceScores
	conferenceScoresAgainst
	conference
	conferenceSeed
	}
}	  
`

// PLL
type PLL struct {
	token  string
	client *graphql.Client
}

// NewPLL creates a new value of PLL with an initialized
// GraphQL client using the given token.
func NewPLL(token string) *PLL {
	return &PLL{
		token:  token,
		client: graphql.NewClient(graphqlEndpoint),
	}
}

// Standings
func (p *PLL) Standings(ctx context.Context, year int, champSeries bool) (*StandingsResponse, error) {
	req := graphql.NewRequest(standingsQuery)
	req.Var("year", year)
	req.Var("champSeries", champSeries)
	req.Header.Set("Authorization", "Bearer "+p.token)

	var res StandingsResponse
	if err := p.client.Run(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func main() {
	ctx := context.Background()

	token := os.Getenv("PLL_BEARER_TOKEN")

	pll := NewPLL(token)

	standings, err := pll.Standings(ctx, 2023, false)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("%#+v\n", standings)

	os.Exit(0)
}
