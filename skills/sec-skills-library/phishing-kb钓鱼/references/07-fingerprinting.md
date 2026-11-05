# TLS Fingerprinting & Client Telemetry Reference

## JARM Fingerprinting
JARM generates a hash from a server's TLS handshake responses to 10 crafted client hello packets. Different server software produces different JARM hashes, enabling identification of web server technology.

**Mitigation**: Modify Apache SSL configuration (cipher order, supported versions) to change the JARM hash. Place a reverse proxy (Caddy/Nginx) in front to mask the real server.

## JA3/JA3S
JA3 fingerprints TLS clients by hashing: TLS version, cipher suites, extensions, elliptic curves, EC point formats. JA3S fingerprints the server response similarly.

**Implication**: Security scanners have known JA3 fingerprints. Blocking known scanner JA3s prevents automated analysis.

## JA4 Fingerprinting
JA4 is the evolution of JA3 with better collision resistance. Format: `{protocol}{version}{SNI}{cipher_count}{ext_count}{ALPN}_{cipher_hash}_{ext+sig_hash}`

### HAProxy + Lua JA4 Plugin
The Lua plugin calculates JA4 fingerprints from TLS handshake data and passes them as HTTP headers to the backend.

**haproxy.cfg:**
```
global
    lua-load /etc/haproxy/lua/ja4.lua

frontend https
    bind *:443 ssl crt /etc/haproxy/ssl/cert.pem
    http-request lua.fingerprint_ja4
    http-request set-header X-JA4 %[var(txn.fingerprint_ja4)]
    default_backend web

backend web
    server backend1 127.0.0.1:8080
```

### PHP: Extract JA4 Header
```php
<?php
$ja4 = isset($_SERVER['HTTP_X_JA4']) ? $_SERVER['HTTP_X_JA4'] : 'unknown';
echo "JA4: " . htmlspecialchars($ja4);
?>
```

### JA4 Blacklisting
Maintain a map file of known bot JA4 fingerprints:
```
# ja4_bots.map
t13d1517h2_8daaf6152771_b0da82dd1658 curl
t13d1715h2_5b57614c22b0_bde62ecc39b4 python-requests
```

**HAProxy map lookup:**
```
http-request set-var(txn.bot_id) var(txn.fingerprint_ja4),map(/etc/haproxy/ja4_bots.map)
http-request set-header X-Bot-Id %[var(txn.bot_id)]
```

**PHP: Check for bot:**
```php
<?php
$botId = isset($_SERVER['HTTP_X_BOT_ID']) ? $_SERVER['HTTP_X_BOT_ID'] : '';
if (!empty($botId)) {
    header('Location: /blocked.php'); exit;
}
?>
```

### JA4 Partial Blacklisting
Split JA4 into parts (ja4_a, ja4_b, ja4_c) and check each independently:
```php
<?php
$ja4 = $_SERVER['HTTP_X_JA4'] ?? '';
$parts = explode('_', $ja4);
$ja4_a = $parts[0] ?? '';
$ja4_b = $parts[1] ?? '';
$ja4_c = $parts[2] ?? '';

$blacklist = json_decode(file_get_contents('/var/www/blacklist_ja4.json'), true);
foreach ($blacklist as $entry) {
    if ((!empty($entry['ja4_a']) && $entry['ja4_a'] === $ja4_a) ||
        (!empty($entry['ja4_b']) && $entry['ja4_b'] === $ja4_b) ||
        (!empty($entry['ja4_c']) && $entry['ja4_c'] === $ja4_c)) {
        header('Location: /blocked.php'); exit;
    }
}
?>
```

---

## Client Telemetry Logging Library

Comprehensive JavaScript library that collects 30+ browser properties for visitor profiling and bot detection.

```javascript
class ClientLogger {
    constructor(endpoint) {
        this.endpoint = endpoint;
        this.data = {};
    }

    collect() {
        var n = navigator;
        this.data = {
            userAgent: n.userAgent,
            platform: n.platform,
            language: n.language,
            languages: n.languages ? n.languages.join(',') : '',
            cookieEnabled: n.cookieEnabled,
            doNotTrack: n.doNotTrack,
            hardwareConcurrency: n.hardwareConcurrency || 0,
            maxTouchPoints: n.maxTouchPoints || 0,
            webdriver: n.webdriver,
            vendor: n.vendor,

            screenWidth: screen.width,
            screenHeight: screen.height,
            screenColorDepth: screen.colorDepth,
            screenPixelDepth: screen.pixelDepth,

            innerWidth: window.innerWidth,
            innerHeight: window.innerHeight,
            outerWidth: window.outerWidth,
            outerHeight: window.outerHeight,
            devicePixelRatio: window.devicePixelRatio,

            timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
            timezoneOffset: new Date().getTimezoneOffset(),

            webglVendor: this.getWebGLVendor(),
            webglRenderer: this.getWebGLRenderer(),

            connectionType: (n.connection || {}).effectiveType || 'unknown',
            connectionRTT: (n.connection || {}).rtt || 0,

            localStorageAvailable: this.testStorage('localStorage'),
            sessionStorageAvailable: this.testStorage('sessionStorage'),
            indexedDBAvailable: !!window.indexedDB,

            serviceWorkerAvailable: 'serviceWorker' in n,
            notificationPermission: typeof Notification !== 'undefined' ? Notification.permission : 'unavailable',
        };
    }

    getWebGLVendor() {
        try {
            var c = document.createElement('canvas');
            var gl = c.getContext('webgl') || c.getContext('experimental-webgl');
            var ext = gl.getExtension('WEBGL_debug_renderer_info');
            return ext ? gl.getParameter(ext.UNMASKED_VENDOR_WEBGL) : 'N/A';
        } catch(e) { return 'N/A'; }
    }

    getWebGLRenderer() {
        try {
            var c = document.createElement('canvas');
            var gl = c.getContext('webgl') || c.getContext('experimental-webgl');
            var ext = gl.getExtension('WEBGL_debug_renderer_info');
            return ext ? gl.getParameter(ext.UNMASKED_RENDERER_WEBGL) : 'N/A';
        } catch(e) { return 'N/A'; }
    }

    testStorage(type) {
        try { var s = window[type]; s.setItem('t', 't'); s.removeItem('t'); return true; }
        catch(e) { return false; }
    }

    send() {
        this.collect();
        fetch(this.endpoint, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(this.data)
        });
    }
}

var logger = new ClientLogger('/analytics.php');
logger.send();
```

### Server-Side Analytics Handler
```php
<?php
$data = json_decode(file_get_contents('php://input'), true);
if ($data) {
    $data['server_ip'] = $_SERVER['REMOTE_ADDR'];
    $data['server_ua'] = $_SERVER['HTTP_USER_AGENT'];
    $data['timestamp'] = date('Y-m-d H:i:s');
    $logFile = '/var/www/telemetry.json';
    file_put_contents($logFile, json_encode($data) . "\n", FILE_APPEND);
}
?>
```
