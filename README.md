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
MERGE (:Tier_Diamond {name: "Diamond Tier", discount: "Free shipping on all orders and free lifetime warranty on applicable products"})
MERGE (:Tier_Gold {name: "Gold Tier", discount: "Free shipping on orders above $100 and discounted warranty on applicable products"})
MERGE (:Tier_Silver {name: "Silver Tier", discount: "Free shipping on orders above $200"})
MERGE (:Tier_Member {name: "Member Tier"})
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

#### Create a random sample of customers for demo purposes
```
MATCH (c:customer)
WITH c, rand() AS randomValue
ORDER BY randomValue
LIMIT 50
WITH COLLECT(c.customer_id) AS sample_customers
UNWIND sample_customers AS customer_id

MERGE (s:sample_customer {customer_id: customer_id})
```

```
MATCH (c:sample_customer), (o:order)
WHERE c.customer_id = o.customer_id
MERGE (c)-[:placed]->(o)
RETURN c, o
```
```
%%oc

UNWIND [
"018d3918d60c2e821b654b8eb4cfde55",
"2ad313a8287a197d9b303c380a6981db",
"46efeb419cd0b58214d62fb6f13623f7",
"40e2d996e2e7b68e453e264a377ba7f8",
"0197368f9bed3b6ac07506404f78a155",
"32b1d8089a0d72d4c54f9617d432ddbd",
"8153bf38caa53372d335a611651d7ae1",
"bf16ebb1ed47bce228a5b9f57fd4ff82",
"c6afce1779e2ee8b819a6130d120a48e",
"de74037ebc271b8781a9b7d703924ac3",
"24cf4347382b8e190a512b2de912f89c",
"2069f7dac27ea08a26fbdcdb18ae6142",
"3516087073944b772c2f8b3d0f93a31b",
"ff423802c40cadf7c338845c137fe304",
"4058d823a6156dd8f9a068750310bac5",
"b2ce968dfdaf1da5b510bc7f310f8dd9",
"1392fdcad98720de6dcb4296794f204b",
"5b6e97fea8528cf0060d86dbd82f9c4a",
"df2554dfedbe850dff2ae6c179fb21b2",
"2db3e03b1db3a836ed63120f3ce8e362",
"1a5967d84c2fb8d6c22c6eff14643058",
"bdcc20055a51ea3a85d7b9087f0a53ef",
"9c338ea8093192e203bc16add78c123c",
"7099bc9e000fed5fe3cace34788e7714",
"477ba7008fa296dafcedb37b8bd9b702",
"ea50d78b023d96a5d45a86d2059348e1",
"7f59e4e8a71ab50abe6f08288e94480e",
"4907ba2deda6f3cf96409a181c097ef5",
"180404f910942ed7c4b5f3e952007686",
"c4eb325091a03f5a95a8b188eea38273",
"fd4d78d4ac99d34e9c5d0f66a4540d8f",
"4ab4db2f93c68d5914dc5eb566dc486c",
"c9cb57b640e67ea0437a79903e2d2fcb",
"09ce0754a6b5bcbb8c24d38a6ce54543",
"1fc6b3289b080e774c436bacf707eeb7",
"1c0f65288ca605e8f359d5dc62043aed",
"606aca2f93f152b1d2a86dca7c556b5d",
"a684554378131af35310b25179278c1c",
"4fa00e989992cc755f0e1fe2a1b89ee7",
"dfcde7971c143052874fa2bf1623a3ab",
"5a4b754ec98d34c53658d10de7b1f620",
"ab4e4ce50bd32ebffeca76ac9ade7044",
"55cae1e7c9b2dd0420cf1f95699d77a6",
"267894d0c4c7fa22ff3b7eefca26a46b",
"77e188b1127d58981db82ce19e0b601d",
"69ca21797c2f506e1776ff086dd987f9",
"33829a5ab9c9fbe1c5a55943b73250ad",
"8b6fb564288d7bee8cd7234845cfb0e6",
"4f1f2b13805c2ab2ce70a6cad8001b18",
"bb465a223e03add1dea0f9b32822f59c"
] AS customer_id
MERGE (t:demo_set_customer {customer_id: customer_id})
RETURN t
```




