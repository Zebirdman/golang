.PHONY: all binary install clean
export GOPATH:=$(CURDIR)/Godeps/_workspace:$(GOPATH)

LIBDIR=${DESTDIR}/lib/systemd/system
BINDIR=${DESTDIR}/usr/libexec/docker/

all: man binary

binary:
	go build  -o zeft-auth .

install:
	install -d -m 0755 ${LIBDIR}
	install -m 644 systemd/zeft-auth.service ${LIBDIR}
	install -d -m 0755 ${LIBDIR}
	install -m 644 systemd/zeft-auth.socket ${LIBDIR}
	install -d -m 0755 ${BINDIR}
	install -m 755 docker-novolume-plugin ${BINDIR}
	install -m 644 docker-novolume-plugin.8 ${MANINSTALLDIR}/man8/

clean:
	rm -f zeft-auth
