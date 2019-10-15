# Catapult
Catapult aims to provide a somewhat full-featured management system to manage firecracker VMs.

# Prerequisites
You need to have postgresql installed, and a ```db.toml``` file configured in a way that points to it.

Create a user
```sql
  create user catapult password 'catapult';
```

Create a database
```sql
  create database catapult owner catapult template template0
encoding 'UTF8' lc_collate 'en_US.UTF-8' lc_ctype 'en_US.UTF-8';

```

Install pgcrypto extention in our database
```sql
  \c catapult
  CREATE EXTENSION "pgcrypto";
```

# Run
```shell script
  mv db.example.toml db.toml # Only if you run a local pg server
  ./catapult migrate init # Only on first run
  ./catapult migrate
  ./catapult serve
```

# FAQ
> Will we succeed?

Unclear at the moment. Only time will tell.
