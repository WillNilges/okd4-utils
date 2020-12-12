package main

import "fmt"

// func main() {
// 	services := Host{"boot", "1.0.69.101"}
// 	bootstrap := Host{"boot", "1.0.1.102"}
// 	masters := []Host{Host{"host1", "1.1.1.111"}, Host{"host2", "1.1.1.112"}, Host{"host3", "1.1.1.113"}}
// 	workers := []Host{Host{"work1", "1.2.1.111"}, Host{"work2", "1.2.1.112"}, Host{"work3", "1.2.1.113"}}

// 	generateOkdDomainBindConfig("postave.us", "willard_cluster", services, bootstrap, masters, workers)
// 	generateSubnetBindConfig("postave.us", "willard_cluster", services, bootstrap, masters, workers)
// 	generateLocalBindConfig("postave.us", services)
// 	generateBindConfig(services)
// }

func main() {
	fmt.Println("Welcome to QuickShift") // Working title.
	fmt.Println("This program is designed to make it easier to configure an OpenShift cluster by automatically configuring your services machine, DNS config files, and HAProxy config file.")

	// Get domain name. Important for generating okd4 config yaml
	var domain string
	fmt.Print("\nEnter your domain name: ")
	fmt.Scanln(&domain)

	var cluster string
	fmt.Print("\n\nEnter your cluster name: ")
	fmt.Scanln(&cluster)

	// Acquire network details of the service machine. Auto-detects devices and IP addresses.
	fmt.Println("\n\nAcquiring service host networking information...")

	// Get service machine IP
	serviceIP, err := chooseIP()
	if err != nil {
		fmt.Println(err)
	}

	// Get service machine hostnames
	serviceHostname, err := getHostname()
	if err != nil {
		fmt.Println(err)
	}
	service := Host{serviceHostname, serviceIP}
	fmt.Println("Your service host is:\n", service.Hostname, "\n", service.Ipaddr, "\n")

	// Get bootstrap information. Manual input.
	var bootstrap Host
	fmt.Print("Enter your bootstrap node's information: ")
	bootstrap.Hostname = inputHostname()
	bootstrap.Ipaddr = inputIPAddr()

	// Acquire master information. Quantity and details. Manual input.
	fmt.Println("\n\n=== Masters ===")

	var mastersAsWorkersS string
	var mastersAsWorkers bool
	for {
		fmt.Print("Would you like to use masters as workers (y/n)? ")
		_, err := fmt.Scanf("%s", &mastersAsWorkersS)
		if err == nil && mastersAsWorkersS == "y" {
			mastersAsWorkers = true
			break
		} else if err == nil && mastersAsWorkersS == "y" {
			mastersAsWorkers = false
			break
		} else {
			fmt.Println("Invalid input detected.")
		}
	}

	var masterCount int
	for {
		fmt.Print("Enter number of masters: ")
		_, err := fmt.Scanf("%d", &masterCount)
		if err == nil || masterCount < 1 {
			break
		} else {
			fmt.Println("Invalid input detected.")
		}
	}
	masters := nodeDetails(masterCount)
	fmt.Println(masters)

	fmt.Println("\n\n=== Workers ===")
	var workerCount int
	for {
		fmt.Print("Enter number of workers: ")
		_, err := fmt.Scanf("%d", &workerCount)
		if err == nil {
			break
		} else {
			fmt.Println("Invalid input detected.")
		}
	}
	workers := nodeDetails(workerCount)
	fmt.Println(workers)

	generateHAProxyConfig(bootstrap, masters, workers, mastersAsWorkers)
	generateOkdDomainBindConfig(domain, cluster, service, bootstrap, masters, workers)
	generateSubnetBindConfig(domain, cluster, service, bootstrap, masters, workers)
	generateLocalBindConfig(domain, service)
	generateBindConfig(service)
}
