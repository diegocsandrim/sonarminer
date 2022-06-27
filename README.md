# sonarminer

SonarMiner is a tool to automate SonarQube analysis over time in Git repositories.


# Setup

go build .
docker-compose up -d

```sh
sudo sysctl -w vm.max_map_count=262144 # or edit /etc/sysctl.conf
docker-compose up
```

# Usage

```sh
./sonarminer analyse diegocsandrim/sonarminer
```

## Data access

Basic data can be accessed with SQL:

```
cat queries/basic.sql | docker exec -i sonarminer_db_1 psql -U sonar
```


# Clean up

```
docker-compose down
```
