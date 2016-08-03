#!/bin/bash

export GO=$(shell which go)
export GOINSTALL="${GO} install"
export GOBUILD="${GO} build"
export GOCLEAN="${GO} clean"
export GOGET="${GO} get"

#:<<'COMMENT1'
export  EXENAME=go-stomp-server
#BUILDPATH ="/usr/local/$(EXENAME)"
export BUILDPATH=${PWD}

export logdir=/var/log/${EXENAME}
export bindir=/usr/bin
export demondir=/etc/init
export confdir=/etc/${EXENAME}

export CONF=${EXENAME}.config
export LOGCONF=${EXENAME}.logconfig
export DEMONF=${EXENAME}.conf

export GOFLAGS = "-o go-stomp-server \
	  -ldflags \"-X github.com/KristinaEtc/slflog.configLogFile=${confdir}/${LOGCONF} \
	  -X main.configFile=${confdir}/${CONF}\" "
#COMMENT1
