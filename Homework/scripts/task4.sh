#!/usr/bin/env bash
hosts=(192.168.0.1 173.194.222.113 87.250.250.242)
port=80
while ((1==1))
do
	for addr in ${hosts[@]}
	do
		curl "$addr:$port"
		if [ $? != 0 ]
		then
			echo -n $i "|" $addr":"$port $(date "+%D %T")": failed" >> error.log
			exit
		else
			echo -n $i "|" $addr":"$port $(date "+%D %T")": DONE" >> ping.log
		fi
	done
done
