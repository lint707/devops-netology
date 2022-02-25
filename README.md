Hello 25.02.2022
Test text
Second test
Third test

Файлы, которые будут в локальном каталоге ".terraform", включая вложенные каталоги в него. Сам каталог может быть любой вложенности.

Исключает все файлы с расширением: *.tfstate
Исключает все файлы, что содержат: *.tfstate.*,  

Исключает файл лога: crash.log и файлы логов в названии которых присутствует: crash.*.log

Исключает файлы содаржащие конфиденциальные данные, с расширением: *.tfvars 
и с: *.tfvars.json 

Исключает конкретные файлы и производные от них:
override.tf
override.tf.json
*_override.tf
*_override.tf.json

файлы конфигурации CLI
.terraformrc
terraform.rc

one
two
hird
