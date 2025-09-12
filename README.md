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

#### Typical environment file
... would look like this:
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

To create the environment file, you simply run `pgtools env create`. You will be prompted for the values to put

#### Other `env` command outputs
- List environment files
```bash
[23:08:41|jfgratton@london:src]: /opt/bin/pgtools env ls
Number of environment files: 2
┏━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━┓
┃ Environment file ┃ File size ┃ Modification time   ┃
┣━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━┫
┃ defaultEnv.json  ┃ 184       ┃ 2025/09/10 22:32:52 ┃
┃ test.json        ┃ 244       ┃ 2025/09/10 21:12:51 ┃
┗━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━┛
```
- Get details on an environment file named `test.json`
```bash
[23:08:48|jfgratton@london:src]: /opt/bin/pgtools env info test
┏━━━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━┳━━━━━━━━━━━┳━━━━━━━━━━━━━┳━━━━━━━━━━━━━┳━━━━━━━━━━━━┳━━━━━━━━━━━━┳━━━━━━━━━━━━━━━━━━━━━━━┓
┃ Environment file ┃ DB server host ┃ DB server port ┃ DB user   ┃ DB password ┃ SSL enabled ┃ SSL cert   ┃ SSL key    ┃ Description           ┃
┣━━━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━╋━━━━━━━━━━━╋━━━━━━━━━━━━━╋━━━━━━━━━━━━━╋━━━━━━━━━━━━╋━━━━━━━━━━━━╋━━━━━━━━━━━━━━━━━━━━━━━┫
┃ test.json        ┃ london         ┃ 5432           ┃ jfgratton ┃ *ENCODED*   ┃ require     ┃ london.crt ┃ london.key ┃ test environment file ┃
┗━━━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━┻━━━━━━━━━━━┻━━━━━━━━━━━━━┻━━━━━━━━━━━━━┻━━━━━━━━━━━━┻━━━━━━━━━━━━┻━━━━━━━━━━━━━━━━━━━━━━━┛
```
- Remove the `test.json` environment file:
```bash
[23:15:28|jfgratton@london:src]: /opt/bin/pgtools env rm test
test.json removed succesfully
```
Again, more information is available with `pgtools env -h` or `pgtools env $SUBCOMMAND -h`.

### Backup one or many databases

*TODO*
complete doc