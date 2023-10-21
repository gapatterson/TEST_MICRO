package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gapat/goMicro/service"
	"github.com/google/uuid"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func GetAllowedCountry() http.HandlerFunc {

	return func(rw http.ResponseWriter, r *http.Request) {
		var allowedRequest service.AllowedCountryRequest

		err := decoder.Decode(&allowedRequest, r.URL.Query())
		if err != nil { // Decoder should not fail.
			log.Printf("[ERROR] Returned bad request: %v", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		// TODO: upgrade to context passing
		// TODO: add context-dependent logging
		if allowedRequest.RequestId == "" {
			allowedRequest.RequestId = uuid.New().String()
		}

		var missingFields []string

		if allowedRequest.Ip == "" {
			missingFields = append(missingFields, "ip")
		}
		if allowedRequest.CustomerId == 0 {
			missingFields = append(missingFields, "customer_id")
		}
		if len(missingFields) > 0 {
			log.Printf("[INFO] %v Returned bad request: missing fields %v", allowedRequest.RequestId, missingFields)
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(fmt.Sprintf("Bad Request: Missing fields: %v", strings.Join(missingFields[:], ", "))))
		}

		allowed, err := service.CustomerAllowedIp(allowedRequest)
		if err == service.ErrorInvalidIp {
			log.Printf("[INFO] %v Returned bad request: invalid ip %v", allowedRequest.RequestId, allowedRequest.Ip)
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(fmt.Sprintf("Bad Request: Invalid IP: %v", allowedRequest.Ip)))
			return
		} else if err == service.ErrorInvalidCustomerId {
			// Do not inform the client about the customer id being invalid.  They should not be aware of internal data.
			log.Printf("[INFO] %v Returned forbidden", allowedRequest.RequestId)
			rw.WriteHeader(http.StatusForbidden)
			return
		} else if err != nil {
			log.Printf("[INFO] %v Returned internal server error", allowedRequest.RequestId)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		if allowed {
			log.Printf("[INFO] %v Returned OK", allowedRequest.RequestId)
			rw.WriteHeader(http.StatusOK)
			return
		} else {
			log.Printf("[INFO] %v Returned forbidden", allowedRequest.RequestId)
			rw.WriteHeader(http.StatusForbidden)
			return
		}

	}
}
