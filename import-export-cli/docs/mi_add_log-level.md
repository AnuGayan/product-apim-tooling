## mi add log-level

Add new Logger to a Micro Integrator

### Synopsis

Add new Logger named [logger-name] to the [class-name] with log level [log-level] specified by the command line arguments to a Micro Integrator in the environment specified by the flag --environment, -e

```
mi add log-level [logger-name] [class-name] [log-level] [flags]
```

### Examples

```
To add a new logger
   mi add log-level synapse-api org.apache.synapse.rest.API DEBUG -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment of the micro integrator to which a new logger should be added
  -h, --help                 help for log-level
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [mi add](mi_add.md)	 - Add new users or loggers to a Micro Integrator instance

