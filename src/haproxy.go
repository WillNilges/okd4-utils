package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

//HAProxyInfo is ...
type HAProxyInfo struct {
	Bootstrap Host
	Masters   []Host
	Workers   []Host
}

// Generate an haproxy config file.
func generateHAProxyConfig(bootstrap Host, masters []Host, workers []Host) {
	input, err := ioutil.ReadFile("template/haproxy.cfg.template")
	if err != nil {
		fmt.Println(err)
	}

	info := HAProxyInfo{
		bootstrap,
		masters,
		workers,
	}

	/*
		TODO:
		Append Bootstrap to masters list
		Append Masters to masters list
		Append Workers to workers list
		If Masters can be workers
		Append Masters to workers list
	*/

	haproxyTemplate := template.Must(template.New("").Parse(string(input)))
	// nodes := []Host{Host{"host1", "1.1.1.1"}, Host{"host2", "1.1.1.2"}, Host{"host3", "1.1.1.3"}}

	//nodeNames := []string{}
	//nodeIPs := []string{}
	//for _,node := range nodes {
	//	nodeNames = append(nodeNames, node.Hostname)
	//	nodeIPs = append(nodeIPs, node.ipaddr)
	//}
	//fmt.Println(nodeNames)

	os.MkdirAll("output", os.ModePerm)
	f, err := os.Create("output/haproxy.cfg")
	haproxyTemplate.Execute(f, info)

	// err = ioutil.WriteFile("output/haproxy.cfg", []byte(output), 0644)
	if err != nil {
		fmt.Println(err)
	}
}
