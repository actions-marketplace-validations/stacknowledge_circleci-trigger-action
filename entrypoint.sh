#!/bin/sh -l

/circleci-trigger-action run --id $1 --project $2 --branch $3 --token $4 --timeout $5