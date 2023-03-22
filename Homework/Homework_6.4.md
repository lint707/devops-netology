# Домашнее задание к занятию "6.4. PostgreSQL"

## Задача 1  

Используя docker поднял инстанс PostgreSQL (версию 13). Данные БД сохранил в volume.  

Подключился к БД PostgreSQL используя `psql`.  
Воспользовался командой `\?` для вывода подсказки по имеющимся в `psql` управляющим командам:  
![psg2](img/psg2_start.jpg)  

Примеры управляющих команды для:  
- вывода списка БД: `postgres=# \l`  
![psg2](img/psg2_l.jpg)  

- подключения к БД: `postgres=# \c postgres`  
![psg2](img/psg2_c.jpg)  

- вывода списка таблиц: `postgres=# \dtS` или `\dt`  
![psg2](img/psg2_dts.jpg)  

- вывода описания содержимого таблиц: `postgres=# \dS pg_event_trigger` или `\d имя таблицы`
![psg2](img/psg2_ds.jpg)  

- выхода из psql: `postgres=# \q`  
![psg2](img/psg2_q.jpg)  

## Задача 2  
Используя `psql` создал БД `test_database`:  
![psg2](img/psg2_l2.jpg)  

Изучил [бэкап БД](https://github.com/netology-code/virt-homeworks/tree/master/06-db-04-postgresql/test_data).  

Восстановил бэкап БД в `test_database`:  
![psg2](img/psg2_bk.jpg)  

Перешёл в управляющую консоль `psql` внутри контейнера.  
Подключился к восстановленной БД и провёл операцию ANALYZE для сбора статистики по таблице.  
![psg2](img/psg2_cr.jpg)  

Используя таблицу [pg_stats](https://postgrespro.ru/docs/postgresql/12/view-pg-stats), нашёл столбец таблицы `orders`   
с наибольшим средним значением размера элементов в байтах.  
![psg2](img/psg2_pgstats.jpg)  

## Задача 3
Архитектор и администратор БД выяснили, что ваша таблица orders разрослась до невиданных размеров и
поиск по ней занимает долгое время. Вам, как успешному выпускнику курсов DevOps в нетологии предложили
провести разбиение таблицы на 2 (шардировать на orders_1 - price>499 и orders_2 - price<=499).
SQL-транзакции для проведения данной операции:  
![psg2](img/psg2_order.jpg)  
![psg2](img/psg2_order_sl.jpg)  

Можно ли было изначально исключить "ручное" разбиение при проектировании таблицы orders?
 - Ручного разбиения можно было исбежать, если при изначальном проектировании сделать таблицу как секционарованную.

## Задача 4
Используя утилиту `pg_dump` создал бекап БД `test_database`:  
![psg2](img/psg2_bkp.jpg)  

Как бы вы доработали бэкап-файл, чтобы добавить уникальность значения столбца `title` для таблиц `test_database`?
 - Доработать файл бекапа можно следующим образом, добавив UNIQUE(title):
```
CREATE TABLE public.orders (
    id integer NOT NULL,
    title character varying(80) NOT NULL,
    price integer DEFAULT 0,
    UNIQUE(title)
);
```

---
