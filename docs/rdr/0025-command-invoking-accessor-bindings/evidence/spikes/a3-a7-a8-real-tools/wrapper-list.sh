#!/bin/sh
# R2 wrapper: git config --list (k=v lines) -> flat JSON object of strings. {artifact} arrives as $1.
git config --file "$1" --list | awk -F= 'BEGIN{printf "{"} {gsub(/"/,"\\\"",$2); printf "%s\"%s\":\"%s\"", (NR>1?",":""), $1, substr($0,length($1)+2)} END{printf "}\n"}'
