package v1

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/models"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type fakeDeleter struct{ err error }

func (f fakeDeleter) DeleteUserDataById(_, _ string) error { return f.err }

type sentMail struct{ to, subject string }

type fakeMailer struct{ sent chan sentMail }

func (f fakeMailer) SendEmail(to, subject, _ string) error {
	f.sent <- sentMail{to, subject}
	return nil
}

func deleteAccount(t *testing.T, deleteErr error) (int, chan sentMail) {
	t.Helper()
	mailer := fakeMailer{sent: make(chan sentMail, 2)}
	ctrl := &AuthController{userSvc: fakeDeleter{err: deleteErr}, emailService: mailer, logger: zap.NewNop()}

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals(constants.ContextUid, "user-1")
		c.Locals(constants.KratosID, "kratos-1")
		c.Locals(constants.ContextUser, models.User{ID: "user-1", Email: "gone@example.com"})
		return c.Next()
	})
	app.Delete("/user", ctrl.DeleteRegisteredUser)

	res, err := app.Test(httptest.NewRequest(http.MethodDelete, "/user", nil))
	if err != nil {
		t.Fatal(err)
	}
	return res.StatusCode, mailer.sent
}

// GitHub #46: a deleted account gets one confirmation email at its address.
func TestDeleteRegisteredUserSendsConfirmationEmail(t *testing.T) {
	status, sent := deleteAccount(t, nil)
	if status != http.StatusOK {
		t.Fatalf("status %d, want 200", status)
	}
	select {
	case m := <-sent:
		if m.to != "gone@example.com" || m.subject != constants.AccountDeletedEmailSubject {
			t.Fatalf("sent %+v", m)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("no confirmation email after a successful deletion")
	}
	select {
	case m := <-sent:
		t.Fatalf("second email sent: %+v", m)
	case <-time.After(100 * time.Millisecond):
	}
}

func TestDeleteRegisteredUserSendsNoEmailWhenDeletionFails(t *testing.T) {
	status, sent := deleteAccount(t, errors.New("foreign key"))
	if status != http.StatusInternalServerError {
		t.Fatalf("status %d, want 500", status)
	}
	select {
	case m := <-sent:
		t.Fatalf("email sent for a failed deletion: %+v", m)
	case <-time.After(100 * time.Millisecond):
	}
}
