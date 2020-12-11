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

	/*
		TODO:
		Append Bootstrap to masters list
		Append Masters to masters list
		Append Workers to workers list
		If Masters can be workers
		Append Masters to workers list
	*/

	haproxyTemplate := template.Must(template.New("").Parse(string(input)))

	os.MkdirAll("output", os.ModePerm)
	// f, err := os.Create("output/haproxy.cfg")
	fmt.Println("I am outputting!!!!!")
	haproxyTemplate.Execute(os.Stdout, info)

	// err = ioutil.WriteFile("output/haproxy.cfg", []byte(output), 0644)
	if err != nil {
		fmt.Println(err)
	}
}
