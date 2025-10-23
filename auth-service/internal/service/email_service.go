package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type EmailService struct {
	apiKey    string
	fromEmail string
	fromName  string
	baseURL   string
}

type ResendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
	Text    string   `json:"text"`
}

func NewEmailService(apiKey, fromEmail, baseURL string) *EmailService {
	return &EmailService{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		fromName:  "Study Platform",
		baseURL:   baseURL,
	}
}

func (s *EmailService) SendVerificationEmail(email, username, otp, token string) error {
	verificationLink := fmt.Sprintf("%s/auth/verify/%s", s.baseURL, token)

	htmlContent := s.getAcademicEmailHTML(username, otp, verificationLink)
	textContent := s.getAcademicEmailText(username, otp, verificationLink)

	reqBody := ResendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail),
		To:      []string{email},
		Subject: "Verify Your Email - Study Platform",
		HTML:    htmlContent,
		Text:    textContent,
	}

	jsonData, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("resend API error: status %d", resp.StatusCode)
	}

	return nil
}

func (s *EmailService) getAcademicEmailHTML(username, otp, link string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"></head>
<body style="margin:0;padding:0;font-family:'Crimson Text','Times New Roman',serif;background-color:#fefefe;">
    <table width="100%%" cellpadding="0" cellspacing="0" style="padding:40px 20px;">
        <tr><td align="center">
            <table width="600" cellpadding="0" cellspacing="0" style="background:#fff;border-radius:12px;border:1px solid #d8d8dc;box-shadow:0 2px 8px rgba(30,58,138,0.08);">
                <tr><td style="background:linear-gradient(135deg,#1e3a8a 0%%,#3b5998 100%%);padding:40px;text-align:center;border-radius:12px 12px 0 0;">
                    <h1 style="color:#fff;margin:0;font-size:32px;font-weight:600;letter-spacing:-0.5px;">Study Platform</h1>
                    <p style="color:rgba(255,255,255,0.9);margin:8px 0 0;font-size:15px;font-style:italic;">Academic Excellence Portal</p>
                </td></tr>
                <tr><td style="padding:40px;">
                    <p style="color:#1a1a2e;font-size:16px;margin:0 0 10px 0;">Dear %s,</p>
                    <p style="color:#4a4a5e;font-size:15px;line-height:1.7;margin:0 0 30px 0;">
                        Welcome to Study Platform! To complete your registration and access our academic resources, please verify your email using the code below:
                    </p>
                    <div style="background:#f5f5f7;border:2px solid #d8d8dc;border-radius:8px;padding:30px;text-align:center;margin:30px 0;">
                        <p style="color:#666;font-size:13px;margin:0 0 12px 0;text-transform:uppercase;letter-spacing:1.5px;font-family:Inter,sans-serif;">Verification Code</p>
                        <div style="font-size:42px;font-weight:bold;color:#1e3a8a;letter-spacing:12px;font-family:'Courier New',monospace;">%s</div>
                        <p style="color:#666;font-size:12px;margin:20px 0 0;font-family:Inter,sans-serif;">⏱️ Expires in <strong>15 minutes</strong></p>
                    </div>
                    <div style="text-align:center;margin:30px 0;position:relative;">
                        <div style="border-top:1px solid #d8d8dc;"></div>
                        <span style="background:#fff;padding:0 15px;position:relative;top:-12px;color:#666;font-size:13px;font-family:Inter,sans-serif;">OR</span>
                    </div>
                    <div style="text-align:center;margin:30px 0;">
                        <a href="%s" style="display:inline-block;background:#1e3a8a;color:#fff;padding:16px 36px;border-radius:8px;text-decoration:none;font-weight:600;font-size:15px;font-family:Inter,sans-serif;box-shadow:0 2px 4px rgba(30,58,138,0.2);">Verify Email Address</a>
                    </div>
                    <div style="background:#fef5ed;border-left:4px solid #9b4819;padding:16px 20px;margin:30px 0;border-radius:4px;">
                        <p style="color:#7a3a15;font-size:13px;margin:0;line-height:1.6;">
                            🔒 <strong>Security Notice:</strong> If you did not create an account, please disregard this email.
                        </p>
                    </div>
                </td></tr>
                <tr><td style="background:#f5f5f7;padding:30px;text-align:center;border-radius:0 0 12px 12px;border-top:1px solid #d8d8dc;">
                    <p style="color:#666;font-size:13px;margin:0 0 12px;font-family:Inter,sans-serif;">
                        Questions? Contact <a href="mailto:support@study.khoipn.id.vn" style="color:#1e3a8a;text-decoration:none;">support@study.khoipn.id.vn</a>
                    </p>
                    <p style="color:#999;font-size:11px;margin:0;font-family:Inter,sans-serif;">
                        © 2024 Study Platform • study.khoipn.id.vn
                    </p>
                </td></tr>
            </table>
        </td></tr>
    </table>
</body>
</html>`, username, otp, link)
}

func (s *EmailService) getAcademicEmailText(username, otp, link string) string {
	return fmt.Sprintf(`
Study Platform - Verify Your Email Address

Dear %s,

Welcome to Study Platform! To complete your registration, please verify your email address.

Your Verification Code: %s
(Expires in 15 minutes)

Or click this link: %s

Security Notice: If you did not create an account, please disregard this email.

Questions? Contact support@study.khoipn.id.vn
© 2024 Study Platform • study.khoipn.id.vn
`, username, otp, link)
}
