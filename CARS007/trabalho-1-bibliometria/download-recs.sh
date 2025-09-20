#!/bin/bash

source .env

archives_remaining=10078
total_files=21
start_file=12

increase_by=500

start=5511
end=$(($start + increase_by))
#end=$increase_by

archives_remaining=$(($archives_remaining - $start))

for i in $(seq $start_file $total_files); do
  echo "$i - start: $start | end: $end | archives_remaining: $archives_remaining | increase_by: $increase_by"

  curl 'https://www-webofscience-com.ez48.periodicos.capes.gov.br/api/wosnx/indic/export/saveToFile' \
    -X POST \
    -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:142.0) Gecko/20100101 Firefox/142.0' \
    -H 'Accept: application/json, text/plain, */*' \
    -H 'Accept-Language: en-US,en;q=0.5' \
    -H 'Accept-Encoding: gzip, deflate, br, zstd' \
    -H 'X-1P-WOS-SID: USW2EC0A3Achxkxla37hKCXo9Bxkn' \
    -H 'Content-Type: application/json' \
    -H 'Origin: https://www-webofscience-com.ez48.periodicos.capes.gov.br' \
    -H 'DNT: 1' \
    -H 'Connection: keep-alive' \
    -H 'Referer: https://www-webofscience-com.ez48.periodicos.capes.gov.br/wos/woscc/summary/79f1087c-6e6e-44e9-84ce-508be95d38a8-0176057d87/relevance/1(overlay:export/ext)' \
    -H 'Sec-Fetch-Dest: empty' \
    -H 'Sec-Fetch-Mode: cors' \
    -H 'Sec-Fetch-Site: same-origin' \
    -H 'Priority: u=0' \
    -H "Cookie: $AUTH_COOKIES" \
    --data-raw "{\"parentQid\":\"79f1087c-6e6e-44e9-84ce-508be95d38a8-0176057d87\",\"sortBy\":\"relevance\",\"displayTimesCited\":\"true\",\"displayCitedRefs\":\"true\",\"product\":\"UA\",\"colName\":\"WOS\",\"displayUsageInfo\":\"true\",\"fileOpt\":\"othersoftware\",\"action\":\"saveToTab\",\"markFrom\":\"$start\",\"markTo\":\"$end\",\"view\":\"summary\",\"isRefQuery\":\"false\",\"locale\":\"en_US\",\"filters\":\"fullRecordPlus\"}" \
    -o "data/savedrecs-$i.txt"

  start=$(($start + $increase_by))
  archives_remaining=$(($archives_remaining - $increase_by))
  if (($archives_remaining < $increase_by)); then
    increase_by=$archives_remaining
  fi
  end=$(($end + $increase_by))

  echo "%%%%%%%%%%%%%%%%%%%%%%%%%%%%%"
  echo "size: $(cat data/savedrecs-$i.txt | wc -l)"
  echo "%%%%%%%%%%%%%%%%%%%%%%%%%%%%%"
  echo ""

  sleep 5
done
