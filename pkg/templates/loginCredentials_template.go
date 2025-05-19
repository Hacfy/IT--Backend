package templates

func GetVerifyEmailOtpTemplate(password, email string) string {
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login Credintials - IT Management System</title>
</head>
<body style="margin: 0; padding: 0; font-family: 'Poppins', Arial, sans-serif; background-color: #f4f4f4;">
    <table role="presentation" style="width: 100%; border-collapse: collapse;">
        <tr>
            <td style="padding: 0;">
                <table role="presentation" style="width: 100%; max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 8px; margin-top: 20px; margin-bottom: 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
                    <!-- Header -->
                    <tr>
                        <td style="padding: 20px 0; text-align: center; background-color: #4169E1; border-radius: 8px 8px 0 0;">
                            <h1 style="color: #ffffff; font-size: 24px; margin: 0; font-weight: 600; font-family: 'Poppins', Arial, sans-serif;">IT Management System</h1>
                        </td>
                    </tr>
                    <!-- Content -->
                    <tr>
                        <td style="padding: 40px 30px;">
                            <h2 style="color: #333333; font-size: 22px; margin: 0 0 20px 0; font-family: 'Poppins', Arial, sans-serif;">Login Credintials</h2>
                            <p style="color: #666666; font-size: 16px; line-height: 24px; margin: 0 0 20px 0; font-family: 'Poppins', Arial, sans-serif;">
                                Please Login through the below mentioned email and password:
                            </p>
                            <!-- PASSWORD Code Box -->
                            <div style="background-color: #f8f9fa; border: 2px solid #4169E1; border-radius: 6px; padding: 20px; text-align: center; margin: 30px 0;">
                                <span style="font-family: 'Courier New', monospace; font-size: 32px; font-weight: bold; color: #4169E1; letter-spacing: 5px;">
                                    ` + email + `
                                </span>
                            </div>
                            <div style="background-color: #f8f9fa; border: 2px solid #4169E1; border-radius: 6px; padding: 20px; text-align: center; margin: 30px 0;">
                                <span style="font-family: 'Courier New', monospace; font-size: 32px; font-weight: bold; color: #4169E1; letter-spacing: 5px;">
                                    ` + password + `
                                </span>
                            </div>
                            <p style="color: #d84c4c; font-size: 16px; line-height: 24px; margin: 0 0 20px 0; font-family: 'Poppins', Arial, sans-serif; font-weight: bold; text-decoration: underline;">
                                This is a temporary password, please update your password.
                            </p>                            
                        </td>
                    </tr>
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 30px; background-color: #f8f9fa; border-radius: 0 0 8px 8px; text-align: center;">
                            <p style="color: #999999; font-size: 14px; margin: 0; font-family: 'Poppins', Arial, sans-serif;">
                                © 2025 IT Management System. All Rights Reserved.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>`
}
