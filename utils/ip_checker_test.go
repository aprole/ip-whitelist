package utils

import (
	"net"
	"testing"

	"github.com/oschwald/geoip2-golang"
	"github.com/stretchr/testify/require"
)

func TestCheckIP(t *testing.T) {
	testCases := []struct {
		ip              string
		whitelist       []string
		exp_accepted    bool
		exp_country_iso string
		exp_coutry_name string
	}{
		{"75.127.6.164", []string{"US", "RS", "IT", "JP"}, true, "US", "United States"},
		{"101.33.28.0", []string{"AR", "ES", "FR", "PT"}, false, "SG", "Singapore"},
		{"103.167.154.0", []string{"NL", "RS", "IE"}, true, "IE", "Ireland"},
		{"100.42.20.34", []string{}, false, "CA", "Canada"},
	}

	db, err := geoip2.Open("/var/lib/GeoIP/GeoLite2-Country.mmdb")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	for _, tc := range testCases {
		t.Run(tc.ip, func(t *testing.T) {
			allowed, record, err := CheckIP(net.ParseIP(tc.ip), tc.whitelist, db)
			require.NoError(t, err)
			require.Equal(t, tc.exp_accepted, allowed)
			require.Equal(t, tc.exp_country_iso, record.Country.IsoCode)
			require.Equal(t, tc.exp_coutry_name, record.Country.Names["en"])
		})
	}
}
