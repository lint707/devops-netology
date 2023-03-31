# Домашнее задание к занятию "Как работает сеть в K8S"

### Задание 1. Создать сетевую политику (или несколько политик) для обеспечения доступа

1. Создать deployment'ы приложений frontend, backend и cache и соответсвующие сервисы.
2. В качестве образа использовать network-multitool.
3. Разместить поды в namespace app.
4. Создать политики чтобы обеспечить доступ frontend -> backend -> cache. Другие виды подключений должны быть запрещены.
5. Продемонстрировать, что трафик разрешен и запрещен.

![nwp0](img/nwp-0.jpg)
![nwp1](img/nwp-1.jpg)
![nwp3](img/nwp-3.jpg)

Deployment:
[frontend](manifests/10-frontend.yaml)
[backend](manifests/20-backend.yaml)
[cache](manifests/30-cache.yaml)

NetworkPolicy:
[default](manifests/nwp-default.yaml)
[backend](manifests/nwp-backend.yaml)
[cache](manifests/nwp-cache.yaml)


