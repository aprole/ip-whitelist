package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/aprole/ip-whitelist/models"
	"github.com/aprole/ip-whitelist/utils"
	"github.com/oschwald/geoip2-golang"
)

type Country struct {
	ISOCode string
	Name    string
}

func Handler(w http.ResponseWriter, r *http.Request, db *geoip2.Reader) {
	var req models.HTTPRequest
	defer r.Body.Close()
	if err := json.NewDecoder(io.LimitReader(r.Body, 2<<20)).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("error decoding: ", err.Error())
		return
	}

	ipAddress := net.ParseIP(req.IP)
	if ipAddress == nil {
		errMsg := fmt.Sprintf("invalid ip address: %s", req.IP)
		http.Error(w, errMsg, http.StatusBadRequest)
		log.Println(errMsg)
		return
	}

	if req.AllowedCountries == nil {
		errMsg := "bad request: 'allowedCountries' is required"
		http.Error(w, errMsg, http.StatusBadRequest)
		log.Println(errMsg)
		return
	}

	accepted, record, err := utils.CheckIP(ipAddress, req.AllowedCountries, db)
	if err != nil {
		errMsg := fmt.Sprintf("server error: %s", err.Error())
		http.Error(w, errMsg, http.StatusInternalServerError)
		log.Println(errMsg)
		return
	}

	json.NewEncoder(w).Encode(
		models.HTTPResponse{
			Accepted:    accepted,
			IP:          req.IP,
			Country:     record.Country.IsoCode,
			CountryName: record.Country.Names["en"]})

	utils.LogResult("HTTP", accepted, ipAddress, record)
}
