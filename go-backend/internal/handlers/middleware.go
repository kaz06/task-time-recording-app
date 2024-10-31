package handlers

import (
	"encoding/json"
	"errors"
	"fmt"

	"io"

	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

var firebaseCerts map[string]string

func fetchFirebaseCerts() error {
	resp, err := http.Get("https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &firebaseCerts)
	if err != nil {
		return err
	}

	return nil
}

func (c *Container) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		authHeader := ctx.Request().Header.Get("Authorization")
		if authHeader == "" {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Not Found Authorization Header"})
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer"))
		if tokenString == "" {
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Not Found Token"})
		}

		_, claims, err := verifyIDToken(tokenString)
		if err != nil {

			fmt.Printf("Error verifying ID token: %v\n", err)
			return ctx.JSON(http.StatusUnauthorized, map[string]string{"error": "Failed to verify token"})
		}

		user, err := c.getOrCreateUser(ctx, claims)

		if err != nil {
			fmt.Printf("Error getting or creating user: %v\n", err)
			return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get or create user"})
		}
		ctx.Set("user", user)

		return next(ctx)
	}
}

func verifyIDToken(idToken string) (*jwt.Token, jwt.MapClaims, error) {
	if firebaseCerts == nil {
		err := fetchFirebaseCerts()
		if err != nil {
			return nil, nil, err
		}
	}

	token, err := jwt.Parse(idToken, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid not found in toekn header")
		}

		cert, ok := firebaseCerts[kid]
		if !ok {
			err := fetchFirebaseCerts()
			if err != nil {
				return nil, err
			}
			cert, ok = firebaseCerts[kid]
			if !ok {
				return nil, fmt.Errorf("unable to find certificate for kid: %s", kid)
			}
		}

		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		if err != nil {
			return nil, err
		}
		return key, nil

	})

	if err != nil {
		return nil, nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if err := verifyClaims(claims); err != nil {
			return nil, nil, err
		}
		return token, claims, nil
	} else {
		return nil, nil, errors.New("Invalid token")
	}

}

func verifyClaims(claims jwt.MapClaims) error {

	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	if projectID == "" {
		return fmt.Errorf("FIREBASE_PROJECT_ID is not set")
	}

	iss, ok := claims["iss"].(string)
	expectedIss := fmt.Sprintf("https://securetoken.google.com/%s", projectID)
	if !ok || iss != expectedIss {
		return fmt.Errorf("invalid issuer")
	}

	aud, ok := claims["aud"].(string)
	if !ok || aud != projectID {
		return fmt.Errorf("invalid audience")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return fmt.Errorf("invalid expiration time")
	}
	if int64(exp) < time.Now().Unix() {
		return fmt.Errorf("token is expired")
	}
	return nil
}
