#!/bin/bash

if [ -z "${SDK}" ]; then
	echo "You must set the SDK environment variable."
	exit 1
fi

mkdir -p "development-data"
mkdir -p "development-data/blobstore"

exec python -O "${SDK}/dev_appserver.py" -d \
     --show_mail_body \
     --high_replication \
     --blobstore_path=development-data/blobstore \
     --datastore_path=development-data/datastore \
     --history_path=development-data/datastore.history \
     --search_indexes_path=development-data/searchindexes \
     $@ \
     .
