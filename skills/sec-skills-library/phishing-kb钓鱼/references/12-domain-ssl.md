# Domain & SSL Certificate Reference

## Domain Selection Strategy

### Naming
- **Avoid** suspicious keywords (login, password, microsoft, etc.) — triggers DNS hunting
- **Avoid** typosquatting — actively hunted by security teams
- **Best approach**: neutral domain + targeted subdomain (e.g., `healthyfoods.com` with `docusign.healthyfoods.com`)
- Keep domain name under 15 characters
- Use `.com`, `.org`, `.net` TLDs — uncommon TLDs get blocked

### Pre-Purchase Checks
1. **WHOIS history** — check DomainTools for prior abuse
2. **VirusTotal** — search domain for malicious flags
3. **URLScan.io** — check for prior associations
4. **Wayback Machine** — review historical content
5. **AbuseIPDB** — check associated IP reputation

### Domain Add-ons
- **Privacy/WHOIS Guard** — always enable
- **Duration** — 1 year (extend if not burned)
- **Email** — optional but useful for SSL validation

---

## SSL Certificates

### Let's Encrypt (Free) - Apache
```bash
sudo apt install certbot python3-certbot-apache
sudo certbot --apache
```
**Note**: Free SSL triggers more scrutiny. Certificate is logged to CT immediately.

### Let's Encrypt (Free) - Nginx
```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx
```

### Wildcard Certificate (Hides Subdomains)
```bash
sudo certbot certonly --manual -d "*.domain.com" -d domain.com --preferred-challenges dns
```
Requires DNS TXT record `_acme-challenge` with provided value.

CT logs show `*.domain.com` — scanners can't determine active subdomains unless they brute-force.

**Use hard-to-guess subdomains** (e.g., `x7k9m2.domain.com` not `login.domain.com`).

### Paid SSL (Comodo/Sectigo)
```bash
openssl req -new -newkey rsa:2048 -nodes -keyout domain.key -out domain.csr
```

Submit CSR to provider, validate via DNS or email.

**Apache config for paid cert:**
```apache
<VirtualHost *:443>
    ServerName domain.com
    DocumentRoot /var/www/html
    SSLEngine on
    SSLCertificateFile /etc/ssl/certs/domain.crt
    SSLCertificateKeyFile /etc/ssl/private/domain.key
    SSLCertificateChainFile /etc/ssl/certs/domain.ca-bundle.crt
</VirtualHost>
```

**Result**: Less scanner traffic, not logged by all CT providers, higher trust score.

---

## Certificate Transparency Evasion

| Method | CT Exposure | Scanner Traffic |
|--------|-------------|-----------------|
| Single domain cert | Full domain visible | High |
| Subdomain cert | Full subdomain visible | Very high |
| Wildcard cert | Only `*.domain.com` visible | Low |
| Paid wildcard cert | Only `*.domain.com` visible | Very low |

---

## Domain Reputation Building

### Domain Aging
Register domain weeks/months before campaign. Freshly registered domains are flagged.

### Domain Categorization
Submit domain to categorization services to get classified as "Business", "Technology", etc.:
- Symantec/Bluecoat
- Fortiguard
- Palo Alto URL Filtering

Host benign content (blog, landing page) during aging period.

### Web Traffic Simulation
Generate legitimate-looking traffic to build domain reputation:
- Social media shares
- SEO indexing
- Organic-looking visits

---

## Provider Selection

### Hosting
- Use well-known cloud providers (AWS, Azure, DigitalOcean, Vultr)
- Avoid bulletproof hosting (poor IP reputation)
- Check IP on AbuseIPDB, VirusTotal, blacklist checkers before deploying
- Geographic proximity to target reduces suspicion

### Domain Registrar
- Research takedown speed — prefer slower responders
- Use separate registrar from hosting provider
- Use privacy protection on all domains
- Different registrar per campaign to avoid correlation
