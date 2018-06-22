# Getting Started

```bash
APP_SECRET="app-secret" go run bencher/*.go
```

# TODO

* [X] Make one network call be successful
* [X] Make 500 concurrent network calls be successful
* [ ] Accept an arg `-c` for number of concurrency that runs til cancel
* [ ] Output summary report
* [ ] Accept an arg `-t` for how long, in seconds, for it to run
* [ ] Accept an arg `-d` for how long, in seconds, for delay between each request
