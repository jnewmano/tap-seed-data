package seeddata

import (
	"embed"
)

const (
	lastNameDataset        = "lastnames.txt"
	femaleFirstNameDataset = "firstnames-female.txt"
	maleFirstNameDataset   = "firstnames-male.txt"
)

//go:embed namedata
var nameFiles embed.FS

func init() {
	files := []string{
		lastNameDataset,
		femaleFirstNameDataset,
		maleFirstNameDataset,
	}

	for _, v := range files {
		err := loadNameFile(v)
		if err != nil {
			panic(err)
		}
	}
}

func RandomName() (firstName string, lastName string) {
	if randomGenerator.Float64() < 0.5 {
		return RandomMaleName()
	}
	return RandomFemaleName()
}

func RandomFemaleName() (firstName string, lastName string) {
	firstName = randomFemaleFirstName()
	lastName = randomLastName()
	return
}

func RandomMaleName() (firstName string, lastName string) {
	firstName = randomMaleFirstName()
	lastName = randomLastName()
	return
}

func randomMaleFirstName() string {
	return pick(maleFirstNameDataset)
}

func randomFemaleFirstName() string {
	return pick(femaleFirstNameDataset)
}

func randomLastName() string {
	return pick(lastNameDataset)
}
