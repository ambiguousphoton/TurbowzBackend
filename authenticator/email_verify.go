package authenticator

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/smtp"
	"sync"
	"time"
)

type otpEntry struct {
	Code      string
	ExpiresAt time.Time
}

type EmailVerifier struct {
	mu       sync.Mutex
	store    map[string]otpEntry // email -> OTP
	otpTTL   time.Duration

	// SMTP config
	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string
	FromAddr string
}

func NewEmailVerifier(smtpHost, smtpPort, smtpUser, smtpPass, fromAddr string) *EmailVerifier {
	ev := &EmailVerifier{
		store:    make(map[string]otpEntry),
		otpTTL:   5 * time.Minute,
		SMTPHost: smtpHost,
		SMTPPort: smtpPort,
		SMTPUser: smtpUser,
		SMTPPass: smtpPass,
		FromAddr: fromAddr,
	}
	go ev.cleanupLoop()
	return ev
}

func (ev *EmailVerifier) GenerateAndSend(email string) error {
	code, err := generateOTP(6)
	if err != nil {
		return fmt.Errorf("failed to generate OTP: %w", err)
	}

	ev.mu.Lock()
	ev.store[email] = otpEntry{Code: code, ExpiresAt: time.Now().Add(ev.otpTTL)}
	ev.mu.Unlock()

	// Dev mode: skip SMTP if credentials are not configured
	if ev.SMTPUser == "" || ev.SMTPPass == "" {
		log.Printf("EmailVerifier [DEV MODE]: OTP for %s is %s", email, code)
		return nil
	}

	if err := ev.sendEmail(email, code); err != nil {
		log.Printf("EmailVerifier: failed to send OTP to %s - %v", email, err)
		ev.mu.Lock()
		delete(ev.store, email)
		ev.mu.Unlock()
		return fmt.Errorf("failed to send verification email")
	}

	log.Printf("EmailVerifier: OTP sent to %s", email)
	return nil
}

func (ev *EmailVerifier) Verify(email, code string) bool {
	ev.mu.Lock()
	defer ev.mu.Unlock()

	entry, ok := ev.store[email]
	if !ok {
		return false
	}
	if time.Now().After(entry.ExpiresAt) {
		delete(ev.store, email)
		return false
	}
	if entry.Code != code {
		return false
	}
	delete(ev.store, email)
	return true
}

func (ev *EmailVerifier) sendEmail(to, code string) error {
	subject := "Turbowz - Email Verification Code"
	body := fmt.Sprintf("Your verification code is: %s\n\nThis code expires in 5 minutes.", code)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", ev.FromAddr, to, subject, body)

	auth := smtp.PlainAuth("", ev.SMTPUser, ev.SMTPPass, ev.SMTPHost)
	addr := ev.SMTPHost + ":" + ev.SMTPPort
	return smtp.SendMail(addr, auth, ev.FromAddr, []string{to}, []byte(msg))
}

func (ev *EmailVerifier) cleanupLoop() {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		ev.mu.Lock()
		now := time.Now()
		for email, entry := range ev.store {
			if now.After(entry.ExpiresAt) {
				delete(ev.store, email)
			}
		}
		ev.mu.Unlock()
	}
}

func generateOTP(length int) (string, error) {
	otp := ""
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		otp += fmt.Sprintf("%d", n.Int64())
	}
	return otp, nil
}
