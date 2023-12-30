package repositories

import "time"

func validTimeToUnixOrNil(theTime time.Time) any {
	if (theTime == time.Time{}) {
		return nil
	} else {
		return theTime.Unix()
	}
}
