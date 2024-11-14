#!/bin/bash

export AWS_REGION="us-east-2"
export AWS_NEPTUNE_LOADER_FORMAT="opencypher"
export AWS_NEPTUNE_ENDPOINT="db-neptune-aws-reinvent.cluster-cmabrddbmjfm.us-east-2.neptune.amazonaws.com"
export AWS_NEPTUNE_ENDPOINT_PORT=8182
export AWS_NEPTUNE_ARN="arn:aws:iam::317913635185:role/NeptuneLoadFromS3"

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
load_file="s3://${AWS_S3_BUCKET}/${AWS_S3_BUCKET_KEY}"
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
