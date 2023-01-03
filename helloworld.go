package main

import (
	"fmt"
	"encoding/hex"
)

func main() {
	var succIdHex,_=hex.Decode([]byte("AASD08924X"))
	var currIdHex,_=hex.Decode([]byte("AASD08924Z"))
	fmt.Println("La distanza é %d",hex.Encode(succIdHex[]-currIdHex[]))
}
