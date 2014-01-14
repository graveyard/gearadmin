// Package gearadmin provides simple bindings to the gearman admin protocol: http://gearman.org/protocol/.
//
//
// Usage
//
// Here's an example program that tails all operations starting from the time the program is launched:
//
//         package main
//
//         import (
//         	"fmt"
//         	"github.com/Clever/gearadmin"
//         	"net"
//         )
//
//         func main() {
//         	c, err := net.Dial("tcp", "localhost:4730")
//         	if err != nil {
//         		panic(err)
//         	}
//         	defer c.Close()
//         	gearadmin := gearadmin.NewGearmanAdmin(c)
//         	status, _ := gearadmin.Status()
//         	fmt.Printf("%#v\n", status)
//         }
package gearadmin
