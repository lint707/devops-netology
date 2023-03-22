# Домашнее задание к занятию "5.3. Введение. Экосистема. Архитектура. Жизненный цикл Docker контейнера"

---

## Задача 1

Сценарий выполения задачи:

- создайте свой репозиторий на https://hub.docker.com;
- выберете любой образ, который содержит веб-сервер Nginx;
- создайте свой fork образа;
- реализуйте функциональность:
запуск веб-сервера в фоне с индекс-страницей, содержащей HTML-код ниже:
```
<html>
<head>
Hey, Netology
</head>
<body>
<h1>I’m DevOps Engineer!</h1>
</body>
</html>
```
Опубликуйте созданный форк в своем репозитории и предоставьте ответ в виде ссылки на https://hub.docker.com/username_repo.

Ссылка на репозторий: https://hub.docker.com/r/lint707/nginx_netology5.3

### Решение:

Скачиваем образ:
```
vagrant@server1:~$ docker pull ubuntu/nginx:1.18-22.04_edge
1.18-22.04_edge: Pulling from ubuntu/nginx
afad979d238f: Pull complete 
7efc54ee67e0: Pull complete 
9185eb702abc: Pull complete 
70bc7ac6a482: Pull complete 
Digest: sha256:c798b60006732b86a43bca84d9e5059b465ecdcdb305dd3452ed5613b29e30ac
Status: Downloaded newer image for ubuntu/nginx:1.18-22.04_edge
docker.io/ubuntu/nginx:1.18-22.04_edge

vagrant@server1:~$ docker images
REPOSITORY     TAG               IMAGE ID       CREATED      SIZE
ubuntu/nginx   1.18-22.04_edge   c1d6caca3772   9 days ago   144MB
```
Создаём dockerfile:
```
vagrant@server1:~/docker$ vim dockerfile
FROM ubuntu/nginx:1.18-22.04_edge
RUN echo '<html><head>Hey, Netology</head><body><h1>I am DevOps Engineer!</h1></body></html>' > /var/www/html/index.html
```
Делаем fork образа:
```
vagrant@server1:~/docker$ docker build -f dockerfile -t lint707/nginx_netology5.3  .
Sending build context to Docker daemon  2.048kB
Step 1/2 : FROM ubuntu/nginx:1.18-22.04_edge
 ---> c1d6caca3772
Step 2/2 : RUN echo '<html><head>Hey, Netology</head><body><h1>I am DevOps Engineer!</h1></body></html>' > /usr/share/nginx/html/index.html
 ---> Running in 49d94709c9e8
Removing intermediate container 49d94709c9e8
 ---> a61bb2bbbf69
Successfully built a61bb2bbbf69
Successfully tagged lint707/nginx_netology5.3:latest

vagrant@server1:~/docker$ docker images
REPOSITORY                  TAG               IMAGE ID       CREATED              SIZE
lint707/nginx_netology5.3   latest            a61bb2bbbf69   About a minute ago   144MB
ubuntu/nginx                1.18-22.04_edge   c1d6caca3772   9 days ago           144MB
```
Авторизация и загрузка на hud.docker:
```
vagrant@server1:~/docker$ docker login -u lint707
Password: 
Login Succeeded

vagrant@server1:~/docker$ docker push lint707/nginx_netology5.3
Using default tag: latest
The push refers to repository [docker.io/lint707/nginx_netology5.3]
760a4341f34a: Pushed 
e29410d289ba: Mounted from ubuntu/nginx 
9126dae3b3a7: Mounted from ubuntu/nginx 
39767fb85b82: Mounted from ubuntu/nginx 
a790f937a6ae: Mounted from ubuntu/nginx 
latest: digest: sha256:cee73146fb887ef3c7937a9c7a9c483e965b335b3b64acf047cb228ceb76e31a size: 1362
```
Запуск контейнера:
```
vagrant@server1:~/docker$ docker run -d -v /home/vagrant/docker/index.html:/var/www/html/index.html -p 8080:80 lint707/nginx_netology5.3
34fd7120644b6d72cc3d36fc9e044e690a42e6363b9f352f828b3f8282ea343d
vagrant@server1:~/docker$ curl localhost:8080
<html> 
	<head>Hey, Netology</head>
	<body><h1>I am DevOps Engineer!</h1></body>
</html>
```

## Задача 2

Посмотрите на сценарий ниже и ответьте на вопрос:
"Подходит ли в этом сценарии использование Docker контейнеров или лучше подойдет виртуальная машина, физическая машина? Может быть возможны разные варианты?"

Детально опишите и обоснуйте свой выбор.

--

Ответ:
- Высоконагруженное монолитное java веб-приложение:
> Физичиские или виртуальные машины, если в приложениине заложено масштабированиетогла дучше физическая машина,
в случае мастрабируется и может взаимодействовать с балансировщиком, то виртуальная машина.
- Nodejs веб-приложение:
> Docker, удобство масштабирования, легковесность, простота развётрывания.
- Мобильное приложение c версиями для Android и iOS:
> Виртуальные машины, упрощает тестирование и размещение.
- Шина данных на базе Apache Kafka:
> Docker, возможность быстрого отката, в случае проблем на проде, изолированность приложений.
- Elasticsearch кластер для реализации логирования продуктивного веб-приложения - три ноды elasticsearch, два logstash и две ноды kibana:
>Docker, простота при обновлении, миграции, удаления логов, удобнее при кластеризации - меньше времени на запуск контейнеров.
- Мониторинг-стек на базе Prometheus и Grafana:
> Docker, удобство масштабирования, легковесность, простота развётрывания, есть готовые образы.
- MongoDB, как основное хранилище данных для java-приложения:
> Виртуальная машина, не подходящее решение хранить БД в контейнере. 
- Gitlab сервер для реализации CI/CD процессов и приватный (закрытый) Docker Registry:
> Виртуальная машина, удобство бекапов и миграции.

## Задача 3

- Запустите первый контейнер из образа ***centos*** c любым тэгом в фоновом режиме, подключив папку ```/data``` из текущей рабочей директории на хостовой машине в ```/data``` контейнера;
- Запустите второй контейнер из образа ***debian*** в фоновом режиме, подключив папку ```/data``` из текущей рабочей директории на хостовой машине в ```/data``` контейнера;
- Подключитесь к первому контейнеру с помощью ```docker exec``` и создайте текстовый файл любого содержания в ```/data```;
- Добавьте еще один файл в папку ```/data``` на хостовой машине;
- Подключитесь во второй контейнер и отобразите листинг и содержание файлов в ```/data``` контейнера.

### Решение:

```
vagrant@server1:~/docker$ docker images
REPOSITORY                  TAG                      IMAGE ID       CREATED        SIZE
lint707/nginx_netology5.3   latest                   a61bb2bbbf69   2 days ago     144MB
debian                      unstable-20220622-slim   2a8550669e9d   11 days ago    76MB
centos                      centos7.9.2009           eeb6ee3f44bd   9 months ago   204MB
vagrant@server1:~/docker$ docker run -v /data:/data -dt --name debian
"docker run" requires at least 1 argument.
See 'docker run --help'.

Usage:  docker run [OPTIONS] IMAGE [COMMAND] [ARG...]

Run a command in a new container
vagrant@server1:~/docker$ docker run -v /home/vagrant/docker/data:/data -dt --name debian debian
Unable to find image 'debian:latest' locally
latest: Pulling from library/debian
1339eaac5b67: Pull complete 
Digest: sha256:859ea45db307402ee024b153c7a63ad4888eb4751921abbef68679fc73c4c739
Status: Downloaded newer image for debian:latest
50de199aaba04fb0770d6fc58e2594ae8d525a0a1090d167a1129220b31a43b5
vagrant@server1:~/docker$ docker run -v /home/vagrant/docker/data:/data -dt --name centos centos
Unable to find image 'centos:latest' locally
latest: Pulling from library/centos
a1d0c7532777: Pull complete 
Digest: sha256:a27fd8080b517143cbbbab9dfb7c8571c40d67d534bbdee55bd6c473f432b177
Status: Downloaded newer image for centos:latest
eefe838d522feb6951afbc001ed3af6f8994f90a2d170cbdff4b302ce1e9ab47
vagrant@server1:~/docker$ 
```
Создание файла в первом контейнере:
```
vagrant@server1:~/docker$ docker exec -it 3a056480b140 bash
[root@3a056480b140 /]# cd data
[root@3a056480b140 data]# echo "test text" >> test.txt
```
Создание файла на хостовой машине:
```
vagrant@server1:~/docker/data$ vim test.txt
```
Проверка в втором контейнере:
```
vagrant@server1:~/docker/data$ docker exec -it 50de199aaba0 bash
[root@50de199aaba0:/data]# ls -lha
total 16K
drwxrwxr-x 2 1000 1000 4.0K Jul  4 07:52 .
drwxr-xr-x 1 root root 4.0K Jul  4 07:51 ..
-rw-r--r-- 1 root root   10 Jul  4 07:52 test.txt
-rw-rw-r-- 1 1000 1000   10 Jul  4 07:58 test_host.txt
```

## Задача 4 (*)

Воспроизвести практическую часть лекции самостоятельно.

Соберите Docker образ с Ansible, загрузите на Docker Hub и пришлите ссылку вместе с остальными ответами к задачам.

https://hub.docker.com/r/lint707/ansible_ntl

---


Сценарий:



