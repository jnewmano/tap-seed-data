package seeddata

import (
	"strconv"
	"time"
)

func init() {
	loadPopulationData()
}

// RandomBirthday generates a random birthday using
// the age distribution of the United States
func RandomBirthday() time.Time {

	bucket := pick("ages")
	ageBucket, _ := strconv.Atoi(bucket)

	// once we have the age bucket, we can get an age out of that bucket
	// each bucket is 5 years wide
	ar := randomGenerator.Intn(5)
	age := ageBucket*5 + int(ar)

	dayR := randomGenerator.Intn(365)
	birthday := time.Now().AddDate(-age, 0, -dayR)

	return birthday
}

// US census data 2019
var populationDistribution = []struct {
	Desc  string
	Count int
}{
	{"Under 5 years", 19736},
	{"5 to 9 years	", 20212},
	{"10 to 14 years", 20827},
	{"15 to 19 years", 20849},
	{"20 to 24 years", 21254},
	{"25 to 29 years", 23277},
	{"30 to 34 years", 21932},
	{"35 to 39 years", 21443},
	{"40 to 44 years", 19584},
	{"45 to 49 years", 20345},
	{"50 to 54 years", 20355},
	{"55 to 59 years", 21163},
	{"60 to 64 years", 20592},
	{"65 to 69 years", 17356},
	{"70 to 74 years", 14131},
	{"75 to 79 years", 9357},
	{"80 to 84 years", 6050},
	{"85 to 89 years", 5535}, // split out (85-90] into one extra bucket
	{"90 years and over", 358},
}
