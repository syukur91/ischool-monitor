
# Development

## Contributing Guide

Several things that I would like to try using on this project:

1. Git flow
    1. public branch (master and development) should be linear and clean
    1. private branch can be in any type. do pull request to master when ready
    1. rebase and squash before pull request
1. Don't push .env to repo (12-factor violation)
1. Don't push artifact to repo (debug binary, executable)
1. Use sql (no ORM)
1. Use PostgreSQL materialized view for query (if needed)
1. Don't mess with migration after deployed
1. Use normalized table. For 1-n and n-m relationship, use relation table.
1. Use go mod (not dep). Use Go 1.11 up 
1. Write tests

## Project structure

1. `api/controller`: contains HTTP request handler
1. `api/schema`: contains request and response schema. We are trying to separate database structure with request and response structure. So both not necessarily same
1. `migration`: contains SQL migrations. See `Makefile`
1. `pkg`: contains some small libraries
1. `service`: contains logic and data access part

## Create database and user

```
$ createuser school
$ createdb school

$ psql
# alter user school with encrypted password 'school123';
# GRANT ALL PRIVILEGES ON DATABASE school to school;
```

## Migration

Download migrate binary. I put to bin inside the project

```
curl -L https://github.com/golang-migrate/migrate/releases/download/v3.5.2/migrate.darwin-amd64.tar.gz | tar xvz
```
Rename the binary to `migrate`

In project folder
```
// up
make migrate-up

// down
make migrate-down
```

## Running

```
git clone
export GO111MODULE=on
go build
make test
make run
```

## Git workflow

### Main branch

- `master` is production release branch
- `development` is latest development changes to be released in production (merged to master)

Those are public branch and don't create new branch with above name

### Supporting branches

Depending on your task create branch with following convention

- `feature-*` to add new feature. Example `feature-add-tier`
- `release-*` to prepare new release for production. Example `release-2.0.0`
- `hotfix-*` to fix issue in `master`. Example `hotfix-crash-when-delete`

See https://nvie.com/posts/a-successful-git-branching-model

### Merge with pull request

Do rebase and/or squash before doing pull request to development of master branch.

## Tests

Write test for `services` and `pkgs` if that function is new and untested. See https://www.youtube.com/watch?v=a6oP24CSdUg for why we need to write tests

To run test 
```
make test
```

## Code generation

We have code generation that will create `schema`, `controller` and `service` go file. To start, edit [template/cmd/gen.go](template/cmd/gen.go) `models` with model that you want to generate. 

```
cd template
go generate
```

**IMPORTANT** Be careful not to run go generate again to generated file that you have edited as it will overwrite the content. You can set `Skip=true` to disable code generation for specific model

# License

Proprietary