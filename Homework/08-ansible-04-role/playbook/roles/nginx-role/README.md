Role Name: nginx-role
=========

Simple nginx deploy and management role.

Role Variables
--------------

F: You can specify a particular version (or `*` for the latest). Please note that downgrade isn't supported.
```yaml
nginx_version: "1.22.0"
```

Example Playbook
----------------

Including an example of how to use your role (for instance, with variables passed in as parameters) is always nice for users too:

    - hosts: servers
      roles:
         - { role: username.rolename, x: 42 }
```
  - hosts: nginx
    remote_user: root
    vars:
      - nginx_version: "1.22.0"
    roles: 
      - nginx-role
    
    
```

License
-------

BSD

Author Information
------------------

Role by [lint707](https://github.com/lint707).

Dear contributors, thank you.
