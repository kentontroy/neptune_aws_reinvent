# neo4j_neptune_aws_reinvent

```
git clone https://github.com/kentontroy/neo4j_neptune_aws_reinvent
```

```
cd ./src/go/neptune-database-load
go run create-relationship-customer-order.go
```

```
go run upload-to-s3.go \
  --source="./data/bulk-loader-example-opencypher-format/relationship-customer-to-order.csv" \
  --aws_config="317913635185_cldr_poweruser" \
  --aws_region="us-east-2" \
  --aws_bucket="kdavis-bucket" \
  --aws_bucket_key=""data/bulk-loader-example-opencypher-format/relationship-customer-to-order.csv"

```

```
./scripts/load-to-neptune.sh "kdavis-bucket" "data/bulk-loader-example-opencypher-format/node-olist-customers.csv"
./scripts/load-to-neptune.sh "kdavis-bucket" "data/bulk-loader-example-opencypher-format/node-olist-orders.csv"
./scripts/load-to-neptune.sh "kdavis-bucket" "data/bulk-loader-example-opencypher-format/relationship-customer-to-order.csv"
```
