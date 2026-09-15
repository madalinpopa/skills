package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const minSigningKeyBytes = 32

type Config struct {
	SigningKey      []byte
	Issuer          string
	Audience        string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	CookieSecure    bool
	CookieSameSite  http.SameSite
}

func (c Config) Validate() error {
	var err error

	if len(c.SigningKey) < minSigningKeyBytes {
		err = errors.Join(err, fmt.Errorf("signing key must have at least %d bytes", minSigningKeyBytes))
	}

	if strings.TrimSpace(c.Issuer) == "" {
		err = errors.Join(err, errors.New("issuer must not be empty"))
	}

	if strings.TrimSpace(c.Audience) == "" {
		err = errors.Join(err, errors.New("audience must not be empty"))
	}

	if c.AccessTokenTTL <= 0 {
		err = errors.Join(err, errors.New("access token ttl must be positive"))
	}

	if c.RefreshTokenTTL <= 0 {
		err = errors.Join(err, errors.New("refresh token ttl must be positive"))
	}

	if c.CookieSameSite == http.SameSiteNoneMode && !c.CookieSecure {
		err = errors.Join(err, errors.New("cookie same site none requires a secure cookie"))
	}

	return err
}

// ConfigFromEnv reads the module settings. AUTH_SECRET_KEY or the file named
// by AUTH_SECRET_KEY_FILE is required; the rest have defaults.
func ConfigFromEnv() (Config, error) {
	secret, err := secretFromEnv("AUTH_SECRET_KEY")
	if err != nil {
		return Config{}, err
	}

	accessTTL, err := durationFromEnv("AUTH_ACCESS_TOKEN_TTL", 15*time.Minute)
	if err != nil {
		return Config{}, err
	}

	refreshTTL, err := durationFromEnv("AUTH_REFRESH_TOKEN_TTL", 12*time.Hour)
	if err != nil {
		return Config{}, err
	}

	cookieSecure, err := boolFromEnv("AUTH_COOKIE_SECURE", true)
	if err != nil {
		return Config{}, err
	}

	issuer := stringFromEnv("AUTH_ISSUER", "api")

	return Config{
		SigningKey:      []byte(secret),
		Issuer:          issuer,
		Audience:        stringFromEnv("AUTH_AUDIENCE", issuer),
		AccessTokenTTL:  accessTTL,
		RefreshTokenTTL: refreshTTL,
		CookieSecure:    cookieSecure,
		CookieSameSite:  http.SameSiteLaxMode,
	}, nil
}

func secretFromEnv(name string) (string, error) {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value, nil
	}

	path := strings.TrimSpace(os.Getenv(name + "_FILE"))
	if path == "" {
		return "", fmt.Errorf("%s or %s_FILE must be set", name, name)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("reading %s_FILE: %w", name, err)
	}

	value := strings.TrimSpace(string(data))
	if value == "" {
		return "", fmt.Errorf("%s_FILE %q is empty", name, path)
	}

	return value, nil
}

func stringFromEnv(name, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}

	return fallback
}

func durationFromEnv(name string, fallback time.Duration) (time.Duration, error) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return fallback, nil
	}

	value, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("%s: %q is not a duration such as 15m or 12h", name, raw)
	}

	return value, nil
}

func boolFromEnv(name string, fallback bool) (bool, error) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return fallback, nil
	}

	value, err := strconv.ParseBool(raw)
	if err != nil {
		return false, fmt.Errorf("%s: %q is not a boolean", name, raw)
	}

	return value, nil
}
