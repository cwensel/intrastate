#!/bin/sh
# R3 wrapper: stdin JSON object of strings -> one `git config --file $1 key value` per key.
# jq-free on purpose (relies on the C3 shape: flat object of strings, no escaped quotes in values).
{ cat; echo; } | sed -e 's/^{//; s/}$//; s/","/"\n"/g' | sed -e 's/^"\([^"]*\)":"\(.*\)"$/\1 \2/' \
  | while read -r k v; do [ -n "$k" ] && { git config --file "$1" "$k" "$v" || exit $?; }; done
