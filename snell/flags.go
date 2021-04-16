package main

import "flag"

var forwardAddr = flag.String("http_forward_address", "http://localhost:8080/", "address to forward http requests")
