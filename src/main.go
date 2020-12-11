package main

func main() {
	bootstrap := Host{"boot", "1.0.1.1"}
	masters := []Host{Host{"host1", "1.1.1.1"}, Host{"host2", "1.1.1.2"}, Host{"host3", "1.1.1.3"}}
	workers := []Host{Host{"work1", "1.2.1.1"}, Host{"work2", "1.2.1.2"}, Host{"work3", "1.2.1.3"}}
	generateHAProxyConfig(bootstrap, masters, workers, true)
}

// func main() {
// 	fmt.Println("Welcome to QuickShift") // Working title.
// 	fmt.Println("This program is designed to make it easier to configure an OpenShift cluster by automatically configuring your services machine, DNS config files, and HAProxy config file.")

// 	// Get domain name. Important for generating okd4 config yaml
// 	var domain string
// 	fmt.Print("\nEnter your domain name: ")
// 	fmt.Scanln(&domain)

// 	var cluster string
// 	fmt.Print("\n\nEnter your cluster name: ")
// 	fmt.Scanln(&cluster)

// 	// Acquire network details of the service machine. Auto-detects devices and IP addresses.
// 	fmt.Println("\n\nAcquiring service host networking information...")

// 	// Get service machine IP
// 	serviceIP, err := chooseIP()
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	// Get service machine hostnames
// 	serviceHostname, err := getHostname()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	service := Host{serviceHostname, serviceIP}
// 	fmt.Println("Your service host is:\n", service.hostname, "\n", service.ipaddr, "\n")

// 	// Get bootstrap information. Manual input.
// 	var bootstrap Host
// 	fmt.Print("Enter your bootstrap node's information: ")
// 	bootstrap.hostname = inputHostname()
// 	bootstrap.ipaddr = inputIPAddr()

// 	// Acquire master information. Quantity and details. Manual input.
// 	fmt.Println("\n\n=== Masters ===")
// 	var masterCount int
// 	for {
// 		fmt.Print("Enter number of masters: ")
// 		_, err := fmt.Scanf("%d", &masterCount)
// 		if err == nil || masterCount < 1 {
// 			break
// 		} else {
// 			fmt.Println("Invalid input detected.")
// 		}
// 	}
// 	masters := nodeDetails(masterCount)
// 	fmt.Println(masters)

// 	fmt.Println("\n\n=== Workers ===")
// 	var workerCount int
// 	for {
// 		fmt.Print("Enter number of workers: ")
// 		_, err := fmt.Scanf("%d", &workerCount)
// 		if err == nil {
// 			break
// 		} else {
// 			fmt.Println("Invalid input detected.")
// 		}
// 	}
// 	workers := nodeDetails(workerCount)
// 	fmt.Println(workers)

// 	// haproxy_gen(bootstrap, masters, workers)
// 	// bind_gen_subdomain(domain, cluster, service, bootstrap, masters, workers)
// 	// bind_gen_subnet(domain, cluster, service, bootstrap, masters, workers)
// 	// bind_named_conf_gen(service)
// 	// bind_named_conf_local_gen(domain, service)
// 	generateHAProxyConfig(bootstrap, masters, workers)
// }
