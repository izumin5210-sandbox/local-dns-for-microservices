# Local service discovery for microservices development
## Run your local

```
# Install dependencies
$ dep ensure -v -vendor-only

# Start dns, reverse proxy and web servers
$ DOCKER_HOST_IP=$(/sbin/ifconfig en0 inet | tail -n 1 | awk '{ print $2 }') docker-compose up -d
```

### Do requests

```
$ curl foo.izumin.local/ping
pong

$ curl 'bar.izumin.local/echo?message=Hello!'
Hello!

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
