package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

// OkdDomainInfo is ... fuck.
type OkdDomainInfo struct {
	Domain    string
	Cluster   string
	Services  Host
	Bootstrap Host
	Masters   []Host
	Workers   []Host
}

// SubnetInfo is ... fuck.
type SubnetInfo struct {
	Domain      string
	Cluster     string
	Services    Host
	Services24  string
	Bootstrap   Host
	Bootstrap24 string
	Masters     []Host
	Masters24   []string
	Workers     []Host
	Workers24   []string
}

// LocalInfo is ...
type LocalInfo struct {
	Domain        string
	Subnet        string
	SubnetReverse string
}

// NamedInfo is ...
type NamedInfo struct {
	Services      Host
	Subnet        string
	SubnetReverse string
}

func generateOkdDomainBindConfig(domain string, cluster string, services Host, bootstrap Host, masters []Host, workers []Host) {
	input, err := ioutil.ReadFile("template/named/db.okd.local.template")
	if err != nil {
		fmt.Println(err)
	}

	info := OkdDomainInfo{
		domain,
		cluster,
		services,
		bootstrap,
		masters,
		workers,
	}

	haproxyTemplate := template.Must(template.New("").Parse(string(input)))

	os.MkdirAll("output", os.ModePerm)
	f, err := os.Create(fmt.Sprintf("output/db.%s.%s", cluster, domain))
	haproxyTemplate.Execute(f, info)
	// haproxyTemplate.Execute(os.Stdout, info) // DEBUG
}

func generateSubnetBindConfig(domain string, cluster string, services Host, bootstrap Host, masters []Host, workers []Host) {
	input, err := ioutil.ReadFile("template/named/db.subnet.template")
	if err != nil {
		fmt.Println(err)
	}

	masters24 := []string{}
	workers24 := []string{}

	for _, node := range masters {
		masters24 = append(masters24, strings.Split(node.Ipaddr, ".")[3])
	}

	for _, node := range workers {
		workers24 = append(workers24, strings.Split(node.Ipaddr, ".")[3])
	}

	info := SubnetInfo{
		domain,
		cluster,
		services,
		strings.Split(services.Ipaddr, ".")[3],
		bootstrap,
		strings.Split(bootstrap.Ipaddr, ".")[3],
		masters,
		masters24,
		workers,
		workers24,
	}

	haproxyTemplate := template.Must(template.New("").Parse(string(input)))

	os.MkdirAll("output", os.ModePerm)
	f, err := os.Create(fmt.Sprintf("output/db.%s.%s", cluster, domain))
	haproxyTemplate.Execute(f, info)
	// haproxyTemplate.Execute(os.Stdout, info) // DEBUG
}

func generateLocalBindConfig(domain string, services Host) {
	input, err := ioutil.ReadFile("template/named/named.conf.local.template")
	if err != nil {
		fmt.Println(err)
	}

	// Reverse subnet
	s := strings.Split(services.Ipaddr, ".")[0:3]
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	subnetReverse := strings.Join(s, ".")

	info := LocalInfo{
		domain,
		strings.Join(strings.Split(services.Ipaddr, ".")[0:3], "."),
		subnetReverse,
	}

	haproxyTemplate := template.Must(template.New("").Parse(string(input)))

	os.MkdirAll("output", os.ModePerm)
	f, err := os.Create("output/named.conf.local")
	haproxyTemplate.Execute(f, info)
	// haproxyTemplate.Execute(os.Stdout, info) // DEBUG
}

func generateBindConfig(services Host) {
	input, err := ioutil.ReadFile("template/named/named.conf.template")
	if err != nil {
		fmt.Println(err)
	}

	// Reverse subnet
	s := strings.Split(services.Ipaddr, ".")[0:3]
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	subnetReverse := strings.Join(s, ".")

	info := NamedInfo{
		services,
		strings.Join(strings.Split(services.Ipaddr, ".")[0:3], "."),
		subnetReverse,
	}

	haproxyTemplate := template.Must(template.New("").Parse(string(input)))

	os.MkdirAll("output", os.ModePerm)
	f, err := os.Create("output/named.conf")
	haproxyTemplate.Execute(f, info)
	// haproxyTemplate.Execute(os.Stdout, info) // DEBUG
}
