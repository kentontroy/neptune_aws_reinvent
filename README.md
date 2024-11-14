# neo4j_neptune_aws_reinvent

git clone https://github.com/kentontroy/neo4j_neptune_aws_reinvent

./scripts/load-to-neptune.sh "kdavis-bucket" "data/bulk-loader-example-opencypher-format/node-olist-customers.csv"
./scripts/load-to-neptune.sh "kdavis-bucket" "data/bulk-loader-example-opencypher-format/node-olist-orders.csv"
./scripts/load-to-neptune.sh "kdavis-bucket" "data/bulk-loader-example-opencypher-format/relationship-customer-to-order.csv"
