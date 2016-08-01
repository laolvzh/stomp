#!/bin/bash
set -e
autoconf
./configure --prefix=/
make
sudo checkinstall -D --pkgversion=0.6 --pkgname=go-stomp-server \
       --maintainer="Kristina Kovalevskaya isitiriss@gmail.com" --autodoinst=yes \
       --spec=ABOUT.md --provides=""

RETVAL=$?
[ $RETVAL -eq 0 ] && echo Success
[ $RETVAL -ne 0 ] && echo Failure
