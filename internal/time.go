package internal

import "time"

func TimestampByUnit(t time.Time, unit string) int64 {
	switch unit {
	case "ns":
		return t.UnixNano()
	case "us", "µs":
		return t.UnixNano() / int64(time.Microsecond) // support go version less than 1.17
	case "ms":
		return t.UnixNano() / int64(time.Millisecond) // support go version less than 1.17
	case "s":
		return t.Unix()
	}

	return 0
}

func GetTimeByUnit(ts int64, unit string) time.Time {
	switch unit {
	case "ns":
		return time.Unix(0, ts)
	case "us", "µs":
		return time.Unix(0, ts*int64(time.Microsecond)) // support go version less than 1.17
	case "ms":
		return time.Unix(0, ts*int64(time.Millisecond)) // support go version less than 1.17
	case "s":
		return time.Unix(ts, 0)
	}

	return time.Time{}
}
