package main

import (
	nr "github.com/newrelic/go-agent"
	"net/http"
	"os"
)

func main() {
	appName := getEnv("NEW_RELIC_APP_NAME", "HEROKU_APP_NAME")
	licenseKey := getEnv("NEW_RELIC_LICENSE_KEY")
	//logplexDrainToken := getEnv("LOGPLEX_DRAIN_TOKEN")
	authSecret := os.Getenv("AUTH_SECRET")

	config := nr.NewConfig(appName, licenseKey)
	if os.Getenv("APP_ENV") != "production" {
		config.Logger = nr.NewDebugLogger(os.Stdout)
	}
	app, err := nr.NewApplication(config)

	if err != nil {
		panic(err.Error())
	}

	logplexHandler := newLogplexHandler(app)

	http.HandleFunc(nr.WrapHandleFunc(app, "/", func(rw http.ResponseWriter, req *http.Request) {
		/*
			TODO: Things to consider:
			  [X] Drain Token Security
			  [X] NewRelic
				[ ] Tagged logs with Request ID
			  [ ] Be idempotent (Logplex-Frame-Id)
			  [ ] Operations
				    [ ] Memory/CPU usage
				    [ ] GC times?
				    [ ] Go routines?
				[ ] Support for multiple log drain tokens
			  [ ] Logging for future debugging e.g. `User-Agent` since it maps to release version
		*/

		//token := req.Header.Get("Logplex-Drain-Token")

		//if token != logplexDrainToken {
		//logger.Debugf("Bad token password. '%s' does not match '%s'", token, logplexDrainToken)
		////rw.WriteHeader(http.StatusUnauthorized)
		//return
		//}

		if authSecret != "" {
			_, pass, ok := req.BasicAuth()

			if !ok || authSecret != pass {
				logger.Debugf("Bad auth secret. '%s' does not match '%s'", pass, authSecret)
				//rw.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		logplexHandler(rw, req)
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "7000"
	}
	logger.Fatal(http.ListenAndServe(":"+port, nil))
}

func getEnv(mainKey string, otherKeys ...string) string {
	otherKeys = append([]string{mainKey}, otherKeys...)
	for _, k := range otherKeys {
		name := os.Getenv(k)

		if name != "" {
			return name
		}
	}

	panic("Env var for " + mainKey + " is missing. Please make sure that env var is set")
}