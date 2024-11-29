package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"
)

var TWEET_URL_REGEX = regexp.MustCompile(`^https:\/\/(?:twitter\.com|x\.com)\/(?:#!\/)?(\w+)\/status(?:es)?\/(\d+)(?:\/.*)?`)

// normalizeUrl normalizes the URL and removes tracking parameters.
func NormalizeUrl(input string) (string, error) {
	// Parse the URL
	u, err := url.Parse(input)
	if err != nil {
		return "", err
	}

	// Convert scheme and host to lowercase
	u.Scheme = strings.ToLower(u.Scheme)
	u.Host = strings.ToLower(u.Host)

	// Remove default ports (80 for HTTP, 443 for HTTPS)
	if (u.Scheme == "http" && u.Port() == "80") || (u.Scheme == "https" && u.Port() == "443") {
		u.Host = u.Hostname()
	}

	// Remove trailing slash if it's not the root
	if u.Path != "/" && strings.HasSuffix(u.Path, "/") {
		u.Path = strings.TrimSuffix(u.Path, "/")
	}

	// Remove tracking parameters
	u.RawQuery = removeTrackingParams(u.String(), u.RawQuery)

	// Sort query parameters by key
	if len(u.RawQuery) > 0 {
		queryParams := u.Query()
		var keys []string
		for k := range queryParams {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		var sortedParams []string
		for _, k := range keys {
			sortedParams = append(sortedParams, fmt.Sprintf("%s=%s", k, queryParams.Get(k)))
		}
		u.RawQuery = strings.Join(sortedParams, "&")
	}

	// Remove fragment
	u.Fragment = ""

	// Canonicalize the path by decoding percent-encoded characters and removing redundant slashes
	u.Path = decodePath(u.Path)

	// Return the normalized URL
	return u.String(), nil
}

// removeTrackingParams removes common tracking parameters from the query string.
func removeTrackingParams(link string, query string) string {
	trackingParams := []string{"utm_source", "utm_medium", "utm_campaign", "utm_term", "utm_content", "fbclid", "gclid"}

	if TWEET_URL_REGEX.MatchString(link) {
		// remove tracking parameters from tweet links:
		// https://x.com/exhibitSaveSoil/status/1532405039217664002?s=20&t=R91quPajs0E53Yds-fhv2g
		trackingParams = append(trackingParams, "s", "t")
	}

	// Parse query parameters
	queryParams := make(url.Values)
	queryValues, err := url.ParseQuery(query)
	if err != nil {
		return ""
	}
	for k, v := range queryValues {
		// Skip tracking parameters
		if contains(trackingParams, k) {
			continue
		}
		queryParams[k] = v
	}

	return queryParams.Encode()
}

// contains checks if a parameter is in the list of tracking parameters
func contains(params []string, key string) bool {
	for _, param := range params {
		if param == key {
			return true
		}
	}
	return false
}

// decodePath decodes percent-encoded characters and removes redundant slashes.
func decodePath(path string) string {
	// Decode percent-encoded characters
	decodedPath, err := url.PathUnescape(path)
	if err != nil {
		return path
	}

	// Remove redundant slashes and handle `..` or `.` segments
	segments := strings.Split(decodedPath, "/")
	var canonicalSegments []string
	for _, segment := range segments {
		if segment == "" || segment == "." {
			continue
		}
		if segment == ".." {
			if len(canonicalSegments) > 0 {
				canonicalSegments = canonicalSegments[:len(canonicalSegments)-1]
			}
		} else {
			canonicalSegments = append(canonicalSegments, segment)
		}
	}
	return "/" + strings.Join(canonicalSegments, "/")
}
