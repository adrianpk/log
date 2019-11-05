# Package log

A structured and context logger based on [zerolog](https://github.com/rs/zerolog).

## Quick Guide - draft

### TL;DR
```go
// Instantiate
l := log.NewLogger(log.Debug, "app-name", "always-add", "this-field")

// Add context values
l.SetDyna("some-contextual", "text", "another-context-num", 2)

// Log
log.Info("I'm ready to log", "this", "value", "and-this-number", 5, "right?", true)
```

**To get this output**
```json
{"level":"info","name":"app-name","always-add":"this-field","some-contextual":"text","another-contextual-num":2,"right?":true,"this":"value","and-this-number":5,"time":"2019-06-15T19:31:54+02:00","message":"I'm ready to log"}
```


### Static fields

```go
// Instatiate a new logger
name := "my-service"
rev := "5844c5a725"
l := log.NewLogger(level, name, "revision", rev)
```
Level could be values from 0 to 3 or for better readability:
`log.Debug`, `log.Info`, `log.Warn`, `log.Error`.

If name is not empty it will be appended to the logger normal output (i.e.: `"name": "my-service"`)

Additional pair of values will be appended to the logger normal output (i.e.: `"revision": "5844c5a725"`)

In this case `"revision": "5844c5a725"` but it could be any other pair of values: `..."api-version", "v1.0.0", "revision", "5844c5a725", "env", "stage")`.

Name and additional key-value pairs are considered static fields. Once set, it is not possible to modify them until restart.


### Dynamic fields
```go
// Previous steps ommited

// reqID and ref: some sample values to log
reqID := r.Header.Get("X-Request-ID")
if reqID == "" {
  reqID = uuid.New().String()
  r.Header.Set("X-Request-ID", rid)
}

ref := req.Referrer()

l.SetDyna("reqID", reqID, "referer", ref)
```

In addition to static key-values, the logger will add these tuples to the output.
These values can be overwritten and are unique to each instance of the logger.

Particularly useful so that, for example, a middleware can set contextual values such as request ID, session ID, etc.

### Logging
```go
//...
log.Info("Hello!", "name", "Clark", "lastname", "Kent")
```

Will produce an output like this
```json
{"level":"info","name":"my-service","revision":"5844c5a725","referer":"https://www.google.com/search?&q=nice+logger","reqID":"48f14c8c1fc854sa12bfd012a58f90a4","name":"Clark","lastname":"Kent","time":"2019-06-15T16:32:01+02:00","message":"Hello!"}
```

First argument is used as the log message `"message":"Hello!"` then each pair of additional values are added to the output of the log as can be seen in the example.


(cont...)
