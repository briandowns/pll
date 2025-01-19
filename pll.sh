#!/bin/sh

set -x 

TOKEN='Bearer N)eIKy1rZ%/%fm1WhM7tuVcrR*UIsc'
YEAR='2024'
CHAMP_SERIES="false"
QUERY='query($year: Int!, $champSeries: Boolean!) {standings(season: $year, champSeries: $champSeries){team{officialId location locationCode urlLogo fullName} seed wins @skip(if: $champSeries) losses @skip(if: $champSeries) ties @skip(if: $champSeries) scores @skip(if: $champSeries) scoresAgainst @skip(if: $champSeries) scoreDiff @skip(if: $champSeries) csWins @include(if: $champSeries) csLosses @include(if: $champSeries) csTies @include(if: $champSeries) csScores @include(if: $champSeries) csScoresAgainst @include(if: $champSeries) csScoreDiff @include(if: $champSeries) conferenceWins conferenceLosses conferenceTies conferenceScores conferenceScoresAgainst conference conferenceSeed}}'
SCHEMA='query {__schema {types {name}}}'
PAYLOAD="{\"operationName\":null,\"variables\":{\"year\": ${YEAR}, \"champSeries\": ${CHAMP_SERIES}}, \"query\": \"${QUERY}\"}"
SCHEMA_PAYLOAD="{\"operationName\":null,\"variables\":{\"year\": ${YEAR}, \"champSeries\": ${CHAMP_SERIES}}, \"query\": \"${SCHEMA}\"}"
# printf -v PAYLOAD '{
#   "operationName":null,
#   "variables": {
#     "year": ${YEAR},
#     "champSeries": ${CHAMP_SERIES}
#   },
# }'

curl -vvvv -s -XPOST \
     -H "Content-Type: application/json" \
     -H "Accept-encoding: gzip, deflate, br, zstd" \
     -H "Accept-language: en-US,en;q=0.9,it;q=0.8" \
     -H "Authorization: ${TOKEN}" \
     -H "Connection: keep-alive" \
     -d "${SCHEMA_PAYLOAD}" \
     https://api.stats.premierlacrosseleague.com/graphql

exit 0
