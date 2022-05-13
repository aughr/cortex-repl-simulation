# Cortex replication simulations

This repo helps simulate how often Cortex queries will fail under different possible architectures. All of this presumes the presence of 3 availability zones, zonal replication, and sharding by all labels.

## Scenarios:

- `even`: Status quo as of May 2022. Ingested series for a given tenant are written to S shards out of N ingesters. Each series writes to 3 ingesters, one in each AZ, with majority success required. Note that the replication sets for series can overlap, so series A can write to ingesters 1, 2, 3 while series B writes to ingesters 2, 3, 4. At query time, all S shards are queried. If 2 zones have failing ingesters, the query fails, because there could be a series that had those two failing ingesters as their majority.

- `even6`: `even`, but with a 6x replication factor. Ingested series for a given tenant are written to S shards out of N ingesters. Each series writes to 6 ingesters, two in each AZ, requiring all ingesters in 2 AZs to succeed. At query time, all S shards are queried. If 2 zones have at least 2 failing ingesters, the query fails.

- `even6-weaker`: `even6`, but with a slightly weaker ingest guarantee. Instead of requiring both ingesters in two zones for each series to completely succeed at ingest time, a simple majority of ingesters (4) must succeed for each series. This weakens query failure tolerance. Queries can now fail with either 2 ingesters in 2 AZs each failing _or_ 2 ingesters in 1 AZ plus 1 ingester each in the other AZs failing.

- `clumps`: The set of ingesters is consistently partitioned into shards. Each ingester belongs to one and only one shard, and each shard has one ingester in each AZ. Ingested series for a given tenant are written to S shards out of N/3 ingesters. Each series writes to a given shard, requiring a majority of ingesters in the shard to succeed. At query time, all S shards are queried. If, for any shard, a majority of ingesters fail, the query fails.

- `clumps6`: `clumps`, but with a 6x replication factor. The set of ingesters is consistently partitioned into shards. Each ingester belongs to one and only one shard, and each shard has 2 ingesters in each AZ. Ingested series for a given tenant are written to S shards out of N/6 ingesters. Each series writes to a given shard, requiring all ingesters in 2 zones for the shard to succeed. At query time, all S shards are queried. If, for any shard, two AZs entirely fail, the query fails.

- `clumps6-weaker`: `clumps6`, but with a slightly weaker ingest guarantee. Instead of requiring both ingesters in two zones for each shard to completely succeed at ingest time, a simple majority of ingesters in each shard must succeed. This weakens query failure tolerance, since now any 4 failing ingesters in a shard will fail a query.

## A run

via `./run.sh`:

```
-failures 5 -shards 75 -ingesters 99, will 2x all for 6x replication case
even
2022/05/13 13:37:59 78889 successes, 921111 failures == 92.11% failure rate
even6
2022/05/13 13:38:08 84201 successes, 915799 failures == 91.58% failure rate
even6-weaker
2022/05/13 13:38:13 36461 successes, 963539 failures == 96.35% failure rate
clumps
2022/05/13 13:38:15 852479 successes, 147521 failures == 14.75% failure rate
clumps6
2022/05/13 13:38:20 999754 successes, 246 failures == 0.02% failure rate
clumps6-weaker
2022/05/13 13:38:24 998836 successes, 1164 failures == 0.12% failure rate

-failures 5 -shards 300 -ingesters 999, will 2x all for 6x replication case
even
2022/05/13 13:38:45 644829 successes, 355171 failures == 35.52% failure rate
even6
2022/05/13 13:39:24 851861 successes, 148139 failures == 14.81% failure rate
even6-weaker
2022/05/13 13:40:06 486844 successes, 513156 failures == 51.32% failure rate
clumps
2022/05/13 13:40:27 993953 successes, 6047 failures == 0.60% failure rate
clumps6
2022/05/13 13:41:09 1000000 successes, 0 failures == 0.00% failure rate
clumps6-weaker
2022/05/13 13:41:49 999998 successes, 2 failures == 0.00% failure rate

-failures 100 -shards 600 -ingesters 999, will 2x all for 6x replication case
even
2022/05/13 13:42:10 0 successes, 1000000 failures == 100.00% failure rate
even6
2022/05/13 13:42:50 0 successes, 1000000 failures == 100.00% failure rate
even6-weaker
2022/05/13 13:43:30 6459 successes, 993541 failures == 99.35% failure rate
clumps
2022/05/13 13:43:53 1789 successes, 998211 failures == 99.82% failure rate
clumps6
2022/05/13 13:44:42 943570 successes, 56430 failures == 5.64% failure rate
clumps6-weaker
2022/05/13 13:45:30 778756 successes, 221244 failures == 22.12% failure rate
```
