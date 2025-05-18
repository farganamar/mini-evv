package helpers

import "net/http"

func GetRealIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Origin-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarder-For")
	}

	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Real-IP")
	}

	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
