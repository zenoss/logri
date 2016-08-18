# Logri

Logri is a wrapper for [Logrus](https://github.com/Sirupsen/logrus) that
provides **hierarchical, configurable, structured logging**. 

Like Logrus, it's a drop-in replacement for Go's standard logging library, but
it adds the ability to:

* Define loggers that inherit their log levels and output streams from parent loggers
* Configure loggers from a YAML file
* Update configuration on the fly
* Optionally watch a logging configuration file for changes

And, of course, it retains Logrus's excellent structured logging capabilities.

## Usage

You can drop Logri in to replace logging or Logrus very simply:

```go
package main


import (
        "github.com/Sirupsen/logrus"
    log "github.com/zenoss/logri"
)

func main() {

    log.Infof("Logri can replace the %s package", "logging")

    log.WithFields(logrus.Fields{
        "package": "logrus",
    }).Infof("Or another popular logging package")
}
```

### Named loggers

The power of Logri comes in with named hierarchical loggers. Use
`logri.GetLogger(name)` with a dotted name to return an individual logger that
inherits its log level and outputs from its parent, but can add or override its
own.

```go
package main

import (
    "github.com/Sirupsen/logrus"
    "github.com/zenoss/logri"
)

var (
    pkglog = logri.GetLogger("package")
    cmplog = logri.GetLogger("package.component")
    subcmplog = logri.GetLogger("package.component.subcomponent")
)

func main() {

    pkglog.SetLevel(logrus.DebugLevel, true) // Second argument makes it inherited
    // package.component and package.component.subcomponent are also Debug level now

    // Quiet package.component down but leave subcomponent at debug
    cmplog.SetLevel(logrus.ErrorLevel, false) // Second argument false means
                                              // local to this logger only
}
```

Further calls to `logri.GetLogger(name)` will retrieve the same logger
instance, so there's no need to jump through hoops exporting loggers to share
them among packages.

### Configuration via file

You can also configure Logri using a YAML file.