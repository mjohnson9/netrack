@ECHO OFF

SetLocal
SetLocal enableextensions

md "development-data"
md "development-data/blobstore"

python -O "C:\App Engine\Go\dev_appserver.py" --show_mail_body --blobstore_path=development-data/blobstore --datastore_path=development-data/datastore --prospective_search_path=development-data/prospective_search --search_indexes_path=development-data/searchindexes %* .

EndLocal
