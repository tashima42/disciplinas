#!/usr/bin/env bash

source ./ids.sh

for id in "${IDS[@]}"; do
  echo "scraping $id"

  curl "https://www.wikiaves.com.br/$id" \
    -X 'GET' \
    -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8' \
    -H 'Sec-Fetch-Site: none' \
    -H 'Cookie: CHAVE=fa07c33ad0ff1e4bc83ce638a7007bb7; LOGIN=9a8e661d02a03c3; LOGINSPY=9a8e661d02a03c3; mybb[forumread]=a%3A1%3A%7Bi%3A11%3Bi%3A1760619509%3B%7D; mybb[lastactive]=1760619509; mybb[threadread]=a%3A1%3A%7Bi%3A4993%3Bi%3A1760619509%3B%7D; PHPSESSID=5382e914ea0145a2e39d0c2851466706; WIKILANG=pt-br' \
    -H 'Sec-Fetch-Mode: navigate' \
    -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.0 Safari/605.1.15' \
    -H 'Accept-Language: en-GB,en-US;q=0.9,en;q=0.8' \
    -H 'Accept-Encoding: gzip' \
    -H 'Sec-Fetch-Dest: document' \
    -H 'Priority: u=0, i' \
    --compressed \
    --output "dist/$id.html"

  echo "done $id"

  sleep 1
done
