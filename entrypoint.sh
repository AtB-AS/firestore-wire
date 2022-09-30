#!/bin/sh
document="$1"

result=`echo "$1" | /bin/firestore-wire`

echo "::set-output name=result::$result"
