/*-
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2025 Brian J. Downs
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE REGENTS AND CONTRIBUTORS ``AS IS'' AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED.  IN NO EVENT SHALL THE REGENTS OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
 * OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
 * LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
 * OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
 * SUCH DAMAGE.
 */

// package pll provides access to the Premier Lacrosse League's GraphQL
// statistics database.
package pll

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/machinebox/graphql"
)

const graphqlEndpoint = "https://api.stats.premierlacrosseleague.com/graphql"

var seasonSegments = []string{
	"regular",
	"post",
	"champSeries",
}

var PlayerStatistics = []string{
	"points",
	"onePointGoals",
	"twoPointGoals",
	"scoringPoints",
	"assists",
	"shots",
	"pointsPG",
	"onePointGoalsPG",
	"assistsPG",
	"shotsPG",
	"shotPct",
	"touches",
	"faceoffPct",
	"faceoffWinsPG",
	"touchesPG",
	"savesPG",
	"savePct",
	"causedTurnovers",
	"causedTurnoversPG",
	"groundBalls",
	"groundBallsPG",
}

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

// PlayerStatLeader
type PlayerStatLeader struct {
	OfficialID string `json:"officialId"`
	ProfileURL string `json:"profileUrl"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Position   string `json:"position"`
	StatType   string `json:"statType"`
	Slug       string `json:"slug"`
	StatValue  string `json:"statValue"`
	PlayerRank int    `json:"playerRank"`
	JerseyNum  string `json:"jerseyNum"`
	TeamID     string `json:"teamId"`
	Year       int    `json:"year"`
}

// PlayerStatsResponse
type PlayerStatsResponse struct {
	Data struct {
		PlayerStatLeaders []PlayerStatLeader `json:"playerStatLeaders"`
	} `json:"data"`
}

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

// TeamStats
func (p *PLL) PlayerStats(ctx context.Context, year, limit int, seasonSegment string, stats []string) (*PlayerStatsResponse, error) {
	if err := ValidSeasonSegment(seasonSegment); err != nil {
		return nil, err
	}

	if err := ValidStats(stats); err != nil {
		return nil, errors.New("invalid stats")
	}

	req := graphql.NewRequest(playerStatsQuery)
	req.Var("year", year)
	req.Var("seasonSegment", seasonSegment)
	req.Var("statList", strings.Join(stats, ","))
	req.Var("limit", limit)
	req.Header.Set("Authorization", "Bearer "+p.token)

	var res PlayerStatsResponse
	if err := p.client.Run(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// ValidSeasonSegment checks to see if the given season
// segment is valid.
func ValidSeasonSegment(segment string) error {
	if !slices.Contains(seasonSegments, segment) {
		return errors.New("invalid segment: " + segment)
	}

	return nil
}

// ValidStats checks to see if the given stats are valid.
func ValidStats(stats []string) error {
	if slices.Equal(PlayerStatistics, stats) {
		return nil
	}

	for _, stat := range stats {
		if !slices.Contains(PlayerStatistics, stat) {
			return errors.New("invalid stat: " + stat)
		}
	}

	return nil
}

type AutoGenerated struct {
	Data struct {
		AllTeams []struct {
			OfficialID     string `json:"officialId"`
			LocationCode   string `json:"locationCode"`
			Location       string `json:"location"`
			FullName       string `json:"fullName"`
			URLLogo        string `json:"urlLogo"`
			Slogan         string `json:"slogan"`
			TeamWins       int    `json:"teamWins"`
			TeamLosses     int    `json:"teamLosses"`
			TeamTies       int    `json:"teamTies"`
			TeamWinsPost   int    `json:"teamWinsPost"`
			TeamLossesPost int    `json:"teamLossesPost"`
			TeamTiesPost   int    `json:"teamTiesPost"`
			League         string `json:"league"`
			Coaches        []struct {
				Name      string `json:"name"`
				CoachType string `json:"coachType"`
			} `json:"coaches"`
			Stats struct {
				Scores                     int     `json:"scores"`
				FaceoffPct                 float64 `json:"faceoffPct"`
				ShotPct                    float64 `json:"shotPct"`
				TwoPointShotPct            float64 `json:"twoPointShotPct"`
				TwoPointShotsOnGoalPct     float64 `json:"twoPointShotsOnGoalPct"`
				ClearPct                   float64 `json:"clearPct"`
				RidesPct                   float64 `json:"ridesPct"`
				SavePct                    float64 `json:"savePct"`
				ShortHandedPct             int     `json:"shortHandedPct"`
				ShortHandedGoalsAgainstPct float64 `json:"shortHandedGoalsAgainstPct"`
				PowerPlayGoalsAgainstPct   float64 `json:"powerPlayGoalsAgainstPct"`
				ManDownPct                 float64 `json:"manDownPct"`
				ShotsOnGoalPct             float64 `json:"shotsOnGoalPct"`
				OnePointGoals              int     `json:"onePointGoals"`
				ScoresAgainst              int     `json:"scoresAgainst"`
				Saa                        int     `json:"saa"`
				PowerPlayPct               float64 `json:"powerPlayPct"`
				GamesPlayed                int     `json:"gamesPlayed"`
				Goals                      int     `json:"goals"`
				TwoPointGoals              int     `json:"twoPointGoals"`
				Assists                    int     `json:"assists"`
				GroundBalls                int     `json:"groundBalls"`
				Turnovers                  int     `json:"turnovers"`
				CausedTurnovers            int     `json:"causedTurnovers"`
				FaceoffsWon                int     `json:"faceoffsWon"`
				FaceoffsLost               int     `json:"faceoffsLost"`
				Faceoffs                   int     `json:"faceoffs"`
				Shots                      int     `json:"shots"`
				TwoPointShots              int     `json:"twoPointShots"`
				TwoPointShotsOnGoal        int     `json:"twoPointShotsOnGoal"`
				GoalsAgainst               int     `json:"goalsAgainst"`
				TwoPointGoalsAgainst       int     `json:"twoPointGoalsAgainst"`
				NumPenalties               int     `json:"numPenalties"`
				Pim                        float64 `json:"pim"`
				Clears                     int     `json:"clears"`
				ClearAttempts              int     `json:"clearAttempts"`
				Rides                      int     `json:"rides"`
				RideAttempts               int     `json:"rideAttempts"`
				Saves                      int     `json:"saves"`
				Offsides                   int     `json:"offsides"`
				ShotClockExpirations       int     `json:"shotClockExpirations"`
				PowerPlayGoals             int     `json:"powerPlayGoals"`
				PowerPlayShots             int     `json:"powerPlayShots"`
				ShortHandedGoals           int     `json:"shortHandedGoals"`
				ShortHandedShots           int     `json:"shortHandedShots"`
				ShortHandedShotsAgainst    int     `json:"shortHandedShotsAgainst"`
				ShortHandedGoalsAgainst    int     `json:"shortHandedGoalsAgainst"`
				PowerPlayGoalsAgainst      int     `json:"powerPlayGoalsAgainst"`
				PowerPlayShotsAgainst      int     `json:"powerPlayShotsAgainst"`
				TimesManUp                 int     `json:"timesManUp"`
				TimesShortHanded           int     `json:"timesShortHanded"`
				ShotsOnGoal                int     `json:"shotsOnGoal"`
				ScoresPG                   float64 `json:"scoresPG"`
				ShotsPG                    float64 `json:"shotsPG"`
				TotalPasses                int     `json:"totalPasses"`
				Touches                    int     `json:"touches"`
			} `json:"stats"`
			PostStats struct {
				Scores                     int     `json:"scores"`
				FaceoffPct                 float64 `json:"faceoffPct"`
				ShotPct                    float64 `json:"shotPct"`
				TwoPointShotPct            float64 `json:"twoPointShotPct"`
				TwoPointShotsOnGoalPct     float64 `json:"twoPointShotsOnGoalPct"`
				ClearPct                   int     `json:"clearPct"`
				RidesPct                   float64 `json:"ridesPct"`
				SavePct                    float64 `json:"savePct"`
				ShortHandedPct             int     `json:"shortHandedPct"`
				ShortHandedGoalsAgainstPct int     `json:"shortHandedGoalsAgainstPct"`
				PowerPlayGoalsAgainstPct   int     `json:"powerPlayGoalsAgainstPct"`
				ManDownPct                 int     `json:"manDownPct"`
				ShotsOnGoalPct             float64 `json:"shotsOnGoalPct"`
				OnePointGoals              int     `json:"onePointGoals"`
				ScoresAgainst              int     `json:"scoresAgainst"`
				Saa                        float64 `json:"saa"`
				PowerPlayPct               float64 `json:"powerPlayPct"`
				GamesPlayed                int     `json:"gamesPlayed"`
				Goals                      int     `json:"goals"`
				TwoPointGoals              int     `json:"twoPointGoals"`
				Assists                    int     `json:"assists"`
				GroundBalls                int     `json:"groundBalls"`
				Turnovers                  int     `json:"turnovers"`
				CausedTurnovers            int     `json:"causedTurnovers"`
				FaceoffsWon                int     `json:"faceoffsWon"`
				FaceoffsLost               int     `json:"faceoffsLost"`
				Faceoffs                   int     `json:"faceoffs"`
				Shots                      int     `json:"shots"`
				TwoPointShots              int     `json:"twoPointShots"`
				TwoPointShotsOnGoal        int     `json:"twoPointShotsOnGoal"`
				GoalsAgainst               int     `json:"goalsAgainst"`
				TwoPointGoalsAgainst       int     `json:"twoPointGoalsAgainst"`
				NumPenalties               int     `json:"numPenalties"`
				Pim                        float64 `json:"pim"`
				Clears                     int     `json:"clears"`
				ClearAttempts              int     `json:"clearAttempts"`
				Rides                      int     `json:"rides"`
				RideAttempts               int     `json:"rideAttempts"`
				Saves                      int     `json:"saves"`
				Offsides                   int     `json:"offsides"`
				ShotClockExpirations       int     `json:"shotClockExpirations"`
				PowerPlayGoals             int     `json:"powerPlayGoals"`
				PowerPlayShots             int     `json:"powerPlayShots"`
				ShortHandedGoals           int     `json:"shortHandedGoals"`
				ShortHandedShots           int     `json:"shortHandedShots"`
				ShortHandedShotsAgainst    int     `json:"shortHandedShotsAgainst"`
				ShortHandedGoalsAgainst    int     `json:"shortHandedGoalsAgainst"`
				PowerPlayGoalsAgainst      int     `json:"powerPlayGoalsAgainst"`
				PowerPlayShotsAgainst      int     `json:"powerPlayShotsAgainst"`
				TimesManUp                 int     `json:"timesManUp"`
				TimesShortHanded           int     `json:"timesShortHanded"`
				ShotsOnGoal                int     `json:"shotsOnGoal"`
				ScoresPG                   int     `json:"scoresPG"`
				ShotsPG                    float64 `json:"shotsPG"`
				TotalPasses                int     `json:"totalPasses"`
				Touches                    int     `json:"touches"`
			} `json:"postStats"`
			ChampSeries struct {
				TeamWins   int `json:"teamWins"`
				TeamLosses int `json:"teamLosses"`
				TeamTies   int `json:"teamTies"`
				Stats      struct {
					Scores                     int     `json:"scores"`
					FaceoffPct                 float64 `json:"faceoffPct"`
					ShotPct                    float64 `json:"shotPct"`
					TwoPointShotPct            float64 `json:"twoPointShotPct"`
					TwoPointShotsOnGoalPct     float64 `json:"twoPointShotsOnGoalPct"`
					ClearPct                   int     `json:"clearPct"`
					RidesPct                   float64 `json:"ridesPct"`
					SavePct                    float64 `json:"savePct"`
					ShortHandedPct             int     `json:"shortHandedPct"`
					ShortHandedGoalsAgainstPct int     `json:"shortHandedGoalsAgainstPct"`
					PowerPlayGoalsAgainstPct   float64 `json:"powerPlayGoalsAgainstPct"`
					ManDownPct                 float64 `json:"manDownPct"`
					ShotsOnGoalPct             float64 `json:"shotsOnGoalPct"`
					OnePointGoals              int     `json:"onePointGoals"`
					ScoresAgainst              int     `json:"scoresAgainst"`
					Saa                        int     `json:"saa"`
					PowerPlayPct               float64 `json:"powerPlayPct"`
					GamesPlayed                int     `json:"gamesPlayed"`
					Goals                      int     `json:"goals"`
					TwoPointGoals              int     `json:"twoPointGoals"`
					Assists                    int     `json:"assists"`
					GroundBalls                int     `json:"groundBalls"`
					Turnovers                  int     `json:"turnovers"`
					CausedTurnovers            int     `json:"causedTurnovers"`
					FaceoffsWon                int     `json:"faceoffsWon"`
					FaceoffsLost               int     `json:"faceoffsLost"`
					Faceoffs                   int     `json:"faceoffs"`
					Shots                      int     `json:"shots"`
					TwoPointShots              int     `json:"twoPointShots"`
					TwoPointShotsOnGoal        int     `json:"twoPointShotsOnGoal"`
					GoalsAgainst               int     `json:"goalsAgainst"`
					TwoPointGoalsAgainst       int     `json:"twoPointGoalsAgainst"`
					NumPenalties               int     `json:"numPenalties"`
					Pim                        int     `json:"pim"`
					Clears                     int     `json:"clears"`
					ClearAttempts              int     `json:"clearAttempts"`
					Rides                      int     `json:"rides"`
					RideAttempts               int     `json:"rideAttempts"`
					Saves                      int     `json:"saves"`
					Offsides                   int     `json:"offsides"`
					ShotClockExpirations       int     `json:"shotClockExpirations"`
					PowerPlayGoals             int     `json:"powerPlayGoals"`
					PowerPlayShots             int     `json:"powerPlayShots"`
					ShortHandedGoals           int     `json:"shortHandedGoals"`
					ShortHandedShots           int     `json:"shortHandedShots"`
					ShortHandedShotsAgainst    int     `json:"shortHandedShotsAgainst"`
					ShortHandedGoalsAgainst    int     `json:"shortHandedGoalsAgainst"`
					PowerPlayGoalsAgainst      int     `json:"powerPlayGoalsAgainst"`
					PowerPlayShotsAgainst      int     `json:"powerPlayShotsAgainst"`
					TimesManUp                 int     `json:"timesManUp"`
					TimesShortHanded           int     `json:"timesShortHanded"`
					ShotsOnGoal                int     `json:"shotsOnGoal"`
					ScoresPG                   float64 `json:"scoresPG"`
					ShotsPG                    float64 `json:"shotsPG"`
					TotalPasses                int     `json:"totalPasses"`
					Touches                    int     `json:"touches"`
				} `json:"stats"`
			} `json:"champSeries"`
		} `json:"allTeams"`
	} `json:"data"`
}
