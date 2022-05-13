# Cortex replication simulations

This repo helps simulate how often Cortex queries will fail under different possible architectures. All of this presumes the presence of 3 availability zones, zonal replication, and sharding by all labels.

Scenarios:

- `even`: Status quo as of May 2022. Ingested series for a given tenant are written to S shards out of N ingesters. Each series writes to 3 ingesters, one in each AZ, with majority success required. Note that the replication sets for series can overlap, so series A can write to ingesters 1, 2, 3 while series B writes to ingesters 2, 3, 4. At query time, all S shards are queried. If 2 zones have failing ingesters, the query fails, because there could be a series that had those two failing ingesters as their majority.

- `even6`: `even`, but with a 6x replication factor. Ingested series for a given tenant are written to S shards out of N ingesters. Each series writes to 6 ingesters, two in each AZ, requiring all ingesters in 2 AZs to succeed. At query time, all S shards are queried. If 2 zones have at least 2 failing ingesters, the query fails.

- `even6-weaker`: `even6`, but with a slightly weaker ingest guarantee. Instead of requiring both ingesters in two zones for each series to completely succeed at ingest time, a simple majority of ingesters (4) must succeed for each series. This weakens query failure tolerance. Queries can now fail with either 2 ingesters in 2 AZs each failing _or_ 2 ingesters in 1 AZ plus 1 ingester each in the other AZs failing.

- `clumps`: The set of ingesters is consistently partitioned into shards. Each ingester belongs to one and only one shard, and each shard has one ingester in each AZ. Ingested series for a given tenant are written to S shards out of N/3 ingesters. Each series writes to a given shard, requiring a majority of ingesters in the shard to succeed. At query time, all S shards are queried. If, for any shard, a majority of ingesters fail, the query fails.

- `clumps6`: `clumps`, but with a 6x replication factor. The set of ingesters is consistently partitioned into shards. Each ingester belongs to one and only one shard, and each shard has 2 ingesters in each AZ. Ingested series for a given tenant are written to S shards out of N/6 ingesters. Each series writes to a given shard, requiring all ingesters in 2 zones for the shard to succeed. At query time, all S shards are queried. If, for any shard, two AZs entirely fail, the query fails.

- `clumps6-weaker`: `clumps6`, but with a slightly weaker ingest guarantee. Instead of requiring both ingesters in two zones for each shard to completely succeed at ingest time, a simple majority of ingesters in each shard must succeed. This weakens query failure tolerance, since now any 4 failing ingesters in a shard will fail a query.
