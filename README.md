# Golang-Port-ScannerA custom port scanner written entirely in golang. Attempts to connect toports and displays if they can be connected to. Allows a user to scan asingle IP, a subnet, or a hostname.Due to the program's ability to attempt many connections at once, it isfaster than a program such as nmap, however the program is very noisy andshould not be used for a penetration test.## Flags-help: Displays the help message.-CIDR (string): Used if scanning a subnet. Only pass the integer into. -con (int): How many connections should be attempted at once.-host (string): Used to declare what host/subnet your are scanning.-ps (int): Used to declare the starting port that is being scanned.-pe (int): Used to declare the ending port that is being scanned.example: ./portscanner.exe -host 127.0.0.1 -con 100 -ps 1 -pe 1023