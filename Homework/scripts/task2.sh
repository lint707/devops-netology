#!/usr/bin/env bash
while ((1==1))
do
	curl http://localhost
	if (($? != 0))
	then
		date >> curl.log
	else break 
	fi
	sleep 5
done
