package main

import (
	"flag"
	"log"
	"math/rand"
	"time"
)

func main() {
	ingesters := flag.Int("ingesters", 198, "how many ingesters run in the cluster")
	shardcount := flag.Int("shards", 6, "how many shards a client maps to")
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
	if *shardcount%6 != 0 || *ingesters%6 != 0 {
		log.Fatal("Shards and ingesters must be divisible by 6")
	}

	shards := make(map[int]bool, *shardcount/6)
	order := rand.Perm(*ingesters / 6)
	for i := 0; i < *shardcount/6; i++ {
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
	shardfailures := make([]tracker, ingesters/6)

	order := rand.Perm(ingesters)
	for i := 0; i < failures; i++ {
		failed := order[i]
		failedshard := shardFor(failed)
		// if the failed shard isn't in the shards used, continue
		if !shards[failedshard] {
			continue
		}

		if !shardfailures[failedshard].record() {
			return false
		}
	}

	return true
}

func shardFor(i int) int {
	return i / 6
}

type tracker struct {
	failed int
}

func (t *tracker) record() bool {
	t.failed += 1
	return t.failed < 4
}
