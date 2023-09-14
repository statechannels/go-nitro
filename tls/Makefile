MKCERT ?= $(LOCALBIN)/mkcert
MKCERT_VERSION ?= v1.4.4
...
.PHONY: mkcert
mkcert: $(MKCERT)
$(MKCERT): $(LOCALBIN)
	curl -L "https://github.com/FiloSottile/mkcert/releases/download/$(MKCERT_VERSION)/mkcert-$(MKCERT_VERSION)-$(UNAME)-amd64" -o $(LOCALBIN)/mkcert
	chmod +x $(LOCALBIN)/mkcert