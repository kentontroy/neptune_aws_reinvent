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
export AWS_NEPTUNE_LOADER_FORMAT="opencypher"
export AWS_NEPTUNE_ENDPOINT="db-neptune-aws-reinvent.cluster-cmabrddbmjfm.us-east-2.neptune.amazonaws.com"
export AWS_NEPTUNE_ENDPOINT_PORT=8182
export AWS_NEPTUNE_ARN="arn:aws:iam::317913635185:role/NeptuneLoadFromS3"
```

```
git clone https://github.com/kentontroy/neptune_aws_reinvent
cd neptune_aws_reinvent
export NEPTUNE_PROJECT_HOME=$PWD
export NEPTUNE_LOADER_FILE_DIR=${NEPTUNE_PROJECT_HOME}/data/bulk-loader-example-opencypher-format
```

```
cd ./src/neptune-database-load/go

go run upload-to-s3.go \
  --source="${NEPTUNE_LOADER_FILE_DIR}/node-olist-customers.csv" \
  --aws_config="${AWS_CONFIG}" \
  --aws_region="${AWS_REGION}" \
  --aws_bucket="${AWS_BUCKET}" \
  --aws_bucket_key="${AWS_BUCKET_KEY_DIR}/node-olist-customers.csv"

go run upload-to-s3.go \
  --source="${NEPTUNE_LOADER_FILE_DIR}/node-olist-orders.csv" \
  --aws_config="${AWS_CONFIG}" \
  --aws_region="${AWS_REGION}" \
  --aws_bucket="${AWS_BUCKET}" \
  --aws_bucket_key="${AWS_BUCKET_KEY_DIR}/node-olist-orders.csv"

go run create-relationship-customer-order.go

go run upload-to-s3.go \
  --source="${NEPTUNE_LOADER_FILE_DIR}/node-olist-orders.csv" \
  --aws_config="${AWS_CONFIG}" \
  --aws_region="${AWS_REGION}" \
  --aws_bucket="${AWS_BUCKET}" \
  --aws_bucket_key="${AWS_BUCKET_KEY_DIR}/relationship-customer-to-order.csv"

go run upload-to-s3.go \
  --source="${NEPTUNE_LOADER_FILE_DIR}/node-olist-orders.csv" \
  --aws_config="${AWS_CONFIG}" \
  --aws_region="${AWS_REGION}" \
  --aws_bucket="${AWS_BUCKET}" \
  --aws_bucket_key="${AWS_BUCKET_KEY_DIR}/relationship-customer-to-order.csv"

go run upload-to-s3.go \
  --source="${NEPTUNE_LOADER_FILE_DIR}/node-olist-products.csv" \
  --aws_config="${AWS_CONFIG}" \
  --aws_region="${AWS_REGION}" \
  --aws_bucket="${AWS_BUCKET}" \
  --aws_bucket_key="${AWS_BUCKET_KEY_DIR}/node-olist-products.csv"

go run upload-to-s3.go \
  --source="${NEPTUNE_LOADER_FILE_DIR}/relationship-order-to-product.csv" \
  --aws_config="${AWS_CONFIG}" \
  --aws_region="${AWS_REGION}" \
  --aws_bucket="${AWS_BUCKET}" \
  --aws_bucket_key="${AWS_BUCKET_KEY_DIR}/relationship-order-to-product.csv"

```

```
chmod -R 755 ${NEPTUNE_PROJECT_HOME}/scripts
cd ${NEPTUNE_PROJECT_HOME}/scripts
./load-to-neptune.sh "${AWS_BUCKET}" "${AWS_BUCKET_KEY_DIR}/node-olist-customers.csv"
./load-to-neptune.sh "${AWS_BUCKET}" "${AWS_BUCKET_KEY_DIR}/node-olist-orders.csv"
./load-to-neptune.sh "${AWS_BUCKET}" "${AWS_BUCKET_KEY_DIR}/relationship-customer-to-order.csv"
./load-to-neptune.sh "${AWS_BUCKET}" "${AWS_BUCKET_KEY_DIR}/node-olist-products.csv"
./load-to-neptune.sh "${AWS_BUCKET}" "${AWS_BUCKET_KEY_DIR}/relationship-order-to-product"

```
