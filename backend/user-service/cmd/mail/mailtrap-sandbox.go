package test

import (
	"fmt"
	"net/smtp"
)

func main(){
	host := "sandbox.smtp.mailtrap.io"
	port := 2525
	user := "49bbd2a554bd3e"
	password := "4f4bc1b03a1d70"
	from := "noreply@drive-master.com"
	to := "test@example.com"

	auth := smtp.PlainAuth("", user, password, host)

	subject := "Subject: Test Email from Drive Master\r\n"
	mime    := "MIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n"
	body    := `
		<h1>Hello from Drive Master!</h1>
		<p>This is a test email from the user-service.</p>
	`
	msg := []byte(subject + mime + "\r\n" + body)

	var address = fmt.Sprintf("%s:%d", host, port)
	err := smtp.SendMail(address, auth, from, []string{to}, msg)
	if err != nil {
		fmt.Println("Failed to send email:", err)
		return
	}

	fmt.Println("Email sent successfully! Check your Mailtrap inbox.")

}