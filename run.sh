#!/bin/bash

function trial {
  failures=$1
  shards=$2
  ingesters=$3
  echo "-failures $failures -shards $shards -ingesters $ingesters, will 2x all for 6x replication case"

  (( failures2 = failures * 2 ))
  (( shards2 = shards * 2 ))
  (( ingesters2 = ingesters * 2 ))

  run "even" $failures $shards $ingesters
  run "even6" $failures2 $shards2 $ingesters2
  run "even6-weaker" $failures2 $shards2 $ingesters2
  run "clumps" $failures $shards $ingesters
  run "clumps6" $failures2 $shards2 $ingesters2
  run "clumps6-weaker" $failures2 $shards2 $ingesters2

  echo ""
}

function run {
  echo "$1"
  go run "./$1/main.go" -failures $2 -shards $3 -ingesters $4 -trials 1000000
}

trial 5 75 99
trial 5 300 999
trial 100 600 999
