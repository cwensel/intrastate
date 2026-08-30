#!/bin/sh
# R5 wrapper: emit a ~2 MiB JSON object (single key, 2 MiB value).
printf '{"blob":"'; head -c 2097152 /dev/zero | tr '\0' 'x'; printf '"}\n'
