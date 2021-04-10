# ghip
simple server to update and show github url ip

This project is inspired by https://github.com/521xueweihan/GitHub520. However, this is a self-hosted server version.

## get started
 your minimum go version should be 1.16, since ```go:embed``` is used. After compiled it, you have to put the binary to place with page subdir and db subdir. Full dir when running is like this:
 ```bash
 .
├── cache
│   ├── host.html
│   └── index.html
├── db
│   └── ipsaves.sqlite
├── ghip
├── ghip.toml
└── page
    └── index.md

 ```
 
 cache dir can generate automatically. However, db and page dir is essential.
 
 You have to run ```ghip -gen``` to generate a config file with default settings. After that, you can run it simply with:
 ```bash
 ghip
 ```
 
 ## About cron job
 ghip will update ip as a cron job. Cron is using the go library [github.com/robfig/cron](https://github.com/robfig/cron), the cron spec format is:

 Field name   | Mandatory? | Allowed values  | Allowed special characters
----------   | ---------- | --------------  | --------------------------
Seconds      | Yes        | 0-59            | * / , -
Minutes      | Yes        | 0-59            | * / , -
Hours        | Yes        | 0-23            | * / , -
Day of month | Yes        | 1-31            | * / , - ?
Month        | Yes        | 1-12 or JAN-DEC | * / , -
Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?

more details see https://pkg.go.dev/github.com/robfig/cron
