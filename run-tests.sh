sh scripts/database/support/init-db.sh demo_test
PGPASSWORD=canary psql -h localhost -U demo demo_test -f scripts/database/sql/table-scripts.sql
go clean -testcache
DB_HOST=localhost DB_NAME=demo_test TESTDATA_PATH=$(pwd)/testdata gotestsum ./...
