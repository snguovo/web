#!/bin/sh
go build -v -ldflags "-s -w" -o  web_nom main.go;
if [ -f "web" ]; then
rm -rf web
fi
upx -9 -o web web_nom;
rm -rf web_nom;