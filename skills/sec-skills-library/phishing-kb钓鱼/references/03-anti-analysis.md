# Anti-Analysis Techniques Reference

## Cookie Check

### Client-Side: navigator.cookieEnabled
```javascript
if (navigator.cookieEnabled) {
    console.log('Cookies enabled');
} else {
    console.log('Cookies disabled');
}
```

### Client-Side: Set/Read Verification
More reliable than navigator property (which can be spoofed):
```javascript
function checkCookies() {
    const cookieVal = 'test_value';
    document.cookie = `test_cookie=${cookieVal}; path=/;`;
    const cookies = {};
    document.cookie.split('; ').forEach(cookie => {
        const [name, value] = cookie.split('=');
        cookies[name] = value;
    });
    if (cookies.test_cookie === cookieVal) {
        console.log("Cookies enabled");
    } else {
        console.log("Cookies disabled");
    }
}
checkCookies();
```

### Server-Side: PHP Cookie Persistence
**set-cookie.php** (landing page):
```php
<?php
setcookie('session_check', 'valid', time() + 3600, '/', "", false, false);
header('Location: login.php');
exit();
?>
```

**login.php** (check persistence):
```php
<?php
if (!isset($_COOKIE['session_check'])) {
    header('Location: https://google.com');
    exit();
}
?>
```

---

## Headless Browser / WebDriver Detection

```javascript
if (navigator.webdriver) {
    document.write("Automated browser detected");
} else {
}
```

---

## User-Agent Filtering

### PHP: Keyword Blacklist
```php
<?php
$blockedKeywords = array('curl', 'headless', 'python', 'wget', 'bot', 'spider', 'scanner');
$userAgent = strtolower($_SERVER['HTTP_USER_AGENT']);
foreach ($blockedKeywords as $keyword) {
    if (strpos($userAgent, $keyword) !== false) {
        header('Location: blocked.php');
        exit;
    }
}
header('Location: login.php');
?>
```

### PHP: Crawler Database Blacklist
Download `crawler-user-agents.json` from GitHub and place outside docroot:
```php
<?php
$botList = json_decode(file_get_contents('/var/www/crawler-user-agents.json'), true);
$userAgent = $_SERVER['HTTP_USER_AGENT'];
function isBot($ua, $botList) {
    foreach ($botList as $bot) {
        if (preg_match('/' . $bot['pattern'] . '/i', $ua)) return true;
    }
    return false;
}
if (isBot($userAgent, $botList)) {
    header('Location: blocked.php'); exit;
}
header('Location: login.php');
?>
```

### PHP: Whitelist Browsers Only
```php
<?php
$whitelistedKeywords = array('Chrome', 'Firefox', 'Edg', 'Version');
$userAgent = $_SERVER['HTTP_USER_AGENT'];
$allowed = false;
foreach ($whitelistedKeywords as $kw) {
    if (strpos($userAgent, $kw) !== false) { $allowed = true; break; }
}
if (!$allowed) { header('Location: blocked.php'); exit; }
header('Location: login.php');
?>
```

### PHP: Block Old Browser Versions
```php
<?php
$ua = $_SERVER['HTTP_USER_AGENT'];
function getBrowserVersion($ua, $browser) {
    $patterns = [
        'Chrome' => '/Chrome\/(\d+)/', 'Firefox' => '/Firefox\/(\d+)/',
        'Edge' => '/Edg\/(\d+)/', 'Safari' => '/Version\/(\d+)/'
    ];
    if (isset($patterns[$browser]) && preg_match($patterns[$browser], $ua, $m)) return (int)$m[1];
    return null;
}
$cv = getBrowserVersion($ua, 'Chrome');
$fv = getBrowserVersion($ua, 'Firefox');
$ev = getBrowserVersion($ua, 'Edge');
$sv = getBrowserVersion($ua, 'Safari');
if (($cv !== null && $cv >= 100) || ($fv !== null && $fv >= 100) ||
    ($ev !== null && $ev >= 100) || ($sv !== null && $sv >= 16)) {
    header('Location: login.php');
} else {
    header('Location: blocked.php');
}
?>
```

### PHP: Auto-Blacklist IPs from Bad UAs
```php
<?php
$ipBlacklistFile = '/var/www/ip_blacklist.txt';
function isIpBlacklisted($ip, $file) {
    if (file_exists($file)) {
        $ips = file($file, FILE_IGNORE_NEW_LINES | FILE_SKIP_EMPTY_LINES);
        return in_array($ip, $ips);
    }
    return false;
}
function blacklistIp($ip, $file) {
    file_put_contents($file, $ip . PHP_EOL, FILE_APPEND | LOCK_EX);
}
$clientIp = $_SERVER['REMOTE_ADDR'];
if (isIpBlacklisted($clientIp, $ipBlacklistFile)) {
    header('Location: blocked.php'); exit;
}
$blockedKeywords = array('curl', 'headless', 'python');
$ua = strtolower($_SERVER['HTTP_USER_AGENT']);
foreach ($blockedKeywords as $kw) {
    if (strpos($ua, $kw) !== false) {
        blacklistIp($clientIp, $ipBlacklistFile);
        header('Location: blocked.php'); exit;
    }
}
header('Location: login.php');
?>
```

### Apache: Block Missing User-Agent
```apache
RewriteEngine On
RewriteCond %{HTTP_USER_AGENT} ^$
RewriteRule ^ /blocked.php [L,R=302]
```

---

## UA Spoofing Detection

### JS + PHP: Navigator vs HTTP Header Validation
**Client JS:**
```javascript
const data = { userAgent: navigator.userAgent, platform: navigator.platform };
fetch('validate-ua.php', {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body: new URLSearchParams(data).toString()
}).then(r => r.json()).then(d => {
    if (!d.match) console.log('UA spoofing detected');
});
```

**Server PHP:**
```php
<?php
if ($_SERVER['REQUEST_METHOD'] === 'POST') {
    $clientUA = $_POST['userAgent'];
    $clientPlatform = $_POST['platform'];
    $serverUA = $_SERVER['HTTP_USER_AGENT'];
    $uaMatch = stripos($serverUA, $clientUA) !== false;
    $platformMatch = false;
    if (stripos($serverUA, 'windows') !== false)
        $platformMatch = in_array($clientPlatform, ['Windows', 'Win32', 'Win64']);
    elseif (stripos($serverUA, 'mac') !== false)
        $platformMatch = in_array($clientPlatform, ['Macintosh', 'MacIntel', 'MacOS']);
    elseif (stripos($serverUA, 'linux') !== false)
        $platformMatch = in_array($clientPlatform, ['Linux x86_64', 'Linux armv81']);
    echo json_encode(['match' => $uaMatch && $platformMatch]);
}
?>
```

### Browser-Specific Navigator Properties
```javascript
if ('oscpu' in navigator) console.log("Firefox confirmed");
if ('protectedAudience' in navigator) console.log("Chromium confirmed");
if (CSS.supports('-webkit-line-grid','none')) console.log("Safari confirmed");
```

### Mobile UA + Screen Size Mismatch
```javascript
fetch('validate-screen.php', {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body: new URLSearchParams({ screenWidth: screen.width }).toString()
}).then(r => r.json()).then(d => {
    if (!d.match) console.log('Mobile UA spoofing detected');
});
```

### Mouse Movement on Mobile
```javascript
if (/Android|iPhone|iPad/i.test(navigator.userAgent)) {
    window.addEventListener('mousemove', () => console.log('Mouse on mobile = spoofed UA'));
    window.addEventListener('touchstart', () => console.log('Touch = legitimate'));
}
```

### Font-Based OS Detection
```javascript
function detectOSSpoofing() {
    const ua = navigator.userAgent;
    const fontsByOS = {
        'Windows': ['Calibri', 'Segoe UI'],
        'MacOS': ['San Francisco', 'Helvetica Neue'],
        'Linux': ['DejaVu Sans', 'Liberation Sans']
    };
    let uaOS = ua.includes('Windows') ? 'Windows' : ua.includes('Mac') ? 'MacOS' : ua.includes('Linux') ? 'Linux' : 'Unknown';
    let fontOS = 'Unknown';
    const canvas = document.createElement('canvas');
    const ctx = canvas.getContext('2d');
    ctx.font = '72px monospace';
    const baseline = ctx.measureText('A').width;
    for (const [os, fonts] of Object.entries(fontsByOS)) {
        for (const font of fonts) {
            ctx.font = `72px ${font}, monospace`;
            if (ctx.measureText('A').width !== baseline) { fontOS = os; break; }
        }
        if (fontOS !== 'Unknown') break;
    }
    if (fontOS !== 'Unknown' && fontOS !== uaOS) {
        console.log(`OS spoofing: UA=${uaOS}, Font=${fontOS}`);
    }
}
```

### Math Precision Fingerprint
```javascript
Math.hypot(-24.42, -50.519999999999930000)
```

---

## IP Restrictions

### PHP: Country Whitelist
```php
<?php
$whitelistedCountries = ['US', 'DE'];
$ip = $_SERVER['REMOTE_ADDR'];
$response = file_get_contents("https://api.country.is/" . $ip);
$data = json_decode($response, true);
if (!isset($data['country']) || !in_array($data['country'], $whitelistedCountries)) {
    header('Location: /blocked.php'); exit;
}
header('Location: /login.php');
?>
```

### PHP: VPN/Proxy Detection (IPInfo)
```php
<?php
$apiToken = 'YOUR_TOKEN';
$ip = $_SERVER['REMOTE_ADDR'];
$response = file_get_contents("https://ipinfo.io/{$ip}?token={$apiToken}");
$data = json_decode($response, true);
if (isset($data['privacy'])) {
    $p = $data['privacy'];
    if ($p['hosting'] || $p['proxy'] || $p['tor'] || $p['relay'] || $p['vpn']) {
        header('Location: /blocked.php'); exit;
    }
}
header('Location: /login.php');
?>
```

### PHP: Block Specific VPN Providers
```php
<?php
$blockedVPNs = ['NordVPN', 'Mullvad'];
$apiToken = 'YOUR_TOKEN';
$ip = $_SERVER['REMOTE_ADDR'];
$data = json_decode(file_get_contents("https://ipinfo.io/{$ip}?token={$apiToken}"), true);
if (isset($data['privacy']) && $data['privacy']['vpn'] && in_array($data['privacy']['service'], $blockedVPNs)) {
    header('Location: /blocked.php'); exit;
}
header('Location: /login.php');
?>
```

### Apache: Rate Limiting (mod_evasive)
```bash
sudo apt install libapache2-mod-evasive
```
Config (`/etc/apache2/mods-available/evasive.conf`):
```apache
<IfModule mod_evasive20.c>
    DOSHashTableSize    3097
    DOSPageCount        5
    DOSSiteCount        20
    DOSPageInterval     1
    DOSSiteInterval     5
    DOSBlockingPeriod   3600
</IfModule>
```

---

## Window Size Detection

```javascript
if (window.innerWidth < 300 || window.innerHeight < 300) {
    console.log("Suspicious: tiny window");
}
if (window.outerWidth === 0 || window.outerHeight === 0) {
    console.log("Headless browser: zero outer dimensions");
}
if (window.innerWidth === screen.width && window.innerHeight === screen.height) {
    console.log("Fullscreen detected");
}
```

---

## Incognito Mode Detection

Uses storage API quota differences between normal and private browsing:
```javascript
async function detectIncognito() {
    return new Promise((resolve) => {
        if ('storage' in navigator && 'estimate' in navigator.storage) {
            navigator.storage.estimate().then(({quota}) => {
                resolve(quota < 120000000);
            });
        } else {
            resolve(false);
        }
    });
}
```

---

## Ad Blocker Detection

```javascript
function detectAdBlocker() {
    var testAd = document.createElement('div');
    testAd.innerHTML = '&nbsp;';
    testAd.className = 'adsbox ad-banner';
    testAd.style.cssText = 'position:absolute;left:-9999px;';
    document.body.appendChild(testAd);
    setTimeout(function() {
        if (testAd.offsetHeight === 0 || testAd.clientHeight === 0) {
            console.log('Ad blocker detected');
        }
        testAd.remove();
    }, 100);
}
```
