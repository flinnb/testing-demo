db_name=$1

PGPASSWORD=canary

if [ $(psql -d postgres -U postgres -h localhost -t -c "SELECT COUNT(*) FROM pg_user WHERE usename='demo';") -eq 0 ]; then
	echo "Creating DB role 'demo'"
	psql -U postgres -h localhost -c "CREATE USER demo WITH PASSWORD 'canary';";
fi

dropdb -h localhost -U postgres --if-exists $db_name
createdb -h localhost -U postgres -T template0 -O "demo" $db_name
