#!/bin/bash

# AWS アイコンパッケージのリンクを取得
curl -s "https://aws.amazon.com/jp/architecture/icons/" | \
pup 'script[type="application/json"] text{}' | \
grep -o 'https://[^"]*\.zip' | \
sort -u
