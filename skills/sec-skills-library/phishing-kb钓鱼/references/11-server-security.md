# Server Security Hardening Reference

## HTTPS Enforcement

### Redirect HTTP to HTTPS (Apache mod_rewrite)
```apache
RewriteEngine on
RewriteCond %{HTTPS} off
RewriteRule ^ https://%{HTTP_HOST}%{REQUEST_URI} [L,R=301]
```

### Disable HTTP Entirely (UFW)
```bash
sudo ufw deny http
```

### Disable HTTP (Apache)
Comment out `Listen 80` in `/etc/apache2/ports.conf`, then reload.

---

## Block Direct IP Access

### Method 1: Host Header Check (Apache)
```apache
RewriteEngine On
RewriteCond %{HTTP_HOST} !^(www\.)?yourdomain\.com$ [NC]
RewriteCond %{HTTP_HOST} !=""
RewriteRule ^ - [F]
```

### Method 2: Separate Config File
Create `deny-ip.conf`:
```apache
<IfModule mod_ssl.c>
    <VirtualHost *:443>
        ServerName YOUR_SERVER_IP
        Redirect 403 /
        DocumentRoot /dev/null
    </VirtualHost>
</IfModule>
```
Name it so it loads AFTER main config (e.g., `zzz-deny-ip.conf`).

### Nginx
```nginx
server {
    listen 80 default_server;
    return 444; # drop connection
}
```

---

## Remove Verbose Information

### Apache
Add to `/etc/apache2/apache2.conf`:
```apache
ServerTokens Prod
ServerSignature Off
```

### Spoof Server Header (Apache + ModSecurity)
```bash
sudo apt install libapache2-mod-security2
```
```apache
<IfModule security2_module>
    SecRuleEngine on
    ServerTokens Min
    SecServerSignature "LiteSpeed"
</IfModule>
```

### Nginx
In `/etc/nginx/nginx.conf` inside `http {}`:
```nginx
server_tokens off;
```

### SSH Banner
Add to `/etc/ssh/sshd_config`:
```
DebianBanner no
```

---

## SSH Hardening

### Key-Based Authentication
```bash
ssh-keygen -t rsa -b 4096

ssh-copy-id -i ~/.ssh/id_rsa.pub root@server-ip

ssh -i ~/.ssh/id_rsa root@server-ip
```

### Disable Password Auth
In `/etc/ssh/sshd_config`:
```
PasswordAuthentication no
```
Check override files in `/etc/ssh/sshd_config.d/*.conf` too.

### Jump Server Setup
Restrict SSH to jump server only:
```bash
sudo ufw allow from JUMP_SERVER_IP to any port 22
sudo ufw deny 22
sudo ufw reload
```

---

## Apache mod_rewrite Redirector

### Concept
A redirector sits between users and the phishing server. It filters traffic (block bots, allow targets) and reverse-proxies legitimate requests to the backend.

### Basic UA Filtering
```apache
RewriteEngine On
RewriteCond %{HTTP_USER_AGENT} curl [NC]
RewriteRule ^.*$ http://google.com [L,R=302]
```

### IP Blacklist/Whitelist
```apache
RewriteCond %{REMOTE_ADDR} ^1\.2\.3\.4$
RewriteRule ^.*$ http://google.com/ [L,R=302]

RewriteCond %{REMOTE_ADDR} !^1\.2\.3\.4$
RewriteRule ^.*$ http://google.com/ [L,R=302]
```

### Time-Based Filtering
```apache
RewriteCond %{TIME_HOUR}%{TIME_MIN} <0800 [OR]
RewriteCond %{TIME_HOUR}%{TIME_MIN} >1700
RewriteRule ^.*$ http://google.com/ [L,R=302]
```

### Reverse Proxy to Backend (instead of redirect)
```apache
RewriteEngine On
RewriteCond %{HTTP_USER_AGENT} curl [NC]
RewriteRule ^.*$ http://BACKEND_IP/404.html [P]
RewriteRule ^.*$ http://BACKEND_IP/login.html [P]
```
Requires: `sudo a2enmod rewrite proxy proxy_http`

### Referrer-Policy
**Apache:**
```apache
<IfModule mod_headers.c>
    Header set Referrer-Policy "no-referrer"
</IfModule>
```
**Nginx:**
```nginx
add_header Referrer-Policy "no-referrer";
```

---

## Cloudflare WAF

### Setup Steps
1. Add domain to Cloudflare
2. Update NS records to Cloudflare nameservers
3. Enable SSL (Full Strict)
4. Enable Bot Fight Mode + Block AI Bots
5. Set Security Level to High

### WAF Custom Rules
Block by country, ASN, or bot status:
- Condition: `Known Bots = true` OR `Country != US` OR `ASN = 1335`
- Action: Block

### Rate Limiting
- Condition: Path = `/login.php`, 5 requests in 10 seconds
- Action: Block for 1 hour
