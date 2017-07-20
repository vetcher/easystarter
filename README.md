# Easy Starter
Tool to manage your services

## Install
`go get github.com/vetcher/easystarter`

## Commands

| Title              | Commands (means the same)                  | Description                                                                                                                            | Parameters                                                    |   |
|--------------------|--------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------|---------------------------------------------------------------|--:|
| Exit               | exit, e, ext, out, end, break, close, quit | Exit program, services may skill work                                                                                                  |                                                               |   |
| Start service      | start, s, up                               | Start specified service or start all                                                                                                   | `all` or service name and command line args for start service |   |
| Stop service       | stop, kill, k, down                        | Stop specified service or stop all                                                                                                     | `all` or service name                                         |   |
| Reload service     | reload, r                                  | Stop and start service.  If parameter all specified, program will reload configuration from `services.json` file before start services | `all` or service name                                         |   |
| List services      | list, ps                                   | Print all services, their args and status.                                                                                             |                                                               |   |
| List environment   | env, vars                                  | Print environment variables from `env.ini` file or all                                                                                 | `all`                                                         |   |
| Reload environment | reload env, reenv                          | Reload environment variables from `env.ini` file                                                                                       |                                                               |   |
| Help               | help, h                                    | Show some help information                                                                                                             |                                                               |   |