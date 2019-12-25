#!/usr/bin/env bash
mkdir youxibd
cp -rv ./ui ./youxibd
cp -rv ./certs ./youxibd
cp -v ./steamcommunity-proxy.exe ./youxibd
zip -rv youxibd.zip ./youxibd
rm -rvf ./youxibd
