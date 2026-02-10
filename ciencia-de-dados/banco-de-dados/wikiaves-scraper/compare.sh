#!/usr/bin/env sh

f1=$1
f2=$2

# checar que os dois arquivos existem
#
# extrair usando unzip

checksum1=$(md5sum --quiet $f1)
checksum2=$(md5sum --quiet $f2)

echo "$checksum1 = $checksum2"

if [ "$checksum1" == "$checksum2" ]; then
  echo "iguais"
else
  echo "diferentes"
fi
