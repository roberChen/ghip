// Package ipget gets ip of url both locally or from ipaddress.com
package ipget

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/roberChen/ghip/conf"

	"github.com/PuerkitoBio/goquery"
)

const (
	IPADDRESS_SCHEMA = "https://"
	IPADDRESS_PREFIX = ".ipaddress.com"
)

// IPAddrLocal get the real ip for url:port connection via local connection try, the return will be without port info. url would not include scheme, and port should't start with colon
func IPAddrLocal(url, port string) (string, error) {
	urlport := url + ":" + port
	conn, err := net.Dial("tcp", urlport)
	if err != nil {
		return "", fmt.Errorf("connection to %s error: %s", urlport, err)
	}
	ipport := conn.RemoteAddr().String()
	return ipport[:len(ipport)-len(port)-1], nil
}

func IPAddrServer(url string) ([]string, error) {
	ipaddrurl := ipaddressURL(url)
	out := []string{}
	for i := 0; i < conf.RETRY; i++ {
		resp, err := http.Get(ipaddrurl)
		if err != nil {
			log.Printf("try to get ip of `%s` via ipaddress.com failed (retry %d): %s\n", url, i, err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			log.Printf("try to get ip of `%s` via ipaddress.com failed (retry %d): statuc code %d\n", url, i, resp.StatusCode)
			continue
		}
		defer resp.Body.Close()
		d, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Printf("try to get ip of `%s` via ipaddress.com failed (retry %d): %s\n", url, i, err)
			continue
		}
		maxi := -1
		d.Find("body > div.resp.main > main > section:nth-child(10) > table > tbody > tr").Each(func(_ int, s *goquery.Selection) {
			if strings.HasPrefix(s.Find("h3").First().Text(),"What IP address") {
				// right FAQ block
				s.Find("a").Each(func(i int, sa *goquery.Selection) {
					maxi = i
					out = append(out, sa.Text())
				})
			}
		})
		if maxi == -1 {
			return out, fmt.Errorf("try to get ip of `%s` from ippaddress.com failed: wrong selector", url)
		}
		break
	}
	if len(out) == 0 {
		return out, fmt.Errorf("try to get ip of `%s` from ipaddress.com failed, retry used up or incorrect find method", url)
	}
	return out, nil
}

// get url to get ip from ipadress.com
func ipaddressURL(url string) string {
	sp := strings.Split(url, ".")
	if len(sp) < 2 {
		return url
	} else if len(sp) == 2 {
		return IPADDRESS_SCHEMA + url + IPADDRESS_PREFIX
	}
	return IPADDRESS_SCHEMA + strings.Join(sp[len(sp)-2:], ".") + IPADDRESS_PREFIX + "/" + url
}
