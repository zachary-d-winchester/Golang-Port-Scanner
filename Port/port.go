package port // Package port

import (
	log "github.com/sirupsen/logrus" // Adds advanced logging functionality using the logrus package.
	"net"
	"sort"
	"strconv"
	"sync"
	"time"
)

type Result struct { // Creates a structure to be called later for implementation into an array. Allows the array to easily store both the port and state of the port in each element
	HostName string
	Port     int
	Protocol string
	State    string
}

var (
	mutex         sync.Mutex // Mutual Exclusion. Locks a variable so that concurrent functions do not encounter race conditions.
	ClosedCounter int        // Used to count how many ports are determined to be closed so that they do not have to be printed.
)

func worker(jobs <-chan int, resultC chan<- Result, hostname string) {
	for i := range jobs {
		func(port int) {
			result := ScanPort("tcp", hostname, port)
			if result.State == "Closed" {
				mutex.Lock()
				ClosedCounter++ // This needs to be locked because many concurrent functions might try to increment it at once.
				mutex.Unlock()
			}
			resultC <- result
		}(i)
	}
}

func ScanPort(protocol, hostname string, port int) Result { // Function that takes a protocol, hostname, and port. Returns as the Result structure.
	result := Result{Port: port} // Sets the Port element to the port number taken in by this function.
	result.HostName = hostname
	result.Protocol = protocol                                     // Sets the Protocol element to the protocol type taken in by this function.
	address := hostname + ":" + strconv.Itoa(port)                 // Takes the hostname (represented as an ip), concatenates a ':' (signifies a socket) to it, and then turns the port number into a string so that it can be concatenated to the rest of the address. Stores in address.
	conn, err := net.DialTimeout(protocol, address, 1*time.Second) /* DialTimeout is a function in the net package. It takes a protocol, address, and a timeout duration. It attempts to connect to the address
	 									using the given protocol. Returns conn (which reads data from the connection), and error. If there is no error, then it returns nil for error.
		  								conn and error are then assigned to conn and err in this function to be used later.*/

	if err != nil { // If there is an error, run this:
		result.State = "Closed" // Sets the state in this port's element to "Closed".
		return result           // Returns the element.
	}
	defer func(conn net.Conn) { // Waits for surrounding functions to return before this function executes. Tries to close the connection that was established previously.
		err := conn.Close() // Stores any produced errors from trying to close the connection.
		if err != nil {     // If there is an error, print that an error was encountered.
			log.Error("Connection close error") // If the connection fails to close, log an error.
		}
	}(conn) // Immediately Invoked Function. Right after the above function is declared, we invoke it with conn as an argument.
	result.State = "Open" // Should no errors be encountered, set the state attribute of the element in go to Open.
	return result         // Returns the element of the array.
}

func ScanningPorts(hostname string, ps int, pe int, concurrent int) []Result { // Takes a hostname, starting port, ending port, and how many concurrent connections there should be.
	var results []Result
	ports := pe - ps + 1 // Finds how many ports are between the starting and ending port. Inclusive
	if pe == 0 {         // This enables the program to only scan one port.
		ports = 1
	}
	jobs := make(chan int, ports)        // Creates a jobs channel with a buffer size of ports
	resultsC := make(chan Result, ports) // Creates a resultsC channel with a buffer size of ports
	for i := 1; i <= concurrent; i++ {   // Starts a number of worker functions equal to the number of concurrent connections requested
		go worker(jobs, resultsC, hostname)
	}
	for i := 0; i < ports; i++ { // Sends all of the ports needed to be scanned to the jobs channel
		jobs <- ps + i
	}
	close(jobs)
	for i := 1; i <= ports; i++ { // Takes the results stored in the resultsC channel and converts it to an array
		results = append(results, <-resultsC)
	}
	close(resultsC)
	sort.SliceStable(results, func(i, j int) bool { // Sorts the array by their port number, lowest to highest.
		return results[i].Port < results[j].Port
	})
	return results // Return the results array.
}
