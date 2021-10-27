package seeddata

import (
	"embed"
)

const (
	familyNameDataset      = "familynames.txt"
	femaleGivenNameDataset = "givennames-female.txt"
	maleGivenNameDataset   = "givennames-male.txt"
)

//go:embed namedata
var nameFiles embed.FS

func init() {
	files := []string{
		familyNameDataset,
		femaleGivenNameDataset,
		maleGivenNameDataset,
	}

	for _, v := range files {
		err := loadNameFile(v)
		if err != nil {
			panic(err)
		}
	}
}

func RandomName() (givenName string, familyName string) {
	if randomGenerator.Float64() < 0.5 {
		return RandomMaleName()
	}
	return RandomFemaleName()
}

func RandomFemaleName() (givenName string, familyName string) {
	givenName = randomFemaleGivenName()
	familyName = randomFamilyName()
	return
}

func RandomMaleName() (givenName string, familyName string) {
	givenName = randomMaleGivenName()
	familyName = randomFamilyName()
	return
}

func randomMaleGivenName() string {
	return pick(maleGivenNameDataset)
}

func randomFemaleGivenName() string {
	return pick(femaleGivenNameDataset)
}

func randomFamilyName() string {
	return pick(familyNameDataset)
}
