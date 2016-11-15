# TCP/UDP output plugin

Send beats trough a TPC or UDP connection to a remote endpoint

## Configuration

Set `udp_tcp` as the output.

```yaml
output:

  # The wanted output
  type: udp_tcp
```

### File based

Set the network address and port, and the protocol wanted:

* `output.udp_tcp.network`: `tcp` or `udp`
* `output.udp_tcp.raddr`: `<remote_host>:<port>`

If you want to configure a secure connection trough your endpoint, set the `output.udp_tcp.tls.enable` to `true` property 
to true and set:

* `output.udp_tcp.tls.enable`: `false` by default, `true` if you want to set a secure connection.
* `output.udp_tcp.tls.ca_path`: the path to the ca file.
* `output.udp_tcp.tls.cert_path`: the path to the cert file.
* `output.udp_tcp.tls.key_path`: `the path to the key file.


Here is the default configuration used:

```yaml
output:

  # The wanted output
  type: ${BFWD_OUTPUT_TYPE:udp_tcp}
  
   # UDP and TCP settings
    udp_tcp:
  
      # Protocol connection, by default tcp
      network: ${BFWD_UDPTCP_NETWORK:tcp}
  
      # The remote endpoint
      raddr: ${BFWD_UDPTCP_RADDR}
  
      # Secure communication settings, if enabled, configure path for key and certificates.
      tls:
        # By default, TLS is disabled
        enable: ${BFWD_ENABLE_TLS:false}
  
        # Path to the ca file
        ca_path: ${BFWD_ENABLE_TLS_CA_PATH}
  
        # Path to the cert file
        cert_path: ${BFWD_ENABLE_TLS_CERT_PATH}
  
        # Path to the key file
        key_path: ${BFWD_ENABLE_TLS_KEY_PATH}
```

### Environment variable based

Here are the env vars used to configure the plugin:

* Network endpoint: `BFWD_UDPTCP_NETWORK`
* Network protocol: `BFWD_UDPTCP_NETWORK`
* Enable the TLS connection:  `BFWD_ENABLE_TLS`
* Path to the CA file:  `BFWD_ENABLE_TLS_CA_PATH`
* Path to the CERT file:  `BFWD_ENABLE_TLS_CERT_PATH`
* Path to the Key file:  `BFWD_ENABLE_TLS_KEY_PATH`


### Example to Logmatic.io


```
mkdir demo-tcp && cd demo-tcp

# Get files


```