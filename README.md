# mundiclient

Client for the mundi laser printer.

It is setup to communicate over tcp, so you need to convert the serial port to ethernet. In this case we used https://www.antratek.nl/cse-h21

The implementation is done based on the documentation provided, see file `./906-0152-00.pdf`. Post 2.10.001 is used.

Getting started:
```go
package main

import (
  "fmt"
  "github.com/naipath/mundiclient"
)

func main() {
  client, _ := mundiclient.New("192.168.1.171", 1470)
  defer client.Close()
  
  fmt.Println(client.GetVersion())
}
```

Implemented commands:
- Get version
- External field
- Upload logo image
- Select file for marking
- Request list of marking files
- Request current marking file
- Get laser parameter
- Set laser parameter
- Get file (download file from mundi scan)
- Get line settings
- Set line settings
- Get current counters
- Reset current counters
- Incremental counters (Get and Set)
- Get status data
- Get status message
- Get field contents
- Offset field
- System Time and date

Not yet implemented: 
- Alter Logo Text Fields
- Sending marking file
- Get history files