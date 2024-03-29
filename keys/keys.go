package keys

const (
	TG_TOKEN                           = "TG"
	TRANSLATION                        = "translation"
	API_TRANSLATION_PATH               = "/version"
	PROVIDER                           = "provider"
	API_PROVIDER_PATH                  = "/provider"
	PLAN                               = "plan"
	API_PLAN_PATH                      = "/plan"
	READING_DAY                        = "reading_day"
	API_READING_DAY_PATH               = "/day"
	API_READING_DAY_MOVE_BACKWARD_PATH = "/prev"
	API_READING_DAY_MOVE_FORWARD_PATH  = "/next"
)

var SETTINGS_KEYS = [...]string{
	TRANSLATION, PLAN, READING_DAY, PROVIDER,
}
