package middleware

import "net/http"

func Any(next http.HandlerFunc) http.HandlerFunc {
	return traceRequest(
		logIncomingRequests(
			addContentType("application/json",
				enableCors(next))))
}
