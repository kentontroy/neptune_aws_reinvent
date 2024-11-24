# neptune_aws_reinvent

#### Install Golang
```
cd $HOME
wget -c https://go.dev/dl/go1.23.1.linux-amd64.tar.gz
sudo tar -C /usr/local/ -xzf go1.23.1.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
go version
```

#### Install pyvenv on Ubuntu
```
sudo apt install -y make build-essential wget curl
sudo apt install -y libssl-dev
sudo apt install -y zliblg-dev libbz2-dev
sudo apt install -y libreadline-dev libsqlite3-dev
sudo apt install -y llvm libncurses5-dev libncursesw5-dev
sudo apt install -y xz-utils-dev tk-dev libffi-dev liblzms-dev python-openssl git
sudo apt install -y liblzma-dev python-openssl git

python3 -m pip install --upgrade pip

sudo curl https://pyenv.run | bash
cat $HOME/.bash_profile
...
# User specific environment and startup programs
export PYENV_ROOT="$HOME/.pyenv"
[[ -d $PYENV_ROOT/bin ]] && export PATH="$PYENV_ROOT/bin:$PATH"
eval "$(pyenv init -)"
eval "$(pyenv virtualenv-init -)"

pyenv install 3.10.14
pyenv virtualenv 3.10.14 venv
pyenv activate venv

pip install neo4j pandas mlxtend==0.23.1

```

#### Set environment variables
```
export AWS_CONFIG=317913635185_cldr_poweruser
export AWS_REGION=us-east-2
export AWS_BUCKET=kdavis-bucket
export AWS_BUCKET_KEY_DIR=data/bulk-loader-example-opencypher-format
export AWS_NEPTUNE_LOADER_FORMAT="opencypher"
export AWS_NEPTUNE_ENDPOINT=reinvent-db-neptune-1.cluster-cnlncb8narog.us-east-2.neptune.amazonaws.com
export AWS_NEPTUNE_ENDPOINT_PORT=8182
export AWS_NEPTUNE_ARN="arn:aws:iam::317913635185:role/NeptuneLoadFromS3"
```
#### Clone gitrepo
```
git clone https://github.com/kentontroy/neptune_aws_reinvent
cd neptune_aws_reinvent
export NEPTUNE_PROJECT_HOME=$PWD
export NEPTUNE_LOADER_FILE_DIR=${NEPTUNE_PROJECT_HOME}/data/bulk-loader-example-opencypher-format
```

#### Extract geolocation dataset
```
cd ${NEPTUNE_LOADER_FILE_DIR}
gunzip node-olist-geolocation.csv.gz
```

#### Upload data sets to S3. Create relationship between entities using custom code.
The datasets are formatted according to the Cypher bulk-loading specification
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

#### Load the datasets from S3 into Neptune. 
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

#### Create a Rewards Tier with specified discounts
```
MERGE (:tier_diamond {name: "Diamond Tier", discount: "Free shipping on all orders and free lifetime warranty on applicable products"})
MERGE (:tier_gold {name: "Gold Tier", discount: "Free shipping on orders above $100 and discounted warranty on applicable products"})
MERGE (:tier_silver {name: "Silver Tier", discount: "Free shipping on orders above $200"})
MERGE (:tier_member {name: "Member Tier"})
```

#### Calculate the variables used to map Customers to the Rewards Tier level
```
MATCH (o:order)-[r:has_item]->(p:product)
WITH o.order_id AS order_id, ROUND(SUM(r.price) * 100) / 100 AS purchase_amount
WITH AVG(purchase_amount) AS avg_purchase_amount, STDEVP(purchase_amount) AS stddev_purchase_amount

MERGE (l:lifetime_rewards_variable)
SET l.average_purchase_amount = avg_purchase_amount, l.stddev_purchase_amount = stddev_purchase_amount
RETURN l
```

```
MATCH (c:customer)-[:placed]->(o:order)-[r:has_item]->(p:product), (l:lifetime_rewards_variable)
WITH l, c, ROUND(SUM(r.price) * 100) / 100 AS purchase_amount
RETURN
CASE 
    WHEN purchase_amount > l.average_purchase_amount + (2 * l.stddev_purchase_amount) THEN "Diamond"
    WHEN purchase_amount > l.average_purchase_amount + l.stddev_purchase_amount THEN "Gold"
    WHEN purchase_amount >= l.average_purchase_amount THEN "Silver"
    ELSE "Member"
END AS tier

```



