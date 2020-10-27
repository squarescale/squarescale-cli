#!/bin/sh

set -e

if [ -z "$SQSC_TOKEN" ]; then
  echo "SQSC_TOKEN is not set. Quitting."
  exit 1
fi

CMD="/sqsc $@"
echo $CMD
eval $CMD
