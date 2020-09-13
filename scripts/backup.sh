#!/bin/bash
export PATH="/usr/local/bin:/usr/bin:/bin"

AWS_ACCESS_KEY=AKIASEJMBX3ZVP6IVP4L
AWS_SECRET_KEY=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
BUCKET=views-dgraph-backups-development
REGION=eu-west-1
EXPORT_PATH=/exports

dgbr backup \
--AWS_ACCESS_KEY=${AWS_ACCESS_KEY} \
--AWS_SECRET_KEY=${AWS_SECRET_KEY} \
--bucket=${BUCKET}  \
--region=${REGION}  \
--export=${EXPORT_PATH} \
--prefix=ddd