package common

import (
	"fmt"
	"net"
	"net/url"
)

func ParseDomainName(urlParam string) (string, error) {
	u, err := url.Parse(urlParam)
	if err != nil {
		fmt.Printf("Failed to parse url:%+v\n", err)
		return "", err
	}
	if u.Host == "" {
		return u.Path, nil
	}
	host, port, err := net.SplitHostPort(u.Host)
	if err != nil {
		if port == "" {
			host = u.Host
		} else {
			fmt.Printf("Failed to split host and port:%+v\n", err)
		}
	}
	return host, nil
}
