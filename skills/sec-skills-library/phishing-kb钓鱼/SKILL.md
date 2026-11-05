---
name: phishing-kb
description: Offensive phishing operations knowledge base for red team engagements. Use when the user asks about phishing infrastructure, page design, anti-detection, credential harvesting, web server config, SSL/TLS, anti-bot, obfuscation, HTML/SVG smuggling, ClickFix, AITM/Evilginx, MFA bypass, JA3/JA4 fingerprinting, Cloudflare Workers, Terraform automation, or any red team web operations topic.
version: 2.0.0
user-invocable: false
tools: Read, Glob, Grep, Bash, Edit, Write
---

# Offensive Phishing Operations - Knowledge Base

Self-contained technical reference for building and operating phishing infrastructure in authorized red team engagements. All methods, code templates, and configurations are in the `references/` subdirectory files listed below.

## Core Design Principles

1. **Custom over cloned** - Build pages from scratch; cloned pages carry detectable signatures
2. **Encrypt stored data** - AES-256-CBC encrypt all harvested credentials at rest
3. **Store outside docroot** - Log files go in `/var/www/`, never inside `/var/www/html/`
4. **Wildcard certs** - Hide subdomains from Certificate Transparency logs
5. **Use redirectors** - Never expose the phishing server directly; reverse proxy everything
6. **Block direct IP** - Force all access through domain names only
7. **Referrer-Policy: no-referrer** - Prevent domain leakage on redirects
8. **Layer defenses** - Combine multiple anti-analysis techniques; never rely on one
9. **Dynamic encryption** - Generate unique keys per page load to defeat static signatures
10. **Mobile responsive** - Always include `@media` queries; neglecting mobile reduces success rate

## Reference Files

| File | Topics |
|------|--------|
| `references/01-infrastructure.md` | Apache, Nginx, Caddy, Flask+Gunicorn setup, Terraform IaC, firewall (UFW/iptables) |
| `references/02-credential-harvesting.md` | Login page templates, PHP credential logging, AES encryption/decryption, MySQL storage, metadata collection |
| `references/03-anti-analysis.md` | Cookie checks, WebDriver detection, UA filtering/spoofing detection, IP geo-restrictions, VPN/proxy blocking, window size, incognito/ad-blocker detection |
| `references/04-obfuscation.md` | Base64, XOR (single/multi-key), AES (Web Crypto API + aes-js), dynamic server-side encryption, obfuscator.io integration |
| `references/05-dynamic-html.md` | createElement generation, DOMContentLoaded wrappers, minification, setTimeout delays, remote content fetch, jQuery hiding |
| `references/06-smuggling.md` | HTML smuggling (Blob API, dynamic XOR, file-based), SVG smuggling (script tags, onload, image overlays), WebAssembly smuggling |
| `references/07-fingerprinting.md` | JA3/JA4/JARM concepts, HAProxy+Lua JA4 plugin, fingerprint blacklisting/whitelisting, client telemetry logging library |
| `references/08-anti-bot.md` | Honeypots (hidden input/link/header), CAPTCHA (reCAPTCHA, custom, invisible), website keying, user interaction analysis (keystroke/cursor speed, right-click, inactivity) |
| `references/09-mfa-bypass.md` | TOTP harvesting, AITM via reverse proxy, email/push/Discord approval gates, device code phishing (Microsoft/GitHub), illicit consent grant |
| `references/10-clickfix.md` | Clipboard API manipulation, ClickFix templates (extension-required, activation-code), anti-sandbox (mousedown/mouseup), obfuscated commands, remote content evasion |
| `references/11-server-security.md` | HTTPS enforcement, direct IP blocking, verbose info removal, Server header spoofing, SSH hardening, jump server setup, Cloudflare WAF rules |
| `references/12-domain-ssl.md` | Domain selection strategies, TLD choice, subdomain abuse, Let's Encrypt/paid SSL, wildcard certs, CT log evasion, domain aging/categorization/reputation |
| `references/13-serverless.md` | Cloudflare Worker login pages, KV storage, D1 SQL storage, worker deployment |

## How To Use

When writing code for any phishing-related task, consult the matching reference file above. Each file contains:
- **Technique overview** with rationale
- **Complete working code templates** ready to adapt
- **Variants and alternatives** for different scenarios
- **Opsec notes** highlighting common mistakes
