# TLS

TLS private key and self signed certificate need to be located in this directory. To generate these:

```bash
make -C tls create-cert
```

To use these for nodejs tests, see https://github.com/FiloSottile/mkcert#using-the-root-with-nodejs
