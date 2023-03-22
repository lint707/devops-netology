# Домашнее задание к занятию "3.5. Файловые системы"

1. Ознакомился с материалом о [sparse](https://ru.wikipedia.org/wiki/%D0%A0%D0%B0%D0%B7%D1%80%D0%B5%D0%B6%D1%91%D0%BD%D0%BD%D1%8B%D0%B9_%D1%84%D0%B0%D0%B9%D0%BB) (разряженных) файлах.

1. `Hardlink` это ссылка на тот же самый файл и имеет тот же `inode`, соотвественно права будут одни и теже, в качестве примера выполнил прверку с его созданием:
    ```bash
    vagrant@vagrant:~/test2$ ln new.txt new.link
    vagrant@vagrant:~/test2$ ls -ilh
    total 20K
    1179655 -rw-rw-r-- 2 vagrant vagrant 24 Mar 25 07:20 new.link
    1179655 -rw-rw-r-- 2 vagrant vagrant 24 Mar 25 07:20 new.txt
    vagrant@vagrant:~/test2$
    vagrant@vagrant:~/test2$ chmod 0555 new.txt
    vagrant@vagrant:~/test2$ ls -ilh
    total 20K
    1179655 -r-xr-xr-x 2 vagrant vagrant 24 Mar 25 07:20 new.link
    1179655 -r-xr-xr-x 2 vagrant vagrant 24 Mar 25 07:20 new.txt
    vagrant@vagrant:~/test2$
    ```
1. Выполнил `vagrant destroy` на имеющийся инстанс Ubuntu. Заменил содержимое Vagrantfile следующим:

    ```bash
    Vagrant.configure("2") do |config|
      config.vm.box = "bento/ubuntu-20.04"
      config.vm.provider :virtualbox do |vb|
        lvm_experiments_disk0_path = "/tmp/lvm_experiments_disk0.vmdk"
        lvm_experiments_disk1_path = "/tmp/lvm_experiments_disk1.vmdk"
        vb.customize ['createmedium', '--filename', lvm_experiments_disk0_path, '--size', 2560]
        vb.customize ['createmedium', '--filename', lvm_experiments_disk1_path, '--size', 2560]
        vb.customize ['storageattach', :id, '--storagectl', 'SATA Controller', '--port', 1, '--device', 0, '--type', 'hdd', '--medium', lvm_experiments_disk0_path]
        vb.customize ['storageattach', :id, '--storagectl', 'SATA Controller', '--port', 2, '--device', 0, '--type', 'hdd', '--medium', lvm_experiments_disk1_path]
      end
    end
    ```

    Создал новую виртуальную машину с двумя дополнительными неразмеченными дисками по 2.5 Гб.
    ```bash
    vagrant@vagrant:~$ lsblk
    NAME                      MAJ:MIN RM  SIZE RO TYPE MOUNTPOINT
    loop0                       7:0    0 70.3M  1 loop /snap/lxd/21029
    loop1                       7:1    0 55.4M  1 loop /snap/core18/2128
    loop2                       7:2    0 32.3M  1 loop /snap/snapd/12704
    loop3                       7:3    0 43.6M  1 loop /snap/snapd/15177
    loop4                       7:4    0 55.5M  1 loop /snap/core18/2344
    sda                         8:0    0   64G  0 disk
    ├─sda1                      8:1    0    1M  0 part
    ├─sda2                      8:2    0    1G  0 part /boot
    └─sda3                      8:3    0   63G  0 part
      └─ubuntu--vg-ubuntu--lv 253:0    0 31.5G  0 lvm  /    
    sdb                         8:16   0  2.5G  0 disk
    sdc                         8:32   0  2.5G  0 disk
    ```
    
1. Используя `fdisk`, разбил диск `sdb` на 2 раздела: 2 Гб и оставшееся пространство.
    ```bash
    vagrant@vagrant:~$ sudo fdisk /dev/sdb
    Welcome to fdisk (util-linux 2.34).
    Changes will remain in memory only, until you decide to write them.
    Be careful before using the write command.
    Device does not contain a recognized partition table.
    Created a new DOS disklabel with disk identifier 0xf6da8a78.
    
    Command (m for help): n
    Partition type
       p   primary (0 primary, 0 extended, 4 free)
       e   extended (container for logical partitions)
    Select (default p): p   
    Partition number (1-4, default 1): 1
    First sector (2048-5242879, default 2048): 2048
    Last sector, +/-sectors or +/-size{K,M,G,T,P} (2048-5242879, default 5242879): +2G
    Created a new partition 1 of type 'Linux' and of size 2 GiB.
    
    Command (m for help): n
    Partition type
       p   primary (1 primary, 0 extended, 3 free)
       e   extended (container for logical partitions)
    Select (default p): p
    Partition number (2-4, default 2): 2
    First sector (4196352-5242879, default 4196352): 4196352
    Last sector, +/-sectors or +/-size{K,M,G,T,P} (4196352-5242879, default 5242879):
    Created a new partition 2 of type 'Linux' and of size 511 MiB.
    
    Command (m for help): w
    The partition table has been altered.
    Calling ioctl() to re-read partition table.
    Syncing disks.
    ```
    Результат выполения:
    ```bash
    vagrant@vagrant:~$ lsblk
    NAME                      MAJ:MIN RM  SIZE RO TYPE MOUNTPOINT
    sdb                         8:16   0  2.5G  0 disk
    ├─sdb1                      8:17   0    2G  0 part
    └─sdb2                      8:18   0  511M  0 part
    ```
    
1. Используя `sfdisk`, перенёс таблицу разделов диска `sdb` на диск `sdc`.
    ```bash
    vagrant@vagrant:~$ sudo sfdisk -d /dev/sdb | sudo sfdisk --force /dev/sdc
    Checking that no-one is using this disk right now ... OK
    
    Disk /dev/sdc: 2.51 GiB, 2684354560 bytes, 5242880 sectors
    Disk model: VBOX HARDDISK
    Units: sectors of 1 * 512 = 512 bytes
    Sector size (logical/physical): 512 bytes / 512 bytes
    I/O size (minimum/optimal): 512 bytes / 512 bytes
    
    >>> Script header accepted.
    >>> Script header accepted.
    >>> Script header accepted.
    >>> Script header accepted.
    >>> Created a new DOS disklabel with disk identifier 0x1c9f6b42.
    /dev/sdc1: Created a new partition 1 of type 'Linux' and of size 2 GiB.
    /dev/sdc2: Created a new partition 2 of type 'Linux' and of size 511 MiB.
    /dev/sdc3: Done.
    
    New situation:
    Disklabel type: dos
    Disk identifier: 0x1c9f6b42
    
    Device     Boot   Start     End Sectors  Size Id Type
    /dev/sdc1          2048 4196351 4194304    2G 83 Linux
    /dev/sdc2       4196352 5242879 1046528  511M 83 Linux
    
    The partition table has been altered.
    Calling ioctl() to re-read partition table.
    Syncing disks.
    vagrant@vagrant:~$ lsblk
    NAME                      MAJ:MIN RM  SIZE RO TYPE MOUNTPOINT
    sdb                         8:16   0  2.5G  0 disk
    ├─sdb1                      8:17   0    2G  0 part
    └─sdb2                      8:18   0  511M  0 part
    sdc                         8:32   0  2.5G  0 disk
    ├─sdc1                      8:33   0    2G  0 part
    └─sdc2                      8:34   0  511M  0 part
    ```
    
1. Используя `mdadm` собрал RAID1 на паре разделов 2 Гб.
    ```bash
    vagrant@vagrant:~$ sudo mdadm --create --verbose /dev/md0 --level=1 --raid-devices=2 /dev/sdb1 /dev/sdc1
    mdadm: Note: this array has metadata at the start and
    may not be suitable as a boot device.  If you plan to
    store '/boot' on this device please ensure that
    your boot-loader understands md/v1.x metadata, or use
    --metadata=0.90
    mdadm: size set to 2094080K
    Continue creating array? y
    mdadm: Defaulting to version 1.2 metadata
    mdadm: array /dev/md0 started.
    vagrant@vagrant:~$ lsblk
    NAME                      MAJ:MIN RM  SIZE RO TYPE  MOUNTPOINT
    sdb                         8:16   0  2.5G  0 disk
    ├─sdb1                      8:17   0    2G  0 part
    │ └─md0                     9:0    0    2G  0 raid1
    └─sdb2                      8:18   0  511M  0 part
    sdc                         8:32   0  2.5G  0 disk
    ├─sdc1                      8:33   0    2G  0 part
    │ └─md0                     9:0    0    2G  0 raid1
    └─sdc2                      8:34   0  511M  0 part
    ``` 

1. Используя `mdadm` собрал RAID0 на второй паре маленьких разделов.
    ```bash
    vagrant@vagrant:~$ sudo mdadm --create --verbose /dev/md1 --level=0 --raid-devices=2 /dev/sdb2 /dev/sdc2
    mdadm: chunk size defaults to 512K
    mdadm: Defaulting to version 1.2 metadata
    mdadm: array /dev/md1 started.
    vagrant@vagrant:~$ lsblk
    NAME                      MAJ:MIN RM  SIZE RO TYPE  MOUNTPOINT
    sdb                         8:16   0  2.5G  0 disk
    ├─sdb1                      8:17   0    2G  0 part
    │ └─md0                     9:0    0    2G  0 raid1
    └─sdb2                      8:18   0  511M  0 part
      └─md1                     9:1    0 1018M  0 raid0
    sdc                         8:32   0  2.5G  0 disk
    ├─sdc1                      8:33   0    2G  0 part
    │ └─md0                     9:0    0    2G  0 raid1
    └─sdc2                      8:34   0  511M  0 part  
      └─md1                     9:1    0 1018M  0 raid0
    ```

1. Создал 2 независимых PV на получившихся md-устройствах.
    ```bash
    vagrant@vagrant:~$ sudo pvcreate /dev/md0
      Physical volume "/dev/md0" successfully created.
    vagrant@vagrant:~$ sudo pvcreate /dev/md1
      Physical volume "/dev/md1" successfully created.
    ```
    
1. Создал общую volume-group на этих двух PV.
    ```bash
    vagrant@vagrant:~$ sudo vgcreate vg0 /dev/md0 /dev/md1
      Volume group "vg0" successfully created
    vagrant@vagrant:~$ sudo pvdisplay
      --- Physical volume ---
      PV Name               /dev/sda3
      VG Name               ubuntu-vg
      PV Size               <63.00 GiB / not usable 0
      Allocatable           yes
      PE Size               4.00 MiB
      Total PE              16127
      Free PE               8063
      Allocated PE          8064
      PV UUID               sDUvKe-EtCc-gKuY-ZXTD-1B1d-eh9Q-XldxLf
    
      --- Physical volume ---
      PV Name               /dev/md0
      VG Name               vg0
      PV Size               <2.00 GiB / not usable 0
      Allocatable           yes
      PE Size               4.00 MiB
      Total PE              511
      Free PE               511
      Allocated PE          0
      PV UUID               xgFyEm-2cC1-LqLw-ypyb-FjgW-o8Uo-UP78lv
    
      --- Physical volume ---
      PV Name               /dev/md1
      VG Name               vg0
      PV Size               1018.00 MiB / not usable 2.00 MiB
      Allocatable           yes
      PE Size               4.00 MiB
      Total PE              254
      Free PE               254
      Allocated PE          0
      PV UUID               oG4Hnd-KZdd-UZ65-SguQ-cc9f-hkl3-YxlLLq
    ```

1. Создал LV размером 100 Мб, указав его расположение на PV с RAID0.
    ```bash
    vagrant@vagrant:~$ sudo lvcreate -L 100M vg0 /dev/md1
      Logical volume "lvol0" created.
    vagrant@vagrant:~$ lsblk
    NAME                      MAJ:MIN RM  SIZE RO TYPE  MOUNTPOINT
    sdb                         8:16   0  2.5G  0 disk
    ├─sdb1                      8:17   0    2G  0 part
    │ └─md0                     9:0    0    2G  0 raid1
    └─sdb2                      8:18   0  511M  0 part
      └─md1                     9:1    0 1018M  0 raid0
        └─vg0-lvol0           253:1    0  100M  0 lvm
    sdc                         8:32   0  2.5G  0 disk
    ├─sdc1                      8:33   0    2G  0 part
    │ └─md0                     9:0    0    2G  0 raid1
    └─sdc2                      8:34   0  511M  0 part
      └─md1                     9:1    0 1018M  0 raid0
        └─vg0-lvol0           253:1    0  100M  0 lvm
    ```

1. Создайте `mkfs.ext4` ФС на получившемся LV.
    ```bash
    vagrant@vagrant:~$ sudo mkfs.ext4 /dev/vg0/lvol0
    mke2fs 1.45.5 (07-Jan-2020)
    Creating filesystem with 25600 4k blocks and 25600 inodes
    
    Allocating group tables: done
    Writing inode tables: done
    Creating journal (1024 blocks): done
    Writing superblocks and filesystem accounting information: done
    ```

1. Смонтировал этот раздел в директорию `/tmp/vg0`.
    ```bash
    vagrant@vagrant:~$ sudo mkdir /tmp/vg0
    vagrant@vagrant:~$ sudo mount /dev/vg0/lvol0 /tmp/vg0
    ```

1. Поместил тестовый файл: `wget https://mirror.yandex.ru/ubuntu/ls-lR.gz -O /tmp/vg0/test.gz`.
    ```bash
    vagrant@vagrant:/tmp/vg0$ sudo wget https://mirror.yandex.ru/ubuntu/ls-lR.gz -O /tmp/vg0/test.gz
    --2022-04-04 07:59:31--  https://mirror.yandex.ru/ubuntu/ls-lR.gz
    Resolving mirror.yandex.ru (mirror.yandex.ru)... 213.180.204.183, 2a02:6b8::183
    Connecting to mirror.yandex.ru (mirror.yandex.ru)|213.180.204.183|:443... connected.
    HTTP request sent, awaiting response... 200 OK
    Length: 22314849 (21M) [application/octet-stream]
    Saving to: ‘/tmp/vg0/test.gz’
    
    /tmp/vg0/test.gz                          100%[=====================================================================================>]  21.28M  10.5MB/s    in 2.0s
    
    2022-04-04 07:59:34 (10.5 MB/s) - ‘/tmp/vg0/test.gz’ saved [22314849/22314849]
    
    vagrant@vagrant:/tmp/vg0$ ls -l
    total 21808
    drwx------ 2 root root    16384 Apr  4 07:50 lost+found
    -rw-r--r-- 1 root root 22314849 Apr  4 06:49 test.gz
    ```

1. Вывод `lsblk`.
    ```bash
    vagrant@vagrant:~$ lsblk
    NAME                      MAJ:MIN RM  SIZE RO TYPE  MOUNTPOINT
    sdb                         8:16   0  2.5G  0 disk
    ├─sdb1                      8:17   0    2G  0 part
    │ └─md0                     9:0    0    2G  0 raid1
    └─sdb2                      8:18   0  511M  0 part
      └─md1                     9:1    0 1018M  0 raid0
        └─vg0-lvol0           253:1    0  100M  0 lvm   /tmp/vg0
    sdc                         8:32   0  2.5G  0 disk
    ├─sdc1                      8:33   0    2G  0 part
    │ └─md0                     9:0    0    2G  0 raid1
    └─sdc2                      8:34   0  511M  0 part
      └─md1                     9:1    0 1018M  0 raid0
        └─vg0-lvol0           253:1    0  100M  0 lvm   /tmp/vg0
    ```
    
1. Протестируйте целостность файла:
    ```bash
    vagrant@vagrant:~$ sudo gzip -t /tmp/vg0/test.gz
    vagrant@vagrant:~$ sudo echo $?
    0
    ```

1. Используя pvmove, переместил содержимое PV с RAID0 на RAID1.
    ```bash
    vagrant@vagrant:~$ sudo pvmove -n /dev/vg0/lvol0 /dev/md1 /dev/md0
      /dev/md1: Moved: 12.00%
      /dev/md1: Moved: 100.00%
    vagrant@vagrant:~$ lsblk
    NAME                      MAJ:MIN RM  SIZE RO TYPE  MOUNTPOINT
    sdb                         8:16   0  2.5G  0 disk
    ├─sdb1                      8:17   0    2G  0 part
    │ └─md0                     9:0    0    2G  0 raid1
    │   └─vg0-lvol0           253:1    0  100M  0 lvm   /tmp/vg0
    └─sdb2                      8:18   0  511M  0 part
      └─md1                     9:1    0 1018M  0 raid0
    sdc                         8:32   0  2.5G  0 disk
    ├─sdc1                      8:33   0    2G  0 part
    │ └─md0                     9:0    0    2G  0 raid1
    │   └─vg0-lvol0           253:1    0  100M  0 lvm   /tmp/vg0
    └─sdc2                      8:34   0  511M  0 part
      └─md1                     9:1    0 1018M  0 raid0
    ```
    
1. Сделали `--fail` на устройство в вашем RAID1 md.
    ```bash
    vagrant@vagrant:~$ vagrant@vagrant:~$ sudo mdadm /dev/md0 --fail /dev/sdc1
    mdadm: set /dev/sdc1 faulty in /dev/md0
    ```
1. Подтвердите выводом `dmesg`, что RAID1 работает в деградированном состоянии.
    ```bash
    vagrant@vagrant:~$ dmesg | grep md0
    [ 2283.455009] md/raid1:md0: not clean -- starting background reconstruction
    [ 2283.455011] md/raid1:md0: active with 2 out of 2 mirrors
    [ 2283.455033] md0: detected capacity change from 0 to 2144337920
    [ 2283.455929] md: resync of RAID array md0
    [ 2293.766432] md: md0: resync done.
    [ 6030.120992] md/raid1:md0: Disk failure on sdc1, disabling device.
                   md/raid1:md0: Operation continuing on 1 devices.
    ```bash
1. Проверил целостность файла, он доступен:
    ```bash
    vagrant@vagrant:~$ sudo gzip -t /tmp/vg0/test.gz
    vagrant@vagrant:~$ sudo echo $?
    0
    ```

1. Погасите тестовый хост, `vagrant destroy`.
    ```bash
    vagrant@vagrant:~$ exit
    logout
    Connection to 127.0.0.1 closed.
    PS C:\Users\lint\vagrant> vagrant suspend
    ==> default: Saving VM state and suspending execution...
    PS C:\Users\lint\vagrant> vagrant halt
    ==> default: Discarding saved state of VM...
    PS C:\Users\lint\vagrant> vagrant destroy
    default: Are you sure you want to destroy the 'default' VM? [y/N] y
    ==> default: Destroying VM and associated drives...
    PS C:\Users\lint\vagrant>
    ``` 
 ---

