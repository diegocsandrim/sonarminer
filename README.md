# sonarminer

SonarMiner is a tool to automate SonarQube analysis over time in GitHub repositories.


# Requirements

- Go
- Docker
- Docker compose

# Setup

```sh
sudo sysctl -w vm.max_map_count=262144 # or edit /etc/sysctl.conf
go build .
docker compose up -d
```

# Usage

```sh
./sonarminer analyse diegocsandrim/sonarminer
```

## Data access

Basic data can be accessed with SQL:

```
cat queries/basic.sql | docker compose exec -T db psql -U sonar
```


# Clean up

```
docker compose down
```
