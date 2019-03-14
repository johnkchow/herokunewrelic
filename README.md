[![Build Status](https://travis-ci.org/johnkchow/herokunewrelic.svg?branch=master)](https://travis-ci.org/johnkchow/herokunewrelic)

# Getting Started

```bash
dep ensure
go build
export NEW_RELIC_STREAMING_LICENSE_KEY=key
export NEW_RELIC_APP_NAME="Your Heroku App Name"
export AUTH_SECRET=auth_secret
./herokunewrelic
```

To execute tests:

```bash
go test -v
```

**NOTE**: The commands above assumes that you have both `go@1.9.7` and `dep` installed. For further instructions on how to install the dependencies:
* https://golang.github.io/dep/docs/daily-dep.html
* https://github.com/moovweb/gvm

# How to setup on Heroku

Supposed that:

* Your `herokunewrelic` instance is accessible via https://herokunewrelic.herokuapp.com.
* You have a Heroku app name called `heroku-app`

Simply execute the following command:

```bash
APP_HOST=herokunewrelic.herokuapp.com
HEROKU_APP_NAME=heroku-app
heroku drains:add https://$HEROKU_APP_NAME:$AUTH_SECRET@$APP_HOST -a HEROKU_APP_NAME
```

**NOTE**: The username for the basic auth part must will become the `sourceAppName` property for the custom NewRelic events.

# TODO

* [X] NewRelic
* [X] Parse arbitrary metrics in addition to dyno runtime metrics
* [ ] Honor env var LOG_LEVEL
* [ ] Tagged logs with Request ID
* [ ] Be idempotent (Logplex-Frame-Id)
* [ ] Support custom event backends
* [ ] Turn parseFrames into iterator to further reduce memory footprint
* [X] Operations
  * [X] Memory/CPU usage
  * [X] GC times
  * [X] Go routines
* [X] Support for multiple log drain tokens
* [ ] Logging for future debugging e.g. `User-Agent` since it maps to release version
