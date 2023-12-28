package business

import (
	"net/url"
)

type BibleGatewaySource string

func (bg BibleGatewaySource) GetPassageLink(osis string, version string) string {
	linkObj := formURL(string(bg), "/passage/")
	linkObj = attachQuery(linkObj, "search", osis)
	attachQuery(linkObj, "version", version)
	// link := linkObj.String()

	return linkObj.String()
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
