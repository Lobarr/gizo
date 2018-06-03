package p2p

import (
	"net/url"
)

//ParseAddr returns dispatcher url as a map
func ParseAddr(addr string) (map[string]string, error) {
	temp := make(map[string]string)
	parsed, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	temp["pub"] = parsed.User.Username()
	temp["ip"] = parsed.Hostname()
	temp["port"] = parsed.Port()
	return temp, nil
}
