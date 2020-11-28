package main

import (
	// "errors"
	"fmt"
	"net"
	"os"
	"strings"
)

func choose_IP() (IP string, err error) {
	// var IPs []string
	ifaces, err := net.Interfaces()
	// handle err
	if err != nil {
		fmt.Println(err)
	}
	for index, iface := range ifaces {
		// addrs, err := iface.Addrs()
		// handle err
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(index,": ",iface.Name)
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
	addrs, err := ifaces[option].Addrs()
	if err != nil {
		fmt.Println(err)
	}
	for index, addr := range addrs {
		fmt.Println(index,": ",addr)
	}
	var addr_option int
	for {
		fmt.Print("\nSelect your IP address: ")
		_, err := fmt.Scanf("%d", &addr_option)
		if err == nil {
			break
		} else {
			fmt.Println("Invalid input detected.")
		}
	}
	IP = addrs[addr_option].String()
	if strings.Contains(IP, "/") {
		IP = strings.Split(IP, "/")[0]
	}
	return
}

func get_hostname() (hostname string, err error) {
	return os.Hostname()
}