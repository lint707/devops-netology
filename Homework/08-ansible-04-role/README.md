# Домашнее задание к занятию "8.4 Работа с Roles"

## Подготовка к выполнению
1. Создал три пустых публичных репозитория в любом своём проекте: vector-role, nginx-role и lighthouse-role.
```
https://github.com/lint707/lighthouse-role
https://github.com/lint707/nginx-role
https://github.com/lint707/vector-role
```
2. Добавьте публичную часть своего ключа к своему профилю в github.
> Done

## Основная часть

Наша основная цель - разбить наш playbook на отдельные roles. Задача: сделать roles для clickhouse, vector и lighthouse и написать playbook для использования этих ролей. Ожидаемый результат: существуют три ваших репозитория: два с roles и один с playbook.

1. Создал в старой версии playbook файл `requirements.yml` и заполил его следующим содержимым:
`vim requirements.yml`

   ```yaml
   ---
     - src: git@github.com:AlexeySetevoi/ansible-clickhouse.git
       scm: git
       version: "1.11.0"
       name: clickhouse 
   ```

2. При помощи `ansible-galaxy` скачал себе эту роль.
```yaml
user@user-VirtualBox:~/Desktop/playbook$ ansible-galaxy role install -r requirements.yml --force
Starting galaxy role install process
The authenticity of host 'github.com (140.82.121.3)' can't be established.
ECDSA key fingerprint is SHA256:p2QAMXNIC1TJYWeIOttrVc98/R1BUFWu3/LiyKgUfQM.
Are you sure you want to continue connecting (yes/no/[fingerprint])? yes
- extracting clickhouse to /home/user/.ansible/roles/clickhouse
- clickhouse (1.11.0) was installed successfully
```

3. Создал новый каталог с ролью при помощи `ansible-galaxy role init vector-role`.
```yaml
ansible-galaxy role init vector-role
ansible-galaxy role init lighthouse-role
ansible-galaxy role init nginx-role
```

4. На основе tasks из старого playbook заполнил новую role. Разнёс переменные между `vars` и `default`. 
> Done

5. Перенёс нужные шаблоны конфигов в `templates`.
> Done

6. Описал в `README.md` роли и их параметры.
> Done

7. Повторил шаги 3-6 для lighthouse, nginx. 
> Done

8. Выложил все roles в репозитории. Проставил тэги, используя семантическую нумерацию Добавил roles в `requirements.yml` в playbook.
> Done

9. Переработал playbook на использование roles. Не забудьте про зависимости lighthouse и возможности совмещения `roles` с `tasks`.
> Done

10. Выложил playbook в репозиторий.
> Done

11. Ссылки на репозитория с roles и ссылка на репозиторий с playbook:
[lighthouse-role](https://github.com/lint707/lighthouse-role)  
[nginx-role](https://github.com/lint707/nginx-role)  
[vector-role](https://github.com/lint707/vector-role)  
[08-ansible-04-role](https://github.com/lint707/devops-netology-1/tree/main/Homework/08-ansible-04-role)  

---

### Как оформить ДЗ?

Выполненное домашнее задание пришлите ссылкой на .md-файл в вашем репозитории.

---
