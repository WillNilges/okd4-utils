package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"strings"
)

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