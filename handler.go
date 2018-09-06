package urlshort

import (
	"encoding/json"
	"net/http"

	yaml "gopkg.in/yaml.v2"
)

// Redirection represents a path redirection /foo to a defined URL.
type Redirection struct {
	Path string `json:"path", yaml:"path"`
	URL  string `json:"url", yaml:"url"`
}

type Redirections []Redirection

// Map converts a redirection slice into a MapHandler parameters.
func (rs Redirections) Map() map[string]string {
	result := make(map[string]string, len(rs))
	for _, r := range rs {
		result[r.Path] = r.URL
	}
	return result
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if url, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, url, http.StatusPermanentRedirect)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var redirections Redirections
	if err := yaml.Unmarshal(yml, &redirections); err != nil {
		return nil, err
	}
	return MapHandler(redirections.Map(), fallback), nil
}

// JSONHandler will parse provided json and returns a redirection handler.
func JSONHandler(raw []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var redirections Redirections
	if err := json.Unmarshal(raw, &redirections); err != nil {
		return nil, err
	}
	return MapHandler(redirections.Map(), fallback), nil
}
