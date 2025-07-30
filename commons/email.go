package commons

import (
	"encoding/hex"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"unicode"

	"github.com/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type EmailHeaders struct {
	From    string
	To      string
	Subject string
	Headers map[string]string
}

func (e *EmailHeaders) Build() string {
	lines := []string{
		fmt.Sprintf("From: %s", e.From),
		fmt.Sprintf("To: %s", e.To),
		fmt.Sprintf("Subject: %s", e.Subject),
		`Content-Type: text/html; charset="utf-8"`,
	}

	for k, v := range e.Headers {
		lines = append(lines, fmt.Sprintf("%s: %s", k, v))
	}

	return strings.Join(lines, "\r\n") + "\r\n\r\n"
}

func DecodeCloudflareEmail(encoded string) (string, error) {
	data, err := hex.DecodeString(encoded)
	if err != nil || len(data) < 1 {
		return "", err
	}

	key := data[0]
	for i := 1; i < len(data); i++ {
		data[i] ^= key
	}

	return string(data[1:]), nil
}

func EmailToDisplayName(email string) string {
	localPart := strings.Split(email, "@")[0]
	nameParts := strings.FieldsFunc(localPart, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	title := cases.Title(language.Und)
	for i, part := range nameParts {
		nameParts[i] = title.String(strings.ToLower(part))
	}

	return strings.Join(nameParts, " ")
}

func SendLocalHTMLEmail(emailHeader *EmailHeaders, htmlBody string) error {
	var (
		msg = []byte(emailHeader.Build())
	)

	conn, err := net.Dial("tcp", "localhost:25")
	if err != nil {
		return err
	}
	defer conn.Close()

	if client, err := smtp.NewClient(conn, "localhost"); err != nil {
		return errors.Wrapf(err, "NewClient")
	} else if err := client.Mail(emailHeader.From); err != nil {
		return errors.Wrapf(err, "set from")
	} else if err := client.Rcpt(emailHeader.To); err != nil {
		return errors.Wrapf(err, "set to")
	} else if writer, err := client.Data(); err != nil {
		return err
	} else {
		defer writer.Close()

		if _, err := writer.Write(msg); err != nil {
			return err
		}

		client.Quit()
	}

	return nil
}
