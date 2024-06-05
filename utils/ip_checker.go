package utils

import (
	"log"
	"net"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

func CheckIP(ip net.IP, countryWhitelist []string, db *geoip2.Reader) (bool, *geoip2.Country, error) {
	record, err := db.Country(ip)
	if err != nil {
		return false, nil, err
	}

	accepted := false
	for _, country := range countryWhitelist {
		if strings.EqualFold(country, record.Country.IsoCode) {
			accepted = true
			break
		}
	}

	return accepted, record, nil
}

func LogResult(source string, accepted bool, ip net.IP, record *geoip2.Country) {
	var result string
	if accepted {
		result = "accepted"
	} else {
		result = "denied"
	}

	log.Printf("%s request %s: ip: %s, country: %s - %s\n",
		source, result, ip, record.Country.IsoCode, record.Country.Names["en"])
}
