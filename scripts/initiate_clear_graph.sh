#!/bin/bash

##############################################################################################
# References
# https://docs.aws.amazon.com/neptune/latest/userguide/manage-console-fast-reset.html#:~:text=To%20delete%20all%20data%20from%20a%20Neptune%20DB%20cluster%20using%20the%20API&text=You%20do%20this%20by%20sending,to%20specify%20the%20initiateDatabaseReset%20action.&text=The%20token%20remains%20valid%20for,minutes)%20after%20it%20is%20issued.
##############################################################################################

curl -X POST \
  -H 'Content-Type: application/json' \
      https://${AWS_NEPTUNE_ENDPOINT}:${AWS_NEPTUNE_ENDPOINT_PORT}/system \
  -d '{ "action" : "initiateDatabaseReset" }'

