/*-
 * SPDX-License-Identifier: BSD-2-Clause
 *
 * Copyright (c) 2024 Brian J. Downs
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

package pll

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

const playerStatsQuery = `query($year: Int!, $seasonSegment: SeasonSegment, $statList: [String], $limit: Int) {
	playerStatLeaders(year: $year, seasonSegment: $seasonSegment, statList: $statList, limit: $limit) {
	officialId
	profileUrl
	firstName
	lastName
	position
	statType
	slug
	statValue
	playerRank
	jerseyNum
	teamId
	year
	}
}
`
