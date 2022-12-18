#!/usr/bin/env bash
set -euo pipefail

doc_string_pattern='#[[:space:]]*@.+:.+'
fmt_pattern='s/#[[:space:]]*@(.+):[[:space:]]*(.+)/\1%\2/g'

main() {
  local output='Usage: make <target>\n\nTargets:'

  local doc_strings
  doc_strings="$(grep -E "$doc_string_pattern" "$@")"

  local curr_entry
  local prev_entry=''

  while IFS= read -r doc_string; do
    curr_entry="$(cut -d '%' -f1 <<<"$doc_string")"

    if [[ "$curr_entry" != "$prev_entry" ]]; then
      [[ -n "$prev_entry" ]] && output+='\n'
      prev_entry="$curr_entry"
      output+="\n\t${curr_entry}"
    fi

    value="$(cut -d '%' -f2 <<<"$doc_string")"
    output+="\n\t\t${value}"
  done < <(sed -r "$fmt_pattern" <<<"$doc_strings")

  echo -e "$output" | expand -t2
}

main "$@"
