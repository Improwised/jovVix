package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSMTPHelper(t *testing.T) {
	helper := NewSMTPHelper("localhost", "879", "xyz@example.com", "12345")

	assert.Equal(t, helper.Host, "localhost")
	assert.Equal(t, helper.Port, "879")
	assert.Equal(t, helper.Username, "xyz@example.com")
	assert.Equal(t, helper.Password, "12345")

	helper.SetSubject("This is subject")
	assert.Equal(t, helper.MailDetails.Subject, "Subject: This is subject\n\n")

	helper.SetReceivers([]string{
		"abc@example.com",
	})
	assert.Equal(t, helper.MailDetails.Receivers, []string{"abc@example.com"})

	helper.SetSender("sender@example.com")
	assert.Equal(t, helper.MailDetails.Sender, "sender@example.com")

	helper.SetPlainBody([]byte("this is body"))
	assert.Equal(t, helper.MailDetails.Data, []byte("this is body"))

	helper.SetHTMLBody([]byte("<h1>this is html body</h1>"))
	assert.Equal(t, helper.MailDetails.Data, []byte("<h1>this is html body</h1>"))

	err := helper.SendMail()
	assert.Error(t, err)
}
