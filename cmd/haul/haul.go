package main

import (
	"flag"
	"fmt"
	"github.com/kseleven/haul/pkg/resource"
	"github.com/kseleven/haul/pkg/service"
	"strings"
)

var (
	host   string
	domain string
	qType  string
)

func main() {
	flag.StringVar(&host, "h", "127.0.0.1", "dns server ip")
	flag.StringVar(&domain, "r", ".", "domain name")
	flag.StringVar(&qType, "t", "A", "type:one of (a,any,mx,ns,soa,hinfo,axfr,txt,...)")
	flag.Parse()

	qTypeCode, ok := resource.QTypeSet[strings.ToUpper(qType)]
	if !ok {
		qTypeCode = resource.QTypeA
	}

	if err := service.Request(resource.Request{
		Host:  host,
		Port:  53,
		QName: domain,
		QType: qTypeCode,
	}); err != nil {
		fmt.Println(err)
	}
}
