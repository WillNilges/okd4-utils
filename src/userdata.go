package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

// Host is a representation of a Node, be it service, bootstrap, master, or compute
type Host struct {
	Hostname string
	Ipaddr   string
}

// Automatically detect the network devices and IP addresses of the host.
// Interactively let the user choose which address to use.
func chooseIP() (IP string, err error) {
	// Detect interfaces and let the user pick which one they want.
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
	}
	for index, iface := range ifaces {
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(index, ": ", iface.Name)
	}
	var option int
	for {
		fmt.Print("Select your primary interface: ")
		_, err := fmt.Scanf("%d", &option)
		if err == nil {
			break
		} else {
			fmt.Println("Invalid input detected.")
		}
	}

	// With their selection, detect IP addresses and let the user pick which one they want.
	addrs, err := ifaces[option].Addrs()
	if err != nil {
		fmt.Println(err)
	}
	for index, addr := range addrs {
		fmt.Println(index, ": ", addr)
	}
	var addrOption int
	for {
		fmt.Print("\nSelect your IP address: ")
		_, err := fmt.Scanf("%d", &addrOption)
		if err == nil {
			break
		} else {
			fmt.Println("Invalid input detected.")
		}
	}
	IP = addrs[addrOption].String()
	if strings.Contains(IP, "/") {
		IP = strings.Split(IP, "/")[0]
	}
	return
}

// Returns the hostname of, well, the host.
func getHostname() (hostname string, err error) {
	return os.Hostname()
}

// Wrapper function for inputting hostname and IP address in a somewhat
// elegant way.
func nodeDetails(nodeCount int) []Host {
	nodes := []Host{}
	for i := 0; i < nodeCount; i++ {
		if nodeCount > 1 {
			fmt.Println(i+1, "/", nodeCount)
		}
		fmt.Println("Enter your node's information.")
		var host Host
		host.Hostname = inputHostname()
		host.Ipaddr = inputIPAddr()
		nodes = append(nodes, host)
	}
	return nodes
}

// Asks the user for the hostname of a node
func inputHostname() (hostname string) {
	for {
		fmt.Print("\nhostname: ")
		fmt.Scanln(&hostname)
		if len(hostname) == 0 {
			fmt.Println("Error. Hostname cannot be empty.")
		} else {
			return
		}
	}
}

// Asks the user for the IP address of a node.
func inputIPAddr() (ipaddr string) {
	for {
		fmt.Print("\nIP Address: ")
		fmt.Scanln(&ipaddr)
		validIP := regexp.MustCompile(`^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)
		if !validIP.MatchString(ipaddr) {
			fmt.Println("Error. That does not look like a valid IP.")
		} else {
			return
		}
	}
}
