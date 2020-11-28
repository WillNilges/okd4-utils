package main

import (
	"fmt"
	"regexp"
)

type Host struct {
	hostname string
	ipaddr string
}

func input_hostname() (hostname string) {
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

func input_ipaddr() (ipaddr string) {
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