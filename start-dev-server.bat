@ECHO OFF

SetLocal
SetLocal enableextensions

IF "%SDK%"=="" (
	Echo You must set the SDK environment variable.
	Exit /B
)

md "development-data"
md "development-data/blobstore"

python -O "%SDK%/dev_appserver.py" --dev_appserver_log_level debug --show_mail_body --blobstore_path=development-data/blobstore --datastore_path=development-data/datastore --prospective_search_path=development-data/prospective_search --search_indexes_path=development-data/searchindexes %* .

EndLocal
