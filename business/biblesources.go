package business

type BSI interface {
	GetPassageLink(string, string) string // osis, translation code
	String() string
}

type BibleSource string

const (
	BibleGatewayLink BibleSource = "https://biblegateway.com"
	BibleGatewayCode             = "bgway"
)

var BibleSourcesColl = map[string]BibleSource{
	BibleGatewayCode: BibleGatewayLink,
}

func GetProviderInterface(source string) (bsi BSI) {
	switch source {
	case BibleGatewayCode:
		bsi = BibleGatewaySource(BibleGatewayLink)
	default:
		bsi = BibleGatewaySource(BibleGatewayLink)
	}
	return
}
