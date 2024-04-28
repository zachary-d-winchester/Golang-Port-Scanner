package main

import (
	"PortScanner/port"          // Package for connecting to ports
	"PortScanner/scan"          // Package for iterating through IP addresses in a subnet
	"flag"                      // Allows for command line arguments
	"fmt"                       // Formatted input/output
	"github.com/thediveo/netdb" // Package to check the ports and give what services run on them
	"strconv"
	"strings"
	"time"
)

func main() {
	oldHostName := ""
	host := flag.String("host", "127.0.0.1", "What host is being scanned?")
	CIDR := flag.String("CIDR", "32", "If the address is a subnet, what is the CIDR? (Please only list the integer. eg. 24 not /24)")
	help := flag.Bool("help", false, "Tells the program to print out all of the flags and their usages.")
	ps := flag.Int("ps", 1, "What port do you want to be the first you scan?")
	pe := flag.Int("pe", 0, "What port do you want to be the last you scan? (Must be greater or equal to the first port value!) Leave blank if you just want to scan one port.")
	concurrent := flag.Int("con", 100, "How many connections should be attempted at once?")
	flag.Parse()                                           // Above are all of the command line args that the user can set
	if *ps <= 0 || *pe < 0 || *ps > 65535 || *pe > 65535 { // Port validation
		fmt.Println("The numbers you input for ports are invalid. Please check again.")
	} else if *ps > *pe && *pe != 0 { // Makes sure that the starting port is lower than the ending port
		fmt.Println("Your starting port is greater than your ending port.\nPlease make sure that you input a starting port that is lower than or equal to your ending port.")
	} else if *help != false { // Checks the help flag
		flag.PrintDefaults()
	} else {
		res1 := strings.Split(*host, ".")
		_, err := strconv.Atoi(res1[0])
		if err == nil {
			startTime := time.Now() // Gets the time before the program runs.
			results := scan.Scan(*host, *CIDR, *ps, *pe, *concurrent)
			endTime := time.Now()                                          // Gets the time after the program runs.
			executionTime := endTime.Sub(startTime)                        // Subtracts the start time from the end time to get the time the program took to run.
			fmt.Printf(" Finished.\nScan finished in %v\n", executionTime) // The following are outputs for the function to provide the user with information about the scanned IP(s)
			fmt.Printf("There are %d closed ports\n", port.ClosedCounter)
			for _, outer := range results { // Formats the [][]Result structure and prints it out.
				for _, inner := range outer.PortsScanned {
					if inner.State != "Closed" { // Checks to see if the port is closed. Only prints the open ports.
						if oldHostName != inner.HostName { // For use in scanning a subnet. Allows the ports to also be ordered by the IP address scanned.
							fmt.Printf("%v\n", inner.HostName)
							oldHostName = inner.HostName
						}
						service := netdb.ServiceByPort(inner.Port, inner.Protocol)
						if service != nil { // Checks to see if the service can be identified.
							fmt.Printf("%5d/%s %s %s\n", inner.Port, inner.Protocol, inner.State, service.Name)
						} else {
							fmt.Printf("%5d/%s %s\n", inner.Port, inner.Protocol, inner.State)
						}
					}
				}
			}
		} else {
			startTime := time.Now() // Gets the time before the program runs.
			results := port.ScanningPorts(*host, *ps, *pe, *concurrent)
			endTime := time.Now()                                          // Gets the time after the program runs.
			executionTime := endTime.Sub(startTime)                        // Subtracts the start time from the end time to get the time the program took to run.
			fmt.Printf(" Finished.\nScan finished in %v\n", executionTime) // The following are outputs for the function to provide the user with information about the scanned IP(s)
			fmt.Printf("There are %d closed ports\n", port.ClosedCounter)
			for _, res := range results {
				if res.State != "Closed" { // Checks to see if the port is closed. Only prints the open ports.
					if oldHostName != res.HostName { // For use in scanning a subnet. Allows the ports to also be ordered by the IP address scanned.
						fmt.Printf("%v\n", res.HostName)
						oldHostName = res.HostName
					}
					service := netdb.ServiceByPort(res.Port, res.Protocol)
					if service != nil { // Checks to see if the service can be identified.
						fmt.Printf("%5d/%s %s %s\n", res.Port, res.Protocol, res.State, service.Name)
					} else {
						fmt.Printf("%5d/%s %s\n", res.Port, res.Protocol, res.State)
					}
				}
			}
		}
	}
}
