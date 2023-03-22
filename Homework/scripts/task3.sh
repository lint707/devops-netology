#!/usr/bin/env bash
hosts=(192.168.0.1 173.194.222.113 87.250.250.242)
port=80
count=( 1 2 3 4 5)
for i in ${count[@]}
do
	for addr in ${hosts[@]}
	do
		echo -n $i "|" $addr":"$port $(date "+%D %T") >>ping.log
		curl "$addr:$port"
		if [ $? != 0 ]
		then
			echo ": failed" >> ping.log
		else
			echo ": DONE" >> ping.log
		fi
	done
done
