# How To Run Project

## Prerequisites
Install Swag and Goose CLI tools:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
```

First, start the required services using Docker Compose:
```bash
docker compose up -d
```

Then, run the Go application:
```bash
go run main.go
```

### Create new migration file

```bash
goose create <table-name> sql
```

## Database Migration with Goose

### Apply all available migrations

```bash
goose up
# OK    001_basics.sql
# OK    002_next.sql
# OK    003_and_again.go
```

### Migrate up to a specific version

```bash
goose up-to 20170506082420
# OK    20170506082420_create_table.sql
```

### Migrate up by one step

```bash
goose up-by-one
# OK    20170614145246_change_type.sql
```

### Roll back a single migration

```bash
goose down
# OK    003_and_again.go
```

### Roll back to a specific version

```bash
goose down-to 20170506082527
# OK    20170506082527_alter_column.sql
```

### Roll back all migrations (⚠️ careful!)

```bash
goose down-to 0
```

### Show migration status

```bash
goose status
#   Applied At                  Migration
#   =======================================
#   Sun Jan  6 11:25:03 2013 -- 001_basics.sql
#   Sun Jan  6 11:25:03 2013 -- 002_next.sql
#   Pending                  -- 003_and_again.go
```

### Show current database version

```bash
goose version
# goose: version 002
```

## Contribution
  1. นพณัช สาทิพย์พงษ์ besterOz
  2. พงศธร รักงาน prukngan
  3. ภวัต เลิศตระกูลชัย Phawat Loedtrakunchai
  4. ธฤต จันทร์ดี tharitpr

## System architecture
  Microservice





