Stream Control Transmission Protocol (SCTP)
----

THIS IS NOT AN OFFICIAL SCTP PACKAGE!

I merged this to fit my needs from github.com/ishidawataru/sctp and others. 



Examples
----

See `example/sctp.go`

```go
$ cd example
$ go build
$ # run example SCTP server
$ ./example -server -port 1000 -ip 10.10.0.1,10.20.0.1
$ # run example SCTP client
$ ./example -port 1000 -ip 10.10.0.1,10.20.0.1
```
