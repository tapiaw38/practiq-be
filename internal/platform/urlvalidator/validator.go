package urlvalidator

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

var (
	// Blocked IP ranges for SSRF protection
	blockedNetworks = []*net.IPNet{
		// Private IPv4 ranges (RFC 1918)
		parseCIDR("10.0.0.0/8"),
		parseCIDR("172.16.0.0/12"),
		parseCIDR("192.168.0.0/16"),
		// Localhost
		parseCIDR("127.0.0.0/8"),
		// Link-local addresses
		parseCIDR("169.254.0.0/16"),
		// Loopback IPv6
		parseCIDR("::1/128"),
		// Link-local IPv6
		parseCIDR("fe80::/10"),
		// Unique local IPv6
		parseCIDR("fc00::/7"),
	}
)

func parseCIDR(cidr string) *net.IPNet {
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(fmt.Sprintf("invalid CIDR in urlvalidator: %s", cidr))
	}
	return network
}

// ValidateURL validates a URL to prevent SSRF attacks.
// It checks:
// - URL is well-formed
// - Scheme is http or https
// - Hostname is not in allowlist (if provided)
// - Resolved IP is not in blocked private/internal ranges
func ValidateURL(rawURL string, allowedDomains []string) error {
	if rawURL == "" {
		return fmt.Errorf("URL cannot be empty")
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Only allow http and https schemes
	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("unsupported URL scheme: %s (only http and https are allowed)", parsedURL.Scheme)
	}

	hostname := parsedURL.Hostname()
	if hostname == "" {
		return fmt.Errorf("URL must have a hostname")
	}

	// If allowlist is provided, enforce it
	if len(allowedDomains) > 0 {
		allowed := false
		for _, domain := range allowedDomains {
			if matchesDomain(hostname, domain) {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("domain %s is not in the allowlist", hostname)
		}
	}

	// Resolve hostname to IP addresses
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return fmt.Errorf("failed to resolve hostname %s: %w", hostname, err)
	}

	if len(ips) == 0 {
		return fmt.Errorf("hostname %s did not resolve to any IP addresses", hostname)
	}

	// Check all resolved IPs against blocked networks
	for _, ip := range ips {
		if err := validateIP(ip); err != nil {
			return fmt.Errorf("resolved IP %s for hostname %s: %w", ip.String(), hostname, err)
		}
	}

	return nil
}

func validateIP(ip net.IP) error {
	// Check if IP is in any blocked network
	for _, network := range blockedNetworks {
		if network.Contains(ip) {
			return fmt.Errorf("IP address is in blocked range %s (private/internal network)", network.String())
		}
	}
	return nil
}

// matchesDomain checks if hostname matches the allowed domain.
// Supports exact match and wildcard subdomain match (e.g., "*.example.com")
func matchesDomain(hostname, allowedDomain string) bool {
	hostname = strings.ToLower(hostname)
	allowedDomain = strings.ToLower(allowedDomain)

	// Exact match
	if hostname == allowedDomain {
		return true
	}

	// Wildcard subdomain match (e.g., "*.example.com")
	if strings.HasPrefix(allowedDomain, "*.") {
		baseDomain := allowedDomain[2:] // Remove "*."
		// Check if hostname ends with .baseDomain
		if strings.HasSuffix(hostname, "."+baseDomain) {
			return true
		}
		// Also allow exact match with base domain
		if hostname == baseDomain {
			return true
		}
	}

	return false
}
