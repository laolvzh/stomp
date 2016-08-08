#!/bin/bash
set -e
autoconf
./configure --prefix=/
make
fakeroot checkinstall --install=no --reset-uids=yes --fstrans=yes -D --pkgversion=0.6.1 --pkgname=go-stomp-server \
      --maintainer="Kristina Kovalevskaya isitiriss@gmail.com" \
       --spec=ABOUT.md --provides="" --pkgsource=go-stomp-server

RETVAL=$?
[ $RETVAL -eq 0 ] && echo Success
[ $RETVAL -ne 0 ] && echo Failure
