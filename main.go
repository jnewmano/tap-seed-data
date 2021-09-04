package main

import (
	"log"
	"strconv"
	"time"

	"github.com/jnewmano/tap-seed-data/seeddata"
	"github.com/jnewmano/tap-seed-data/singertap"
)

func main() {

	s := singertap.New()

	attributes := map[string]float64{
		"IsActive":  0.80,
		"IsPatient": 1.00,
	}

	copts := []seeddata.ContactOption{
		seeddata.NewContactAttributesOption(attributes),
	}

	for i := 1; i <= 10000; i++ {
		contact := seeddata.GenerateContact(strconv.Itoa(i), copts...)
		err := s.WriteRecord("Contact", &contact)
		if err != nil {
			log.Fatal(err)
		}
	}

	aopts := []seeddata.AppointmentOption{
		seeddata.NewAppointmentDurationOption(45, 60, 75, 90),           // appointments should be 60 minutes long
		seeddata.NewAppointmentContactIDRangeOption(1, 10000),           // randomly generate contact ids
		seeddata.NewAppointmentStartTimeMinStepOption(time.Minute * 15), // schedule at 5 minute intervals
		seeddata.NewAppointmentDateIntervalOption(time.Now().AddDate(-1, 0, 0), time.Now().AddDate(0, 6, 0)),
		seeddata.NewAppointmentTimeFilterOption(allowStartTimeDaysAndBusinessHoursFunc),
		seeddata.NewAppointmentResourceOption(seeddata.AppointmentResourceKeyOperatory, []string{"1", "2", "3", "4"}),
		seeddata.NewAppointmentResourceOption(seeddata.AppointmentResourceKeyDoctor, []string{"1", "3"}),
		seeddata.NewAppointmentResourceOption(seeddata.AppointmentResourceKeyHygenist, []string{"2", "4", "5", "6"}),
	}

	for i := 1; i <= 100000; i++ {
		apt := seeddata.GenerateAppointment(strconv.Itoa(i), aopts...)
		err := s.WriteRecord("Appointment", &apt)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func allowStartTimeDaysAndBusinessHoursFunc(start time.Time, stop time.Time) bool {

	if start.Hour() < 9 || start.Hour() >= 17 || stop.Hour() >= 17 {
		return false
	}

	switch start.Weekday() {
	case time.Saturday, time.Sunday:
		return false
	}

	return true
}
