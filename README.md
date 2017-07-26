# Easy Starter
Tool to manage your services

## Install
`go get github.com/vetcher/easystarter`

## Command line arguments
|Argument      |Description                                                 |
|:------------:|------------------------------------------------------------|
|`-config`     | path to `env.ini` file with environment variables          |
|`-prefix`*    | Prefix of path to `Makefile`                               |
|`-suffix`*    | Suffix of path to `Makefile`                               |
|`-filename`*  | Use this file name instead of `Makefile` when start service|

> Before start new service, tool builds it by instructions described in Makefile's `install` _target_. If you start not configured in `services.json` service,
> program look at `<prefix>/<service name>/<suffix>/<filename>`, that should be makefile with `install` _rule_.

## Commands

|        Title       | Command                   | Description                                                                                                                            | Parameters                                                     | Other                                                           |
|:------------------:|---------------------------|----------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------|-----------------------------------------------------------------|
| Exit               | `exit`                    | Exit program                                                                                                                           |                                                                | Bug: Запущенные сервисы все еще могут работать после завершения |
| Start service      | `start`                   | Start specified service or start all                                                                                                   | `-all` or service name and command line args for start service |                                                                 |
| Stop service       | `stop`                    | Stop specified service or stop all (send `SIGTERM` signal)                                                                             | `-all` or service name                                         |                                                                 |
| Kill service       | `kill`                    | Kill specified service or kill all                                                                                                     | `-all` or service name                                         |                                                                 |
| Restart service    | `restart`                 | Stop and start service.  If flag `-all` specified, program will reload configuration from `services.json` file before start services   | `-all` or service name                                         |                                                                 |
| List services      | `ps`                      | Print all services, their args and status.                                                                                             | `-all`                                                         |                                                                 |
| List environment   | `env`                     | Print environment variables from `env.ini` file or all. With flag `-reload` reloads environment from `env.ini` file                    | `-all` or `-reload`                                            |                                                                 |

## Usage
Program creates `env.ini` file and `logs` folder if it does not exist yet.
You can specify services in file `services.json`, where you may set name, target Makefile with `install` _rule_ and command line arguments for service.
For file structure refer at `services.json` file in repository. __Field `target` required__.
Logs for each service writes to `./logs/<servicename>.log` file.
