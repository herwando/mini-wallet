# SKU Management

[SKU Management as A Service](https://bukalapak.atlassian.net/wiki/spaces/BL20/pages/2207758194/Draft+SKU+Management+as+A+Service)

## Setup
- Run dependency
```
make docker-up
```
- Download database migration tools
```
make tool-migrate
```
- Run migration for each module. See makefile for custom database env config. Also see more argument on [golang-migrate](https://github.com/golang-migrate/migrate)
```
MIGRATE_ARGS=up make migrate
```
- Run golang main on cmd
```
go run cmd/rpc/main.go
```

## Useful Links
- [Core Services Developer Tutorials](https://bukalapak.atlassian.net/wiki/spaces/CIS/pages/2248710958/Developer+Tutorials)
- [Core Services Developer Guidelines](https://bukalapak.atlassian.net/wiki/spaces/CIS/pages/2225867497/Developer+Guidelines)
- [Developing a Domain](https://bukalapak.atlassian.net/wiki/spaces/CIS/pages/2248875394/Developing+a+Domain)
