# Домашнее задание к занятию "Микросервисы: принципы"

Вы работаете в крупной компанию, которая строит систему на основе микросервисной архитектуры.
Вам как DevOps специалисту необходимо выдвинуть предложение по организации инфраструктуры, для разработки и эксплуатации.

## Задача 1: API Gateway 

Предложите решение для обеспечения реализации API Gateway. Составьте сравнительную таблицу возможностей различных программных решений. На основе таблицы сделайте выбор решения.

Решение должно соответствовать следующим требованиям:
- Маршрутизация запросов к нужному сервису на основе конфигурации
- Возможность проверки аутентификационной информации в запросах
- Обеспечение терминации HTTPS

#### Ответ:
https://geekflare.com/api-gateway/
https://coderlessons.com/articles/programmirovanie/mikroservisy-dlia-razrabotchikov-java-shliuzy-i-agregatory-api
https://www.moesif.com/blog/technical/api-gateways/How-to-Choose-The-Right-API-Gateway-For-Your-Platform-Comparison-Of-Kong-Tyk-Apigee-And-Alternatives/

| Products |	Kong | Tyk.io |	APIGee |	AWS Gateway |	Azure Gateway |	Express Gateway |
|----------|-------|--------|--------|--------------|---------------|-----------------|
| Deployment Complexity	| Single node	| Single node	| Many nodes with different roles	| Cloud vendor PaaS	| Cloud vendor PaaS |	Flexible |
| Data Stores Required	| Cassandra or Postgres	| Redis |	Cassandra, Zookeeper, and Postgres	| Cloud vendor PaaS	| Cloud vendor PaaS |	Redis |
| Open Source |	Yes, Apache 2.0 |	Yes, MPL	| No	| No	| No	| Yes, Apache 2.0 |
| Core Technology	| NGINX/Lua |	GoLang |	Java |	Not open |	Not open |	Node.js Express |
| On Premise	| Yes	| Yes	| Yes	| No	| Mo	| Yes |
| Community/Extensions | Large	| Medium	| No	| No	| No	| Small |
| Authorization/API Keys |	Yes	| Yes	| Yes	| Yes	| Yes	| Yes |
| Rate Limiting |	Yes	| Yes	| Yes	| Yes	| Yes	| Yes |
| Data Transformation	| HTTP	| HTTP	| Yes	| No	| No	| No |
| Integrated Billing	| No	| No	| Yes	| No	| No	| No |
| Маршрутизация запросов	| 	| 	| 	| 	| 	|  |
| Проверка аутентификации	| 	| 	| 	| 	| 	|  |
| Терминации HTTPS	| 	| 	| 	| 	| 	|  |



## Задача 2: Брокер сообщений

Составьте таблицу возможностей различных брокеров сообщений. На основе таблицы сделайте обоснованный выбор решения.

Решение должно соответствовать следующим требованиям:
- Поддержка кластеризации для обеспечения надежности
- Хранение сообщений на диске в процессе доставки
- Высокая скорость работы
- Поддержка различных форматов сообщений
- Разделение прав доступа к различным потокам сообщений
- Протота эксплуатации

#### Ответ:

| Критерий | [RabbitMQ](https://www.rabbitmq.com/) | [Apache Kafka](https://kafka.apache.org/) | [Qpid](https://qpid.apache.org/components/cpp-broker/index.html) | [SwiftMQ](https://www.swiftmq.com/)	
|----------|----------|----------|----------|-----------|
| Кластеризации | + | + | + | + |
| Хранение сообщений на диске в процессе доставки | + | + | + | + |
| Cкорость работы | + | + | - | - | 
| Поддержка различных форматов сообщений | + | + | - | - |
| Разделение прав доступа | + | + | + | + | 
| Простота эксплуатации | + | + | - | - |

Я бы выбрал Apache Kafka, который обладает следующими преимуществами:
- позволяет масштабировать  систему до бесконечности;
- может достичь пропускной способности в миллионы сообщений в секунду даже при ограниченных ресурсах;
- сообщения в Kafka не удаляются брокерами по мере их обработки консьюмерами — данные в Kafka могут храниться днями, неделями, годами. Благодаря этому одно и то же сообщение может быть обработано сколько угодно раз разными консьюмерами и в разных контекстах.
- поддерживает ACL;
- гарантирует порядок доставки сообщений.
---
