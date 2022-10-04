sh scripts/database/support/init-db.sh connectrn_test
PGPASSWORD=canary psql -h localhost -U connectrn connectrn_test -f scripts/database/sql/table-scripts.sql
go clean -testcache
DB_HOST=localhost DB_NAME=connectrn_test TESTDATA_PATH=$(pwd)/testdata gotestsum ./...
