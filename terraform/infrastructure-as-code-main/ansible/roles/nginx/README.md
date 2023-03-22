                                                                                        ```nginx```
                                                                                ```Роль для установки NGINX .```

``Aдаптирована для установки в Ubuntu, Debian, Centos.```
```
Requirements
Нужна роль common

Role Variables
```
В defaults/main.yml объявлены переменые:
```
config_nginx: true  # true создаем конфиг с нашими переменнами # false просто ставим nginx
nginx_worker_processes: 'auto'  # лучше ставить auto nginx (сам определит сколько воркеров нужно запускать)
nginx_enable_ipv6: false
nginx_access_log: '/var/log/nginx/access.log'
nginx_error_log: '/var/log/nginx/error.log'
nginx_events_block:
  - 'multi_accept on' #Если multi_accept включён, рабочие процессы будут принимать новые соединения по очереди.
  - 'worker_connections 768' 

# custom headers
nginx_headers: []
  # - 'set_real_ip_from 192.168.0.0/24'  # ip  адресс балансировщика (openstack)
  # - 'real_ip_header X-Real-IP'
  # - 'real_ip_recursive on'

nginx_http_block:
  basic_settings:
    - 'keepalive_timeout 65'
    - 'sendfile on'
    - 'server_names_hash_bucket_size 64'
    - 'server_name_in_redirect off'
    - 'server_tokens off'
    - 'tcp_nodelay on'
    - 'tcp_nopush on'
    - 'types_hash_max_size 2048'
    - 'include /etc/nginx/mime.types'
    - 'default_type application/octet-stream'
    - 'open_file_cache max=200000 inactive=20s'
    - 'open_file_cache_valid 30s'
    - 'open_file_cache_min_uses 2'
    - 'open_file_cache_errors on'
    - 'keepalive_requests 1000'
    - 'reset_timedout_connection on'
    - 'client_body_timeout 10'
    - 'client_max_body_size 64m'
    - 'send_timeout 2'
    - 'etag on'
    - 'limit_req_zone $binary_remote_addr zone=requests:10m rate=20r/s'
    - 'limit_req zone=requests burst=40'


  gzip_settings:
    - 'gzip on'
    - 'gzip_disable "msie6"'
    - 'gzip_vary on'
    - 'gzip_proxied any'
    - 'gzip_comp_level 6'
    - 'gzip_buffers 16 8k'
    - 'gzip_http_version 1.1'
    - 'gzip_types text/plain text/css application/json application/x-javascript text/xml application/xml application/xml+rss text/javascript'
  logging_settings:
    - 'access_log {{ nginx_access_log }}'
    - 'error_log {{ nginx_error_log }}'
  vhost_configs:
    - 'include /etc/nginx/conf.d/*.conf'
    - 'include /etc/nginx/sites-enabled/*'
nginx_listen_port: 80

nginx_default_sites_block:
    - 'server_name  _'
    - 'return 444'
    - 'access_log  off'
    - 'error_log off'

```
```Example Playbook
- hosts: servers
  roles:
     - common
     - { role: nginx, config_nginx: true }
```
