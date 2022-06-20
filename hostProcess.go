package main

type (
	hostProcess struct {
		api string
		ex  []Expiry
		op  []Feature
	}
	option  int    // usable option
	Expiry  option // expiry options
	Feature option // featured options
)

var (
	host = hostProcess{"url",
		[]Expiry{Hour, Day, Week, Month, Year, Never},
		[]Feature{Burn}}
)

const (
	Hour       Expiry  = iota // expires after 1 hour
	Day                       // expires after 1 day
	Week                      // expires after 1 week
	Month                     // expires after 1 month
	Year                      // expires after 1 year
	Never                     // expires `"never"`
	Burn       Feature = iota // delete after reading once
	UploadFile                // upload a file
)

func (e Expiry) String() string {
	switch e {
	case Hour:
		return "1hour"
	case Day:
		return "1day"
	case Week:
		return "1week"
	case Month:
		return "1month"
	case Year:
		println("Defaulting to 1 month, the maximum allowed.")
		return "1month"
	case Never:
		println("Defaulting to 1 month, the maximum allowed.")
		return "1month"
	}
	return ""
}
