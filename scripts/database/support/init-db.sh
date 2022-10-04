db_name=$1

PGPASSWORD=canary

if [ $(psql -d postgres -U postgres -h localhost -t -c "SELECT COUNT(*) FROM pg_user WHERE usename='connectrn';") -eq 0 ]; then
	echo "Creating DB role 'connectrn'"
	psql -U postgres -h localhost -c "CREATE USER connectrn WITH PASSWORD 'canary';";
fi

dropdb -h localhost -U postgres --if-exists $db_name
createdb -h localhost -U postgres -T template0 -O "connectrn" $db_name
