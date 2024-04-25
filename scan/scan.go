package scan

import (
	"PortScanner/port"
	"bytes"
	"fmt"
	gen "github.com/korylprince/ipnetgen"
	"math"
	"net"
	"os"
	"sort"
	"strconv"
)

type Scans struct {
	HostName     string
	PortsScanned []port.Result
}

var (
	results []Scans
)

func worker(jobs <-chan string, resultsC chan<- Scans, ps int, pe int, concurrent int) {
	for ipAddress := range jobs {
		workerResult := Scans{}
		workerResult.PortsScanned = port.ScanningPorts(ipAddress, ps, pe, concurrent)
		workerResult.HostName = ipAddress
		resultsC <- workerResult
	}
}

func Scan(host string, CIDR string, ps int, pe int, concurrent int) []Scans {
	intCIDR, _ := strconv.Atoi(CIDR)
	freeBits := float64(32 - intCIDR)
	exponent := int(math.Pow(2, freeBits))
	if exponent != 1 {
		exponent -= 2
	}
	jobs := make(chan string, exponent)
	resultsC := make(chan Scans, exponent)
	for i := 1; i <= concurrent; i++ {
		go worker(jobs, resultsC, ps, pe, concurrent)
	}
	fullCIDR := host + "/" + CIDR // Builds out the CIDR formatted ip address block
	ipGen, err := gen.New(fullCIDR)
	if err != nil { // CIDR validation
		fmt.Printf("Your CIDR is invalid. Please put in an integer ranging from 0 to 32.\n")
		os.Exit(1)
	}
	fmt.Printf("Port Scanning...")
	for ip := ipGen.Next(); ip != nil; ip = ipGen.Next() { // Iterates through the list of IP addresses in the given block and uses them as inputs for a function to scan them.
		ipAddress := net.IP.String(ip) // Turns the net.IP ip variable into a string that can be used in the port.InitialScan function.
		if ip[3] != 0 && ip[3] != 255 {
			jobs <- ipAddress
		}
	}
	close(jobs)
	for i := 0; i < exponent; i++ {
		results = append(results, <-resultsC)
	}
	close(resultsC)
	sort.SliceStable(results, func(i, j int) bool {
		ip1 := net.ParseIP(results[i].HostName)
		ip2 := net.ParseIP(results[j].HostName)
		return bytes.Compare(ip1, ip2) < 0
	})
	return results
}
