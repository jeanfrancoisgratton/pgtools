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
- server inspection (list, activity, locks, extensions, etc.)
- show or edit PostgreSQL configuration

## How to use the tool
All commands follow the same pattern : `pgtools [global flags] command subcommand [subcommand flags] [extra parameters]`.
You can get help on all commands with the help command: `pgtools command -h`, or `pgtools command subcommand -h`

### The environment file
Before using the tool, we need to create an _environment file_. What is that file? It's a way to keep the tool usable across multiple servers/platforms/environments. The file offers a brief description of the environment, along with credentials.

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
All environment files are located in `$HOME/.config/JFG/pgtools/`; the default file is named `defaultEnv.json`.

To create the environment file, simply run:
`pgtools env create`
You will be prompted for the values to put.

#### Other env command outputs
- List environment files : `pgtools env ls`
- Get details on an environment file named `test.json` : `pgtools env info test`
- Remove the `test.json` environment file : `pgtools env rm test`

Again, more information is available with `pgtools env -h` or `pgtools env $SUBCOMMAND -h`.

### Backup one or many databases
Backups are SQL-based, not binary dumps. They can be saved as raw .sql, .sql.gz, or .sql.tgz.

- Backup a single database : `pgtools db backup mydb backup.sql`
- Backup multiple databases : `pgtools db backup db1 db2 db3 alldbs.sql.gz`
- Backup all databases at once : `pgtools db backup -a everything.sql`
- Backup users / roles : `pgtools db backup -u myusers.sql`

Note that the `-a` and `-u` options are mutually exclusive. If the filename ends with .gz the output is gzip-compressed automatically.

### Restore one or many databases
Restore works in reverse. If the archive was compressed, pgtools decompresses automatically.

- Restore one database from a file : `pgtools db restore mydb backup.sql.gz`
- Restore multiple databases : `pgtools db restore db1 db2 alldbs.sql`

If the target database already exists, pgtools will drop and recreate it before restoring, unless you specify flags to change that behavior.

### Roles management
The role command family lets you inspect and modify PostgreSQL roles.

- List all roles : `pgtools role list`
- Create a new role : `pgtools role create myuser --password=secret --createdb`
- Modify an existing role : `pgtools role edit myuser --superuser`
- Remove a role : `pgtools role drop myuser`

### Server (srv) commands
The srv group exposes introspection of server state. This command will be expanded in a upcoming version.

The main reason it was brought in existence was to have a way to reload the server config without firing a server restart.
This is the way: `pgtools srv reload`

Check current commands with the help facility: `pgtools srv help`


### Show commands
The show group contains commands to list schemas, tables, databases, some stats.

More info: `pgtools show -h`

### Config (conf) commands
The conf group wraps around PostgreSQL configuration (SHOW ALL, ALTER SYSTEM).

- List all config parameters (truncated view) : `pgtools conf list`
- Show full values (fields are not truncated) : `pgtools conf list --full`
- Get specific keys : `pgtools conf get KEY1 KEY2 ... KEYn`<br>(ex: `pgtools conf get work_mem shared_buffers`)
- Set a parameter (requires proper privileges) :`pgtools conf set work_mem = '64MB'`

The set command uses ALTER SYSTEM and automatically reloads the configuration.

---

## Shell completion
Cobra’s completion is available for Bash and Zsh.

### One shot, non-persistent
- Bash : `source <(pgtools completion bash)`
- Zsh : `source <(pgtools completion zsh)`

### Persistent config
- Bash : `pgtools completion bash | sudo tee /etc/bash_completion.d/pgtools > /dev/null`
- Zsh persistent with ~/.zsh.completion.d :
```shell
  mkdir ~/.zsh.completion.d && pgtools completion zsh > ~/.zsh.completion.d/_pgtools
  autoload -Uz compinit && compinit
  (this assumes that fpath is properly set in your .zshrc file)
```

---

## Notes
- This tool is intentionally minimal and only covers common admin tasks.
- All commands use pure Go code through the pgx driver; no external binaries like psql or pg_dump are required.

## Building and installing the tool

1. Clone the repo; we'll assume that the cloned repo is now at /repos/pgtools
2. Look at /repos/pgtools/go.version, that specific GO version is the minimal version needed to compile the tool
3. Helper scripts can be found in /repos/pgtools/src. `updateBuildDeps.sh` can be used to update all the packages consumed by this tool.
The main script is `build.sh` . Inspect it until you are satisfied with its settings. This should be enough to build pgtools.