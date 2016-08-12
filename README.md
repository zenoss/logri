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