# v3 Improve Error Reporting – Port Access Fix

## Problem

The `*selenium.Service` type does not have a public `Port()` method in the version of the selenium library being used. The field `port` is unexported (lowercase), so it cannot be accessed directly. This causes the compilation error:

```
browser/browser.go:59:88: service.Port undefined (type *selenium.Service has no field or method Port, but does have unexported field port)
```

## Changes Made

### 1. `e2e-tests/browser/browser.go`

- **Lines 23-30:** Added a `port` variable and set it based on the browser type before creating the service:

```go
var port int

switch cfg.Browser {
case "chrome":
    port = 9515
    service, err = selenium.NewChromeDriverService(
        "/usr/local/bin/chromedriver",
        port,
        nil,
        selenium.Output(nil),
    )
case "firefox":
    port = 4444
    service, err = selenium.NewGeckoDriverService(
        "/usr/local/bin/geckodriver",
        port,
        nil,
        selenium.Output(nil),
    )
}
```

- **Line 59:** Changed `service.Port()` to `port`:

```go
wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
```

This avoids calling the non-existent `Port()` method and uses the port number we already know from the service creation.

## Result

The code no longer tries to access the unexported `port` field, resolving the compilation error.
