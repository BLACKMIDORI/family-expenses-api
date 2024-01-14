package controllers

import (
	"encoding/json"
	"family-expenses-api/core"
	"family-expenses-api/services"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"runtime/debug"
	"strings"
)

type Filters map[string][]string

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (filters Filters) Get(key string) string {
	values := filters[key]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
func sendInternalServerErrorOnPanic(responseWriter http.ResponseWriter) {
	if err := recover(); err != nil {
		replyAsJson(responseWriter, 500, map[string]any{
			"error": "Internal Server Error",
		})
		log.Println("Internal Server Error:", fmt.Errorf("%v", err))
		fmt.Println("\n" + string(debug.Stack()))
	}
}

func match(regex string, path string) bool {
	return regexp.MustCompile(regex).MatchString(path)
}

func routeParams(request *http.Request, pattern string) map[string]string {
	matches := regexp.MustCompile("{([^{}]+)}").FindAllStringSubmatch(pattern, -1)
	var keys []string
	for _, match := range matches {
		keys = append(keys, match[1])
	}
	newRegex := "^" + pattern + "$"
	for _, key := range keys {
		newRegex = strings.Replace(newRegex, "{"+key+"}", "(.*)", -1)
	}
	values := regexp.MustCompile(newRegex).FindStringSubmatch(request.URL.Path)[1:]

	if len(keys) != len(values) {
		panic("Could not parse route params. " + fmt.Sprint("Keys:", keys, "Values:", values))
	}
	params := make(map[string]string)
	for i, key := range keys {
		params[key] = values[i]
	}
	return params
}

func filtersFromQuery(request *http.Request) (filters Filters) {
	filters = make(Filters)
	for key, values := range request.URL.Query() {
		for _, value := range values {
			if key == "filter" && strings.Contains(value, "__") {
				parts := strings.SplitN(value, "__", 2)
				filterKey := parts[0]
				filterValue := parts[1]
				filters[filterKey] = append(filters[filterKey], filterValue)
			}
		}
	}
	return
}

func bodyJson(request *http.Request) map[string]any {
	var body map[string]any
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		panic(err)
	}
	return body
}

func replyAsJson(responseWriter http.ResponseWriter, status int, body any) {
	responseWriter.WriteHeader(status)
	_ = json.NewEncoder(responseWriter).Encode(body)
}

func parseUser(request *http.Request) (core.AuthenticatedUser, error) {
	tokens := request.Header["Authorization"]
	if len(tokens) == 0 {
		return core.AuthenticatedUser{}, nil
	}
	bearerToken := tokens[0]
	token := strings.ReplaceAll(bearerToken, "Bearer ", "")

	return services.CreateTokenService(nil).VerifyAccessToken(token)
}

func parseUserAsync(request *http.Request) (userChan chan core.AuthenticatedUser) {
	userChan = make(chan core.AuthenticatedUser)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Could not parse user:", err)
				userChan <- core.AuthenticatedUser{}
			}
		}()
		user, err := parseUser(request)
		if err != nil {
			log.Println("Could not parse user:", err)
		}
		userChan <- user
	}()
	return
}
