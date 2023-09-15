YELLOW=\033[0;33m
NC=\033[0m # No Color

.PHONY: create-cert
create-cert:
	mkcert -cert-file statechannels.org.pem -key-file statechannels.org_key.pem statechannels.org localhost 127.0.0.1 ::1

.PHONY: install-mkcert-mac
install-mkcert-mac:
	brew install mkcert
	brew install nss # if you use Firefox
	
	mkcert -install

	@printf "${YELLOW}Remember to set NODE_EXTRA_CA_CERTS environment variable for your shell\n"
	@printf " see  https://github.com/FiloSottile/mkcert#using-the-root-with-nodejs ${NC}\n"
