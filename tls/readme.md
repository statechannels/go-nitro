# TLS

TLS private key and self signed certificate need to be located in this directory. To generate these:

```bash
mkcert -cert-file tls/statechannels.org.pem -key-file tls/statechannels.org_key.pem statechannels.org localhost 127.0.0.1 ::1
```

To use these for nodejs tests, see https://github.com/FiloSottile/mkcert#using-the-root-with-nodejs
