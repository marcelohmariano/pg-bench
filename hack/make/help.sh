#!/usr/bin/env bash
set -e

echo -e 'Usage: make <target>\n'
echo 'Targets:'

doc_string_pattern='#[[:space:]]*@.+:.+'
doc_strings="$(grep -E "$doc_string_pattern" "$@")"

fmt_pattern='s/#[[:space:]]*@(.+):[[:space:]]*(.+)/  \1\t\2/g'
sed -r "$fmt_pattern" <<<"$doc_strings" | expand -t5
