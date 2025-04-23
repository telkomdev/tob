package sslstatus

import (
	"crypto/tls"
	"fmt"
	"log"
	"strings"
	"time"
)

func checkSSLExpiry(domain string, logger *log.Logger) string {
	conn, err := tls.Dial("tcp", domain+":443", &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return fmt.Sprintf("error: %s - %v\n", strings.TrimPrefix(domain, "*."), err)
	}
	defer conn.Close()

	cleanDomain := strings.TrimPrefix(domain, "*.")
	certs := conn.ConnectionState().PeerCertificates

	if len(certs) > 0 {
		cert := certs[0]

		status := "Info"

		// checking whether the certificate matches the domain being checked
		if err := cert.VerifyHostname(cleanDomain); err != nil {
			status = "Danger"
			return fmt.Sprintf("%s: invalid SSL: cert for %s does not match (%s)\n",
				status,
				cleanDomain,
				cert.Subject.CommonName)
		}

		issuer := cert.Issuer.CommonName
		sub := cert.Issuer.CommonName
		logger.Println(cleanDomain, " | ", issuer, " | ", sub)

		expiredDate := cert.NotAfter.Format(time.RFC1123)

		if cert.NotAfter.Before(time.Now()) {
			status = "Danger"
			return fmt.Sprintf("%s: SSL %s expired on %s\n",
				status,
				cleanDomain,
				expiredDate)
		}

		daysLeft := int(time.Until(cert.NotAfter).Hours() / 24)

		if daysLeft <= 7 {
			status = "Critical"
		} else if daysLeft <= 15 {
			status = "Danger"
		} else if daysLeft <= 30 {
			status = "Warning"
		}

		return fmt.Sprintf("%s: SSL %s will expire in %d days (%s)\n",
			status,
			cleanDomain,
			daysLeft,
			expiredDate)
	}

	return fmt.Sprintf("failed to perform a TLS handshake for the domain: %s\n", cleanDomain)
}

func checkSSLExpiryMulti(domains []string, logger *log.Logger) string {
	var sb strings.Builder

	for _, domain := range domains {
		sb.WriteString(checkSSLExpiry(domain, logger))
	}
	return sb.String()
}
