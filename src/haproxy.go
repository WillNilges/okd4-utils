package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"text/template"
)

//HAProxyInfo is ...
type HAProxyInfo struct {
	Bootstrap        Host
	Masters          []Host
	Workers          []Host
	MastersAsWorkers bool
}

// Generate an haproxy config file.
func generateHAProxyConfig(bootstrap Host, masters []Host, workers []Host, mastersAsWorkers bool) {
	input, err := ioutil.ReadFile("template/haproxy.cfg.template")
	if err != nil {
		fmt.Println(err)
	}

	info := HAProxyInfo{
		bootstrap,
		masters,
		workers,
		mastersAsWorkers,
	}

	haproxyTemplate := template.Must(template.New("").Parse(string(input)))

	os.MkdirAll("output", os.ModePerm)
	f, err := os.Create("output/haproxy.cfg")
	haproxyTemplate.Execute(f, info)
}
