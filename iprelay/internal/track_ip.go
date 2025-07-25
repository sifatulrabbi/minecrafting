package internal

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func GetClientIP(r *http.Request) string {
	headers := []string{
		"CF-Connecting-IP",    // Cloudflare
		"True-Client-IP",      // Cloudflare Enterprise
		"X-Real-IP",           // Nginx proxy
		"X-Forwarded-For",     // Standard proxy header
		"X-Client-IP",         // Apache mod_remoteip
		"X-Forwarded",         // Less common
		"X-Cluster-Client-IP", // Cluster environments
		"Forwarded-For",       // RFC 7239 (less common)
		"Forwarded",           // RFC 7239
	}
	for _, header := range headers {
		if ip := getIPFromHeader(r, header); ip != "" {
			fmt.Println("Got the ip from the header.")
			return ip
		}
	}
	return getIPFromRemoteAddr(r.RemoteAddr)
}

func getIPFromHeader(r *http.Request, headerName string) string {
	headerValue := strings.TrimSpace(r.Header.Get(headerName))
	if headerValue == "" {
		return ""
	}
	if headerName == "Forwarded" {
		fmt.Println("Got the ip from the Forwarded header.")
		return parseForwardedHeader(headerValue)
	}
	if strings.Contains(headerValue, ",") {
		ips := strings.Split(headerValue, ",")
		for _, ip := range ips {
			if cleanIP := validateAndCleanIP(strings.TrimSpace(ip)); cleanIP != "" {
				return cleanIP
			}
		}
		return ""
	}
	return validateAndCleanIP(headerValue)
}

// parseForwardedHeader handles RFC 7239 Forwarded header format
// Example: "for=192.0.2.60;proto=http;by=203.0.113.43"
func parseForwardedHeader(header string) string {
	// Split by semicolon and look for "for=" parameter
	parts := strings.Split(header, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(strings.ToLower(part), "for=") {
			forValue := part[4:] // Remove "for="
			// Remove quotes if present
			forValue = strings.Trim(forValue, "\"")
			return validateAndCleanIP(forValue)
		}
	}
	return ""
}

// validateAndCleanIP validates and cleans an IP address
func validateAndCleanIP(ip string) string {
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return ""
	}
	if host, _, err := net.SplitHostPort(ip); err == nil {
		ip = host
	}
	ip = strings.Trim(ip, "[]")
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return ""
	}
	return parsedIP.String()
}

func getIPFromRemoteAddr(remoteAddr string) string {
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		// If SplitHostPort fails, it might be just an IP without port
		if parsedIP := net.ParseIP(remoteAddr); parsedIP != nil {
			return parsedIP.String()
		}
		return ""
	}
	if parsedIP := net.ParseIP(ip); parsedIP != nil {
		return parsedIP.String()
	}
	return ""
}
