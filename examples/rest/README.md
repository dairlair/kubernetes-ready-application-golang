# REST Service example

This service provides a single "/welcome" endpoint.

Build it
```shell script
mame build
# Or for macOS
GOOS=darwin mame build
```

Run it
```shell script
./rest
```

Test it
```shell script
# Hit the application
curl http://localhost/welcome
# Hit the liveness and readiness probes
curl http://localhost:81/healthz -i
curl http://localhost:81/readyz -i
```