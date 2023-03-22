# Домашнее задание к занятию "08.03 Использование Yandex Cloud"

## Подготовка к выполнению

1. (Необязательно) Познакомтесь с [lighthouse](https://youtu.be/ymlrNlaHzIY?t=929)
2. Подготовьте в Yandex Cloud три хоста: для `clickhouse`, для `vector` и для `lighthouse`.

Ссылка на репозиторий LightHouse: https://github.com/VKCOM/lighthouse

## Основная часть
1. Дописал playbook: ещё один play, который устанавливает и настраивает lighthouse.
```yaml
- name: Install Light
  hosts: lighthouse
  handlers:
    - name: reload-nginx
      become: true
      command: nginx -s reload
  tasks:
    - name: Download Nginx
      become: true
      ansible.builtin.get_url:
        url: "https://nginx.org/packages/centos/7/x86_64/RPMS/nginx-1.22.0-1.el7.ngx.x86_64.rpm"
        dest: "./nginx-1.22.0-1.el7.rpm"
    - name: Install Nginx
      become: true
      ansible.builtin.yum:
        name: "nginx-1.22.0-1.el7.rpm"
    - name: NGINX | Create general config
      become: true
      template:
        src: templates/nginx.conf.j2
        dest: /etc/nginx/nginx.conf
        mode: 0755
      notify: reload-nginx
    - name: Lighthouse | Create lh config
      become: true
      template:
        src: "templates/default.conf.j2"
        dest: /etc/nginx/conf.d/default.conf
        mode: 0644
      notify: reload-nginx
      become: true
    - name: Lighthouse | Copy git
      git: 
        repo: "{{ lh_vcs }}"
        version: master
        dest: "{{ lh_loc_dir }}"
      register: copy_remote
```

2. При создании tasks использовал модули: `get_url`, `template`, `yum`.
3. Tasks должны: скачать статику lighthouse, установить nginx или любой другой webserver, настроить его конфиг для открытия lighthouse, запустить webserver.

4. Приготовил свой собственный inventory файл `prod.yml`.
```yaml
---
all:
  children:
    clickhouse:
      hosts:
        centos7:
          ansible_connection: ssh
          ansible_host: 84.201.157.123
    vector:
      hosts:
        centos7:
          ansible_connection: ssh
          ansible_host: 178.154.223.135
    lighthouse:
      hosts:
        centos7:
          ansible_connection: ssh
          ansible_host: 178.154.223.254
```
<<<<<<< HEAD
user@user-VirtualBox:~/Desktop/playbook$ ansible-lint site.yml
[WARNING]: While constructing a mapping from /home/user/Desktop/playbook/site.yml, line 88, column 7, found a duplicate
=======

5. Запустил `ansible-lint site.yml` и исправьте ошибки, если они есть.
```yaml
user@user-VirtualBox:~/Desktop/playbook$ ansible-lint site.yml
[WARNING]: While constructing a mapping from /home/user/Desktop/playbook/site.yml, line 88, column 7, found a duplicate
>>>>>>> c924229ca5c219a566b34ff3721b0b29195fe393
dict key (become). Using last defined value only.
[WARNING]: While constructing a mapping from <unicode string>, line 88, column 7, found a duplicate dict key (become). Using
last defined value only.
Traceback (most recent call last):
  File "/usr/bin/ansible-lint", line 11, in <module>
    load_entry_point('ansible-lint==4.2.0', 'console_scripts', 'ansible-lint')()
  File "/usr/lib/python3/dist-packages/ansiblelint/__main__.py", line 187, in main
    matches.extend(runner.run())
  File "/usr/lib/python3/dist-packages/ansiblelint/__init__.py", line 286, in run
    matches.extend(self.rules.run(file, tags=set(self.tags),
  File "/usr/lib/python3/dist-packages/ansiblelint/__init__.py", line 177, in run
    matches.extend(rule.matchtasks(playbookfile, text))
  File "/usr/lib/python3/dist-packages/ansiblelint/__init__.py", line 87, in matchtasks
    yaml = ansiblelint.utils.append_skipped_rules(yaml, text, file['type'])
  File "/usr/lib/python3/dist-packages/ansiblelint/utils.py", line 596, in append_skipped_rules
    yaml_skip = _append_skipped_rules(pyyaml_data, file_text, file_type)
  File "/usr/lib/python3/dist-packages/ansiblelint/utils.py", line 607, in _append_skipped_rules
    ruamel_data = yaml.load(file_text)
  File "/usr/lib/python3/dist-packages/ruamel/yaml/main.py", line 331, in load
    return constructor.get_single_data()
  File "/usr/lib/python3/dist-packages/ruamel/yaml/constructor.py", line 111, in get_single_data
    return self.construct_document(node)
  File "/usr/lib/python3/dist-packages/ruamel/yaml/constructor.py", line 121, in construct_document
    for _dummy in generator:
  File "/usr/lib/python3/dist-packages/ruamel/yaml/constructor.py", line 1543, in construct_yaml_map
    self.construct_mapping(node, data, deep=True)
  File "/usr/lib/python3/dist-packages/ruamel/yaml/constructor.py", line 1448, in construct_mapping
    value = self.construct_object(value_node, deep=deep)
  File "/usr/lib/python3/dist-packages/ruamel/yaml/constructor.py", line 174, in construct_object
    for _dummy in generator:
  File "/usr/lib/python3/dist-packages/ruamel/yaml/constructor.py", line 1535, in construct_yaml_seq
    data.extend(self.construct_rt_sequence(node, data))
  File "/usr/lib/python3/dist-packages/ruamel/yaml/constructor.py", line 1298, in construct_rt_sequence
    ret_val.append(self.construct_object(child, deep=deep))
  File "/usr/lib/python3/dist-packages/ruamel/yaml/constructor.py", line 174, in construct_object
    for _dummy in generator:
  File "/usr/lib/python3/dist-packages/ruamel/yaml/constructor.py", line 1543, in construct_yaml_map
    self.construct_mapping(node, data, deep=True)
  File "/usr/lib/python3/dist-packages/ruamel/yaml/constructor.py", line 1449, in construct_mapping
    if self.check_mapping_key(node, key_node, maptyp, key, value):
  File "/usr/lib/python3/dist-packages/ruamel/yaml/constructor.py", line 285, in check_mapping_key
    raise DuplicateKeyError(*args)
ruamel.yaml.constructor.DuplicateKeyError: while constructing a mapping
  in "<unicode string>", line 88, column 7:
        - name: Lighthouse | Create lh config
          ^ (line: 88)
found duplicate key "become" with value "True" (original value: "True")
  in "<unicode string>", line 96, column 7:
          become: true
          ^ (line: 96)

To suppress this check see:
    http://yaml.readthedocs.io/en/latest/api.html#duplicate-keys

Duplicate keys will become an error in future releases, and are errors
by default when using the new API.

```
<<<<<<< HEAD
6. Попробуйте запустить playbook на этом окружении с флагом `--check`.
```
user@user-VirtualBox:~/Desktop/playbook$ ansible-playbook -i inventory/prod.yml site.yml --check
=======

6. Запустил playbook на этом окружении с флагом `--check`.
```yaml
user@user-VirtualBox:~/Desktop/playbook$ ansible-playbook -i inventory/prod.yml site.yml --check
>>>>>>> c924229ca5c219a566b34ff3721b0b29195fe393

PLAY [Install Clickhouse] ***************************************************************************************************

TASK [Gathering Facts] ******************************************************************************************************
ok: [centos7]

TASK [Get clickhouse distrib] ***********************************************************************************************
ok: [centos7] => (item=clickhouse-client)
ok: [centos7] => (item=clickhouse-server)
failed: [centos7] (item=clickhouse-common-static) => {"ansible_loop_var": "item", "changed": false, "dest": "./clickhouse-common-static-22.3.3.44.rpm", "elapsed": 0, "gid": 1000, "group": "user", "item": "clickhouse-common-static", "mode": "0664", "msg": "Request failed", "owner": "user", "response": "HTTP Error 404: Not Found", "secontext": "unconfined_u:object_r:user_home_t:s0", "size": 246310036, "state": "file", "status_code": 404, "uid": 1000, "url": "https://packages.clickhouse.com/rpm/stable/clickhouse-common-static-22.3.3.44.noarch.rpm"}

TASK [Get clickhouse distrib] ***********************************************************************************************
ok: [centos7]

TASK [Install clickhouse packages] ******************************************************************************************
ok: [centos7]

TASK [Flush handlers] *******************************************************************************************************

TASK [Create database] ******************************************************************************************************
skipping: [centos7]

PLAY [Install Vector] *******************************************************************************************************

TASK [Gathering Facts] ******************************************************************************************************
ok: [centos7]

TASK [Download Vector] ******************************************************************************************************
ok: [centos7]

TASK [Install Vector] *******************************************************************************************************
ok: [centos7]

PLAY [Install Light] ********************************************************************************************************

TASK [Gathering Facts] ******************************************************************************************************
ok: [centos7]

TASK [Download Nginx] *******************************************************************************************************
changed: [centos7]

TASK [Install Nginx] ********************************************************************************************************
ok: [centos7]

TASK [NGINX | Create general config] ****************************************************************************************
ok: [centos7]

TASK [Lighthouse - Create lh config] ****************************************************************************************
ok: [centos7]

TASK [Lighthouse - Copy git] ************************************************************************************************
ok: [centos7]

PLAY RECAP ******************************************************************************************************************
centos7                    : ok=12   changed=1    unreachable=0    failed=0    skipped=1    rescued=1    ignored=0   
```

<<<<<<< HEAD
7. Запустите playbook на `prod.yml` окружении с флагом `--diff`. Убедитесь, что изменения на системе произведены.
```
user@user-VirtualBox:~/Desktop/playbook$ ansible-playbook -i inventory/prod.yml site.yml --diff
=======
7. Запустил playbook на `prod.yml` окружении с флагом `--diff`. Убедился, что изменения на системе произведены.
```yaml
user@user-VirtualBox:~/Desktop/playbook$ ansible-playbook -i inventory/prod.yml site.yml --diff
>>>>>>> c924229ca5c219a566b34ff3721b0b29195fe393

PLAY [Install Clickhouse] ***************************************************************************************************

TASK [Gathering Facts] ******************************************************************************************************
ok: [centos7]

TASK [Get clickhouse distrib] ***********************************************************************************************
ok: [centos7] => (item=clickhouse-client)
ok: [centos7] => (item=clickhouse-server)
failed: [centos7] (item=clickhouse-common-static) => {"ansible_loop_var": "item", "changed": false, "dest": "./clickhouse-common-static-22.3.3.44.rpm", "elapsed": 0, "gid": 1000, "group": "user", "item": "clickhouse-common-static", "mode": "0664", "msg": "Request failed", "owner": "user", "response": "HTTP Error 404: Not Found", "secontext": "unconfined_u:object_r:user_home_t:s0", "size": 246310036, "state": "file", "status_code": 404, "uid": 1000, "url": "https://packages.clickhouse.com/rpm/stable/clickhouse-common-static-22.3.3.44.noarch.rpm"}

TASK [Get clickhouse distrib] ***********************************************************************************************
ok: [centos7]

TASK [Install clickhouse packages] ******************************************************************************************
ok: [centos7]

TASK [Flush handlers] *******************************************************************************************************

TASK [Create database] ******************************************************************************************************
ok: [centos7]

PLAY [Install Vector] *******************************************************************************************************

TASK [Gathering Facts] ******************************************************************************************************
ok: [centos7]

TASK [Download Vector] ******************************************************************************************************
ok: [centos7]

TASK [Install Vector] *******************************************************************************************************
ok: [centos7]

PLAY [Install Light] ********************************************************************************************************

TASK [Gathering Facts] ******************************************************************************************************
ok: [centos7]

TASK [Download Nginx] *******************************************************************************************************
ok: [centos7]

TASK [Install Nginx] ********************************************************************************************************
ok: [centos7]

TASK [NGINX | Create general config] ****************************************************************************************
ok: [centos7]

TASK [Lighthouse - Create lh config] ****************************************************************************************
ok: [centos7]

TASK [Lighthouse - Copy git] ************************************************************************************************
ok: [centos7]

PLAY RECAP ******************************************************************************************************************
centos7                    : ok=13   changed=0    unreachable=0    failed=0    skipped=0    rescued=1    ignored=0   
```
<<<<<<< HEAD
user@user-VirtualBox:~/Desktop/playbook$ ansible-playbook -i inventory/prod.yml site.yml --diff
=======

8. Повторно запустил playbook с флагом `--diff` и убедился, что playbook идемпотентен.
```yaml
user@user-VirtualBox:~/Desktop/playbook$ ansible-playbook -i inventory/prod.yml site.yml --diff
>>>>>>> c924229ca5c219a566b34ff3721b0b29195fe393

PLAY [Install Clickhouse] ***************************************************************************************************

TASK [Gathering Facts] ******************************************************************************************************
ok: [centos7]

TASK [Get clickhouse distrib] ***********************************************************************************************
ok: [centos7] => (item=clickhouse-client)
ok: [centos7] => (item=clickhouse-server)
failed: [centos7] (item=clickhouse-common-static) => {"ansible_loop_var": "item", "changed": false, "dest": "./clickhouse-common-static-22.3.3.44.rpm", "elapsed": 0, "gid": 1000, "group": "user", "item": "clickhouse-common-static", "mode": "0664", "msg": "Request failed", "owner": "user", "response": "HTTP Error 404: Not Found", "secontext": "unconfined_u:object_r:user_home_t:s0", "size": 246310036, "state": "file", "status_code": 404, "uid": 1000, "url": "https://packages.clickhouse.com/rpm/stable/clickhouse-common-static-22.3.3.44.noarch.rpm"}

TASK [Get clickhouse distrib] ***********************************************************************************************
ok: [centos7]

TASK [Install clickhouse packages] ******************************************************************************************
ok: [centos7]

TASK [Flush handlers] *******************************************************************************************************

TASK [Create database] ******************************************************************************************************
ok: [centos7]

PLAY [Install Vector] *******************************************************************************************************

TASK [Gathering Facts] ******************************************************************************************************
ok: [centos7]

TASK [Download Vector] ******************************************************************************************************
ok: [centos7]

TASK [Install Vector] *******************************************************************************************************
ok: [centos7]

PLAY [Install Light] ********************************************************************************************************

TASK [Gathering Facts] ******************************************************************************************************
ok: [centos7]

TASK [Download Nginx] *******************************************************************************************************
ok: [centos7]

TASK [Install Nginx] ********************************************************************************************************
ok: [centos7]

TASK [NGINX | Create general config] ****************************************************************************************
ok: [centos7]

TASK [Lighthouse - Create lh config] ****************************************************************************************
ok: [centos7]

TASK [Lighthouse - Copy git] ************************************************************************************************
ok: [centos7]

PLAY RECAP ******************************************************************************************************************
centos7                    : ok=13   changed=0    unreachable=0    failed=0    skipped=0    rescued=1    ignored=0   
```

9. Подготовил README.md файл по своему playbook. В нём должно быть описано: что делает playbook, какие у него есть параметры и теги.
10. Готовый playbook выложите в свой репозиторий, поставьте тег `08-ansible-03-yandex` на фиксирующий коммит, в ответ предоставьте ссылку на него.

---

### Как оформить ДЗ?

Выполненное домашнее задание пришлите ссылкой на .md-файл в вашем репозитории.

---
