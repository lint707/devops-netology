# Домашнее задание к занятию "08.05 Тестирование Roles"

## Подготовка к выполнению
1. Установите molecule: `pip3 install "molecule==3.5.2"`
2. Выполните `docker pull aragast/netology:latest` -  это образ с podman, tox и несколькими пайтонами (3.7 и 3.9) внутри

## Основная часть

Наша основная цель - настроить тестирование наших ролей. Задача: сделать сценарии тестирования для vector. Ожидаемый результат: все сценарии успешно проходят тестирование ролей.

### Molecule

1. Запустите  `molecule test -s centos7` внутри корневой директории clickhouse-role, посмотрите на вывод команды.
```yaml
user@user-VirtualBox:~/Desktop/playbook/roles/clickhouse$ molecule test -s centos_8
INFO     centos_8 scenario test matrix: dependency, lint, cleanup, destroy, syntax, create, prepare, converge, idempotence, side_effect, verify, cleanup, destroy
INFO     Performing prerun with role_name_check=0...
INFO     Set ANSIBLE_LIBRARY=/home/user/.cache/ansible-compat/7e099f/modules:/home/user/.ansible/plugins/modules:/usr/share/ansible/plugins/modules
INFO     Set ANSIBLE_COLLECTIONS_PATH=/home/user/.cache/ansible-compat/7e099f/collections:/home/user/.ansible/collections:/usr/share/ansible/collections
INFO     Set ANSIBLE_ROLES_PATH=/home/user/.cache/ansible-compat/7e099f/roles:/home/user/.ansible/roles:/usr/share/ansible/roles:/etc/ansible/roles
INFO     Using /home/user/.cache/ansible-compat/7e099f/roles/alexeysetevoi.clickhouse symlink to current repository in order to enable Ansible to find the role using its expected full name.
INFO     Inventory /home/user/Desktop/playbook/roles/clickhouse/molecule/centos_8/../resources/inventory/hosts.yml linked to /home/user/.cache/molecule/clickhouse/centos_8/inventory/hosts
INFO     Inventory /home/user/Desktop/playbook/roles/clickhouse/molecule/centos_8/../resources/inventory/group_vars/ linked to /home/user/.cache/molecule/clickhouse/centos_8/inventory/group_vars
INFO     Inventory /home/user/Desktop/playbook/roles/clickhouse/molecule/centos_8/../resources/inventory/host_vars/ linked to /home/user/.cache/molecule/clickhouse/centos_8/inventory/host_vars
INFO     Running centos_8 > dependency
WARNING  Skipping, missing the requirements file.
WARNING  Skipping, missing the requirements file.
INFO     Inventory /home/user/Desktop/playbook/roles/clickhouse/molecule/centos_8/../resources/inventory/hosts.yml linked to /home/user/.cache/molecule/clickhouse/centos_8/inventory/hosts
INFO     Inventory /home/user/Desktop/playbook/roles/clickhouse/molecule/centos_8/../resources/inventory/group_vars/ linked to /home/user/.cache/molecule/clickhouse/centos_8/inventory/group_vars
INFO     Inventory /home/user/Desktop/playbook/roles/clickhouse/molecule/centos_8/../resources/inventory/host_vars/ linked to /home/user/.cache/molecule/clickhouse/centos_8/inventory/host_vars
INFO     Running centos_8 > lint
/bin/bash: yamllint: command not found
fatal: not a git repository (or any of the parent directories): .git
Traceback (most recent call last):
  File "/usr/bin/ansible-lint", line 11, in <module>
    load_entry_point('ansible-lint==4.2.0', 'console_scripts', 'ansible-lint')()
  File "/usr/lib/python3/dist-packages/ansiblelint/__main__.py", line 153, in main
    args = get_playbooks_and_roles(options=options)
  File "/usr/lib/python3/dist-packages/ansiblelint/utils.py", line 772, in get_playbooks_and_roles
    files = OrderedDict.fromkeys(sorted(subprocess.check_output(
  File "/usr/lib/python3.8/subprocess.py", line 415, in check_output
    return run(*popenargs, stdout=PIPE, timeout=timeout, check=True,
  File "/usr/lib/python3.8/subprocess.py", line 516, in run
    raise CalledProcessError(retcode, process.args,
subprocess.CalledProcessError: Command '['git', 'ls-files', '*.yaml', '*.yml']' returned non-zero exit status 128.
/bin/bash: line 2: flake8: command not found
WARNING  Retrying execution failure 127 of: y a m l l i n t   . 
 a n s i b l e - l i n t 
 f l a k e 8 

CRITICAL Lint failed with error code 127
WARNING  An error occurred during the test sequence action: 'lint'. Cleaning up.
INFO     Inventory /home/user/Desktop/playbook/roles/clickhouse/molecule/centos_8/../resources/inventory/hosts.yml linked to /home/user/.cache/molecule/clickhouse/centos_8/inventory/hosts
INFO     Inventory /home/user/Desktop/playbook/roles/clickhouse/molecule/centos_8/../resources/inventory/group_vars/ linked to /home/user/.cache/molecule/clickhouse/centos_8/inventory/group_vars
INFO     Inventory /home/user/Desktop/playbook/roles/clickhouse/molecule/centos_8/../resources/inventory/host_vars/ linked to /home/user/.cache/molecule/clickhouse/centos_8/inventory/host_vars
INFO     Running centos_8 > cleanup
WARNING  Skipping, cleanup playbook not configured.
INFO     Inventory /home/user/Desktop/playbook/roles/clickhouse/molecule/centos_8/../resources/inventory/hosts.yml linked to /home/user/.cache/molecule/clickhouse/centos_8/inventory/hosts
INFO     Inventory /home/user/Desktop/playbook/roles/clickhouse/molecule/centos_8/../resources/inventory/group_vars/ linked to /home/user/.cache/molecule/clickhouse/centos_8/inventory/group_vars
INFO     Inventory /home/user/Desktop/playbook/roles/clickhouse/molecule/centos_8/../resources/inventory/host_vars/ linked to /home/user/.cache/molecule/clickhouse/centos_8/inventory/host_vars
INFO     Running centos_8 > destroy
INFO     Sanity checks: 'docker'

PLAY [Destroy] *****************************************************************

TASK [Set async_dir for HOME env] **********************************************
ok: [localhost]

TASK [Destroy molecule instance(s)] ********************************************
changed: [localhost] => (item=centos_8)

TASK [Wait for instance(s) deletion to complete] *******************************
ok: [localhost] => (item=centos_8)

TASK [Delete docker networks(s)] ***********************************************

PLAY RECAP *********************************************************************
localhost                  : ok=3    changed=1    unreachable=0    failed=0    skipped=1    rescued=0    ignored=0

INFO     Pruning extra files from scenario ephemeral directory
```
2. Перейдите в каталог с ролью vector-role и создайте сценарий тестирования по умолчанию при помощи `molecule init scenario --driver-name docker`.
> done

3. Добавьте несколько разных дистрибутивов (centos:8, ubuntu:latest) для инстансов и протестируйте роль, исправьте найденные ошибки, если они есть.
```yaml
user@user-VirtualBox:~/Desktop/playbook/roles/vector-role$ molecule test
INFO     default scenario test matrix: dependency, lint, cleanup, destroy, syntax, create, prepare, converge, idempotence, side_effect, verify, cleanup, destroy
INFO     Performing prerun with role_name_check=0...
INFO     Set ANSIBLE_LIBRARY=/home/user/.cache/ansible-compat/f5bcd7/modules:/home/user/.ansible/plugins/modules:/usr/share/ansible/plugins/modules
INFO     Set ANSIBLE_COLLECTIONS_PATH=/home/user/.cache/ansible-compat/f5bcd7/collections:/home/user/.ansible/collections:/usr/share/ansible/collections
INFO     Set ANSIBLE_ROLES_PATH=/home/user/.cache/ansible-compat/f5bcd7/roles:/home/user/.ansible/roles:/usr/share/ansible/roles:/etc/ansible/roles
ERROR    Computed fully qualified role name of vector-role does not follow current galaxy requirements.
Please edit meta/main.yml and assure we can correctly determine full role name:

galaxy_info:
role_name: my_name  # if absent directory name hosting role is used instead
namespace: my_galaxy_namespace  # if absent, author is used instead

Namespace: https://galaxy.ansible.com/docs/contributing/namespaces.html#galaxy-namespace-limitations
Role: https://galaxy.ansible.com/docs/contributing/creating_role.html#role-names

As an alternative, you can add 'role-name' to either skip_list or warn_list.

Traceback (most recent call last):
  File "/home/user/.local/bin/molecule", line 8, in <module>
    sys.exit(main())
  File "/home/user/.local/lib/python3.8/site-packages/click/core.py", line 1130, in __call__
    return self.main(*args, **kwargs)
  File "/home/user/.local/lib/python3.8/site-packages/click/core.py", line 1055, in main
    rv = self.invoke(ctx)
  File "/home/user/.local/lib/python3.8/site-packages/click/core.py", line 1657, in invoke
    return _process_result(sub_ctx.command.invoke(sub_ctx))
  File "/home/user/.local/lib/python3.8/site-packages/click/core.py", line 1404, in invoke
    return ctx.invoke(self.callback, **ctx.params)
  File "/home/user/.local/lib/python3.8/site-packages/click/core.py", line 760, in invoke
    return __callback(*args, **kwargs)
  File "/home/user/.local/lib/python3.8/site-packages/click/decorators.py", line 26, in new_func
    return f(get_current_context(), *args, **kwargs)
  File "/home/user/.local/lib/python3.8/site-packages/molecule/command/test.py", line 176, in test
    base.execute_cmdline_scenarios(scenario_name, args, command_args, ansible_args)
  File "/home/user/.local/lib/python3.8/site-packages/molecule/command/base.py", line 112, in execute_cmdline_scenarios
    scenario.config.runtime.prepare_environment(
  File "/home/user/.local/lib/python3.8/site-packages/ansible_compat/runtime.py", line 380, in prepare_environment
    self._install_galaxy_role(
  File "/home/user/.local/lib/python3.8/site-packages/ansible_compat/runtime.py", line 540, in _install_galaxy_role
    raise InvalidPrerequisiteError(msg)
ansible_compat.errors.InvalidPrerequisiteError: Computed fully qualified role name of vector-role does not follow current galaxy requirements.
Please edit meta/main.yml and assure we can correctly determine full role name:

galaxy_info:
role_name: my_name  # if absent directory name hosting role is used instead
namespace: my_galaxy_namespace  # if absent, author is used instead

Namespace: https://galaxy.ansible.com/docs/contributing/namespaces.html#galaxy-namespace-limitations
Role: https://galaxy.ansible.com/docs/contributing/creating_role.html#role-names

As an alternative, you can add 'role-name' to either skip_list or warn_list.
```
>after   role_name: vector

```yaml
user@user-VirtualBox:~/Desktop/playbook/roles/vector-role$ molecule test
INFO     default scenario test matrix: dependency, lint, cleanup, destroy, syntax, create, prepare, converge, idempotence, side_effect, verify, cleanup, destroy
INFO     Performing prerun with role_name_check=0...
INFO     Set ANSIBLE_LIBRARY=/home/user/.cache/ansible-compat/f5bcd7/modules:/home/user/.ansible/plugins/modules:/usr/share/ansible/plugins/modules
INFO     Set ANSIBLE_COLLECTIONS_PATH=/home/user/.cache/ansible-compat/f5bcd7/collections:/home/user/.ansible/collections:/usr/share/ansible/collections
INFO     Set ANSIBLE_ROLES_PATH=/home/user/.cache/ansible-compat/f5bcd7/roles:/home/user/.ansible/roles:/usr/share/ansible/roles:/etc/ansible/roles
INFO     Using /home/user/.cache/ansible-compat/f5bcd7/roles/vector_namespace.vector symlink to current repository in order to enable Ansible to find the role using its expected full name.
INFO     Running default > dependency
WARNING  Skipping, missing the requirements file.
WARNING  Skipping, missing the requirements file.
INFO     Running default > lint
INFO     Lint is disabled.
INFO     Running default > cleanup
WARNING  Skipping, cleanup playbook not configured.
INFO     Running default > destroy
INFO     Sanity checks: 'docker'

PLAY [Destroy] *****************************************************************

TASK [Set async_dir for HOME env] **********************************************
ok: [localhost]

TASK [Destroy molecule instance(s)] ********************************************
changed: [localhost] => (item=instance)
changed: [localhost] => (item=instance_ub)

TASK [Wait for instance(s) deletion to complete] *******************************
ok: [localhost] => (item=instance)
ok: [localhost] => (item=instance_ub)

TASK [Delete docker networks(s)] ***********************************************

PLAY RECAP *********************************************************************
localhost                  : ok=3    changed=1    unreachable=0    failed=0    skipped=1    rescued=0    ignored=0

INFO     Running default > syntax

playbook: /home/user/Desktop/playbook/roles/vector-role/molecule/default/converge.yml
INFO     Running default > create

PLAY [Create] ******************************************************************

TASK [Set async_dir for HOME env] **********************************************
ok: [localhost]

TASK [Log into a Docker registry] **********************************************
skipping: [localhost] => (item=None) 
skipping: [localhost] => (item=None) 
skipping: [localhost]

TASK [Check presence of custom Dockerfiles] ************************************
ok: [localhost] => (item={'image': 'quay.io/centos/centos:stream8', 'name': 'instance', 'pre_build_image': True})
ok: [localhost] => (item={'image': 'pycontribs/ubuntu:latest', 'name': 'instance_ub', 'pre_build_image': True})

TASK [Create Dockerfiles from image names] *************************************
skipping: [localhost] => (item={'image': 'quay.io/centos/centos:stream8', 'name': 'instance', 'pre_build_image': True})
skipping: [localhost] => (item={'image': 'pycontribs/ubuntu:latest', 'name': 'instance_ub', 'pre_build_image': True})

TASK [Synchronization the context] *********************************************
skipping: [localhost] => (item={'image': 'quay.io/centos/centos:stream8', 'name': 'instance', 'pre_build_image': True})
skipping: [localhost] => (item={'image': 'pycontribs/ubuntu:latest', 'name': 'instance_ub', 'pre_build_image': True})

TASK [Discover local Docker images] ********************************************
ok: [localhost] => (item={'changed': False, 'skipped': True, 'skip_reason': 'Conditional result was False', 'item': {'image': 'quay.io/centos/centos:stream8', 'name': 'instance', 'pre_build_image': True}, 'ansible_loop_var': 'item', 'i': 0, 'ansible_index_var': 'i'})
ok: [localhost] => (item={'changed': False, 'skipped': True, 'skip_reason': 'Conditional result was False', 'item': {'image': 'pycontribs/ubuntu:latest', 'name': 'instance_ub', 'pre_build_image': True}, 'ansible_loop_var': 'item', 'i': 1, 'ansible_index_var': 'i'})

TASK [Build an Ansible compatible image (new)] *********************************
skipping: [localhost] => (item=molecule_local/quay.io/centos/centos:stream8) 
skipping: [localhost] => (item=molecule_local/pycontribs/ubuntu:latest) 

TASK [Create docker network(s)] ************************************************

TASK [Determine the CMD directives] ********************************************
ok: [localhost] => (item={'image': 'quay.io/centos/centos:stream8', 'name': 'instance', 'pre_build_image': True})
ok: [localhost] => (item={'image': 'pycontribs/ubuntu:latest', 'name': 'instance_ub', 'pre_build_image': True})

TASK [Create molecule instance(s)] *********************************************
changed: [localhost] => (item=instance)
changed: [localhost] => (item=instance_ub)

TASK [Wait for instance(s) creation to complete] *******************************
FAILED - RETRYING: [localhost]: Wait for instance(s) creation to complete (300 retries left).
changed: [localhost] => (item={'failed': 0, 'started': 1, 'finished': 0, 'ansible_job_id': '112610383379.12481', 'results_file': '/home/user/.ansible_async/112610383379.12481', 'changed': True, 'item': {'image': 'quay.io/centos/centos:stream8', 'name': 'instance', 'pre_build_image': True}, 'ansible_loop_var': 'item'})
changed: [localhost] => (item={'failed': 0, 'started': 1, 'finished': 0, 'ansible_job_id': '167012885025.12509', 'results_file': '/home/user/.ansible_async/167012885025.12509', 'changed': True, 'item': {'image': 'pycontribs/ubuntu:latest', 'name': 'instance_ub', 'pre_build_image': True}, 'ansible_loop_var': 'item'})

PLAY RECAP *********************************************************************
localhost                  : ok=6    changed=2    unreachable=0    failed=0    skipped=5    rescued=0    ignored=0

INFO     Running default > prepare
WARNING  Skipping, prepare playbook not configured.
INFO     Running default > converge

PLAY [Converge] ****************************************************************

TASK [Gathering Facts] *********************************************************
ok: [instance]
ok: [instance_ub]

TASK [Include vector-role] *****************************************************

TASK [vector-role : Download Vector] *******************************************
changed: [instance_ub]
changed: [instance]

TASK [vector-role : Install Vector] ********************************************
fatal: [instance]: FAILED! => {"changed": false, "module_stderr": "/bin/sh: sudo: command not found\n", "module_stdout": "", "msg": "MODULE FAILURE\nSee stdout/stderr for the exact error", "rc": 127}
fatal: [instance_ub]: FAILED! => {"ansible_facts": {"pkg_mgr": "apt"}, "changed": false, "msg": ["Could not detect which major revision of yum is in use, which is required to determine module backend.", "You should manually specify use_backend to tell the module whether to use the yum (yum3) or dnf (yum4) backend})"]}

PLAY RECAP *********************************************************************
instance                   : ok=2    changed=1    unreachable=0    failed=1    skipped=0    rescued=0    ignored=0
instance_ub                : ok=2    changed=1    unreachable=0    failed=1    skipped=0    rescued=0    ignored=0

WARNING  Retrying execution failure 2 of: ansible-playbook --inventory /home/user/.cache/molecule/vector-role/default/inventory --skip-tags molecule-notest,notest /home/user/Desktop/playbook/roles/vector-role/molecule/default/converge.yml
CRITICAL Ansible return code was 2, command was: ['ansible-playbook', '--inventory', '/home/user/.cache/molecule/vector-role/default/inventory', '--skip-tags', 'molecule-notest,notest', '/home/user/Desktop/playbook/roles/vector-role/molecule/default/converge.yml']
WARNING  An error occurred during the test sequence action: 'converge'. Cleaning up.
INFO     Running default > cleanup
WARNING  Skipping, cleanup playbook not configured.
INFO     Running default > destroy

PLAY [Destroy] *****************************************************************

TASK [Set async_dir for HOME env] **********************************************
ok: [localhost]

TASK [Destroy molecule instance(s)] ********************************************
changed: [localhost] => (item=instance)
changed: [localhost] => (item=instance_ub)

TASK [Wait for instance(s) deletion to complete] *******************************
FAILED - RETRYING: [localhost]: Wait for instance(s) deletion to complete (300 retries left).
changed: [localhost] => (item=instance)
changed: [localhost] => (item=instance_ub)

TASK [Delete docker networks(s)] ***********************************************

PLAY RECAP *********************************************************************
localhost                  : ok=3    changed=2    unreachable=0    failed=0    skipped=1    rescued=0    ignored=0

INFO     Pruning extra files from scenario ephemeral directory

```

4. Добавьте несколько assert'ов в verify.yml файл для  проверки работоспособности vector-role (проверка, что конфиг валидный, проверка успешности запуска, etc). Запустите тестирование роли повторно и проверьте, что оно прошло успешно.

5. Добавьте новый тег на коммит с рабочим сценарием в соответствии с семантическим версионированием.

### Tox

1. Добавьте в директорию с vector-role файлы из [директории](./example)
2. Запустите `docker run --privileged=True -v <path_to_repo>:/opt/vector-role -w /opt/vector-role -it aragast/netology:latest /bin/bash`, где path_to_repo - путь до корня репозитория с vector-role на вашей файловой системе.
3. Внутри контейнера выполните команду `tox`, посмотрите на вывод.
5. Создайте облегчённый сценарий для `molecule` с драйвером `molecule_podman`. Проверьте его на исполнимость.
6. Пропишите правильную команду в `tox.ini` для того чтобы запускался облегчённый сценарий.
8. Запустите команду `tox`. Убедитесь, что всё отработало успешно.
9. Добавьте новый тег на коммит с рабочим сценарием в соответствии с семантическим версионированием.

После выполнения у вас должно получится два сценария molecule и один tox.ini файл в репозитории. Ссылка на репозиторий являются ответами на домашнее задание. Не забудьте указать в ответе теги решений Tox и Molecule заданий.

## Необязательная часть

1. Проделайте схожие манипуляции для создания роли lighthouse.
2. Создайте сценарий внутри любой из своих ролей, который умеет поднимать весь стек при помощи всех ролей.
3. Убедитесь в работоспособности своего стека. Создайте отдельный verify.yml, который будет проверять работоспособность интеграции всех инструментов между ними.
4. Выложите свои roles в репозитории. В ответ приведите ссылки.

---

### Как оформить ДЗ?

Выполненное домашнее задание пришлите ссылкой на .md-файл в вашем репозитории.

---
