# Домашнее задание к занятию "08.02 Работа с Playbook"

## Подготовка к выполнению

1. (Необязательно) Изучите, что такое [clickhouse](https://www.youtube.com/watch?v=fjTNS2zkeBs) и [vector](https://www.youtube.com/watch?v=CgEhyffisLY)
2. Создайте свой собственный (или используйте старый) публичный репозиторий на github с произвольным именем.
3. Скачайте [playbook](./playbook/) из репозитория с домашним заданием и перенесите его в свой репозиторий.
4. Подготовьте хосты в соответствии с группами из предподготовленного playbook.

## Основная часть

1. Приготовил свой собственный inventory файл `prod.yml`.

```yaml
---
clickhouse:
  hosts:
    centos7:
      ansible_connection: ssh
      ansible_host: 178.154.223.132

      
vector:
  hosts:
    centos7:
      ansible_connection: ssh
      ansible_host: 51.250.68.146
```

2. Дописал playbook: сделать ещё один play, который устанавливает и настраивает [vector](https://vector.dev).

```yaml
- name: Install Vector
  hosts: vector
  tasks:
    - name: Download Vector
      ansible.builtin.get_url:
        url: "https://packages.timber.io/vector/{{ vector_version }}/vector-{{ vector_version }}-1.x86_64.rpm"
        dest: "./vector-{{ vector_version }}-1.x86_64.rpm"
    - name: Install Vector
      become: true
      ansible.builtin.yum:
        name: "vector-{{ vector_version }}-1.x86_64.rpm"

```

3. При создании tasks рекомендую использовать модули: `get_url`, `template`, `unarchive`, `file`.
> Использовал - `get_url`

4. Tasks должны: скачать нужной версии дистрибутив, выполнить распаковку в выбранную директорию, установить vector.
> Установил вектор.

5. Запустил `ansible-lint site.yml` и исправил ошибки, если они есть.
```yaml
user@user-VirtualBox:~/Desktop/08-ansible-02-playbook/playbook$ ansible-lint site.yml
[201] Trailing whitespace
site.yml:31
      meta: flush_handlers 
```
> Исправил ошибку с лишним пробелом.

6. Запустил playbook на этом окружении с флагом `--check`.
```yaml
user@user-VirtualBox:~/Desktop/08-ansible-02-playbook/playbook$ ansible-playbook -i inventory/prod.yml site.yml --check

PLAY [Install Clickhouse] **********************************************************************************************************

TASK [Gathering Facts] *************************************************************************************************************
ok: [centos7]

TASK [Get clickhouse distrib] ******************************************************************************************************
ok: [centos7] => (item=clickhouse-client)
ok: [centos7] => (item=clickhouse-server)
failed: [centos7] (item=clickhouse-common-static) => {"ansible_loop_var": "item", "changed": false, "dest": "./clickhouse-common-static-22.3.3.44.rpm", "elapsed": 0, "gid": 1000, "group": "user", "item": "clickhouse-common-static", "mode": "0664", "msg": "Request failed", "owner": "user", "response": "HTTP Error 404: Not Found", "secontext": "unconfined_u:object_r:user_home_t:s0", "size": 246310036, "state": "file", "status_code": 404, "uid": 1000, "url": "https://packages.clickhouse.com/rpm/stable/clickhouse-common-static-22.3.3.44.noarch.rpm"}

TASK [Get clickhouse distrib] ******************************************************************************************************
ok: [centos7]

TASK [Install clickhouse packages] *************************************************************************************************
ok: [centos7]

TASK [Flush handlers] **************************************************************************************************************

TASK [Create database] *************************************************************************************************************
skipping: [centos7]

PLAY [Install Vector] **************************************************************************************************************

TASK [Gathering Facts] *************************************************************************************************************
ok: [centos7]

TASK [Download Vector] *************************************************************************************************************
ok: [centos7]

TASK [Install Vector] **************************************************************************************************************
ok: [centos7]

PLAY RECAP *************************************************************************************************************************
centos7                    : ok=6    changed=0    unreachable=0    failed=0    skipped=1    rescued=1    ignored=0   

```

7. Запустил playbook на `prod.yml` окружении с флагом `--diff`. Убедился, что изменения на системе произведены.
```yaml
user@user-VirtualBox:~/Desktop/08-ansible-02-playbook/playbook$ ansible-playbook -i inventory/prod.yml site.yml --diff

PLAY [Install Clickhouse] **********************************************************************************************************

TASK [Gathering Facts] *************************************************************************************************************
ok: [centos7]

TASK [Get clickhouse distrib] ******************************************************************************************************
ok: [centos7] => (item=clickhouse-client)
ok: [centos7] => (item=clickhouse-server)
failed: [centos7] (item=clickhouse-common-static) => {"ansible_loop_var": "item", "changed": false, "dest": "./clickhouse-common-static-22.3.3.44.rpm", "elapsed": 0, "gid": 1000, "group": "user", "item": "clickhouse-common-static", "mode": "0664", "msg": "Request failed", "owner": "user", "response": "HTTP Error 404: Not Found", "secontext": "unconfined_u:object_r:user_home_t:s0", "size": 246310036, "state": "file", "status_code": 404, "uid": 1000, "url": "https://packages.clickhouse.com/rpm/stable/clickhouse-common-static-22.3.3.44.noarch.rpm"}

TASK [Get clickhouse distrib] ******************************************************************************************************
ok: [centos7]

TASK [Install clickhouse packages] *************************************************************************************************
ok: [centos7]

TASK [Flush handlers] **************************************************************************************************************

TASK [Create database] *************************************************************************************************************
ok: [centos7]

PLAY [Install Vector] **************************************************************************************************************

TASK [Gathering Facts] *************************************************************************************************************
ok: [centos7]

TASK [Download Vector] *************************************************************************************************************
ok: [centos7]

TASK [Install Vector] **************************************************************************************************************
ok: [centos7]

PLAY RECAP *************************************************************************************************************************
centos7                    : ok=7    changed=0    unreachable=0    failed=0    skipped=0    rescued=1    ignored=0   

```

8. Повторно запустил playbook с флагом `--diff` и убедился, что playbook идемпотентен.
```yaml
user@user-VirtualBox:~/Desktop/08-ansible-02-playbook/playbook$ ansible-playbook -i inventory/prod.yml site.yml --diff

PLAY [Install Clickhouse] **********************************************************************************************************

TASK [Gathering Facts] *************************************************************************************************************
ok: [centos7]

TASK [Get clickhouse distrib] ******************************************************************************************************
ok: [centos7] => (item=clickhouse-client)
ok: [centos7] => (item=clickhouse-server)
failed: [centos7] (item=clickhouse-common-static) => {"ansible_loop_var": "item", "changed": false, "dest": "./clickhouse-common-static-22.3.3.44.rpm", "elapsed": 0, "gid": 1000, "group": "user", "item": "clickhouse-common-static", "mode": "0664", "msg": "Request failed", "owner": "user", "response": "HTTP Error 404: Not Found", "secontext": "unconfined_u:object_r:user_home_t:s0", "size": 246310036, "state": "file", "status_code": 404, "uid": 1000, "url": "https://packages.clickhouse.com/rpm/stable/clickhouse-common-static-22.3.3.44.noarch.rpm"}

TASK [Get clickhouse distrib] ******************************************************************************************************
ok: [centos7]

TASK [Install clickhouse packages] *************************************************************************************************
ok: [centos7]

TASK [Flush handlers] **************************************************************************************************************

TASK [Create database] *************************************************************************************************************
ok: [centos7]

PLAY [Install Vector] **************************************************************************************************************

TASK [Gathering Facts] *************************************************************************************************************
ok: [centos7]

TASK [Download Vector] *************************************************************************************************************
ok: [centos7]

TASK [Install Vector] **************************************************************************************************************
ok: [centos7]

PLAY RECAP *************************************************************************************************************************
centos7                    : ok=7    changed=0    unreachable=0    failed=0    skipped=0    rescued=1    ignored=0   

```

9. Подготовил README.md файл по своему playbook. В нём должно быть описано: что делает playbook, какие у него есть параметры и теги.
10. Готовый playbook выложил в свой репозиторий, поставил тег `08-ansible-02-playbook` на фиксирующий коммит.

---

### Как оформить ДЗ?

Выполненное домашнее задание пришлите ссылкой на .md-файл в вашем репозитории.
---
