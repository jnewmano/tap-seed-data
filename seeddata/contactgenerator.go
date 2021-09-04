package seeddata

import "fmt"

func GenerateContact(contactID string) Contact {
	g := randomGenerator.Intn(2)
	gender := "Female"
	randomName := RandomFemaleName
	if g == 1 {
		gender = "Male"
		randomName = RandomMaleName
	}

	birthday := RandomBirthday().Format("2006-01-02")
	contactMethods := randomContactMethods()

	isActive := true
	if randomGenerator.Float64() < .2 {
		isActive = false
	}

	deceased := false
	if randomGenerator.Float64() < 0.05 {
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

		Attributes: map[string]bool{
			"IsPatient": true,
			"Active":    isActive,
		},
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
		Type:     "Phone",
		Tags:     []string{types[randomGenerator.Intn(len(types))]},
		Priority: randomGenerator.Intn(100),
		Address:  RandomPhoneNumber(),
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
