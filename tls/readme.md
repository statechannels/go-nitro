# TLS

TLS private key and certificate need to be located in this directory. To generate these, install [mkcert](https://github.com/FiloSottile/mkcert#macos). If on a mac, to install mkcert:

```bash
make install-mkcert-mac
```

On other platforms, run [the following command](https://github.com/FiloSottile/mkcert#mkcert) after installing mkcert:

```bash
mkcert -install
```

To create a new certificate, run:

```bash
make create-cert
```

To use these for nodejs tests, see https://github.com/FiloSottile/mkcert#using-the-root-with-nodejs
