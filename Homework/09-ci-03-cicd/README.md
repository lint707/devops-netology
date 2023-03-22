# Домашнее задание к занятию "09.03 CI\CD"

## Подготовка к выполнению

1. Создаём 2 VM в yandex cloud со следующими параметрами: 2CPU 4RAM Centos7(остальное по минимальным требованиям)
> Done
2. Прописываем в [inventory](./infrastructure/inventory/cicd/hosts.yml) [playbook'a](./infrastructure/site.yml) созданные хосты
> Done
3. Добавляем в [files](./infrastructure/files/) файл со своим публичным ключом (id_rsa.pub). Если ключ называется иначе - найдите таску в плейбуке, которая использует id_rsa.pub имя и исправьте на своё
> Done
4. Запускаем playbook, ожидаем успешного завершения
```
user@user-VirtualBox:~/Desktop/09-ci-03-cicd/infrastructure$ ansible-playbook -i inventory/cicd/hosts.yml site.yml
```
> Done
5. Проверяем готовность Sonarqube через [браузер](http://localhost:9000)
> Done
6. Заходим под admin\admin, меняем пароль на свой
> Done
7.  Проверяем готовность Nexus через [бразуер](http://localhost:8081)
> Done
8. Подключаемся под admin\admin123, меняем пароль, сохраняем анонимный доступ
> Done

## Знакомоство с SonarQube

### Основная часть

1. Создаём новый проект, название произвольное
> Done
 
2. Скачиваем пакет sonar-scanner, который нам предлагает скачать сам sonarqube
>Done

3. Делаем так, чтобы binary был доступен через вызов в shell (или меняем переменную PATH или любой другой удобный вам способ)
```
user@user-VirtualBox:~/Desktop/09-ci-03-cicd/example/sonar-scanner-4.7.0.2747-linux/bin$ export PATH=$(pwd):$PATH
```

4. Проверяем `sonar-scanner --version`
```
user@user-VirtualBox:~/Desktop/09-ci-03-cicd/example$ sonar-scanner --version
INFO: Scanner configuration file: /home/user/Desktop/09-ci-03-cicd/example/sonar-scanner-4.7.0.2747-linux/conf/sonar-scanner.properties
INFO: Project root configuration file: NONE
INFO: SonarScanner 4.7.0.2747
INFO: Java 11.0.14.1 Eclipse Adoptium (64-bit)
INFO: Linux 5.15.0-48-generic amd64
```

5. Запускаем анализатор против кода из директории [example](./example) с дополнительным ключом `-Dsonar.coverage.exclusions=fail.py`
```
user@user-VirtualBox:~/Desktop/09-ci-03-cicd/example$ sonar-scanner \
>   -Dsonar.projectKey=Netology \
>   -Dsonar.sources=. \
>   -Dsonar.host.url=http://178.154.220.228:9000 \
>   -Dsonar.login=780f408971631e7e8d7177ca89d89f73722b2798 \
>   -Dsonar.coverage.exclusions=fail.py
INFO: Scanner configuration file: /home/user/Desktop/09-ci-03-cicd/example/sonar-scanner-4.7.0.2747-linux/conf/sonar-scanner.properties
INFO: Project root configuration file: NONE
INFO: SonarScanner 4.7.0.2747
INFO: Java 11.0.14.1 Eclipse Adoptium (64-bit)
INFO: Linux 5.15.0-48-generic amd64
INFO: User cache: /home/user/.sonar/cache
INFO: Scanner configuration file: /home/user/Desktop/09-ci-03-cicd/example/sonar-scanner-4.7.0.2747-linux/conf/sonar-scanner.properties
INFO: Project root configuration file: NONE
INFO: Analyzing on SonarQube server 9.1.0
INFO: Default locale: "en_US", source code encoding: "UTF-8" (analysis is platform dependent)
INFO: Load global settings
INFO: Load global settings (done) | time=317ms
INFO: Server id: 9CFC3560-AYOxoBQPMkGU3vOMMIKo
INFO: User cache: /home/user/.sonar/cache
INFO: Load/download plugins
INFO: Load plugins index
INFO: Load plugins index (done) | time=121ms
INFO: Load/download plugins (done) | time=271857ms
INFO: Process project properties
INFO: Process project properties (done) | time=38ms
INFO: Execute project builders
INFO: Execute project builders (done) | time=11ms
INFO: Project key: Netology
INFO: Base dir: /home/user/Desktop/09-ci-03-cicd/example
INFO: Working dir: /home/user/Desktop/09-ci-03-cicd/example/.scannerwork
INFO: Load project settings for component key: 'Netology'
INFO: Load project settings for component key: 'Netology' (done) | time=86ms
INFO: Load quality profiles
INFO: Load quality profiles (done) | time=144ms
INFO: Load active rules
INFO: Load active rules (done) | time=3942ms
WARN: SCM provider autodetection failed. Please use "sonar.scm.provider" to define SCM of your project, or disable the SCM Sensor in the project settings.
INFO: Indexing files...
INFO: Project configuration:
INFO:   Excluded sources for coverage: fail.py
INFO: 267 files indexed
INFO: Quality profile for py: Sonar way
INFO: ------------- Run sensors on module Netology
INFO: Load metrics repository
INFO: Load metrics repository (done) | time=141ms
INFO: Sensor Python Sensor [python]
WARN: Your code is analyzed as compatible with python 2 and 3 by default. This will prevent the detection of issues specific to python 2 or python 3. You can get a more precise analysis by setting a python version in your configuration via the parameter "sonar.python.version"
INFO: Starting global symbols computation
INFO: 1 source file to be analyzed
INFO: Load project repositories
INFO: Load project repositories (done) | time=82ms
INFO: 1/1 source file has been analyzed
INFO: Starting rules execution
INFO: 1 source file to be analyzed
INFO: 1/1 source file has been analyzed
INFO: Sensor Python Sensor [python] (done) | time=3053ms
INFO: Sensor Cobertura Sensor for Python coverage [python]
INFO: Sensor Cobertura Sensor for Python coverage [python] (done) | time=77ms
INFO: Sensor PythonXUnitSensor [python]
INFO: Sensor PythonXUnitSensor [python] (done) | time=30ms
INFO: Sensor CSS Rules [cssfamily]
INFO: No CSS, PHP, HTML or VueJS files are found in the project. CSS analysis is skipped.
INFO: Sensor CSS Rules [cssfamily] (done) | time=5ms
INFO: Sensor JaCoCo XML Report Importer [jacoco]
INFO: 'sonar.coverage.jacoco.xmlReportPaths' is not defined. Using default locations: target/site/jacoco/jacoco.xml,target/site/jacoco-it/jacoco.xml,build/reports/jacoco/test/jacocoTestReport.xml
INFO: No report imported, no coverage information will be imported by JaCoCo XML Report Importer
INFO: Sensor JaCoCo XML Report Importer [jacoco] (done) | time=23ms
INFO: Sensor C# Project Type Information [csharp]
INFO: Sensor C# Project Type Information [csharp] (done) | time=7ms
INFO: Sensor C# Analysis Log [csharp]
INFO: Sensor C# Analysis Log [csharp] (done) | time=77ms
INFO: Sensor C# Properties [csharp]
INFO: Sensor C# Properties [csharp] (done) | time=0ms
INFO: Sensor JavaXmlSensor [java]
INFO: Sensor JavaXmlSensor [java] (done) | time=9ms
INFO: Sensor HTML [web]
INFO: Sensor HTML [web] (done) | time=4ms
INFO: Sensor VB.NET Project Type Information [vbnet]
INFO: Sensor VB.NET Project Type Information [vbnet] (done) | time=1ms
INFO: Sensor VB.NET Analysis Log [vbnet]
INFO: Sensor VB.NET Analysis Log [vbnet] (done) | time=26ms
INFO: Sensor VB.NET Properties [vbnet]
INFO: Sensor VB.NET Properties [vbnet] (done) | time=0ms
INFO: ------------- Run sensors on project
INFO: Sensor Zero Coverage Sensor
INFO: Sensor Zero Coverage Sensor (done) | time=4ms
INFO: SCM Publisher No SCM system was detected. You can use the 'sonar.scm.provider' property to explicitly specify it.
INFO: CPD Executor Calculating CPD for 1 file
INFO: CPD Executor CPD calculation finished (done) | time=29ms
INFO: Analysis report generated in 184ms, dir size=103.0 kB
INFO: Analysis report compressed in 64ms, zip size=14.1 kB
INFO: Analysis report uploaded in 65ms
INFO: ANALYSIS SUCCESSFUL, you can browse http://178.154.220.228:9000/dashboard?id=Netology
INFO: Note that you will be able to access the updated dashboard once the server has processed the submitted analysis report
INFO: More about the report processing at http://178.154.220.228:9000/api/ce/task?id=AYOxy3mdMkGU3vOMMNPv
INFO: Analysis total time: 16.737 s
INFO: ------------------------------------------------------------------------
INFO: EXECUTION SUCCESS
INFO: ------------------------------------------------------------------------
INFO: Total time: 5:41.305s
INFO: Final Memory: 7M/27M
INFO: ------------------------------------------------------------------------

```
6. Смотрим результат в интерфейсе
![sonar1](https://github.com/lint707/devops-netology-1/blob/main/Homework/img/sonar-1.jpg)

7. Исправляем ошибки, которые он выявил(включая warnings)
>done

8. Запускаем анализатор повторно - проверяем, что QG пройдены успешно
```
user@user-VirtualBox:~/Desktop/09-ci-03-cicd/example$ sonar-scanner   -Dsonar.projectKey=Netology   -Dsonar.sources=.   -Dsonar.host.url=http://178.154.220.228:9000   -Dsonar.login=780f408971631e7e8d7177ca89d89f73722b2798   -Dsonar.coverage.exclusions=fail.py
INFO: Scanner configuration file: /home/user/Desktop/09-ci-03-cicd/example/sonar-scanner-4.7.0.2747-linux/conf/sonar-scanner.properties
INFO: Project root configuration file: NONE
INFO: SonarScanner 4.7.0.2747
INFO: Java 11.0.14.1 Eclipse Adoptium (64-bit)
INFO: Linux 5.15.0-48-generic amd64
INFO: User cache: /home/user/.sonar/cache
INFO: Scanner configuration file: /home/user/Desktop/09-ci-03-cicd/example/sonar-scanner-4.7.0.2747-linux/conf/sonar-scanner.properties
INFO: Project root configuration file: NONE
INFO: Analyzing on SonarQube server 9.1.0
INFO: Default locale: "en_US", source code encoding: "UTF-8" (analysis is platform dependent)
INFO: Load global settings
INFO: Load global settings (done) | time=175ms
INFO: Server id: 9CFC3560-AYOxoBQPMkGU3vOMMIKo
INFO: User cache: /home/user/.sonar/cache
INFO: Load/download plugins
INFO: Load plugins index
INFO: Load plugins index (done) | time=97ms
INFO: Load/download plugins (done) | time=290ms
INFO: Process project properties
INFO: Process project properties (done) | time=16ms
INFO: Execute project builders
INFO: Execute project builders (done) | time=5ms
INFO: Project key: Netology
INFO: Base dir: /home/user/Desktop/09-ci-03-cicd/example
INFO: Working dir: /home/user/Desktop/09-ci-03-cicd/example/.scannerwork
INFO: Load project settings for component key: 'Netology'
INFO: Load project settings for component key: 'Netology' (done) | time=61ms
INFO: Load quality profiles
INFO: Load quality profiles (done) | time=92ms
INFO: Load active rules
INFO: Load active rules (done) | time=2509ms
WARN: SCM provider autodetection failed. Please use "sonar.scm.provider" to define SCM of your project, or disable the SCM Sensor in the project settings.
INFO: Indexing files...
INFO: Project configuration:
INFO:   Excluded sources for coverage: fail.py
INFO: 267 files indexed
INFO: Quality profile for py: Sonar way
INFO: ------------- Run sensors on module Netology
INFO: Load metrics repository
INFO: Load metrics repository (done) | time=72ms
INFO: Sensor Python Sensor [python]
WARN: Your code is analyzed as compatible with python 2 and 3 by default. This will prevent the detection of issues specific to python 2 or python 3. You can get a more precise analysis by setting a python version in your configuration via the parameter "sonar.python.version"
INFO: Starting global symbols computation
INFO: 1 source file to be analyzed
INFO: Load project repositories
INFO: Load project repositories (done) | time=39ms
INFO: 1/1 source file has been analyzed
INFO: Starting rules execution
INFO: 1 source file to be analyzed
INFO: 1/1 source file has been analyzed
INFO: Sensor Python Sensor [python] (done) | time=1137ms
INFO: Sensor Cobertura Sensor for Python coverage [python]
INFO: Sensor Cobertura Sensor for Python coverage [python] (done) | time=50ms
INFO: Sensor PythonXUnitSensor [python]
INFO: Sensor PythonXUnitSensor [python] (done) | time=15ms
INFO: Sensor CSS Rules [cssfamily]
INFO: No CSS, PHP, HTML or VueJS files are found in the project. CSS analysis is skipped.
INFO: Sensor CSS Rules [cssfamily] (done) | time=5ms
INFO: Sensor JaCoCo XML Report Importer [jacoco]
INFO: 'sonar.coverage.jacoco.xmlReportPaths' is not defined. Using default locations: target/site/jacoco/jacoco.xml,target/site/jacoco-it/jacoco.xml,build/reports/jacoco/test/jacocoTestReport.xml
INFO: No report imported, no coverage information will be imported by JaCoCo XML Report Importer
INFO: Sensor JaCoCo XML Report Importer [jacoco] (done) | time=15ms
INFO: Sensor C# Project Type Information [csharp]
INFO: Sensor C# Project Type Information [csharp] (done) | time=13ms
INFO: Sensor C# Analysis Log [csharp]
INFO: Sensor C# Analysis Log [csharp] (done) | time=29ms
INFO: Sensor C# Properties [csharp]
INFO: Sensor C# Properties [csharp] (done) | time=0ms
INFO: Sensor JavaXmlSensor [java]
INFO: Sensor JavaXmlSensor [java] (done) | time=19ms
INFO: Sensor HTML [web]
INFO: Sensor HTML [web] (done) | time=10ms
INFO: Sensor VB.NET Project Type Information [vbnet]
INFO: Sensor VB.NET Project Type Information [vbnet] (done) | time=2ms
INFO: Sensor VB.NET Analysis Log [vbnet]
INFO: Sensor VB.NET Analysis Log [vbnet] (done) | time=17ms
INFO: Sensor VB.NET Properties [vbnet]
INFO: Sensor VB.NET Properties [vbnet] (done) | time=0ms
INFO: ------------- Run sensors on project
INFO: Sensor Zero Coverage Sensor
INFO: Sensor Zero Coverage Sensor (done) | time=15ms
INFO: SCM Publisher No SCM system was detected. You can use the 'sonar.scm.provider' property to explicitly specify it.
INFO: CPD Executor Calculating CPD for 1 file
INFO: CPD Executor CPD calculation finished (done) | time=19ms
INFO: Analysis report generated in 165ms, dir size=102.8 kB
INFO: Analysis report compressed in 95ms, zip size=13.8 kB
INFO: Analysis report uploaded in 65ms
INFO: ANALYSIS SUCCESSFUL, you can browse http://178.154.220.228:9000/dashboard?id=Netology
INFO: Note that you will be able to access the updated dashboard once the server has processed the submitted analysis report
INFO: More about the report processing at http://178.154.220.228:9000/api/ce/task?id=AYOx1UZhMkGU3vOMMNPw
INFO: Analysis total time: 8.171 s
INFO: ------------------------------------------------------------------------
INFO: EXECUTION SUCCESS
INFO: ------------------------------------------------------------------------
INFO: Total time: 10.397s
INFO: Final Memory: 7M/30M
INFO: ------------------------------------------------------------------------

```

9. Делаем скриншот успешного прохождения анализа, прикладываем к решению ДЗ
![sonar2](https://github.com/lint707/devops-netology-1/blob/main/Homework/img/sonar-2.jpg)

## Знакомство с Nexus

### Основная часть

1. В репозиторий `maven-releases` загружаем артефакт с GAV параметрами:
   1. groupId: netology
   2. artifactId: java
   3. version: 8_282
   4. classifier: distrib
   5. type: tar.gz
2. В него же загружаем такой же артефакт, но с version: 8_102
> Done

3. Проверяем, что все файлы загрузились успешно
> Done

4. В ответе присылаем файл `maven-metadata.xml` для этого артефекта
[maven-metadata.xml](https://github.com/lint707/devops-netology-1/blob/main/Homework/09-ci-03-cicd/maven-metadata.xml)

### Знакомство с Maven

### Подготовка к выполнению

1. Скачиваем дистрибутив с [maven](https://maven.apache.org/download.cgi)
> Done

2. Разархивируем, делаем так, чтобы binary был доступен через вызов в shell (или меняем переменную PATH или любой другой удобный вам способ)
> Done

3. Удаляем из `apache-maven-<version>/conf/settings.xml` упоминание о правиле, отвергающем http соединение( раздел mirrors->id: my-repository-http-blocker)
> Done

4. Проверяем `mvn --version`
> Done

5. Забираем директорию [mvn](./mvn) с pom
> Done

### Основная часть

1. Меняем в `pom.xml` блок с зависимостями под наш артефакт из первого пункта задания для Nexus (java с версией 8_282)
> Done

2. Запускаем команду `mvn package` в директории с `pom.xml`, ожидаем успешного окончания
```
[WARNING] JAR will be empty - no content was marked for inclusion!
[INFO] Building jar: /home/user/Desktop/09-ci-03-cicd/mvn/target/simple-app-1.0-SNAPSHOT.jar
[INFO] ------------------------------------------------------------------------
[INFO] BUILD SUCCESS
[INFO] ------------------------------------------------------------------------
[INFO] Total time:  27.589 s
[INFO] Finished at: 2022-10-07T14:09:22+03:00
[INFO] ------------------------------------------------------------------------
```

3. Проверяем директорию `~/.m2/repository/`, находим наш артефакт
```
user@user-VirtualBox:/usr/share/java$ ls -lh ~/.m2/repository/netology/
total 4,0K
drwxrwxr-x 3 user user 4,0K окт  7 14:08 java
```

4. В ответе присылаем исправленный файл `pom.xml`
[pom.xml](https://github.com/lint707/devops-netology-1/blob/main/Homework/09-ci-03-cicd/mvn/pom.xml)

---

### Как оформить ДЗ?

Выполненное домашнее задание пришлите ссылкой на .md-файл в вашем репозитории.

---
