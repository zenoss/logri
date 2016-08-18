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
    log "github.com/zenoss/logri"
)

func main() {

    log.Infof("Logri can replace the %s package", "logging")

    log.WithFields(log.Fields{
        "package": "logrus",
    }).Infof("Or another popular logging package")
}
```