# gh-fetch ![badge](https://github.com/karthikraobr/gh-fetch/workflows/Go/badge.svg)


This repository contains the backend code which fetches the public github repositories of a user/orgnization and its commits and stores it in a datastore. It also caches the github api results for upto 60 seconds before being invalidated.

### Requirements
- go
- docker-compose

### Setup
The make file contains all the commands needed to build/run/mock/test the application. Some commands which are not so obvious:

- `make compose` - one command to rule it all. Build and run the application. Spins up the application and the database server.

- `make mock` - (re)generate the mocks required for testing.

### URLs
- `/user/:username/repositories` - Fetches the public repositories of a user. Optionally query paramaters `page` and `perpage` can be supplied to paginate results. Please note that the paginations works properly only when `cache` is empty. 
e.g. - http://localhost:8000/user/karthikraobr/repositories
- `/user/:username/repository/:repository/commits` - Fetches the commits of a particular repository. Optionally query paramaters `page` and `perpage` can be supplied to paginate results. Since this endpoint does not use a cache or a datastore, the pagination works all the time.
e.g. - http://localhost:8000/user/karthikraobr/repository/gqlgen/commits
- `/user/:username/top20` - Fetches Top 20 recently accessed repositories based on the `last_access` column.


### What is missing?
- Frontend
- Caching/Fetching from the datastore does not paginate the results. Due to time constraints I could not add pagination during these scenarios.
- `package Store` does not contain tests.
- `CI` could have been better.

### What might have gone wrong?
- First time `gin-gonic` and `gorm` user, hence best practices might have taken a backseat.
