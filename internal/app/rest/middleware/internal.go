package middleware

import (
	"net"

	"github.com/gin-gonic/gin"
)

// isTrustedSubnet checks if ip is in subnet
func isTrustedSubnet(subnet string, ip string) bool {
	subnetIP, subnetMask, err := net.ParseCIDR(subnet)
	if err != nil {
		return false
	}

	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return false
	}

	return subnetIP.Equal(ipAddr.Mask(subnetMask.Mask))
}

// Internal checks if request is internal
func Internal(subnet string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if subnet == "" {
			c.AbortWithStatus(403)
			return
		}

		xRealIP := c.GetHeader("X-Real-IP")

		if xRealIP == "" {
			c.AbortWithStatus(403)
			return
		}

		if !isTrustedSubnet(subnet, xRealIP) {
			c.AbortWithStatus(403)
			return
		}

		c.Next()
	}
}
