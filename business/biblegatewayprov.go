package business

import (
	"fmt"
	"net/url"
	"strings"
)

type BibleGatewaySource string

func (bg BibleGatewaySource) GetPassageLink(osisColl []string, version string) string {
	osisQuery := strings.Join(osisColl, ", ")

	linkObj := formURL(string(bg), "/passage/")
	linkObj = attachQuery(linkObj, "search", osisQuery)
	linkObj = attachQuery(linkObj, "version", version)

	linkMsg := fmt.Sprintf("\n[%v](%v)", osisQuery, linkObj.String())
	return linkMsg
}

func (bg BibleGatewaySource) String() string {
	return "Bible Gateway"
}

func formURL(base, path string) *url.URL {
	u, err := url.Parse(base)
	if err != nil {
		return nil
	}
	u = u.JoinPath(path)
	return u
}

func attachQuery(u *url.URL, key, val string) *url.URL {
	q := u.Query()
	q.Add(key, val)
	u.RawQuery = q.Encode()
	return u
}
