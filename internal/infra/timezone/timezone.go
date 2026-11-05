package timezone

import (
	"time"
	_ "time/tzdata" // 내장 시간대 데이터베이스를 바이너리에 포함
)

func MustSetLocalTimezone(tz string) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		panic("failed to load local timezone [" + tz + "]: " + err.Error())
	}
	time.Local = loc
}
