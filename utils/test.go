package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

func MakeTestServeMux(issuesLength int, shouldBePossiblySynced bool, issueStatusCode int) *http.ServeMux {
	sb := new(strings.Builder)
	sb.WriteString("[")

	mux := http.NewServeMux()

	for i := range make([]int, issuesLength) {
		issueKey := fmt.Sprintf("TRADE-%d", i)
		sb.WriteString(fmt.Sprintf(
			`{
				"categoryId": 1, 
				"categoryName": "Wren Kitchens",
				"fromDate": "2023-03-11T00:00:00Z",
				"toDate": "2023-03-11T00:00:00Z",
				"comment": "%s:Dev:Comment",
				"minutes": 30
			}`,
			issueKey,
		))

		if i != issuesLength-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString("]")

	mux.Handle("/api/time/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		Must(w.Write([]byte(sb.String())))
	}))

	mux.Handle("/rest/api/2/issue/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		splited := strings.Split(r.URL.String(), "/")
		b := fmt.Sprintf("{\"key\": \"%s\"}", splited[len(splited)-1])

		w.WriteHeader(issueStatusCode)
		Must(w.Write([]byte(b)))
	}))

	mux.Handle("/rest/api/2/search", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var b string
		if shouldBePossiblySynced {
			b = `{
				"issues": [{"key": "TRADE-1"}],
				"startAt": 0,
				"maxResults": 1,
				"total": 1
			}`
		} else {
			b = `{
				"issues": [],
				"startAt": 0,
				"maxResults": 0,
				"total": 0
			}`
		}

		w.WriteHeader(200)
		Must(w.Write([]byte(b)))
	}))

	return mux
}

func Host(ts *httptest.Server) string {
	return ts.URL + "/"
}
