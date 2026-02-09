package httpapi

import (
	"doctormakarhina/lumos/internal/core/notify"
	"doctormakarhina/lumos/internal/core/payments"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type paymentRegTrial struct {
	srv     payments.Service
	notifer notify.Service
}

func (s *paymentRegTrial) Handle(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	phone := r.FormValue("phone")
	daysDuration := r.FormValue("daysDuration")

	if email == "" {
		s.notifer.ForAdmin(fmt.Sprintf("[TrialFormHandler] recieve email form without email field specified, name = %s, phone = %s", name, phone))
		w.WriteHeader(http.StatusCreated)
		// http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	parsedDuration, err := strconv.Atoi(daysDuration)
	if err != nil {
		s.notifer.ForAdmin(fmt.Sprintf("[TrialFormHanlder] recieve email form without valid duration field specified, email = %s, duration = %s", email, daysDuration))
		w.WriteHeader(http.StatusCreated)
		// http.Error(w, "invalid duration", http.StatusBadRequest)
		return
	}

	trialDuration := time.Duration(parsedDuration) * time.Hour * 24

	err = s.srv.RegisterFromTrial(r.Context(), email, name, phone, trialDuration)
	if err != nil {
		s.notifer.ForAdmin(fmt.Sprintf("[TrialFormHandler] failed to register user from trial, email = %s, name = %s, phone = %s, duration = %d", email, name, phone, trialDuration))
		w.WriteHeader(http.StatusCreated)
		// http.Error(w, "failed to register", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
