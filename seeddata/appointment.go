package seeddata

import (
	"time"
)

// Appointments are events that have been scheduled
// Use reminders for services or events that should be scheduled in the future
type Appointment struct {
	AppointmentID string
	ContactID     string // primary contact or responsible party for the appointment
	ReceiverID    string // person or thing that will be receiving the services for the appointment
	ReceiverType  string // type of thing being serviced at the appointment

	StartTime         time.Time
	EndTime           time.Time
	ServiceLocationID int // id of the place where the service will be provided

	Services  []AppointmentService  // list of services that are scheduled to be received during the appointment
	Resources []AppointmentResource // list of resources assigned for the appointment, operatories, personnel, equipment, etc

	Events []struct {
		Type    string
		EventID string
		Time    time.Time
	} // list of events and times that they occured

	Attributes     map[string]bool   // boolean values that give additional information about an appointment
	AdditionalData map[string]string // key value pairs with additional data about the appointment
}

type AppointmentService struct {
}

type AppointmentResource struct {
	Type       string
	ResourceID string
}

const (
	ReceiverTypeContact = "Contact"

	AppointmentEventKeyConfirmed = "Confirmed"

	AppointmentResourceKeyOperatory = "Operatory"
	AppointmentResourceKeyDoctor    = "Doctor"
	AppointmentResourceKeyHygenist  = "Hygenist"
)

func (a *Appointment) HasEvent(k string) bool {
	for _, v := range a.Events {
		if v.Type == k {
			return true
		}
	}
	return false
}

func (a *Appointment) GetResource(k string) []AppointmentResource {
	var list []AppointmentResource
	for _, v := range a.Resources {
		if v.Type == k {
			list = append(list, v)
		}
	}
	return list
}
