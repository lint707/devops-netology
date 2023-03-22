# Домашнее задание к занятию "08.01 Введение в Ansible"

## Подготовка к выполнению
1. Установил `ansible` версии 2.10 или выше:
```
user@user-VirtualBox:~/ansible_ntl/devops-netology-1/Homework/08-ansible-01-base/playbook$ ansible --version
ansible [core 2.12.8]
  config file = /etc/ansible/ansible.cfg
  configured module search path = ['/home/user/.ansible/plugins/modules', '/usr/share/ansible/plugins/modules']
  ansible python module location = /usr/lib/python3/dist-packages/ansible
  ansible collection location = /home/user/.ansible/collections:/usr/share/ansible/collections
  executable location = /usr/bin/ansible
  python version = 3.8.10 (default, Jun 22 2022, 20:18:18) [GCC 9.4.0]
  jinja version = 2.10.1
  libyaml = True
```
2. Создал свой собственный публичный репозиторий на github с произвольным именем.  
3. Скачал [playbook](./playbook/) из репозитория с домашним заданием и перенесите его в свой репозиторий.  

## Основная часть  
1. Запустил playbook на окружении из `test.yml`, зафиксировал какое значение имеет факт `some_fact` для указанного хоста при выполнении playbook'a.  
```bash
user@user-VirtualBox:~/ansible_ntl/devops-netology-1/Homework/08-ansible-01-base/playbook$ ansible-playbook -i inventory/test.yml site.yml

PLAY [Print os facts] **********************************************************************************************************

TASK [Gathering Facts] *********************************************************************************************************
ok: [localhost]

TASK [Print OS] ****************************************************************************************************************
ok: [localhost] => {
    "msg": "Ubuntu"
}

TASK [Print fact] **************************************************************************************************************
ok: [localhost] => {
    "msg": 12
}

PLAY RECAP *********************************************************************************************************************
localhost                  : ok=3    changed=0    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0   
```

2. Нашёл файл с переменными *(group_vars)* в котором задаётся найденное в первом пункте значение и поменял его на `all default fact`.  
`/group_vars/all` - "msg": `12` -  поменяйте его на `all default fact`.  

3. Подготовил окружение, используя `docker`, для проведения дальнейших испытаний.  
```bash
user@user-VirtualBox:~/ansible_ntl/compose$ docker ps 
CONTAINER ID   IMAGE                 COMMAND            CREATED          STATUS          PORTS     NAMES
ea7dc9f5cf1c   pycontribs/centos:7   "sleep infinity"   30 seconds ago   Up 23 seconds             centos7
e29c011000e1   pycontribs/ubuntu     "sleep infinity"   30 seconds ago   Up 23 seconds             ubuntu
```

4. Провёл запуск `playbook` на окружении из `prod.yml`. Зафиксировал полученные значения `some_fact` для каждого из `managed host`.
```bash
user@user-VirtualBox:~/ansible_ntl/devops-netology-1/Homework/08-ansible-01-base/playbook$ sudo ansible-playbook -i inventory/prod.yml site.yml

PLAY [Print os facts] **********************************************************************************************************

TASK [Gathering Facts] *********************************************************************************************************
ok: [ubuntu]
ok: [centos7]

TASK [Print OS] ****************************************************************************************************************
ok: [centos7] => {
    "msg": "CentOS"
}
ok: [ubuntu] => {
    "msg": "Ubuntu"
}

TASK [Print fact] **************************************************************************************************************
ok: [centos7] => {
    "msg": "el"
}
ok: [ubuntu] => {
    "msg": "deb"
}

PLAY RECAP *********************************************************************************************************************
centos7                    : ok=3    changed=0    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0   
ubuntu                     : ok=3    changed=0    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0   
```

5. Добавил факты в `group_vars` каждой из групп хостов так, чтобы для `some_fact` получились следующие значения:  
для `deb` - `deb default fact`   
для `el` - `el default fact`  

6.  Повторил запуск `playbook` на окружении `prod.yml`. Убедился, что выдаются корректные значения для всех хостов.
```bash
PLAY [Print os facts] **********************************************************************************************************

TASK [Gathering Facts] *********************************************************************************************************
ok: [ubuntu]
ok: [centos7]

TASK [Print OS] ****************************************************************************************************************
ok: [centos7] => {
    "msg": "CentOS"
}
ok: [ubuntu] => {
    "msg": "Ubuntu"
}

TASK [Print fact] **************************************************************************************************************
ok: [ubuntu] => {
    "msg": "deb default fact"
}
ok: [centos7] => {
    "msg": "el default fact"
}

PLAY RECAP *********************************************************************************************************************
centos7                    : ok=3    changed=0    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0   
ubuntu                     : ok=3    changed=0    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0   
```

7. При помощи `ansible-vault` зашифровал факты в `group_vars/deb` и `group_vars/el` с паролем `netology`.  
`ansible-vault encrypt group_vars/deb/examp.yml`  
`ansible-vault encrypt group_vars/el/examp.yml`  

8. Запустил playbook на окружении `prod.yml`. При запуске `ansible` должен запросить у вас пароль. Убедился в работоспособности.  
```bash
user@user-VirtualBox:~/ansible_ntl/devops-netology-1/Homework/08-ansible-01-base/playbook$ sudo ansible-playbook -i inventory/prod.yml site.yml

PLAY [Print os facts] **********************************************************************************************************
ERROR! Attempting to decrypt but no vault secrets found
user@user-VirtualBox:~/ansible_ntl/devops-netology-1/Homework/08-ansible-01-base/playbook$ sudo ansible-playbook -i inventory/prod.yml site.yml --ask-vault-pass
Vault password: 

PLAY [Print os facts] **********************************************************************************************************

TASK [Gathering Facts] *********************************************************************************************************
ok: [ubuntu]
ok: [centos7]

TASK [Print OS] ****************************************************************************************************************
ok: [centos7] => {
    "msg": "CentOS"
}
ok: [ubuntu] => {
    "msg": "Ubuntu"
}

TASK [Print fact] **************************************************************************************************************
ok: [centos7] => {
    "msg": "el default fact"
}
ok: [ubuntu] => {
    "msg": "deb default fact"
}

PLAY RECAP *********************************************************************************************************************
centos7                    : ok=3    changed=0    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0   
ubuntu                     : ok=3    changed=0    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0   
```

9. Посмотрел при помощи `ansible-doc` список плагинов для подключения. Выбрал подходящий для работы на `control node`:  
`local` - execute on controller  

10. В `prod.yml` добавил новую группу хостов с именем  `local`, в ней разместил `localhost` с необходимым типом подключения.  
```
  inside:
    hosts:
      localhost:
        ansible_connection: local
```

11. Запустил `playbook` на окружении `prod.yml`. При запуске `ansible` должен запросить у вас пароль. Убедился что факты `some_fact` для каждого из хостов определены из верных `group_vars`.
```bash
user@user-VirtualBox:~/ansible_ntl/devops-netology-1/Homework/08-ansible-01-base/playbook$ sudo ansible-playbook -i inventory/prod.yml site.yml --ask-vault-pass
Vault password: 

PLAY [Print os facts] **********************************************************************************************************

TASK [Gathering Facts] *********************************************************************************************************
ok: [localhost]
ok: [ubuntu]
ok: [centos7]

TASK [Print OS] ****************************************************************************************************************
ok: [centos7] => {
    "msg": "CentOS"
}
ok: [ubuntu] => {
    "msg": "Ubuntu"
}
ok: [localhost] => {
    "msg": "Ubuntu"
}

TASK [Print fact] **************************************************************************************************************
ok: [centos7] => {
    "msg": "el default fact"
}
ok: [ubuntu] => {
    "msg": "deb default fact"
}
ok: [localhost] => {
    "msg": "all default fact"
}

PLAY RECAP *********************************************************************************************************************
centos7                    : ok=3    changed=0    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0   
localhost                  : ok=3    changed=0    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0   
ubuntu                     : ok=3    changed=0    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0   
```

12. Заполнил `README.md` ответами на вопросы. Сделал `git push` в ветку `master`. В ответе отправьте ссылку на ваш открытый репозиторий с изменённым `playbook` и заполненным `README.md`.

## Необязательная часть

1. При помощи `ansible-vault` расшифровал все зашифрованные файлы с переменными.
```bash
user@user-VirtualBox:~/ansible_ntl/devops-netology-1/Homework/08-ansible-01-base/playbook$ ansible-vault decrypt --ask-vault-password group_vars/deb/* group_vars/el/*
Vault password: 
Decryption successful
```

2. Зашифровал отдельное значение `PaSSw0rd` для переменной `some_fact` паролем `netology`. Добавил полученное значение в `group_vars/all/exmp.yml`.
```bash
user@user-VirtualBox:~/ansible_ntl/devops-netology-1/Homework/08-ansible-01-base/playbook$ ansible-vault encrypt_string "PaSSw0rd"
New Vault password: 
Confirm New Vault password: 
!vault |
          $ANSIBLE_VAULT;1.1;AES256
          32303966363338336638306664343239363361363435643134656266663937313039656535333732
          6138643030613631626162633662623562366434626431380a386364643664363137643736646635
          62613131633839356138356565643732373131343137633236616338633335663166313037326434
          3762373039616235320a623332623438316435393561636530316639326132346233346465396364
          3336
Encryption successful
```

3. Запустил `playbook`, убедился, что для нужных хостов применился новый `fact`.
```bash
user@user-VirtualBox:~/ansible_ntl/devops-netology-1/Homework/08-ansible-01-base/playbook$ sudo ansible-playbook -i inventory/prod.yml site.yml --ask-vault-pass
Vault password: 

PLAY [Print os facts] **********************************************************************************************************

TASK [Gathering Facts] *********************************************************************************************************
ok: [localhost]
ok: [ubuntu]
ok: [centos7]

TASK [Print OS] ****************************************************************************************************************
ok: [ubuntu] => {
    "msg": "Ubuntu"
}
ok: [centos7] => {
    "msg": "CentOS"
}
ok: [localhost] => {
    "msg": "Ubuntu"
}

TASK [Print fact] **************************************************************************************************************
ok: [centos7] => {
    "msg": "el default fact"
}
ok: [ubuntu] => {
    "msg": "deb default fact"
}
ok: [localhost] => {
    "msg": "PaSSw0rd"
}

PLAY RECAP *********************************************************************************************************************
centos7                    : ok=3    changed=0    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0   
localhost                  : ok=3    changed=0    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0   
ubuntu                     : ok=3    changed=0    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0 
```

4. Добавьте новую группу хостов `fedora`, самостоятельно придумайте для неё переменную. В качестве образа можно использовать [этот](https://hub.docker.com/r/pycontribs/fedora).
5. Напишите скрипт на bash: автоматизируйте поднятие необходимых контейнеров, запуск ansible-playbook и остановку контейнеров.
6. Все изменения должны быть зафиксированы и отправлены в вашей личный репозиторий.

---

### Как оформить ДЗ?

Выполненное домашнее задание пришлите ссылкой на .md-файл в вашем репозитории.

---
