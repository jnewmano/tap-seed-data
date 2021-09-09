package seeddata

import (
	"fmt"
)

type contactOptions struct {
	attributes           map[string]bool
	genderDistribution   float64
	deceasedDistribution float64
}

func newContactOptions() contactOptions {
	return contactOptions{
		attributes:           make(map[string]bool),
		genderDistribution:   0.5,
		deceasedDistribution: 0.05,
	}
}

type ContactOption interface {
	apply(*contactOptions)
}

// funcContactOption wraps a function that modifies contactOptions
type funcContactOption struct {
	f func(*contactOptions)
}

func (f *funcContactOption) apply(do *contactOptions) {
	f.f(do)
}

func newFuncContactOption(f func(*contactOptions)) *funcContactOption {
	return &funcContactOption{
		f: f,
	}
}

// NewContactAttributesOptions takes a map of attributes and the
// probability that the attribute should be true
func NewContactAttributesOption(attributes map[string]float64) ContactOption {
	return newFuncContactOption(func(a *contactOptions) {
		for k, v := range attributes {
			set := false
			if randomGenerator.Float64() < v {
				set = true
			}
			a.attributes[k] = set
		}
	})
}

func GenerateContact(contactID string, opts ...ContactOption) Contact {

	co := newContactOptions()

	for _, v := range opts {
		v.apply(&co)
	}

	gender := "Female"
	randomName := RandomFemaleName
	if randomGenerator.Float64() < co.genderDistribution {
		gender = "Male"
		randomName = RandomMaleName
	}

	birthday := RandomBirthday().Format("2006-01-02")
	contactMethods := randomContactMethods()

	deceased := false
	if randomGenerator.Float64() < co.deceasedDistribution {
		deceased = true
	}

	firstName, lastName := randomName()
	c := Contact{
		ContactID: contactID,

		FirstName:     firstName,
		LastName:      lastName,
		MiddleName:    "",
		PreferredName: "",
		Company:       "",

		Gender: gender,

		Birthdate: birthday,
		Deceased:  deceased,

		Address: Address{},

		ContactMethods: contactMethods,

		Attributes:     co.attributes,
		AdditionalData: map[string]string{},
	}

	return c
}

func randomContactMethods() []ContactMethod {
	// only a 10% chance of not having any available contact methods
	if randomGenerator.Float64() < .1 {
		return []ContactMethod{}
	}
	// guarantee at least one contact method
	n := 1 + randomGenerator.Intn(4)
	list := make([]ContactMethod, n)
	for i := 0; i < n; i++ {
		c := RandomPhoneContactMethod()
		list[i] = c
	}

	// add an email address 70% of the time
	if randomGenerator.Float64() < .7 {
		c := RandomEmailContactMethod()

		list = append(list, c)
	}
	return list
}

func RandomPhoneContactMethod() ContactMethod {
	types := []string{"Home", "Mobile", "Work", "Other"}
	c := ContactMethod{
		Type:    "Phone",
		Tags:    []string{types[randomGenerator.Intn(len(types))]},
		Address: RandomPhoneNumber(),
	}

	return c
}

func RandomEmailContactMethod() ContactMethod {
	c := ContactMethod{
		Type:    "Email",
		Address: RandomEmailAddress(),
	}
	return c
}

func RandomEmailAddress() string {
	return randomString(12) + "@example.com"
}

// RandomPhoneNumber returns a random US phone number
// with 555 area code in E.164 format
func RandomPhoneNumber() string {
	r := randomGenerator.Intn(10000000)
	n := fmt.Sprintf("+1555%07d", r)
	return n
}
