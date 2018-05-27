# Local service discovery for microservices development
## Run your local
### Setup for macos users

```
# use localhost DNS
$ sudo mkdir /etc/resolver && sudo echo 'nameserver 127.0.0.1' > /etc/resolver/local
```

### Start servers
```
# Install dependencies
$ make setup

# Build binaries
$ make

# Start discoverer and web servers
$ sudo make run
```

### Do requests

```
$ curl web1.services.local/ping
pong

$ curl 'web2.services.local/echo?message=Hello!'
Hello!

# server1 connects with server2
$ curl 'web1.services.local/delegate?url=http://web3.services.local/ping'
pong

$ curl -I google.com
HTTP/1.1 301 Moved Permanently
Location: http://www.google.com/
Content-Type: text/html; charset=UTF-8
Date: Thu, 10 May 2018 00:57:39 GMT
Expires: Sat, 09 Jun 2018 00:57:39 GMT
Cache-Control: public, max-age=2592000
Server: gws
Content-Length: 219
X-XSS-Protection: 1; mode=block
X-Frame-Options: SAMEORIGIN
```
