package background

import (
	"time"

	"github.com/wufe/polo/pkg/models"
	"github.com/wufe/polo/pkg/storage"
)

type SessionStartWorker struct {
	sessionStorage *storage.Session
	mediator       *Mediator
}

func NewSessionStartWorker(
	sessionStorage *storage.Session,
	mediator *Mediator,
) *SessionStartWorker {
	worker := &SessionStartWorker{
		sessionStorage: sessionStorage,
		mediator:       mediator,
	}
	worker.startAcceptingSessionStartRequests()
	return worker
}

func (w *SessionStartWorker) startAcceptingSessionStartRequests() {
	go func() {
		for {
			session := <-w.mediator.StartSession.Chan
			w.MarkSessionAsStarted(session)
		}
	}()
}

func (w *SessionStartWorker) MarkSessionAsStarted(session *models.Session) {
	app := session.Application
	session.SetStatus(models.SessionStatusStarted)
	session.ResetStartupRetriesCount()
	if app.Watch.Contains(session.Checkout) {
		session.SetMaxAge(-1)
	} else {
		session.SetMaxAge(session.Application.Recycle.InactivityTimeout)
		if session.GetMaxAge() > 0 {
			w.startSessionInactivityTimer(session)
		}
	}

	w.sessionStorage.Update(session)
}

func (w *SessionStartWorker) startSessionInactivityTimer(session *models.Session) {
	session.SetInactiveAt(time.Now().Add(time.Second * time.Duration(session.Application.Recycle.InactivityTimeout)))
	go func() {
		for {
			status := session.GetStatus()
			if status != models.SessionStatusStarted {
				return
			}

			if time.Now().After(session.GetInactiveAt()) {
				w.mediator.DestroySession.Enqueue(session, nil)
				return
			}
			session.DecreaseMaxAge()
			time.Sleep(1 * time.Second)
		}
	}()
}