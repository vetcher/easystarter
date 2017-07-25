# Easy Starter
Tool to manage your services

## Install
`go get github.com/vetcher/easystarter`

## Commands

|        Title       | Commands (means the same) | Description                                                                                                                            | Parameters                                                     | Other                                                           |
|:------------------:|---------------------------|----------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------|-----------------------------------------------------------------|
| Exit               | `exit`                    | Exit program                                                                                                                           |                                                                | Bug: Запущенные сервисы все еще могут работать после завершения |
| Start service      | `start`                   | Start specified service or start all                                                                                                   | `-all` or service name and command line args for start service |                                                                 |
| Stop service       | `stop`                    | Stop specified service or stop all                                                                                                     | `-all` or service name                                         |                                                                 |
| Restart service    | `restart`                 | Stop and start service.  If flag `-all` specified, program will reload configuration from `services.json` file before start services   | `-all` or service name                                         |                                                                 |
| List services      | `ps`                      | Print all services, their args and status.                                                                                             | `-all`                                                         |                                                                 |
| List environment   | `env`                     | Print environment variables from `env.ini` file or all. With flag `-reload` reloads environment from `env.ini` file                    | `-all` or `-reload`                                            |                                                                 |

## Usage
Program creates `env.ini` file and `logs` folder if it does not exist yet.
You can specify services in file `services.json`, where you may set name, target file with `main()` function and command line arguments for service. For file structure look at `services.json` file in repository.