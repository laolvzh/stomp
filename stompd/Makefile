all:
	go build -ldflags "-X github.com/KristinaEtc/slflog.configLogFile=/usr/share/go-stomp-server/go-stomp-server.logconfig \
	-X main.configFile=/usr/share/go-stomp-server/go-stomp-server.config" go-stomp-server.go
	cp go-stomp-server /home/k/work/go/src/Nominatim/deb-package/go-stomp-server/usr/bin/
