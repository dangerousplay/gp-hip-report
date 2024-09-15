package utils

import "gp-hip-report/internal/constants"

func BoolToString(value bool) string {
	if value {
		return constants.Yes
	} else {
		return constants.No
	}
}
