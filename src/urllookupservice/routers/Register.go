package routers

import (
	"net/http"

	"urllookupservice/server"
)

func CreateRouter(s *server.Server) {
	//v1
	s.Router.HandleFunc("/v1/lookupservice/urls/{url}", s.IsMalwareURL).
		Methods(http.MethodGet)
	s.Router.HandleFunc("/v1/lookupservice/urls/{url}", s.AddUrlToMalwareList).
		Methods(http.MethodPost, http.MethodPut)
	s.Router.HandleFunc("/v1/urlookupservicels/urls/{url}", s.RemoveUrlFromMalwareList).
		Methods(http.MethodDelete)

	//v2
	s.Router.HandleFunc("/v2/lookupservice/urls/{url}", s.IsMalwareURLv2).
		Methods(http.MethodGet)
	s.Router.HandleFunc("/v2/lookupservice/urls/{url}", s.AddMalwareURL).
		Methods(http.MethodPost, http.MethodPut)
	s.Router.HandleFunc("/v2/lookupservice/urls/{url}", s.RemoveMalwareURL).
		Methods(http.MethodDelete)
}

