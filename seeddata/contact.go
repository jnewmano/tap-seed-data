package seeddata

type Contact struct {
	ContactID string

	GivenName     string
	FamilyName    string
	MiddleName    string
	PreferredName string
	Company       string

	Gender string

	Birthdate string
	Deceased  bool

	Address        Address
	ContactMethods []ContactMethod

	Attributes     map[string]bool   // boolean values that give additional information about a person
	AdditionalData map[string]string // key value pairs with additional data about the contact
}

func (c *Contact) KeyProperties() []string {
	return []string{"ContactID"}
}

type Address struct {
	Address  string
	Address2 string
	City     string
	State    string
	PostCode string
	Country  string
}

type ContactMethod struct {
	Type    string   // Phone/email/social handle
	Tags    []string // arbitrary tags that add additional information about the contact method, then home, mobile, work, fax, etc
	Address string   // phone number, email, social handle, etc
}

func (c *Contact) BestContactMethod(t string, attributes ...string) ContactMethod {
	return BestContactMethod(c.ContactMethods, t, attributes...)
}

func BestContactMethod(cs []ContactMethod, t string, attributes ...string) ContactMethod {
	for _, v := range cs {
		if v.Type != t {
			continue
		}

		found := false
		for _, a := range attributes {
			if has(a, v.Tags) {
				found = true
				break
			}
		}
		if !found {
			continue
		}

		return v
	}
	return ContactMethod{}
}
