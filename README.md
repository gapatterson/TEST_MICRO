# TEST_MICRO
An example microservice written in ~4 hours.

## Building
**IMPORTANT**: geolocation data for ip addresses comes from: https://dev.maxmind.com/geoip/geoip2/geolite2/
It requires an account to download the db, but it doesn't cost anything.  The GeoLite2 Country `.mmdb` file is used by this project, and must be inserted into the `/db` directory.  I would include it, but that would go against the License.

There is a dockerfile can be used to build and run the microservice.  At the moment, it is hardcoded to port 8080.

```
docker build -t microservice .
docker run -p 8080:8080 microservice
```

## Interaction
Input: 
	`GET` request on `/AllowedCountry?ip=<ip_address>&customer_id=<int>`
Output: 
* `200` header code if ip is in whitelisted country for customer
* `400` header code if missing or malformed query params
* `403` if ip not in whitelist for customer
* `500` if something has gone wrong

The current customers and their allowed regions:
```
1: [US]
2: [US, FR]
3: [ES]
```

Some quick sample ip addresses:
* US:
  * 1.0.0.0
  * 1.0.0.1
* FR:
  * 1.179.112.0
  * 1.179.112.1
* ES:
  * 1.178.224.0
  * 1.178.224.1
---

## Notes

Happy with:
* Basic Functionality (MVP)
* Dockerfile

Unhappy with (in order of personal priority):
* Logging is garbage (doesn't write to anything, not in searchable format), paired with:
* Not using context (request_id should be stored in the context and passed around.  The logger might be as well.)
* Tests missing
* Kubernetes YAML
* Mapping Data plan
* gRCP

Concessions:
Embedding the databases used by the service into the service itself is poor form.  There are env variables set in the dockerfile that point to the db locations, we would want to have them point towards actual db addresses.

This was the first time I've used a Windows machine to use docker, so that took some setup.

Things learned:
sqlite will create imaginary tables if the path is wrong.

there is an sqlite package that does not rely on CGO: modernc.org/sqlite

I definitely prefer to work on a non-windows machine.
