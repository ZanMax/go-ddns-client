# Go DDNS Client
Client for DDNS written on GO

### Providers
- Clourflare

### Configuration
- Edit config.json
- Add domain or subdomain
- Add API TOKEN

### Build

```bash
./build.sh
```

### Configuration Linux service
Edit file service/ddns.service and change path to binary file

service ddns start
service ddns restart
service ddns enable
service ddns status

### Configuration supervisor

```bash
sudo apt update
sudo apt install -y supervisor
sudo service supervisor start
sudo supervisorctl status
```

```bash
sudo vim /etc/supervisor/conf.d/ddns.conf
```

```bash
mkdir /var/log/ddns
```

```
[program:ddns]
directory=/usr/local
command=/usr/local/bin/ddns
autostart=true
autorestart=true
stderr_logfile=/var/log/ddns/err.log
stdout_logfile=/var/log/ddns/out.log
```

```bash
sudo supervisorctl reload
```

### IP source
By default, the utility uses the site http://checkip.amazonaws.com/ to use the IP address, but you can change this.