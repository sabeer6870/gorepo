package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"urllookupservice/common"
)

func TestIsMalwareUrl(t *testing.T) {
	t.Run("returns Pepper's score", func(t *testing.T) {
		s := NewServer("8080")
		s.urlMap.Store("git-scm.com", true)
		testCases := []struct{
			name string
			url string
			expected string
		} {
			{
				"happy path test",
				"https://git-scm.com:9090/download/win",
				"true",
			},
			{
				"happy path test 2",
				"https://www.google.com/",
				"false",
			},
		}

		for _, testCase := range testCases {
			tu := url.QueryEscape(testCase.url)
			request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/v1/lookupservice/urls/%s", tu),nil)
			response := httptest.NewRecorder()
			handler := http.HandlerFunc(s.IsMalwareURL)
			req := mux.SetURLVars(request, map[string]string{
				"url": testCase.url,
			})
			handler.ServeHTTP(response, req)
			got := response.Body.String()
			if response.Result().StatusCode != http.StatusOK {
				t.Errorf("Got status:%d, expected %d", response.Result().StatusCode, http.StatusOK)
			}
			if got != testCase.expected {
				t.Errorf("testcase:%s got %q, want %q", testCase.name, got, testCase.expected)
			}
		}
	})
}

func TestAddUrl(t *testing.T) {
	t.Run("returns Pepper's score", func(t *testing.T) {
		s := NewServer("8080")
		testCases := []struct{
			name string
			url string
			expected string
		} {
			{
				"happy path test",
				"https://git-scm.com:9090/download/win",
				"true",
			},
			{
				"happy path test 2",
				"https://www.google.com/",
				"true",
			},
		}

		for _, testCase := range testCases {
			tu := url.QueryEscape(testCase.url)
			request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/v1/lookupservice/urls/%s", tu),nil)
			response := httptest.NewRecorder()
			handler := http.HandlerFunc(s.AddUrlToMalwareList)
			req := mux.SetURLVars(request, map[string]string{
				"url": testCase.url,
			})
			handler.ServeHTTP(response, req)
			if response.Result().StatusCode != http.StatusOK {
				t.Errorf("Got status:%d, expected %d", response.Result().StatusCode, http.StatusOK)
			}
			d, _ := common.ParseDomainName(testCase.url)
			v, ok := s.urlMap.Load(d)
			if !ok || v.(bool) != true {
				t.Errorf("ok:%v, v :%v", ok, v)
			}
		}
	})
}


func TestRemoveUrl(t *testing.T) {
	t.Run("returns Pepper's score", func(t *testing.T) {
		s := NewServer("8080")
		s.urlMap.Store("git-scm.com", true)
		testCases := []struct{
			name string
			url string
			expected string
		} {
			{
				"happy path test",
				"https://git-scm.com:9090/download/win",
				"true",
			},
		}

		for _, testCase := range testCases {
			tu := url.QueryEscape(testCase.url)
			request, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/lookupservice/url/%s", tu),nil)
			response := httptest.NewRecorder()
			handler := http.HandlerFunc(s.RemoveUrlFromMalwareList)
			req := mux.SetURLVars(request, map[string]string{
				"url": testCase.url,
			})
			handler.ServeHTTP(response, req)
			if response.Result().StatusCode != http.StatusOK {
				t.Errorf("Got status:%d, expected %d", response.Result().StatusCode, http.StatusOK)
			}
			d, _ := common.ParseDomainName(tu)
			v, ok := s.urlMap.Load(d)
			if ok || v != nil {
				t.Errorf("ok:%v, v :%v", ok, v)
			}
		}
	})
}