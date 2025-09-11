# pgtools : the PostgreSQL swiss-army knife
___

## Overview
The goal of this tool is to provide a very basic PostgreSQL client to perform very specific tasks while not having to go through the installation of multiple binary packages (RPM, DEB, APK) and their dependencies.
This tool is a pure, standalone GO tool. No need for anything else.

## What does this tool do ?
Basic functions, really. Originally, this tool was designed to backup and restore databases, but has expanded a bit since then.
The basic operations are:
- backup one or multiple databases
- restore one or multiple databases
- create an empty database (table creation is outside the scope of this tool)
- delete (drop) one or multiple databases
- roles management (list, create, remove, edit)

## How to use the tool
All commands follow the same pattern : `pgtools [global flags] command subcommand [subcommand flags] [extra parameters]`.
You can get help on all commands with the help command: `pgtools command -h`, or `pgtools command subcommand -h`

### The environment file
Before using the tool, we need to create an _environment file_. What is that file ? It's a way to keep the tool usable across multiple servers/platforms/environments. The file offers a brief description of the environment, along with credentials.
Typically, it looks like this:
```json
{
  "description": "optional comment or description",
  "host": "pgsql server fqdn",
  "port": 5432,
  "user": "my_username",
  "password": "encrypted_password",
  "sslmode": "require",
  "sslcert": "path_to_crt_file",
  "sslkey": "path_to_key_file"
}
```
The last two members are omitted if `sslmode` is set to "disable".
All environment files are located in `$HOME/.config/JFG/pgtools/`; the default file is named `defaultEnv.json`

To create the environment file, you simply run `pgtools env create`.

Again, more information is available with `pgtools env -h` or `pgtools env $SUBCOMMAND -h`. There are other commands to list/delete/describe those environment files.

### Backup one or many databases
