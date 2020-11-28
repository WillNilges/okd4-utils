package main

import "fmt"

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
	service_ip, err := choose_IP()
	if err != nil {
		fmt.Println(err)
	}

	// Get service machine hostnames 
	service_hostname, err := get_hostname()
	if err != nil {
		fmt.Println(err)
	}
	service := Host{ service_hostname, service_ip }
	fmt.Println("Your service host is:\n",service.hostname,"\n",service.ipaddr,"\n")

	// Get bootstrap information. Manual input.
	var bootstrap Host
	fmt.Print("Enter your bootstrap node's information: ")
	bootstrap.hostname = input_hostname()
	bootstrap.ipaddr = input_ipaddr()

	// Acquire master information. Quantity and details. Manual input.
	fmt.Println("\n\n=== Masters ===")
	var master_count int
	for {
		fmt.Print("Enter number of masters: ")
		_, err := fmt.Scanf("%d", &master_count)
		if err == nil || master_count < 1 {
			break
		} else {
			fmt.Println("Invalid input detected.")
		}
	}
	masters := node_details(master_count)
	fmt.Println(masters)

	fmt.Println("\n\n=== Workers ===")
	var worker_count int
	for {
		fmt.Print("Enter number of workers: ")
		_, err := fmt.Scanf("%d", &worker_count)
		if err == nil {
			break
		} else {
			fmt.Println("Invalid input detected.")
		}
	}
	workers := node_details(worker_count)
	fmt.Println(workers)

	haproxy_gen(bootstrap, masters, workers)
	bind_gen(domain, cluster, service, bootstrap, masters, workers)
}