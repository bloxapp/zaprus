# zaprus

Ever had a 3rd-party dependency requiring a [logrus](https://github.com/sirupsen/logrus), but you're using [zap](https://github.com/uber-go/zap)?

`zaprus` provides a `logrus.Hook` that makes a `logrus.(Entry|Logger)` replicate it's logs to a given `zap.Logger`.

## Usage

```go
// Create a muted logrus.Logger
lr := logrus.StandardLogger()
lr.SetOutput(io.Discard) // Prevent double-logging

// Forward it's logs to a zap.Logger
hook := zaprus.NewHook(anyZapLogger)
lr.AddHook(hook)

// Pass it to the 3rd-party dependency
an_evil_3rd_party_dependency.New(lr)
```

Or simply...

```go
// Create a muted & zaprus.Hook'd logrus.Logger
lr := zaprus.NewProxyLogrus(anyZapLogger)

// Pass it to the 3rd-party dependency
an_evil_3rd_party_dependency.New(lr)
```