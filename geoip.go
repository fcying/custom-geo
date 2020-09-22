package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"github.com/golang/protobuf/proto"
	"v2ray.com/core/app/router"
	"v2ray.com/core/infra/conf"
	"v2ray.com/core/common"
)

func getIPsList(fileName string) []*router.CIDR {
	d, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	ips := strings.Split(string(d), "\n")

	cidr := make([]*router.CIDR, 0, len(ips))

	for _, ip := range ips {
		if (ip != "") {
			c, err := conf.ParseIP(strings.TrimSpace(ip))
			common.Must(err)
			cidr = append(cidr, c)
		}
	}
	return cidr
}

func GenGeoIP() {
	dir := "ip"

	protoList := new(router.GeoIPList)

	rulefiles, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, rf := range rulefiles {
		filename := rf.Name()
		protoList.Entry = append(protoList.Entry, &router.GeoIP{
			CountryCode: strings.ToUpper(filename),
			Cidr:      getIPsList(dir + "/" + filename),
		})
	}

	protoBytes, err := proto.Marshal(protoList)
	if err != nil {
		fmt.Println("Error marshalling geoip list:", err)
		os.Exit(1)
	}

	if err := ioutil.WriteFile("ip.dat", protoBytes, 0644); err != nil {
		fmt.Println("Error writing geoip to file:", err)
		os.Exit(1)
	} else {
		fmt.Println("ip.dat has been generated successfully in the directory.")
	}
}
