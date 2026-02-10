#!/usr/bin/env sh

f1=$1
f2=$2

# checar que os dois arquivos existem
#
# extrair usando unzip

diff=$(diff -y --suppress-common-lines $f1 $f2 | wc -l)

if [ "$diff" -eq 0 ]; then
  echo "iguais"
else
  echo "diferentes"
fi
