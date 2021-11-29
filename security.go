package chidoc

import "errors"

// AuthType const to know type of authorization documentation
type AuthType string

const (
	// AuthBasic for Authentication basic
	AuthBasic AuthType = "basic"
	// AuthAPIKey for Authentication Basic
	AuthAPIKey = "apiKey"
	// AuthOAuth2 for Authentication Auth2
	AuthOAuth2 = "oauth2"
	// AuthBearer for Authentication Bearer
	AuthBearer = "http"
)

// InType consts authorization input type
type InType string

const (
	// InHeader for APIKey authentication in HTTP Header
	InHeader InType = "header"
	// InQuery for APIKey authentication in URL query
	InQuery = "query"
	// InHttp for Bearer authentication in HTTP
	InHttp = "http"
)

// Auth structs to define authorization documentation
type Auth struct {
	Name        string
	Description string
	Type        AuthType
	Scopes      map[string]string
	UrlAuth     string
	// only case SecurityType was APIKey
	In InType
	// cnly APIKey
	ParameterName string
}

// NewAuthAPIKey creates a security for APIKey
func NewAuthAPIKey(name, description, parameter string, inType InType) Auth {
	return Auth{
		Name:          name,
		Description:   description,
		ParameterName: parameter,
		Type:          AuthAPIKey,
		In:            inType,
	}
}

// NewAuthBasic creates a security for Basic Authentication
func NewAuthBasic(name, description string) Auth {
	return Auth{
		Name:        name,
		Description: description,
		Type:        AuthBasic,
		In:          InHeader,
	}
}

// NewAuthBasic creates a security for Basic Authentication
func NewAuthBearer(name, description string) Auth {
	return Auth{
		Name:        name,
		Description: description,
		Type:        AuthBearer,
	}
}

// NewOAuth
func NewAuthOAuth(name, url, description string, scopes map[string]string) Auth {
	return Auth{
		Name:        name,
		Description: description,
		Scopes:      scopes,
		Type:        AuthOAuth2,
		UrlAuth:     url,
	}
}

// it's working, soon late
//func NewAuthOAuth2()

// Decode security to opeanapi(YAML) parameters
func (a Auth) Decode(ptr map[string]interface{}) (err error) {
	if ptr == nil {
		return errors.New("ptr cannot be nil")
	}

	if _, exists := ptr[a.Name]; exists {
		return errors.New("security already exists")
	}

	auth := make(map[string]interface{})
	auth["type"] = a.Type
	auth["description"] = a.Description

	switch a.Type {
	case AuthBasic:
		break
	case AuthOAuth2:
		auth["name"] = a.Name
		auth["flows"] = map[string]interface{}{
			"clientCredentials": map[string]interface{}{
				"tokenUrl": a.UrlAuth,
				"scopes":   a.Scopes,
			},
		}
	case AuthAPIKey:
		auth["in"] = a.In
		auth["name"] = a.ParameterName
	case AuthBearer:
		auth["scheme"] = "bearer"
		auth["bearerFormat"] = "jwt"
	default:
		return errors.New("SecurityType invalid")
	}
	ptr[a.Name] = auth
	return err
}
