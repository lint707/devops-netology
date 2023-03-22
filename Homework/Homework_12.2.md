# Домашнее задание к занятию "Базовые объекты K8S"

### Задание 1. Создать Pod с именем "hello-world"

1. Создать манифест (yaml-конфигурацию) Pod  
![img1](img/pod_hello.jpg)  
2. Использовать image - gcr.io/kubernetes-e2e-test-images/echoserver:2.2  
3. Подключиться локально к Pod с помощью `kubectl port-forward hello-world 31180:8080` и вывести значение (curl или в браузере)  
![img1](img/get_pods.jpg)  
![img1](img/curl_pod.jpg)  
------

### Задание 2. Создать Service и подключить его к Pod

1. Создать Pod с именем "netology-web"  
![img1](img/pod_web.jpg)  
2. Использовать image - gcr.io/kubernetes-e2e-test-images/echoserver:2.2  
3. Создать Service с именем "netology-svc" и подключить к "netology-web"  
![img1](img/svc_web1.jpg)  
4. Подключиться локально к Service с помощью `kubectl port-forward svc/netology-svc 32280:8080` и вывести значение (curl или в браузере)  
![img1](img/get_svc.jpg)  
![img1](img/curl_svc.jpg)  
------





