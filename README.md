# neptune_aws_reinvent

```
Navigate to the Neptune Cluster just provisioned. Create a notebook that can be used to access the graph database.
```

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
cd ${NEPTUNE_LOADER_FILE_DIR}
gunzip node-olist-geolocation.csv.gz
```

```
cd ${NEPTUNE_PROJECT_HOME}/src/neptune-database-load/go

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

go run upload-to-s3.go \
  --source="${NEPTUNE_LOADER_FILE_DIR}/node-olist-geolocation.csv" \
  --aws_config="${AWS_CONFIG}" \
  --aws_region="${AWS_REGION}" \
  --aws_bucket="${AWS_BUCKET}" \
  --aws_bucket_key="${AWS_BUCKET_KEY_DIR}/node-olist-geolocation.csv"
```

```
chmod -R 755 ${NEPTUNE_PROJECT_HOME}/scripts
cd ${NEPTUNE_PROJECT_HOME}/scripts
./load-to-neptune.sh "${AWS_BUCKET}" "${AWS_BUCKET_KEY_DIR}/node-olist-customers.csv"
./load-to-neptune.sh "${AWS_BUCKET}" "${AWS_BUCKET_KEY_DIR}/node-olist-orders.csv"
./load-to-neptune.sh "${AWS_BUCKET}" "${AWS_BUCKET_KEY_DIR}/relationship-customer-to-order.csv"
./load-to-neptune.sh "${AWS_BUCKET}" "${AWS_BUCKET_KEY_DIR}/node-olist-products.csv"
./load-to-neptune.sh "${AWS_BUCKET}" "${AWS_BUCKET_KEY_DIR}/relationship-order-to-product.csv"
./load-to-neptune.sh "${AWS_BUCKET}" "${AWS_BUCKET_KEY_DIR}/node-olist-geolocation.csv"
```

```
MATCH (c:customer)-[:ordered]->(o:order)-[r:has_item]->(p:product)
WITH c.customer_id AS customer_id, ROUND(SUM(r.price) * 100) / 100 as purchase_amount
ORDER BY purchase_amount DESC
LIMIT 50
WITH COLLECT(customer_id) AS top_customers
UNWIND top_customers AS customer_id
MATCH (c:customer {customer_id: customer_id})-[i:ordered]->(o:order)-[r:has_item]->(p:product)
RETURN customer_id, 
    COLLECT({
        year: i.order_purchase_timestamp_year, month: i.order_purchase_timestamp_month, 
        product: p.product_category_name, amount: ROUND(r.price * 100) / 100
    }) AS purchased_items
```

```
MATCH (o:order)-[r:has_item]->(p:product)
WITH o.order_id AS order_id, ROUND(SUM(r.price) * 100) / 100 AS purchase_amount
WITH AVG(purchase_amount) AS avg_purchase_amount, STDEVP(purchase_amount) AS stddev_purchase_amount

MATCH (c:customer)-[:ordered]->(o:order)-[r:has_item]->(p:product)
WITH avg_purchase_amount, stddev_purchase_amount, c.customer_id AS customer_id, ROUND(SUM(r.price) * 100) / 100 AS purchase_amount
RETURN
CASE 
    WHEN purchase_amount > avg_purchase_amount + (2 * stddev_purchase_amount) THEN "Diamond"
    WHEN purchase_amount > avg_purchase_amount + stddev_purchase_amount THEN "Gold"
    WHEN purchase_amount >= avg_purchase_amount THEN "Silver"
    ELSE "Member"
END AS tier
LIMIT 30
```


```
MATCH (c:customer), (g:geolocation)
WHERE c.geolocation_zip_code_prefix = g.geolocation_zip_code_prefix
MERGE (c)-[:located_at]->(g)

{
  "detailedMessage": "Operation terminated (out of memory)",
  "code": "MemoryLimitExceededException",
  "requestId": "7d2e5760-7114-4c69-ad0d-02f52aa74019",
  "message": "Operation terminated (out of memory)"
}
```
