package service

import (
	"database/sql"
	"errors"
	"log"
	"net"
	"os"

	"github.com/oschwald/geoip2-golang"
	_ "modernc.org/sqlite"
)

var geoipDB = os.Getenv("GEOIP_DB_PATH")
var customersDB = os.Getenv("CUSTOMERS_DB_PATH")

type AllowedCountryRequest struct {
	Ip         string `schema:"ip"`
	CustomerId int    `schema:"customer_id"`
	RequestId  string `schema:"request_id"`
}

var ErrorInvalidIp = errors.New("invalid ip")
var ErrorInvalidCustomerId = errors.New("invalid customer id")

func CustomerAllowedIp(allowedRequest AllowedCountryRequest) (bool, error) {
	log.Printf("[INFO] %v Checking if customer %v is allowed to access from %v", allowedRequest.RequestId, allowedRequest.CustomerId, allowedRequest.Ip)

	country, err := GetCountry(allowedRequest)
	if err != nil {
		return false, err
	}

	log.Printf("[INFO] %v Customer %v checking ip from %v", allowedRequest.RequestId, allowedRequest.CustomerId, country)

	allowedCountries, err := GetAllowed(allowedRequest)
	if err != nil {
		return false, err
	}

	for _, allowed := range allowedCountries {
		if allowed == country {
			return true, nil
		}
	}

	return false, nil
}

func GetCountry(request AllowedCountryRequest) (string, error) {
	db, err := geoip2.Open(geoipDB)
	if err != nil {
		log.Printf("[ERROR] %v error opening country database %v", request.RequestId, err)
		return "", err
	}
	defer db.Close()

	ip := net.ParseIP(request.Ip)
	if ip == nil {
		log.Printf("[ERROR] %v error parsing ip %v", request.RequestId, request.Ip)
		return "", ErrorInvalidIp
	}

	record, err := db.Country(net.ParseIP(request.Ip))
	if err != nil {
		log.Printf("[ERROR] %v error looking up ip %v in country db: %v", request.RequestId, ip, err)
		return "", err
	}

	return record.Country.IsoCode, nil
}

func GetAllowed(request AllowedCountryRequest) ([]string, error) {
	db, err := sql.Open("sqlite", customersDB)
	if err != nil {
		log.Printf("[ERROR] %v error opening customers database %v", request.RequestId, err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT country_code FROM country_whitelist WHERE customer_id = ?", request.CustomerId)
	if err != nil {
		log.Printf("[ERROR] %v error querying customers database %v", request.RequestId, err)
		return nil, err
	}

	var countries []string
	for rows.Next() {
		var country string
		err = rows.Scan(&country)
		if err != nil {
			return nil, err
		}
		countries = append(countries, country)
	}

	if len(countries) == 0 {
		log.Printf("[ERROR] %v customer %v not found", request.RequestId, request.CustomerId)
		return nil, ErrorInvalidCustomerId
	}

	return countries, nil
}
