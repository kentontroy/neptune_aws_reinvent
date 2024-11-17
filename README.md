# neptune_aws_reinvent

```
cd $HOME
wget -c https://go.dev/dl/go1.23.1.linux-amd64.tar.gz
sudo tar -C /usr/local/ -xzf go1.23.1.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
go version
```

```
export AWS_CONFIG=317913635185_cldr_poweruser
export AWS_REGION=us-east-2
export AWS_BUCKET=kdavis-bucket
export AWS_BUCKET_KEY_DIR=data/bulk-loader-example-opencypher-format

```

```
git clone https://github.com/kentontroy/neptune_aws_reinvent
cd neptune_aws_reinvent
export NEPTUNE_PROJECT_HOME=$PWD
export NEPTUNE_PROJECT_FILE_DIR=data/bulk-loader-example-opencypher-format
```

```
cd ./src/neptune-database-load/go

go run upload-to-s3.go \
  --source="${NEPTUNE_PROJECT_HOME}/${NEPTUNE_PROJECT_FILE_DIR}/node-olist-customers.csv" \
  --aws_config="${AWS_CONFIG}" \
  --aws_region="${AWS_REGION}" \
  --aws_bucket="${AWS_BUCKET}" \
  --aws_bucket_key="${AWS_BUCKET_KEY_DIR}/node-olist-customers.csv"

go run upload-to-s3.go \
  --source="${NEPTUNE_PROJECT_HOME}/${NEPTUNE_PROJECT_FILE_DIR}/node-olist-orders.csv" \
  --aws_config="${AWS_CONFIG}" \
  --aws_region="${AWS_REGION}" \
  --aws_bucket="${AWS_BUCKET}" \
  --aws_bucket_key="${AWS_BUCKET_KEY_DIR}/data/bulk-loader-example-opencypher-format/node-olist-orders.csv"

go run create-relationship-customer-order.go

go run upload-to-s3.go \
  --source="${NEPTUNE_PROJECT_HOME}/${NEPTUNE_PROJECT_FILE_DIR}/node-olist-orders.csv" \
  --aws_config="${AWS_CONFIG}" \
  --aws_region="${AWS_REGION}" \
  --aws_bucket="${AWS_BUCKET}" \
  --aws_bucket_key="${AWS_BUCKET_KEY_DIR}/relationship-customer-to-order.csv"

```
```
cd ${NEPTUNE_PROJECT_HOME}/scripts
./load-to-neptune.sh "kdavis-bucket" "data/bulk-loader-example-opencypher-format/node-olist-customers.csv"
./load-to-neptune.sh "kdavis-bucket" "data/bulk-loader-example-opencypher-format/node-olist-orders.csv"
./load-to-neptune.sh "kdavis-bucket" "data/bulk-loader-example-opencypher-format/relationship-customer-to-order.csv"
```
