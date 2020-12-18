package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/kv0s/openhpsdr"
)

// discover, a program to discover SDR radios using openhpsdr protocol 1 & 2
// and tangerinesdr protocol.
//
// This prgram is intended to be use in a script to list the available radios on the
// network.

//Listinterface is a convenience function to print interface data
func Listinterface(itr openhpsdr.Intface, verbose bool) {
	if verbose {
		fmt.Printf("          Computer: (%v)\n", itr.MAC)
		fmt.Printf("                OS: %s (%s) %d CPU(s)\n", runtime.GOOS, runtime.GOARCH, runtime.NumCPU())
		fmt.Printf("              IPV4: %v\n", itr.Ipv4)
		fmt.Printf("              Mask: %d\n", itr.Mask)
		fmt.Printf("           Network: %v\n", itr.Network)
		fmt.Printf("              IPV6: %v\n\n", itr.Ipv6)
	} else {
		fmt.Printf("%s (%v) %s %s %d ", itr.Intname, itr.MAC, runtime.GOOS, runtime.GOARCH, runtime.NumCPU())
		fmt.Printf("%v %d %v %v\n", itr.Ipv4, itr.Mask, itr.Network, itr.Ipv6)
	}
}

// Listboard is a convenience function to print board data
func Listboard(str openhpsdr.Hpsdrboard, verbose bool) {
	if str.Macaddress != "0:0:0:0:0:0" {
		if verbose {
			fmt.Printf("       HPSDR Board: (%s)\n", str.Macaddress)
			fmt.Printf("              IPV4: %s\n", str.Baddress)
			fmt.Printf("              Port: %s\n", str.Bport)
			fmt.Printf("              Type: %s\n", str.Board)
			fmt.Printf("          Firmware: %s\n", str.Firmware)
			fmt.Printf("            Status: %s\n\n", str.Status)
			fmt.Printf("            PC    : %s\n\n", str.Pcaddress)
		} else {
			fmt.Printf("(%s) %s %s", str.Macaddress, str.Baddress, str.Bport)
			fmt.Printf(" %s %s %s %s\n", str.Board, str.Protocol, str.Firmware, str.Status)
		}
	}
}

const version string = "0.1.0"
const started string = "2020-10-27"
const update string = "2020-10-27"

func program() {
	fmt.Printf("discover  version:(%s)\n", version)
	fmt.Printf("    By Dave KV0S, %s, GPL3 \n", started)
	fmt.Printf("    Last Updated: %s \n\n", update)
}

func usage() {
	fmt.Printf("discover program to find SDR radios on your network.\n")
}

func main() {
	ifn := flag.String("interface", "none", "Select one interface")
	veb := flag.Bool("verbose", false, "Select true or false")

	flag.Parse()

	if flag.NFlag() < 1 {
		program()
		usage()
	}

	intf, _ := openhpsdr.Interfaces()

	if flag.NFlag() < 1 {
		fmt.Printf("\nInterfaces on this Computer: \n")
	}

	for i := range intf {
		if flag.NFlag() < 1 {
			// if no flags list the interfaces in short form
			Listinterface(intf[i], *veb)
		} else if (flag.NFlag() == 1) && (*ifn == "none") {
			// if one flag and it is debug = none, list the interface in short form
			Listinterface(intf[i], *veb)
		}

		var adr string
		var bcadr string
		adr = intf[i].Ipv4 + ":1024"
		bcadr = intf[i].Ipv4Bcast + ":1024"

		// perform a discovery
		str, err := openhpsdr.Discover(adr, bcadr, 1, "none")
		if err != nil {
			fmt.Println("Error ", err)
		}

		//var bdid int
		//bdid = 0
		//loop throught the list of discovered HPSDR boards
		for i := 0; i < len(str); i++ {
			//if fg.SelectMAC == str[i].Macaddress {
			// if a MAC is selected
			fmt.Printf("      Selected MAC: %s\n", str[i].Board)
			Listboard(str[i], *veb)
			//bdid = i
			//{}
		}

		fmt.Println()
	}
}
