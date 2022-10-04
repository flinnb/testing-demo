sh scripts/database/support/init-db.sh connectrn_test
PGPASSWORD=canary psql -h postgres -U connectrn connectrn_test -f scripts/database/sql/table-scripts.sql
go clean -testcache
DB_HOST=postgres DB_NAME=connectrn_test TESTDATA_PATH=$(pwd)/testdata gotestsum ./...
