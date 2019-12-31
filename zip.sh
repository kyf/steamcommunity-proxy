#!/usr/bin/env bash
rm -v youxibd.zip
mkdir youxibd
cp -v ./youxibd.ico ./ui/static/
cp -rv ./ui ./youxibd
cp -rv ./certs ./youxibd
cp -v ./steamcommunity-proxy.exe ./youxibd/游戏便当Steam社区加速.exe
zip -rv youxibd.zip ./youxibd
rm -rvf ./youxibd
