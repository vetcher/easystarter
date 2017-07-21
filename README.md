# Easy Starter
Tool to manage your services

## Install
`go get github.com/vetcher/easystarter`

## Commands

|        Title       | Commands (means the same)                                  | Description                                                                                                                            | Parameters                                                     | Other                                                           |
|:------------------:|------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------|-----------------------------------------------------------------|
| Exit               | `exit`, `e`, `ext`, `out`, `end`, `break`, `close`, `quit` | Exit program                                                                                                                           |                                                                | Bug: Запущенные сервисы все еще могут работать после завершения |
| Start service      | `start`, `s`, `up`                                         | Start specified service or start all                                                                                                   | `-all` or service name and command line args for start service |                                                                 |
| Stop service       | `stop`, `kill`, `k`, `down`                                | Stop specified service or stop all                                                                                                     | `-all` or service name                                         |                                                                 |
| Reload service     | `reload`, `r`                                              | Stop and start service.  If parameter all specified, program will reload configuration from `services.json` file before start services | `-all` or service name                                         |                                                                 |
| List services      | `list`, `ps`, `ls`                                         | Print all services, their args and status.                                                                                             |                                                                |                                                                 |
| List environment   | `env`, `vars`                                              | Print environment variables from `env.ini` file or all                                                                                 | `-all`                                                         |                                                                 |
| Reload environment | `reenv`, `env reload`                                      | Reload environment variables from `env.ini` file                                                                                       |                                                                | Bug: неуказанные переменные не удаляются                        |
| Help               | `help`, `h`                                                | Show some help information                                                                                                             |                                                                |                                                                 |

## Usage
Program creates `env.ini` file and `logs` folder if it does not exist yet.
You can specify services in file `services.json`, where you may set name, target file with `main()` function and command line arguments for service. For file structure look at `services.json` file in repository.