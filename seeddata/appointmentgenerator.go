package seeddata

import (
	"strconv"
	"time"
)

type appointmentOptions struct {
	contactID string

	timeWindowStart    time.Time
	timeWindowEnd      time.Time
	startTimeAlignment time.Duration // rule for aligning the start time (start time minute mod this must be zero)

	duration  time.Duration         // list of appointment durations to select from
	resources []AppointmentResource // resources to assign to the generated appointment

	allowTimeFunc func(time.Time, time.Time) bool
}

type AppointmentOption interface {
	apply(*appointmentOptions)
}

// funcAppointmentOption wraps a function that modifies appointmentOptions
type funcAppointmentOption struct {
	f func(*appointmentOptions)
}

func (f *funcAppointmentOption) apply(do *appointmentOptions) {
	f.f(do)
}

func newFuncAppointmentOption(f func(*appointmentOptions)) *funcAppointmentOption {
	return &funcAppointmentOption{
		f: f,
	}
}

func NewAppointmentDateIntervalOption(minDate time.Time, maxDate time.Time) AppointmentOption {
	return newFuncAppointmentOption(func(a *appointmentOptions) {
		a.timeWindowStart = minDate
		a.timeWindowEnd = maxDate
	})
}

func NewAppointmentTimeFilterOption(f func(time.Time, time.Time) bool) AppointmentOption {
	return newFuncAppointmentOption(func(a *appointmentOptions) {
		a.allowTimeFunc = f
	})
}

func NewAppointmentStartTimeMinStepOption(step time.Duration) AppointmentOption {
	return newFuncAppointmentOption(func(a *appointmentOptions) {
		a.startTimeAlignment = step
	})
}

func NewAppointmentContactIDRangeOption(min int, max int) AppointmentOption {
	return newFuncAppointmentOption(func(a *appointmentOptions) {
		a.contactID = strconv.Itoa(randomGenerator.Intn(max-min+1) + min)
	})
}

func NewAppointmentContactIDOption(id string) AppointmentOption {
	return newFuncAppointmentOption(func(a *appointmentOptions) {
		a.contactID = id
	})
}

func NewAppointmentDurationOption(durations ...time.Duration) AppointmentOption {
	return newFuncAppointmentOption(func(a *appointmentOptions) {
		idx := randomGenerator.Intn(len(durations))
		a.duration = durations[idx]
	})
}

func NewAppointmentResourceOption(resourceType string, options []string) AppointmentOption {
	return newFuncAppointmentOption(func(a *appointmentOptions) {
		idx := randomGenerator.Intn(len(options))
		a.resources = append(a.resources,
			AppointmentResource{
				Type:       resourceType,
				ResourceID: options[idx],
			})
	})
}

func GenerateAppointment(appointmentID string, options ...AppointmentOption) Appointment {

	var ao appointmentOptions
	for _, v := range options {
		v.apply(&ao)
	}

	if ao.timeWindowStart.IsZero() {
		ao.timeWindowStart = time.Now().AddDate(-1, 0, 0)
	}
	if ao.timeWindowEnd.IsZero() {
		ao.timeWindowEnd = time.Now().AddDate(1, 0, 0)
	}

	atsRange := ao.timeWindowEnd.Sub(ao.timeWindowStart) / time.Second

	var ats, ate time.Time
	for {
		ats = ao.timeWindowStart.Add(time.Second * time.Duration(randomGenerator.Intn(int(atsRange))))

		// round the appointment minute component down, remove seconds and sub second components down
		if ao.startTimeAlignment == 0 {
			ao.startTimeAlignment = time.Minute
		}

		ats = ats.Truncate(ao.startTimeAlignment)
		ate = ats.Add(time.Duration(ao.duration) * time.Minute)

		if ao.allowTimeFunc == nil {
			break // no filtering needs to happen, break
		}
		if ao.allowTimeFunc(ats, ate) {
			break
		}
	}

	var a Appointment
	a.AppointmentID = appointmentID
	a.ContactID = ao.contactID
	a.ReceiverID = ao.contactID // person receiving the services is the same as the primary contact
	a.ReceiverType = ReceiverTypeContact

	a.StartTime = ats
	a.EndTime = ate

	a.Resources = ao.resources

	a.Events = []AppointmentEvent{
		{
			Type: AppointmentEventKeyConfirmed,
			Time: ats.Add(-time.Hour),
		},
	} // list of events and times that they occured

	return a
}
