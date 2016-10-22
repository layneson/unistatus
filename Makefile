.PHONY: lib unistatus

export CGO_CFLAGS=-I../../../../../rpi-ws281x
export CGO_LDFLAGS=-L../../../../../rpi-ws281x -lm -lws2811
export CGO_ENABLED=1

export LIBRARY_PATH := $(LIBRARY_PATH):./rpi-ws281x

export GOOS=linux
export GOARCH=arm

all: unistatus

lib:
	make -C ./rpi-ws281x/ lib

unistatus: lib
	gb build

install:
	mkdir /opt/unistatus
	cp bin/* /opt/unistatus/
	ln -s /opt/unistatus/unistatus-linux-arm /usr/bin/local/unistatus
	cp unistatus.service /etc/systemd/system/
	systemctl enable unistatus.service

uninstall:
	rm -rf /opt/unistatus
	rm /usr/bin/local/unistatus
	rm /etc/systemd/system/unistatus.service
	rm /etc/systemd/system/multi-user.target.wants/unistatus.service