# Домашнее задание к занятию "4.3. Языки разметки JSON и YAML"


## Обязательная задача 1
Мы выгрузили JSON, который получили через API запрос к нашему сервису:
```
   { "info" : "Sample JSON output from our service\t",
        "elements" :[
            { "name" : "first",
            "type" : "server",
            "ip" : 7175 
            }
            { "name" : "second",
            "type" : "proxy",
            "ip : 71.78.22.43
            }
        ]
    }
```
  Нужно найти и исправить все ошибки, которые допускает наш сервис.
  
  ### Ваш скрипт:
```json
    { "info" : "Sample JSON output from our service\t",
        "elements" : [
            { "name" : "first",
            "type" : "server",
            "ip" : 7175 #не ip-адрес
            },
            { "name" : "second",
            "type" : "proxy",
            "ip" : "71.78.22.43"
            }
        ]
    }
```
## Обязательная задача 2
В прошлый рабочий день мы создавали скрипт, позволяющий опрашивать веб-сервисы и получать их IP. К уже реализованному функционалу нам нужно добавить возможность записи JSON и YAML файлов, описывающих наши сервисы. Формат записи JSON по одному сервису: `{ "имя сервиса" : "его IP"}`. Формат записи YAML по одному сервису: `- имя сервиса: его IP`. Если в момент исполнения скрипта меняется IP у сервиса - он должен так же поменяться в yml и json файле.

### Ваш скрипт:
```python
#!/usr/bin/env python3
import socket
import time
import json
import yaml

srv = {'drive.google.com':'0.0.0.0', 'mail.google.com':'0.0.0.0', 'google.com':'0.0.0.0'}
while 1 == 1:
  for url, ip in srv.items():
    new_ip = socket.gethostbyname(url)
    if new_ip != ip:
      print(' [ERROR] ' + str(url) +' IP mistmatch: '+srv[url]+' '+new_ip)
      srv[url]=new_ip
    else:
     # srv[url] = new_ip
      print(str(url) + ' - ' + ip)
  with open('srv.json', 'w') as json_file:
    json_data= json.dumps(srv, indent=2)
    json_file.write(json_data)
  with open('srv.yaml', 'w') as yaml_file:
    yaml_data= yaml.dump(srv, explicit_start=True, explicit_end=True)
    yaml_file.write(yaml_data)
```

### Вывод скрипта при запуске при тестировании:
```
vagrant@vagrant:~/devops-netology-1$ ./script3.py
 [ERROR] drive.google.com IP mistmatch: 0.0.0.0 64.233.165.194
 [ERROR] mail.google.com IP mistmatch: 0.0.0.0 216.58.210.165
 [ERROR] google.com IP mistmatch: 0.0.0.0 216.58.210.174
drive.google.com - 64.233.165.194
mail.google.com - 216.58.210.165
google.com - 216.58.210.174
```

### json-файл(ы), который(е) записал ваш скрипт:
```json
vagrant@vagrant:~/devops-netology-1$ cat srv.json
{
  "drive.google.com": "64.233.165.194",
  "mail.google.com": "216.58.210.165",
  "google.com": "216.58.210.174"
}
```

### yml-файл(ы), который(е) записал ваш скрипт:
```yaml
vagrant@vagrant:~/devops-netology-1$ cat srv.yaml
---
drive.google.com: 64.233.165.194
google.com: 216.58.210.174
mail.google.com: 216.58.210.165
...
```

## Дополнительное задание (со звездочкой*) - необязательно к выполнению

Так как команды в нашей компании никак не могут прийти к единому мнению о том, какой формат разметки данных использовать: JSON или YAML, нам нужно реализовать парсер из одного формата в другой. Он должен уметь:
   * Принимать на вход имя файла
   * Проверять формат исходного файла. Если файл не json или yml - скрипт должен остановить свою работу
   * Распознавать какой формат данных в файле. Считается, что файлы *.json и *.yml могут быть перепутаны
   * Перекодировать данные из исходного формата во второй доступный (из JSON в YAML, из YAML в JSON)
   * При обнаружении ошибки в исходном файле - указать в стандартном выводе строку с ошибкой синтаксиса и её номер
   * Полученный файл должен иметь имя исходного файла, разница в наименовании обеспечивается разницей расширения файлов

### Ваш скрипт:
```python
???
```

### Пример работы скрипта:
???
