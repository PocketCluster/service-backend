# Nginx Setup  

**Basic Setup (09/12/2016)**

Nginx and GO has four different op modes.  

1. Go HTTP standalone (as the control group)  
2. Nginx proxy to Go HTTP  
3. Nginx fastcgi to Go TCP FastCGI  
4. Nginx fastcgi to Go Unix Socket FastCGI  

We're to go with `Nginx proxy to Go HTTP` as it serves to be the fastest according to [a benchmark](https://gist.github.com/hgfischer/7965620).

The current configuration (09/12/2016) follows below

```sh
# GOLANG-GOJI SERVER
upstream go_http {
    server 127.0.0.1:8000;
    keepalive 300;
}

# HTTP Server with redirect from port 80 to 443 since we are now using TLS for our website
server {
    listen         80;
    server_name    index.pocketcluster.io;
    return 302 https://$server_name$request_uri;
}

# HTTPS server
server {
    listen         443;
    server_name    index.pocketcluster.io;

    # SSL configurations
    ssl on;
    ssl_protocols TLSv1.1 TLSv1.2;
    ssl_ciphers ECDH+AESGCM:DH+AESGCM:ECDH+AES256:DH+AES256:ECDH+AES128:DH+AES:ECDH+3DES:DH+3DES:RSA+AESGCM:RSA+AES:RSA+3DES:!aNULL:!MD5:!DSS;
    ssl_prefer_server_ciphers on;

    ssl_certificate /etc/letsencrypt/live/index.pocketcluster.io/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/index.pocketcluster.io/privkey.pem;

    # enable session resumption to improve https performance
    # http://vincent.bernat.im/en/blog/2011-ssl-session-reuse-rfc5077.html
    ssl_session_cache shared:SSL:50m;
    ssl_session_timeout 5m;

    # Diffie-Hellman parameter for DHE ciphersuites, recommended 2048 bits
    ssl_dhparam /etc/nginx/ssl.crt/dhparams.pem;

    # HTTP Strict Transport Security
    add_header Strict-Transport-Security "max-age=31536000; includeSubdomains;";

    location /robot.txt {
         root      /www-static;
         try_files $uri $uri/ =404;
    }

    location /theme {
        root      /www-static;
        # First attempt to serve request as file, then
        # as directory, then fall back to displaying a 404.
        try_files $uri $uri/ =404;
        autoindex on;
        # Uncomment to enable naxsi on this location
        # include /etc/nginx/naxsi.rules
    }

    location / {
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $remote_addr;
        proxy_set_header Host $host;
        proxy_pass http://go_http;
        proxy_http_version 1.1;
        proxy_set_header Connection "";
    }
}
```

### Pointers

- We need to forward Real IP addresses so set the header.

  ```sh
  proxy_set_header X-Real-IP $remote_addr;
  proxy_set_header X-Forwarded-For $remote_addr;
  proxy_set_header Host $host;
  ```
- `Keepalive` is enabled to reduce # of TCP connection to Nginx

  ```sh
  proxy_pass             http://your_upstream;
  
  # Default is HTTP/1, keepalive is only enabled in HTTP/1.1
  proxy_http_version 1.1;
  
  # Remove the Connection header if the client sends it,
  # it could be "close" to close a keepalive connection
  proxy_set_header Connection "";
  ```
- <sup>*</sup>We can futher reduce # of `location` block, using regex.
  
  
> Reference

- [Benchmarking Nginx with Go](Benchmarking Nginx with Go.pdf)
- [Enable Keepalive connections in Nginx Upstream proxy configurations](Enable Keepalive connections in Nginx Upstream proxy configurations.pdf)  
- [Understanding Nginx Server and Location Block Selection Algorithms](Understanding Nginx Server and Location Block Selection Algorithms _ DigitalOcean.pdf)