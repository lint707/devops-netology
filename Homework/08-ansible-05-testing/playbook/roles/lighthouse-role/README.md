Role Name: lighthouse-role
=========

Simple lighthouse deploy and management role.
    
Example Playbook
----------------

Including an example of how to use your role (for instance, with variables passed in as parameters) is always nice for users too:

    - hosts: servers
      roles:
         - { role: username.rolename, x: 42 }
```
  - hosts: lighthouse
    remote_user: root
    vars:
      lh_loc_dir: /home/user/appserv/lighthouse
      lh_vcs: https://github.com/VKCOM/lighthouse.git
    roles:
      - lighthouse-role
```

License
-------

BSD

Author Information
------------------

Role by [lint707](https://github.com/lint707).

Dear contributors, thank you.
