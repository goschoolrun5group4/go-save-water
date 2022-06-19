package middleware

import "net/http"

func ContentTypeHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Header.Get("Content-type") != "application/json" {
			res.WriteHeader(http.StatusUnprocessableEntity)
			res.Write([]byte("422 - Invalid Content-type"))
			return
		}
		h.ServeHTTP(res, req)
	})
}
