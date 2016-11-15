# Syslog output plugin

Send beats to a local or remote syslog daemon

## Configuration
Set `syslog` as the output.

```yaml
output:

  # The wanted output
  type: syslog
```

### File based

If no network protocol is provided, then the plugin try to send beats to the local service.
You can choose a remote destination by settings the protocol and the remote address as follow:

* `output.syslog.network`: `tcp` or `udp`
* `output.syslog.raddr`: `<remote_host>:<port>`

Do not forget to set an app trough the `output.syslog.tag` property

Here is the default configuration used:

```yaml
output:

  # The wanted output
  type: ${BFWD_OUTPUT_TYPE:syslog}

  # Syslog specific settings
  syslog:
    # Tag or application reported for each log
    tag: ${BFWD_SYSLOG_TAG:beats-fowarder}

    # Protocol connection (tcp or udp), if empty, the local syslog is used
    network: ${BFWD_SYSLOG_NETWORK}

    # The remote endpoint, only if a protocol has been specified
    raddr: ${BFWD_SYSLOG_RADDR}
```

### Environment variable based

Here are the env vars used to configure the plugin:

* Tag or appname:   `BFWD_SYSLOG_TAG`
* Network endpoint: `BFWD_SYSLOG_RADDR`
* Network protocol: `BFWD_SYSLOG_NETWORK`