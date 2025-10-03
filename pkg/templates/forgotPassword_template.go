package templates

func GetForgotPasswordTemplate(email, otp string) string {
	return `
	
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>HACFY - Password Reset OTP</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    /* Reset & base */
    body, table, td, div, p, a { -webkit-text-size-adjust:100%; -ms-text-size-adjust:100%; }
    body {
      font-family: Arial, sans-serif;
      background-color: #1E275A;
      margin: 0;
      padding: 0;
    }

    .wrapper {
      width: 100%;
      table-layout: fixed;
      background-color: #1E275A;
      padding-bottom: 40px;
    }
    .logo{
       justify-items: center;
    }

    .outer {
      margin: 0 auto;
      width: 100%;
      max-width: 600px;
      background: #ffffff;
      border-radius: 40px 40px 0 0;
      overflow: hidden;
    }

    .header {
      text-align: center;
      padding: 30px 20px 20px;
      background: #1E275A;
    }

    .header img.logo {
      height: 70px;
      margin-bottom: 15px;
      border-radius: 8px;
    }

    .header img.banner {
      display: block;
      margin: 0 auto 15px;
      max-width: 320px;
      width: 100%;
      border-radius: 12px;
    }

    .header h1 {
      color: #ffffff;
      font-size: 24px;
      margin: 10px 0 0;
    }

    .main {
      padding: 30px 20px;
      text-align: center;
    }

    .main h2 {
      color: #1E275A;
      font-size: 20px;
      margin-bottom: 10px;
    }

    .main p {
      color: #666;
      font-size: 15px;
      margin-bottom: 20px;
      line-height: 1.5;
    }

    .form-box {
      background: #2A3B6B;
      padding: 20px;
      border-radius: 16px;
      margin: 0 auto;
      max-width: 400px;
      color: #fff;
      text-align: left;
    }

    .info-label {
      font-size: 13px;
      color: #eee;
      margin-bottom: 6px;
    }

    .otp-code {
      display: block;
      width: 100%;
      padding: 14px;
      border-radius: 10px;
      background: #ffffff;
      color: #1E275A;
      font-size: 22px;
      font-weight: 700;
      letter-spacing: 4px;
      text-align: center;
      box-sizing: border-box;
      margin-bottom: 15px;
    }

    .login-btn {
      display: block;
      text-align: center;
      width: 100%;
      padding: 12px;
      background: #ffffff;
      color: #1E275A !important;
      font-size: 16px;
      font-weight: bold;
      border-radius: 8px;
      text-decoration: none;
      margin-top: 10px;
    }

    .info-text {
      margin-top: 15px;
      font-size: 12px;
      color: #ddd;
      line-height: 1.4;
    }

    .footer {
      background: #f5f5f5;
      text-align: center;
      padding: 15px;
      font-size: 12px;
      color: #666;
    }

    /* Mobile tweaks */
    @media only screen and (max-width: 480px) {
      .header h1 { font-size: 18px !important; }
      .main h2 { font-size: 18px !important; }
      .main p { font-size: 14px !important; }
      .otp-code { font-size: 18px !important; letter-spacing: 3px !important; padding: 12px !important; }
      .login-btn { font-size: 14px !important; padding: 10px !important; }
    }
  </style>
</head>
<body>
  <center class="wrapper">
    <div class="outer">
      <!-- Header -->
      <div class="header">
        <img src="pkg/templates/image.png" alt="HACFY Logo" class="logo">
        <img src="pkg/templates/image copy.png" alt="HACFY Banner" class="banner">
        <h1>Welcome to IT Inventory!</h1>
      </div>

      <!-- Main Content -->
      <div class="main">
        <h2>Password Reset OTP</h2>
        <p>Please use the one-time code below to reset your password. This code is valid for a short time only.</p>

        <div class="form-box">
          <div>
            <div class="info-label">Email</div>
            <div class="otp-code">` + email + `</div>
          </div>

          <div>
            <div class="info-label">Your 6-digit OTP</div>
            <div class="otp-code">` + otp + `</div>
          </div>

          <a href="https://your-login-url.com" class="login-btn">Go to Login</a>

          <p class="info-text">
            This OTP is valid for 10 minutes. If you did not request this, please ignore this email or contact support.
          </p>
        </div>
      </div>

      <!-- Footer -->
      <div class="footer">
        © 2025 IT Management System. All Rights Reserved
      </div>
    </div>
  </center>
</body>
</html>


	`
}
