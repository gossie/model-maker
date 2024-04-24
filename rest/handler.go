package rest

import (
	"net/http"

	"github.com/gossie/modelling-service/views"
)

type ServiceHandler func(view *views.View) http.HandlerFunc
