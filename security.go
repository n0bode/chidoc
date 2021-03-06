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
)

// InType consts authorization input type
type InType string

const (
	// InHeader for APIKey authentication in HTTP Header
	InHeader InType = "header"
	// InQuery for APIKey authentication in URL query
	InQuery = "query"
)

// Auth structs to define authorization documentation
type Auth struct {
	Name        string
	Description string
	Type        AuthType
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

// it's working, soon late
//func NewAuthOAuth2()

// Decode security to opeanapi(YAML) parameters
func (a Auth) Decode(ptr map[string]interface{}) (err error) {
	if ptr == nil {
		return errors.New("Ptr cannot be nil")
	}

	if _, exists := ptr[a.Name]; exists {
		return errors.New("Security already exists")
	}

	auth := make(map[string]interface{})
	auth["type"] = a.Type
	auth["description"] = a.Description

	switch a.Type {
	case AuthBasic:
		break
	case AuthOAuth2:
		// it's comming soon
		break
	case AuthAPIKey:
		auth["in"] = a.In
		auth["name"] = a.ParameterName
		break
	default:
		return errors.New("SecurityType invalid")
	}
	ptr[a.Name] = auth
	return err
}
