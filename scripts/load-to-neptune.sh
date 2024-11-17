#!/bin/bash

# Check if bucket name and key are provided
if [ $# -ne 2 ]; then
    echo "Usage: $0 <bucket-name> <bucket-key>"
    exit 1
fi

# Assign arguments to variables
AWS_S3_BUCKET=$1
AWS_S3_BUCKET_KEY=$2

##############################################################################################
# References
# https://docs.aws.amazon.com/neptune/latest/userguide/load-api-reference-load-examples.html
##############################################################################################
load_file="s3://$1/$2"
echo "$load_file"

curl -X POST \
    -H 'Content-Type: application/json' \
    https://${AWS_NEPTUNE_ENDPOINT}:${AWS_NEPTUNE_ENDPOINT_PORT}/loader -d '
    {
        "source" : "'"$load_file"'",
        "format" : "'"${AWS_NEPTUNE_LOADER_FORMAT}"'",
        "iamRoleArn" : "'"${AWS_NEPTUNE_ARN}"'",
        "region" : "'"${AWS_REGION}"'",
        "failOnError" : "FALSE",
        "parallelism" : "MEDIUM",
        "updateSingleCardinalityProperties" : "FALSE",
        "queueRequest" : "TRUE"
    }'
