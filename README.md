# mundiclient

This is a client for the mundi laser printer.

It is setup to communicate over tcp, so you need to convert the serial port to ethernet. In this case we used https://www.antratek.nl/cse-h21

Getting started:
```go
package main

import (
  "fmt"
  "github.com/naipath/mundiclient"
)

func main() {
  client := mundiclient.New("192.168.1.171", 1470)
  
  fmt.Println(client.GetVersion())
  
  client.Close()
}

```
