# nomadcore
## Running
### Local Build
```
docker build -t nomadcore:latest .
docker run -p 8080:8080 -v $(pwd)/config.yml:/app/config.yml nomadcore:latest
```

### GHCR
```
docker run -d -p 8081:8081 -v $(pwd)/config.yml:/app/config.yml ghcr.io/wanderfusion/nomadcore:main
```

## DB
- Postgres 15
- Migrations using https://github.com/jackc/tern
