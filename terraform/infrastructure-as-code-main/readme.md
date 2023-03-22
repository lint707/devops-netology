# Infrastructure As Code

> Для начала работы установите следующие env переменные, получить их можно при установке yandex cloud cli

- export TF_VAR_yc_token= - тот же токен что при установке cli
- export TF_VAR_yc_cloud_id= взять из ГУИ облака

Если возникли проблемы при работе приложения, нет доступа к новостному API то выполните следующие действия:

- Зарегистрируйтесь и получите API ключ [News API](https://newsapi.org/register)
- Полученный ключ укажите в переменной **news_app_api_key** в файле **ansible/inventory/demo/group_vars/news.yml**
- Разверните приложение заново используя **make reconfig**
### Для деплоя в облако:

Для работы с terraform workspace необходимо указать переменную окружения ENV
Например:
* ENV=prod - для prod окружения
* ENV=stage - для stage окружения

Ansible для раскатки использует ENV для поиска в облаке виртуальных машин

```shell
make all
```

### Для удаления из облака:

```shell
make destroy && make clean
```


### TODO

* Добавить демо null_resource
* Зарефакторить makefile убрать cd ..
* Добавить рендеринг ansible_inventory с template