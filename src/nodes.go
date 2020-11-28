package main

import (
	"fmt"
	"net"
	"os"
	"io/ioutil"
	"strings"
	"regexp"
)

type Host struct {
	hostname string
	ipaddr string
}

// Automatically detect the network devices and IP addresses of the host.
// Interactively let the user choose which address to use.
func choose_IP() (IP string, err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
	}
	for index, iface := range ifaces {
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

// Returns the hostname of, well, the host.
func get_hostname() (hostname string, err error) {
	return os.Hostname()
}

// Wrapper function for acquiring hostname and IP address in a somewhat
// elegant way.
func node_details(node_count int) ([]Host){
	nodes := []Host{}
	for i := 0; i < node_count; i++ {
		if node_count > 1 {
			fmt.Println(i+1,"/",node_count)
		}
		fmt.Println("Enter your node's information.")
		var host Host
		host.hostname = input_hostname()
		host.ipaddr = input_ipaddr()
		nodes = append(nodes, host)
	}
	return nodes
}

// Asks the user for the hostname of a node
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

// Asks the user for the IP address of a node.
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

// Generate an haproxy config file.
func haproxy_gen(bootstrap Host, masters []Host, workers []Host) {
	input, err := ioutil.ReadFile("template/haproxy.conf")
	if err != nil {
			fmt.Println(err)
	}

	lines := strings.Split(string(input), "\n")

	offset := 0
	for i, line := range lines {
		var port string
		if strings.Contains(line, "# okd4_k8s_api_be nodes") {
			port = "6443"
		} else if strings.Contains(line, "# okd4_machine_config_server_be nodes") {
			port = "22623"
		} else if strings.Contains(line, "# okd4_http_ingress_traffic_be nodes") {
			port = "80"
		} else if strings.Contains(line, "# okd4_https_ingress_traffic_be nodes") {
			port = "443"
		}

		// Slot the masters and bootstrap into place
		if strings.Contains(line, "# okd4_k8s_api_be nodes") || strings.Contains(line, "# okd4_machine_config_server_be nodes"){
				// Add the bootstrap node
				element := "    server    " + bootstrap.hostname + "    " + bootstrap.ipaddr + ":" + port + " check"
				lines = append(lines, "") // Step 1, make room
				copy(lines[i+offset+1:], lines[i+offset:]) // Step 2, shove everything over
				lines[i+offset] = element // Step 3, insert new item
				offset++
				
				// Add the masters
				for _, node := range masters {
					element := "    server    " + node.hostname + "    " + node.ipaddr + ":" + port + " check"
					lines = append(lines, "")   // Step 1, make room
					copy(lines[i+offset+1:], lines[i+offset:]) // Step 2, shove everything over
					lines[i+offset] = element // Step 3, insert new item
					offset++
				}
		}
		
		// Slot the workers into place
		if strings.Contains(line, "# okd4_http_ingress_traffic_be nodes") || strings.Contains(line, "# okd4_https_ingress_traffic_be nodes"){
			// Add the worker node
			for _, node := range workers {
				element := "    server    " + node.hostname + "    " + node.ipaddr + ":" + port + " check"
				lines = append(lines, "") // Step 1, make room
				copy(lines[i+offset+1:], lines[i+offset:]) // Step 2, shove everything over
				lines[i+offset] = element // Step 3, insert new item
				offset++
			}
		}
	}
	output := strings.Join(lines, "\n")
	os.MkdirAll("output", os.ModePerm)
	err = ioutil.WriteFile("output/haproxy.conf", []byte(output), 0644)
	if err != nil {
		fmt.Println(err)
	}	
}

func bind_gen(domain string, cluster string, services Host, bootstrap Host, masters []Host, workers []Host) {
	lines := open_file("template/named/db.okd.local")

	// Fill out all the single-item stuff
	for i,line := range lines {
		mod_line_1 := strings.Replace(line, "{domain}", domain, -1)
		mod_line_2 := strings.Replace(mod_line_1, "{sub-domain}", cluster, -1)
		mod_line_3 := strings.Replace(mod_line_2, "{bootstrap-ip}", bootstrap.ipaddr, -1)
		mod_line_4 := strings.Replace(mod_line_3, "{services-ip}", services.ipaddr, -1)
		// fmt.Println(mod_line_4)
		lines[i] = mod_line_4
	}

	// Fill out the possible multitudes of masters
	for i,line := range lines {
		if strings.Contains(line, "okd4-master-{master-index}") {
			master_lines := []string{}
			for _, master := range masters {
				master_mod_1 := strings.Replace(line, "okd4-master-{master-index}", master.hostname, -1)
				master_mod_2 := strings.Replace(master_mod_1, "{master-ip}", master.ipaddr, -1)
				master_lines = append(master_lines, master_mod_2)
			}
			lines[i] = strings.Join(master_lines, "\n")
		}
		if strings.Contains(line, "etcd-{master-index}") {
			master_lines := []string{}
			for _, master := range masters {
				master_mod_1 := strings.Replace(line, "{master-index}", master.hostname, -1)
				master_mod_2 := strings.Replace(master_mod_1, "{master-ip}", master.ipaddr, -1)
				master_lines = append(master_lines, master_mod_2)
			}
			lines[i] = strings.Join(master_lines, "\n")
		}
	}

	// Fill out the possible multitudes of workers
	for i,line := range lines {
		if strings.Contains(line, "okd4-worker-{worker-index}") {
		worker_lines := []string{}
			for _, worker := range workers {
				worker_mod_1 := strings.Replace(line, "okd4-worker-{worker-index}", worker.hostname, -1)
				worker_mod_2 := strings.Replace(worker_mod_1, "{worker-ip}", worker.ipaddr, -1)
			worker_lines = append(worker_lines, worker_mod_2)
			}
			lines[i] = strings.Join(worker_lines, "\n")
		}
	}

	os.MkdirAll("output", os.ModePerm)
	write_file("output/db.okd.local", lines)
}

func open_file(path string) ([]string) {
	input, err := ioutil.ReadFile(path)
	if err != nil {
			fmt.Println(err)
	}

	return strings.Split(string(input), "\n")
}

func write_file(path string, lines []string) {
	output := strings.Join(lines, "\n")
	err := ioutil.WriteFile(path, []byte(output), 0644)
	if err != nil {
		fmt.Println(err)
	}
}