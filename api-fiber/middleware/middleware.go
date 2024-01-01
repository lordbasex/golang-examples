package middleware

import (
	"api-fiber/config"
	"api-fiber/models"
	"crypto/x509"
	"encoding/pem"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

var PrivateKey string
var PublicKey string

// CreateToken generates a signed JWT token for the provided user information.
func CreateToken(user models.User) (string, error) {
	// Parse the private key
	privateKeyBytes, _ := pem.Decode([]byte(PrivateKey))
	parsedPrivateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	// Get the current time
	now := time.Now()

	// Set default expiration duration in case of error or absence of configuration
	var expireAt time.Duration = time.Duration(24)

	// Get token expiration hours from the environment variable
	expireHours, err := strconv.Atoi(config.Config("JWT_SECRET_KEY_EXPIRE_HOUR_COUNT"))
	if err == nil {
		// Convert hours to a time duration
		expireAt = time.Duration(expireHours)
	}

	// Prepare JWT claims
	claims := models.CustomClaims{
		UserInfo: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(time.Hour * expireAt).Unix(), // Token valid for the specified hours
			IssuedAt:  now.Unix(),                           // Issuance time
			NotBefore: now.Unix(),                           // Token not valid before issuance time
		},
	}

	// Create a signed token using RS512 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	tokenString, err := token.SignedString(parsedPrivateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken is a middleware function for validating JWT tokens in the "Authorization" header.
func ValidateToken(c *fiber.Ctx) error {
	// Retrieve the "Authorization" header
	authorizationHeader := c.Get("Authorization")
	if authorizationHeader == "" || len(authorizationHeader) <= 7 {
		// Return unauthorized status if the header is invalid or missing
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Authorization header"})
	}

	// Extract the token from the header (remove "Bearer " prefix)
	tokenString := c.Get("Authorization")[7:]

	// Parse the public key
	publicKeyBytes, _ := pem.Decode([]byte(PublicKey))
	parsedPublicKey, err := x509.ParsePKIXPublicKey(publicKeyBytes.Bytes)
	if err != nil {
		// Log fatal error if public key parsing fails
		log.Fatal(err)
	}

	// Validate the token
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return parsedPublicKey, nil
	})

	if err != nil {
		// Return unauthorized status if token validation fails
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Verify token claims
	claims, ok := token.Claims.(*models.CustomClaims)
	if !ok || !token.Valid {
		// Return unauthorized status if claims verification fails
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Store a value in the context
	c.Locals("value1", "Master")

	// Store another value in the context
	c.Locals("value2", "Hello Fiber")

	// Store another user in the context
	c.Locals("user", claims)

	// Continue with the next middleware or handler
	return c.Next()
}

// BlockIP is a middleware function that blocks specified IP addresses.
// It returns a Fiber handler that checks the client's IP address against a list of IPs to block.
// If the client's IP matches any blocked IP, it responds with a 403 Forbidden status.
// If debug mode is enabled, the client's IP address is logged.
func BlockIP(ipsToBlock []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the client's IP address
		clientIP := c.IP()

		// Log the client's IP address if debug mode is enabled
		if config.DebugConfiguration {
			log.Print(clientIP)
		}

		// Check if the client's IP is in the list of IPs to block
		for _, blockedIP := range ipsToBlock {
			if clientIP == blockedIP {
				// Respond with a 403 Forbidden status if the client's IP is blocked
				return c.Status(fiber.StatusForbidden).SendString("Non Authoritative Information")
			}
		}

		// Continue to the next middleware or handler if the client's IP is not blocked
		return c.Next()
	}
}
