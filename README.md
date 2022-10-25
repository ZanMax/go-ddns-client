# Go DDNS Client
Client for DDNS written on GO

### Provider support
- Clourflare

### Configuration Linux service
Edit file service/ddns.service and change path to binary file

service ddns start
service ddns restart
service ddns enable
service ddns status

### IP source
By default the utility uses the site http://checkip.amazonaws.com/ to use the IP address, but you can change this.