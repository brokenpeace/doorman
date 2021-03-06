// Package authn is in charge authenticating requests.

// JWT validators will be instantiated per issuer, and OpenID configuration
// and keys will be cached between requests.

package authn

import (
	"fmt"
	"net/http"
	"strings"
)

// JWTValidator is the interface in charge of extracting JWT claims from request.
type JWTValidator interface {
	ValidateRequest(*http.Request) (*Claims, error)
}

var jwtValidators map[string]JWTValidator

func init() {
	jwtValidators = map[string]JWTValidator{}
}

// NewJWTValidator instantiates or reuses an existing JWT validator for the specified issuer.
func NewJWTValidator(issuer string) (JWTValidator, error) {
	if !strings.HasPrefix(issuer, "https://") {
		return nil, fmt.Errorf("issuer %q not supported or has bad format", issuer)
	}

	// Reuse JWT validators instances among configs if they are for the same issuer.
	v, ok := jwtValidators[issuer]
	if !ok {
		v = newJWTGenericValidator(issuer)
		jwtValidators[issuer] = v
	}
	return v, nil
}
