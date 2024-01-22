# Read-Uncommitted


## Test

```
m test -w txn-rw-register --bin ./kvstore --node-count 2 --concurrency 2n --time-limit 20 --rate 1000 --consistency-models read-uncommitted

m test -w txn-rw-register --bin ./kvstore --node-count 2 --concurrency 2n --time-limit 20 --rate 1000 --consistency-models read-uncommitted --availability total --nemesis partition
```