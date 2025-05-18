package helpers

import (
	"net/http"
	"strings"
)

func GetDeviceInfo(r *http.Request) string {
	userAgent := r.UserAgent()

	if userAgent == "" {
		return "Unknown Device"
	}

	// Check for common mobile devices
	if strings.Contains(userAgent, "iPhone") {
		// Extract iPhone model if possible
		parts := strings.Split(userAgent, "iPhone")
		if len(parts) > 1 {
			modelPart := parts[1]
			// Extract the version number after "iPhone"
			if idx := strings.Index(modelPart, ";"); idx > 0 {
				return "iPhone" + modelPart[:idx]
			}
		}
		return "iPhone"
	}

	if strings.Contains(userAgent, "iPad") {
		return "iPad"
	}

	if strings.Contains(userAgent, "Android") {
		// Try to extract Android device model
		if idx := strings.Index(userAgent, "Android"); idx > 0 {
			remain := userAgent[idx:]
			if modelIdx := strings.Index(remain, ";"); modelIdx > 0 {
				remain = remain[modelIdx+1:]
				if endIdx := strings.Index(remain, ";"); endIdx > 0 {
					model := strings.TrimSpace(remain[:endIdx])
					return model
				}
			}
		}
		return "Android Device"
	}

	// Check for desktop browsers
	if strings.Contains(userAgent, "Windows") {
		return "Windows"
	}

	if strings.Contains(userAgent, "Macintosh") || strings.Contains(userAgent, "Mac OS") {
		return "Mac"
	}

	if strings.Contains(userAgent, "Linux") && !strings.Contains(userAgent, "Android") {
		return "Linux"
	}

	// Default to returning a shortened user agent if it's too long
	if len(userAgent) > 50 {
		return userAgent[:47] + "..."
	}

	return userAgent
}
