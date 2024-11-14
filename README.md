# neo4j_neptune_aws_reinvent

```
git clone https://github.com/kentontroy/neo4j_neptune_aws_reinvent
```

```

```

```
cd ./src/go/neptune-database-load
go run create-relationship-customer-order.go
go run upload-to-s3.go 
```

```
./scripts/load-to-neptune.sh "kdavis-bucket" "data/bulk-loader-example-opencypher-format/node-olist-customers.csv"
./scripts/load-to-neptune.sh "kdavis-bucket" "data/bulk-loader-example-opencypher-format/node-olist-orders.csv"
./scripts/load-to-neptune.sh "kdavis-bucket" "data/bulk-loader-example-opencypher-format/relationship-customer-to-order.csv"
```
