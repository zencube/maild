# Maild
## Modern mail server for SMTP and IMAP

## Warning

**This is not ready yet!**
We're in the process of building this mail server.
Currently it only accepts SMTP connections and does neither send nor store emails it receives.
It also does not support ESMTP, so no TLS support is present *yet*.

## Prerequisites
* Postgresql.
* Go 1.2+

## Running the service
If you are on a linux x64 host, just run:

```shell
  $ ./maild
```

## Building from source
With Go 1.2+ available run:

```shell
  $ make build
```

## Running it for development
```shell
  $ make run
```

## License
This software is MIT licensed. See LICENSE for details.
