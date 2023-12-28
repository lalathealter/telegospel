package business

type BSI interface {
	GetPassageLink(string, string) string
	String() string
}

type BibleSource string

const (
	BibleGatewayLink BibleSource = "https://biblegateway.com"
)

func ChooseProvider(source BibleSource) (bsi BSI) {
	switch source {
	case BibleGatewayLink:
		bsi = BibleGatewaySource(source)
	default:
		bsi = BibleGatewaySource(BibleGatewayLink)
	}
	return
}
