package templates

func GetForgotPasswordTemplate(email string) string {
	return `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Reset Your Password</title>
	</head>
	<body>
		<h1>Reset Your Password</h1>
		<p>Hello {{.Email}},</p>
		<p>We have received a request to reset your password. Please click the button below to reset your password.</p>
		<a href="https://hackfy.com/reset-password?token={{.Token}}">Reset Password</a>
		<p>If you did not request a password reset, please ignore this email.</p>
	</body>
	</html>
	`
}
