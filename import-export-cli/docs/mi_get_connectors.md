## mi get connectors

Get information about connectors deployed in a Micro Integrator

### Synopsis

List all the connectors deployed in a Micro Integrator in the environment specified by the flag --environment, -e

```
mi get connectors [flags]
```

### Examples

```
To list all the connectors
   mi get connectors -e dev
NOTE: The flag (--environment (-e)) is mandatory
```

### Options

```
  -e, --environment string   Environment to be searched
      --format string        Pretty-print using Go Templates. Use "{{ jsonPretty . }}" to list all fields
  -h, --help                 help for connectors
```

### Options inherited from parent commands

```
  -k, --insecure   Allow connections to SSL endpoints without certs
      --verbose    Enable verbose mode
```

### SEE ALSO

* [mi get](mi_get.md)	 - Get information about artifacts deployed in a Micro Integrator instance

