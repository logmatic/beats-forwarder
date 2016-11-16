# Logmatic output plugin

Send beats to Logmatic.io.

## Configuration
Basically, all you need to do is setting your API KEY and choose `logmatic` as the output.

```yaml
output:

  # The wanted output
  type: logmatic
```


### File based

To set the Logmatic.io key use `output.logmatic.key`.
By default, the plugin use a SSL connection trough the 10515 port. You can override 
these settings via the follow properties:

* `output.logmatic.network`: `tcp` or `ssl`
* `output.logmatic.raddr`: `<remote_host>:<port>`

Here is the default configuration used:

```yaml
output:

  # The wanted output
  type: ${BFWD_OUTPUT_TYPE:logmatic}

  # Logmatic specific settings
  logmatic:

    # The Logmatic API Key for authentification
    key: ${BFWD_LOGMATIC_API_KEY}

    # Protocol connection (tcp|ssl), by default ssl
    network: ${BFWD_LOGMATIC_NETWORK:ssl}

    # The remote endpoint
    raddr: ${BFWD_LOGMATIC_RADDR:"api.logmatic.io:10515"}
```

### Environment variable based

Here are the env vars used to configure the plugin:

* Logmatic.io key:  `BFWD_LOGMATIC_API_KEY`
* Network endpoint: `BFWD_LOGMATIC_RADDR`
* Network protocol: `BFWD_LOGMATIC_NETWORK`