package sslstatus

import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"
)

func checkSSLExpiry(domain string) string {
	conn, err := tls.Dial("tcp", domain+":443", &tls.Config{})
	if err != nil {
		return fmt.Sprintf("error: %s - %v\n", strings.TrimPrefix(domain, "*."), err)
	}

	defer conn.Close()

	cleanDomain := strings.TrimPrefix(domain, "*.")

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) > 0 {
		cert := certs[0]
		daysLeft := int(time.Until(cert.NotAfter).Hours() / 24)

		status := "Info"
		if daysLeft <= 30 {
			status = "Warning"
		} else if daysLeft <= 15 {
			status = "Danger"
		} else if daysLeft <= 7 {
			status = "Critical"
		}

		return fmt.Sprintf("%s: %s will expire in %d days (%s)\n",
			status,
			cleanDomain,
			daysLeft,
			cert.NotAfter.Format(time.RFC1123))
	}

	return fmt.Sprintf("failed to perform a TLS handshake for the domain: %s \n", cleanDomain)
}

func checkSSLExpiryMulti(domains []string) string {
	var sb strings.Builder

	for _, domain := range domains {
		sb.WriteString(checkSSLExpiry(domain))
	}
	return sb.String()
}
