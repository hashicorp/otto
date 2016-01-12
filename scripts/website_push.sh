#!/bin/bash
set -e

PROJECT="otto"
FASTLY_SERVICE_ID="7GrxRJP3PVBuqQbyxYQ0MV"

# Ensure the proper AWS environment variables are set
if [ -z "$AWS_ACCESS_KEY_ID" ]; then
  echo "Missing AWS_ACCESS_KEY_ID!"
  exit 1
fi

if [ -z "$AWS_SECRET_ACCESS_KEY" ]; then
  echo "Missing AWS_SECRET_ACCESS_KEY!"
  exit 1
fi

# Ensure the proper Fastly keys are set
if [ -z "$FASTLY_API_KEY" ]; then
  echo "Missing FASTLY_API_KEY!"
  exit 1
fi

# Ensure we have s3cmd installed
if ! command -v "s3cmd" >/dev/null 2>&1; then
  echo "Missing s3cmd!"
  exit 1
fi

# Get the parent directory of where this script is and change into our website
# directory
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

# Upload the files to S3 - we disable mime-type detection by the python library
# and just guess from the file extension because it's surprisingly more
# accurate, especially for CSS and javascript. We also tag the uploaded files
# with the proper Surrogate-Key, which we will later purge in our API call to
# Fastly.
if [ -z "$NO_UPLOAD" ]; then
  echo "Uploading to S3..."

  # Check that the site has been built
  if [ ! -d "$DIR/website/build" ]; then
    echo "Missing compiled website! Running `make build` to compile!"
    exit 1
  fi

  s3cmd \
    --quiet \
    --guess-mime-type \
    --no-mime-magic \
    --acl-public \
    --recursive \
    --add-header="Cache-Control: max-age=31536000" \
    --add-header="x-amz-meta-surrogate-key: site-$PROJECT" \
    put "$DIR/website/build/" "s3://hc-sites/$PROJECT/latest/"
fi

# Perform a soft-purge of the surrogate key.
if [ -z "$NO_PURGE" ]; then
  echo "Purging Fastly cache..."
  curl \
    --fail \
    --silent \
    --output /dev/null \
    --request "POST" \
    --header "Accept: application/json" \
    --header "Fastly-Key: $FASTLY_API_KEY" \
    --header "Fastly-Soft-Purge: 1" \
    "https://api.fastly.com/service/$FASTLY_SERVICE_ID/purge/site-$PROJECT"
fi
