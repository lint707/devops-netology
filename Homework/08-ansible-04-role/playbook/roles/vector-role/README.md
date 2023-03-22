Role Name: vector-role
=========

Simple vector deploy and management role.

Role Variables
--------------

F: You can specify a particular version (or `*` for the latest). Please note that downgrade isn't supported.
```yaml
vector_version: "0.21.0"
```

Example Playbook
----------------

Including an example of how to use your role (for instance, with variables passed in as parameters) is always nice for users too:

    - hosts: servers
      roles:
         - { role: username.rolename, x: 42 }
```
  - hosts: vector
    remote_user: root
    vars:
      - vector_version: "0.21.0"
    roles: 
      - vector-role    
```

License
-------

BSD

Author Information
------------------

Role by [lint707](https://github.com/lint707).

Dear contributors, thank you.
