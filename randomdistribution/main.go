package main

import (
	"flag"
	"log"
	"math/rand"
	"time"
)

func main() {
	ingesters := flag.Int("ingesters", 100, "how many ingesters run in the cluster")
	shardcount := flag.Int("shards", 4, "how many shards a client maps to")
	failures := flag.Int("failures", 0, "how many ingesters have failed")
	trials := flag.Int("trials", 100_000, "how many query trials to run")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	if *shardcount > *ingesters {
		log.Fatal("Cannot have more shards than ingesters")
	}
	if *failures > *ingesters {
		log.Fatal("Cannot have more failures than ingesters")
	}

	shards := make(map[int]bool, *shardcount)
	zones := make([]int, *ingesters)
	order := rand.Perm(*ingesters)
	for i := 0; i < *ingesters; i++ {
		zones[order[i]] = i % 3
	}
	for i := 0; i < *shardcount; i++ {
		shards[order[i]] = true
	}

	failcount := 0
	successcount := 0

	for i := 0; i < *trials; i++ {
		if runTrial(*ingesters, *failures, shards, zones) {
			successcount++
		} else {
			failcount++
		}
	}

	log.Printf("%d successes, %d failures == %0.2f%% failure rate", successcount, failcount, float64(failcount)/float64(successcount+failcount)*100)
}

func runTrial(ingesters int, failures int, shards map[int]bool, zones []int) bool {
	order := rand.Perm(ingesters)
	failedzone := -1
	for i := 0; i < failures; i++ {
		failed := order[i]
		// if the failed node isn't in the shards used, continue
		if !shards[failed] {
			continue
		}

		// if it is, see if we've already failed a zone
		if failedzone == -1 {
			failedzone = zones[failed]
			continue
		}

		if failedzone != zones[failed] {
			return false
		}
	}

	for i := 0; i < len(shards); i++ {
	}
	return true
}
