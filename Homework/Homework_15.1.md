# Домашнее задание к занятию «Организация сети»

---
### Задание 1. Yandex Cloud 

**Что нужно сделать**

1. Создать пустую VPC. Выбрать зону.
2. Публичная подсеть:
 - Создать в VPC subnet с названием public, сетью 192.168.10.0/24.
 - Создать в этой подсети NAT-инстанс, присвоив ему адрес 192.168.10.254. В качестве image_id использовать fd80mrhj8fl2oe87o4e1.
 - Создать в этой публичной подсети виртуалку с публичным IP, подключиться к ней и убедиться, что есть доступ к интернету.

```
user@user-VirtualBox:~/trf$ ssh user@62.84.124.9
The authenticity of host '62.84.124.9 (62.84.124.9)' can't be established.
ECDSA key fingerprint is SHA256:xorjLSa7faqg3+QxIBApQmXw8iCWG3+PUUFtsK3Hu7c.
Are you sure you want to continue connecting (yes/no/[fingerprint])? yes
Warning: Permanently added '62.84.124.9' (ECDSA) to the list of known hosts.
Welcome to Ubuntu 20.04.6 LTS (GNU/Linux 5.4.0-148-generic x86_64)

 * Documentation:  https://help.ubuntu.com
 * Management:     https://landscape.canonical.com
 * Support:        https://ubuntu.com/advantage

The programs included with the Ubuntu system are free software;
the exact distribution terms for each program are described in the
individual files in /usr/share/doc/*/copyright.

Ubuntu comes with ABSOLUTELY NO WARRANTY, to the extent permitted by
applicable law.
```
```
user@publicvm:~$ ip -br a
lo               UNKNOWN        127.0.0.1/8 ::1/128 
eth0             UP             192.168.10.34/24 fe80::d20d:1fff:fe86:80ce/64 
```
```
user@publicvm:~$ ping 8.8.8.8
PING 8.8.8.8 (8.8.8.8) 56(84) bytes of data.
64 bytes from 8.8.8.8: icmp_seq=1 ttl=61 time=21.1 ms
64 bytes from 8.8.8.8: icmp_seq=2 ttl=61 time=20.7 ms
64 bytes from 8.8.8.8: icmp_seq=3 ttl=61 time=20.7 ms
64 bytes from 8.8.8.8: icmp_seq=4 ttl=61 time=20.8 ms
64 bytes from 8.8.8.8: icmp_seq=5 ttl=61 time=20.7 ms
^C
--- 8.8.8.8 ping statistics ---
5 packets transmitted, 5 received, 0% packet loss, time 4007ms
rtt min/avg/max/mdev = 20.652/20.772/21.107/0.171 ms
```

3. Приватная подсеть:
 - Создать в VPC subnet с названием private, сетью 192.168.20.0/24.
 - Создать route table. Добавить статический маршрут, направляющий весь исходящий трафик private сети в NAT-инстанс.
 - Создать в этой приватной подсети виртуалку с внутренним IP, подключиться к ней через виртуалку, созданную ранее, и убедиться, что есть доступ к интернету.
```
user@publicvm:~$ ssh user@192.168.20.13
The authenticity of host '192.168.20.13 (192.168.20.13)' can't be established.
ECDSA key fingerprint is SHA256:ffsZ/dt2UGni2fyLZIIytKvXudaFz8Uh/1KMtRTCwAo.
Are you sure you want to continue connecting (yes/no/[fingerprint])? yes
Warning: Permanently added '192.168.20.13' (ECDSA) to the list of known hosts.
Welcome to Ubuntu 20.04.6 LTS (GNU/Linux 5.4.0-148-generic x86_64)

 * Documentation:  https://help.ubuntu.com
 * Management:     https://landscape.canonical.com
 * Support:        https://ubuntu.com/advantage

The programs included with the Ubuntu system are free software;
the exact distribution terms for each program are described in the
individual files in /usr/share/doc/*/copyright.

Ubuntu comes with ABSOLUTELY NO WARRANTY, to the extent permitted by
applicable law.
```
```
user@privatevm:~$ ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host 
       valid_lft forever preferred_lft forever
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc mq state UP group default qlen 1000
    link/ether d0:0d:12:1f:a5:5c brd ff:ff:ff:ff:ff:ff
    inet 192.168.20.13/24 brd 192.168.20.255 scope global eth0
       valid_lft forever preferred_lft forever
    inet6 fe80::d20d:12ff:fe1f:a55c/64 scope link 
       valid_lft forever preferred_lft forever
```
```
user@privatevm:~$ ping 8.8.8.8
PING 8.8.8.8 (8.8.8.8) 56(84) bytes of data.
64 bytes from 8.8.8.8: icmp_seq=1 ttl=59 time=19.6 ms
64 bytes from 8.8.8.8: icmp_seq=2 ttl=59 time=18.3 ms
64 bytes from 8.8.8.8: icmp_seq=3 ttl=59 time=18.5 ms
64 bytes from 8.8.8.8: icmp_seq=4 ttl=59 time=18.4 ms
^C
--- 8.8.8.8 ping statistics ---
4 packets transmitted, 4 received, 0% packet loss, time 3005ms
rtt min/avg/max/mdev = 18.320/18.679/19.559/0.510 ms
```

[main.tf](file/main.tf)
[variables.tf](file/variables.tf)
[versions.tf](file/versions.tf)
[vpc.tf](file/vpc.tf)

---
