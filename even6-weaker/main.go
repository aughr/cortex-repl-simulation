package main

import (
	"flag"
	"log"
	"math/rand"
	"time"
)

func main() {
	ingesters := flag.Int("ingesters", 99, "how many ingesters run in the cluster")
	shardcount := flag.Int("shards", 3, "how many shards a client maps to")
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
	order := rand.Perm(*ingesters)
	for i := 0; i < *shardcount; i++ {
		shards[order[i]] = true
	}

	failcount := 0
	successcount := 0

	for i := 0; i < *trials; i++ {
		if runTrial(*ingesters, *failures, shards) {
			successcount++
		} else {
			failcount++
		}
	}

	log.Printf("%d successes, %d failures == %0.2f%% failure rate", successcount, failcount, float64(failcount)/float64(successcount+failcount)*100)
}

func runTrial(ingesters int, failures int, shards map[int]bool) bool {
	var failedzones [3]int

	order := rand.Perm(ingesters)
	for i := 0; i < failures; i++ {
		failed := order[i]
		// if the failed node isn't in the shards used, continue
		if !shards[failed] {
			continue
		}

		// note this failure
		zone := zoneFor(failed)
		failedzones[zone]++

		left := failedzones[wrapmod(i-1, 3)]
		right := failedzones[wrapmod(i+1, 3)]
		if failedzones[zone] == 1 {
			// we have failed the query if another zone is at 2 and another zone is at 1
			if (left == 1 && right == 2) || (left == 2 && right == 1) {
				return false
			}
		} else if failedzones[zone] == 2 {
			// we have failed the query if another zone is at 2 _or_ both other zones are at 1
			if (left == 1 && right == 1) || left == 2 || right == 2 {
				return false
			}
		}
	}

	return true
}

func zoneFor(i int) int {
	return i % 3
}

func wrapmod(i int, m int) int {
	v := i % m
	if v >= 0 {
		return v
	}
	return m + v
}
