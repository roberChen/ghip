package ipget

import (
	"fmt"
	"testing"

	"github.com/roberChen/ghip/conf"
)

func TestIpaddressURL(t *testing.T) {
	cases := []struct {
		raw string
		out string
	}{
		{"central.github.com", "https://github.com.ipaddress.com/central.github.com"},
		{"github.com", "https://github.com.ipaddress.com"},
		{"github.io", "https://github.io.ipaddress.com"},
		{"github", "github"},
	}

	for _, c := range cases {
		if c.out != ipaddressURL(c.raw) {
			t.Errorf("failed for %s: %s != %s", c.raw, c.out, ipaddressURL(c.raw))
		}
	}
}

func TestIpaddressServer(t *testing.T) {
	conf.RETRY = 5
	cases := []string{
		"github.com",
		"gitee.com",
		"github.io",
		"central.github.com",
		"avatars1.githubusercontent.com",
		"github-cloud.s3.amazonaws.com",
		"favicons.githubusercontent.com",
		"media.githubusercontent.com",
		"desktop.githubusercontent.com",
	}

	for i, raw := range cases {
		fmt.Printf("testing case `%s`\n", raw)
		o, err := IPAddrServer(raw)
		if err != nil {
			t.Errorf("case %d: %s", i, err)
		}
		fmt.Println(raw, ":", o)
	}
}
