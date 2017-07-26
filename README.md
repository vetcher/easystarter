# Easy Starter
Tool to manage your services

## Install
`go get github.com/vetcher/easystarter`

## Command line arguments
|Argument      |Description                                                                                                    |
|:------------:|---------------------------------------------------------------------------------------------------------------|
|`-config <path-to-env.ini>`  | path to `env.ini` file with environment variables.                                             |
|`-filename <filename>`       | Use this file name instead of `Makefile` when start service.                                   |
|`-s={true|false}`            | This flag means start all services after startup. Same as enter `start -all` after run program.|

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
Program creates `logs` folder if it does not exist yet.
Logs for each service writes to `./logs/<servicename>.log` file.

## Service configuration
You can specify services in file `services.json`, where you may set name, target Makefile with `install` _rule_, custom directory (absolute or relative) to service folder and command line arguments for service.
If `services.json` not in current directory, program use file from `$HOME` folder.
For file structure refer at `services.json` file in repository. __Field `target` required__.
