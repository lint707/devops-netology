# Домашнее задание к занятию "3.3. Операционные системы, лекция 1"

1. Выполнил команду `strace /bin/bash -c 'cd /tmp'`, системный вызов для `cd`: 
    ```
    chdir("/tmp")                           = 0
    ```
    
1. Выполнил команду `file` на объектах разных типов на файловой системе. Например:
    ```bash
    vagrant@vagrant:~$ file /dev/tty
    /dev/tty: character special (5/0)
    vagrant@vagrant:~$ file /dev/sda
    /dev/sda: block special (8/0)
    vagrant@vagrant:~$ file /dev/pts
    /dev/pts: directory
    vagrant@vagrant:~$ file /bin/bash
    /bin/bash: ELF 64-bit LSB shared object, x86-64, version 1 (SYSV), dynamically linked, interpreter /lib64/ld-linux-x86-64.so.2, BuildID[sha1]=a6cb40078351e05121d46daa768e271846d5cc54, for GNU/Linux 3.2.0, stripped
    ```
    База данных `file` находится в: `"/usr/share/misc/magic.mgc"`.
    ```
    stat("/home/vagrant/.magic.mgc", 0x7ffc4a138800) = -1 ENOENT (No such file or directory)
    stat("/home/vagrant/.magic", 0x7ffc4a138800) = -1 ENOENT (No such file or directory)
    openat(AT_FDCWD, "/etc/magic.mgc", O_RDONLY) = -1 ENOENT (No such file or directory)
    stat("/etc/magic", {st_mode=S_IFREG|0644, st_size=111, ...}) = 0
    openat(AT_FDCWD, "/etc/magic", O_RDONLY) = 3
    fstat(3, {st_mode=S_IFREG|0644, st_size=111, ...}) = 0
    read(3, "# Magic local data for file(1) c"..., 4096) = 111
    read(3, "", 4096)                       = 0
    close(3)                                = 0
    openat(AT_FDCWD, "/usr/share/misc/magic.mgc", O_RDONLY) = 3
    ```
    
1. Предположим, приложение пишет лог в текстовый файл. Этот файл оказался удален (deleted в lsof), однако возможности сигналом сказать приложению переоткрыть файлы или просто перезапустить приложение – нет. Так как приложение продолжает писать в удаленный файл, место на диске постепенно заканчивается. Основываясь на знаниях о перенаправлении потоков предложите способ обнуления открытого удаленного файла (чтобы освободить место на файловой системе). </br>
    Терминал 1. Создание файла и настройка логирования:
    ```
    vagrant@vagrant:~/tping$ echo " " > ping.log
    vagrant@vagrant:~/tping$ exec 5> ping.log
    vagrant@vagrant:~/tping$ ping 127.0.0.1 >&5
    ```
    Терминал 2. Удаление файла:
    ```
    vagrant@vagrant:~/tping$ rm ping.log
    vagrant@vagrant:~/tping$ ls -l
    total 0
    ```
    Терминал 2. Поиск: 
    ```
    vagrant@vagrant:~/tping$ lsof | grep ping
    bash      1358                        vagrant  cwd       DIR              253,0     4096    1181766 /home/vagrant/tping
    bash      1358                        vagrant    5w      REG              253,0    14904    1181768 /home/vagrant/tping/ping.log (deleted)
    vagrant@vagrant:~/tping$ sudo lsof -p 1358 | grep ping
    bash    1358 vagrant  cwd    DIR  253,0     4096 1181766 /home/vagrant/tping
    bash    1358 vagrant    5w   REG  253,0    14904 1181768 /home/vagrant/tping/ping.log (deleted)
    ```
    Терминал 2. Очистка:
    ```
    vagrant@vagrant:~/tping$ cat /dev/null | sudo tee /proc/1358/fd/5
    vagrant@vagrant:~/tping$ sudo lsof -p 1358 | grep ping
    bash    1358 vagrant  cwd    DIR  253,0     4096 1181766 /home/vagrant/tping
    bash    1358 vagrant    5w   REG  253,0        0 1181768 /home/vagrant/tping/ping.log (deleted)
    vagrant@vagrant:~/tping$
    ```

1. Зомби-процессы не занимают ресурсы в ОС (CPU, RAM, IO),  но не освобождают запись в таблице процессов.
    
1. В iovisor BCC есть утилита `opensnoop`:
    ```
    vagrant@vagrant:~$ dpkg -L bpfcc-tools | grep sbin/opensnoop
    /usr/sbin/opensnoop-bpfcc
    vagrant@vagrant:~$ sudo /usr/sbin/opensnoop-bpfcc
    PID    COMM               FD ERR PATH
    882    vminfo              4   0 /var/run/utmp
    655    dbus-daemon        -1   2 /usr/local/share/dbus-1/system-services
    655    dbus-daemon        19   0 /usr/share/dbus-1/system-services
    655    dbus-daemon        -1   2 /lib/dbus-1/system-services
    655    dbus-daemon        19   0 /var/lib/snapd/dbus-1/system-services/
    661    irqbalance          6   0 /proc/interrupts
    661    irqbalance          6   0 /proc/stat
    ```
    
1.  `uname -a` использует системный вызов:
    ```
    vagrant@vagrant:~$ strace uname -a > uname.log 2>&1
    vagrant@vagrant:~$ cat uname.log |grep uname
    execve("/usr/bin/uname", ["uname", "-a"], 0x7ffe1e653ea8 /* 24 vars */) = 0
    uname({sysname="Linux", nodename="vagrant", ...}) = 0
    uname({sysname="Linux", nodename="vagrant", ...}) = 0
    uname({sysname="Linux", nodename="vagrant", ...}) = 0
    vagrant@vagrant:~$:
    ```
    Цитата из `man 2 uname` по этому системному вызову:
    ```    
           Part of the utsname information is also accessible via /proc/sys/kernel/{ostype, hostname, osrelease, version,
       domainname}.
    ```
    
1. Оператор `;` выполняет несколько команд одновременно последовательно и обеспечивает вывод без зависимости от успеха и отказа других команд.
    Оператор `&&` (AND оператор) выполнит вторую команду только в том случае, если команда 1 успешно выполнена.
    ```
    vagrant@vagrant:~$ test -d /tmp/some_dir; echo Hi
    Hi
    vagrant@vagrant:~$ test -d /tmp/some_dir&& echo Hi
    vagrant@vagrant:~$
    ```
    `set -e` останавливает выполнение скрипта, если команда или конвейер имеет ошибку, поэтому использование вместе с `&&`, не имеет смысла.


1. Режим bash `set -euxo pipefail` состоит из опций: </br>
    `-e` - прекращает выполнение скрипта если команда завершилась ошибкой, выводит в `stderr` строку с ошибкой. </br>
    `-u` - прекращает выполнение скрипта, если встретилась несуществующая переменная. </br>
    `-x` - выводит выполняемые команды в `stdout` перед выполненинем. </br>
    `-o pipefail` - прекращает выполнение скрипта, даже если одна из частей пайпа завершилась ошибкой. </br>
    Использование данного сценария полезно, в случаеналичия ошибок в скрипте, повышает удобство их отлова, увеличивает детализацию по логам и останавливает выполение скрипта, для последующего исправления. </br>

1. Наиболее часто встречающийся статус у процессов в системе: </br>
    ```
    vagrant@vagrant:~$ ps -o stat 
    STAT
    Ss
    R+
    ```
Дополнительные символы к основной заглавной буквы статуса процессов, нужны для указаия дополнительных характеристик. 
К примену для, `S` - прерываемый режим сна (ожидание завершения события)
`s` - является лидером сеанса
`l` - является многопоточным 
```
PROCESS STATE CODES
       Here are the different values that the s, stat and state output specifiers (header "STAT" or "S") will display
       to describe the state of a process:

               D    uninterruptible sleep (usually IO)
               I    Idle kernel thread
               R    running or runnable (on run queue)
               S    interruptible sleep (waiting for an event to complete)
               T    stopped by job control signal
               t    stopped by debugger during the tracing
               W    paging (not valid since the 2.6.xx kernel)
               X    dead (should never be seen)
               Z    defunct ("zombie") process, terminated but not reaped by its parent

       For BSD formats and when the stat keyword is used, additional characters may be displayed:

               <    high-priority (not nice to other users)
               N    low-priority (nice to other users)
               L    has pages locked into memory (for real-time and custom IO)
               s    is a session leader
               l    is multi-threaded (using CLONE_THREAD, like NPTL pthreads do)
               +    is in the foreground process group
```

---

