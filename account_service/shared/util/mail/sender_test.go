package mail

import (
	"testing"

	"github.com/RoyceAzure/sexy_gpt/account_service/shared/util/config"
	"github.com/stretchr/testify/require"

)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	config, err := config.LoadConfig("../../../../")
	require.NoError(t, err)
	require.NotEmpty(t, config)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	require.NotEmpty(t, sender)

	subject := "test email"
	var contnt string = `
	<h1>Hello world</h1>
	<p>This is a test message from stockinfo</p>
	`
	to := []string{"roycewnag@gmail.com"}
	attachfiles := []string{"../../../test.txt", "../../../環境報告.pdf", "../../../test3.jfif"}
	err = sender.SendEmail(subject, contnt, to, nil, nil, attachfiles)
	require.NoError(t, err)
}
