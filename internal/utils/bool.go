package utils

const (
	Yes = "yes"
	No  = "no"
)

func BoolToString(value bool) string {
	if value {
		return Yes
	} else {
		return No
	}
}
