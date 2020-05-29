package chidoc

import (
	"errors"
	"fmt"
	"strings"
)

type DocSettings struct {
	Title       string
	Description string
	Version     string
	definitions []interface{}
	valuesPath  map[string]interface{}
	auths       []Auth
}

func NewDocSettings(title string) *DocSettings {
	return &DocSettings{
		Title:       title,
		definitions: make([]interface{}, 0),
		valuesPath:  make(map[string]interface{}),
		auths:       make([]Auth, 0),
	}
}

func decodeSetPath(ptr map[string]interface{}, rawPath string, value interface{}) (err error) {
	// in a near future, make a function to decode rawPath, instead of split
	// cos, rawPath is dirt, it may YAML undefined character
	var paths []string = strings.Split(rawPath, ".")

	if len(paths) == 0 {
		return errors.New("Settings path cannot be empty")
	}

	// removes the last elem, cos it's the field
	// dict.dict.dict.[field]
	var field string = paths[len(paths)-1]
	paths = paths[:len(paths)-1]
	for len(paths) > 0 {
		var path string = paths[0]
		fmt.Println(path)
		// pop array
		paths = paths[1:]

		//If exists path, go to it
		if nptr, exists := ptr[path]; exists {
			ptr = nptr.(map[string]interface{})
			continue
		}

		// creates a new dict
		nptr := make(map[string]interface{})
		ptr[path] = nptr
		ptr = nptr
	}

	//Set field value
	ptr[field] = value
	return err
}

// SetDefinitions add model to openAPI YAML
func (s *DocSettings) SetDefinitions(def ...interface{}) {
	s.definitions = def
}

// SetAuths set all authorization for openapi
func (s *DocSettings) SetAuths(auths ...Auth) {
	s.auths = auths
}

// Set a value to openAPI YAML
func (s *DocSettings) Set(name string, value interface{}) {
	s.valuesPath[name] = value
}

// Get returns a path value
func (s *DocSettings) Get(name string) interface{} {
	return s.valuesPath[name]
}

// Decode parse setting for openapi scruct
func (s *DocSettings) Decode(ptr map[string]interface{}) (err error) {
	s.Set("info.title", s.Title)
	s.Set("info.description", s.Description)
	s.Set("info.version", s.Version)
	s.Set("openapi", "3.0.0")

	for path, value := range s.valuesPath {
		if err = decodeSetPath(ptr, path, value); err != nil {
			break
		}
	}
	return err
}
