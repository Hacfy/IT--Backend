package templates

func GetVerifyEmailOtpTemplate(password, email string) string {
	return `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>HACFY - IT Inventory Login</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.0/css/all.min.css">
  <style>
    body {
      margin: 0;
      font-family: Arial, sans-serif;
      background-color: #1E275A;
      color: #000;
    }

    .header {
      text-align: center;
      padding: 40px 20px 20px;
      background-color: #1E275A;
    }

    .header img.logo {
      height: 80px;
      margin-bottom: 10px;
    }

    .header img.banner {
      display: block;
      margin: 0 auto;
      background-color: #ffffff;
      border-radius: 15px;
      height: 60px;
      width: 330px;
      max-width: 90%;
    }

    .header h1 {
      color: white;
      font-size: 24px;
      margin-top: 10px;
    }

    .main {
      background-color: #ffffff;
      border-top-left-radius: 40px;
      border-top-right-radius: 40px;
      padding: 40px 20px;
      text-align: center;
    }

    .main h2 {
      font-size: 28px;
      font-weight: bold;
      margin-top: 20px;
    }

    .main p {
      font-size: 16px;
      color: #333;
    }

    .form-box {
      background-color: #1E275A;
      padding: 30px 50px;
      border-radius: 16px;
      margin-top: 150px;
      max-width: 400px;
      margin-left: auto;
      margin-right: auto;
      box-sizing: border-box;
    }

    .input-group {
      position: relative;
      width: 100%;
      margin-bottom: 20px;
    }

    .input-group i {
      position: absolute;
      top: 50%;
      left: 15px;
      transform: translateY(-50%);
      color: gray;
    }

    .input-group input {
      width: 100%;
      padding: 15px 15px 15px 45px;
      border-radius: 20px;
      border: none;
      font-size: 16px;
      box-sizing: border-box;
    }

    .login-btn {
      background-color: #ffffff;
      color: #1E275A;
      border: none;
      padding: 12px 24px;
      font-size: 16px;
      border-radius: 25px;
      cursor: pointer;
      margin-top: 10px;
      /* width: ; */
      font-weight: bold;
    }

    .login-btn:hover {
      background-color: #ddd;
    }

    .form-box p {
      color: white;
      font-size: 11px;
      margin-top: 20px;
    }

    .footer {
      background-color: #ddd;
      padding: 20px;
      text-align: center;
      border-bottom-left-radius: 30px;
      border-bottom-right-radius: 30px;
    }

    .footer p {
      margin: 0;
      font-size: 14px;
    }

    @media (max-width: 480px) {
      .form-box {
        padding: 20px;
        margin-top: 100px;
      }

      .input-group input {
        font-size: 14px;
        padding-left: 40px;
      }

      .main h2 {
        font-size: 22px;
      }

      .main p {
        font-size: 14px;
      }

      .header h1 {
        font-size: 20px;
      }
    }
  </style>
</head>
<body>

  <div class="header">
    <img src="/home/ashith/Hacfy/IT_INVENTORY/pkg/templates/hacfy png.png" class="logo" alt="Logo">
    <img src="/home/ashith/Hacfy/IT_INVENTORY/pkg/templates/HACFY-C3Z81d0c.webp" class="banner" alt="Logo">
    <h1>Welcome to IT Inventory!</h1>
  </div>

  <div class="main">
    <h2>Login Credentials</h2>
    <p>Please login through the below mentioned email and password</p>

    <div class="form-box">
      <div class="input-group">
        <i class="fas fa-envelope"></i>
        <input type="email" id="email" placeholder="Email" value="` + email + `"/>
      </div>

      <div class="input-group">
        <i class="fas fa-lock"></i>
        <input type="password" id="password" placeholder="Password" value="` + password + `"/>
      </div>

      <button class="login-btn" id="loginBtn">Login</button>

      <p>This is a temporary password. Please update your password</p>
    </div>
  </div>

  <div class="footer">
    <p>2025 IT Management System. All Rights Reserved</p>
  </div>

  <script>
    document.getElementById("loginBtn").addEventListener("click", function () {
      document.getElementById("email").value = ` + email + `;
      document.getElementById("password").value = ` + password + `;
    });
  </script>

</body>
</html>
`
}
