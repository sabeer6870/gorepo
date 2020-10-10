package common

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	var tests = []struct{
		name string
		url string
		expected string
	} {
		{
			"test0",
			"https://git-scm.com:8080/download/win",
			"git-scm.com",
		},
		{
			"test1",
			"https://git-scm.com/download/win",
			"git-scm.com",
		},
		{
			"test2",
			"https://www.chapters.indigo.ca/en-ca/electronics/ekids-paw-patrol-walkie-talkies/092298932637-item.html?gclsrc=aw.ds&gclid=CjwKCAjwq_D7BRADEiwAVMDdHmSp7M8Et_1i_Q-153r_uPy5-2WZ7ga-S1Ih8Ln7ELvZ7rYEgd5X9xoCyckQAvD_BwE&s_campaign=goo-Shopping_Smart_Baby",
			"www.chapters.indigo.ca",
		},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			d, err := ParseDomainName(tst.url)
			if err != nil || d != tst.expected {
				t.Fatal(fmt.Sprintf("d:%s\nError:%+v", d, err))
			}
		})
	}
}
