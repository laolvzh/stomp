#!/bin/bash
sudo service go-stomp-server stop
rm -r logs
sudo rm -r /opt/go-stomp/go-stomp-server/logs
sudo rm /var/log/upstart/go-stomp-server.log
go build
sudo cp stompd /opt/go-stomp/go-stomp-server/stomp
sudo service go-stomp-server start
sudo service go-stomp-server status
echo "------------------"
sudo cat /var/log/upstart/go-stomp-server.log
