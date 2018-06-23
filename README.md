[![Build Status](https://travis-ci.org/johnkchow/herokunewrelic.svg?branch=master)](https://travis-ci.org/johnkchow/herokunewrelic)

# Getting Started

```bash
dep ensure
go build
export NEW_RELIC_LICENSE_KEY=key
export NEW_RELIC_APP_NAME="Your Heroku App Name"
export AUTH_SECRET=auth_secret
./herokunewrelic
```
# TODO

* [X] NewRelic
* [ ] Tagged logs with Request ID
* [ ] Be idempotent (Logplex-Frame-Id)
* [X] Operations
  * [X] Memory/CPU usage
  * [X] GC times
  * [X] Go routines
* [X] Support for multiple log drain tokens
* [ ] Logging for future debugging e.g. `User-Agent` since it maps to release version
