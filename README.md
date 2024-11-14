# neptune_aws_reinvent

```
cd $HOME
wget -c https://go.dev/dl/go1.23.1.linux-amd64.tar.gz
sudo tar -C /usr/local/ -xzf go1.23.1.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
go version
```

```
git clone https://github.com/kentontroy/neptune_aws_reinvent
cd neptune_aws_reinvent
export NEPTUNE_PROJECT_HOME=$PWD
```

```
cd ./src/neptune-database-load/go

go run upload-to-s3.go \
  --source="${NEPTUNE_PROJECT_HOME}/data/bulk-loader-example-opencypher-format/node-olist-customers.csv" \
  --aws_config="317913635185_cldr_poweruser" \
  --aws_region="us-east-2" \
  --aws_bucket="kdavis-bucket" \
  --aws_bucket_key="data/bulk-loader-example-opencypher-format/node-olist-customers.csv"

go run upload-to-s3.go \
  --source="${NEPTUNE_PROJECT_HOME}/data/bulk-loader-example-opencypher-format/node-olist-orders.csv" \
  --aws_config="317913635185_cldr_poweruser" \
  --aws_region="us-east-2" \
  --aws_bucket="kdavis-bucket" \
  --aws_bucket_key="data/bulk-loader-example-opencypher-format/node-olist-orders.csv"

```


```
cd ./src/go/neptune-database-load
go run create-relationship-customer-order.go
```

```
./scripts/load-to-neptune.sh "kdavis-bucket" "data/bulk-loader-example-opencypher-format/node-olist-customers.csv"
./scripts/load-to-neptune.sh "kdavis-bucket" "data/bulk-loader-example-opencypher-format/node-olist-orders.csv"
./scripts/load-to-neptune.sh "kdavis-bucket" "data/bulk-loader-example-opencypher-format/relationship-customer-to-order.csv"
```
