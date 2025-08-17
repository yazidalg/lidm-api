package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretOnce sync.Once
var secretReportDone bool
var lastSecretName string

// getSecret fetches the JWT secret each time (try multiple common env names) and logs once.
func getSecret() []byte {
	names := []string{"SECRET", "JWT_SECRET", "APP_SECRET", "JWT_SIGNING_KEY"}
	var val string
	var used string
	for _, n := range names {
		if v := os.Getenv(n); v != "" {
			val = v
			used = n
			break
		}
	}
	if used == "" { // all empty
		used = "(none)"
	}
	lastSecretName = used
	secretOnce.Do(func() {
		log.Printf("[jwt] secret source=%s length=%d (logged once)", used, len(val))
		if strings.HasPrefix(val, "\"") || strings.HasSuffix(val, "\"") {
			log.Printf("[jwt] NOTE: secret seems to include quotes; consider removing surrounding quotes")
		}
	})
	return []byte(val)
}

func GenerateJwt(userId uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	return token.SignedString(getSecret())
}

func ParseToken(tokenStr string) (uint, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// Ensure signing method HS256
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return getSecret(), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Coba ambil sub dulu, kalau tidak ada ambil id
		var raw interface{}
		if val, ok := claims["sub"]; ok {
			raw = val
		} else if val, ok := claims["id"]; ok {
			raw = val
		}
		switch v := raw.(type) {
		case float64:
			return uint(v), nil
		case int:
			return uint(v), nil
		case int64:
			return uint(v), nil
		case string:
			var id uint
			_, err := fmt.Sscanf(v, "%d", &id)
			if err == nil {
				return id, nil
			}
		default:
			if f, ok := raw.(float64); ok {
				return uint(f), nil
			}
		}
	}
	if err != nil {
		log.Printf("[jwt] parse error: %v (secret_source=%s)", err, lastSecretName)
		// Optional insecure fallback: allow extracting id from payload without verifying signature (DEV ONLY)
		if strings.Contains(err.Error(), "signature is invalid") && os.Getenv("ALLOW_INSECURE_JWT") == "1" {
			parts := strings.Split(tokenStr, ".")
			if len(parts) == 3 {
				if payloadBytes, decErr := base64.RawURLEncoding.DecodeString(parts[1]); decErr == nil {
					var payload map[string]any
					if jsonErr := json.Unmarshal(payloadBytes, &payload); jsonErr == nil {
						// try id / sub keys
						if v, ok := payload["sub"]; ok {
							payload["id"] = v
						}
						if raw, ok := payload["id"]; ok {
							switch vv := raw.(type) {
							case float64:
								return uint(vv), nil
							case int:
								return uint(vv), nil
							case int64:
								return uint(vv), nil
							case string:
								var id uint
								if _, scanErr := fmt.Sscanf(vv, "%d", &id); scanErr == nil {
									return id, nil
								}
							}
						}
					}
				}
			}
		}
	}
	return 0, err
}

// ExtractUnverifiedUserID decodes JWT payload tanpa verifikasi signature (HANYA untuk dev / fallback)
func ExtractUnverifiedUserID(tokenStr string) (uint, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid token format")
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, err
	}
	var payload map[string]any
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return 0, err
	}
	// prefer id then sub
	var raw any
	if v, ok := payload["id"]; ok {
		raw = v
	} else if v, ok := payload["sub"]; ok {
		raw = v
	}
	switch v := raw.(type) {
	case float64:
		return uint(v), nil
	case int:
		return uint(v), nil
	case int64:
		return uint(v), nil
	case string:
		var id uint
		if _, scanErr := fmt.Sscanf(v, "%d", &id); scanErr == nil {
			return id, nil
		}
	}
	return 0, fmt.Errorf("user id not found in payload")
}
