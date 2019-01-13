## What is FIESTA?

FIESTA is a tool created to run side-channel attacks against the HTTPS protocol.
At this moment, only XS-Search and FIESTA attacks have been implementad. However, other attacks such as CRIME, BREACH, TIME and XS-TIME will be implemented in the future.

## How to install FIESTA?

It should be as easy as:

```bash
go get github.com/jselvi/fiesta/cmd/fiesta
```

It hasn't been tested in many platforms, so if you find an issue in your specific choice, please raise an issue.

## How to use FIESTA?

FIESTA tries to mimic Metasploit's look and feel, since it's a well-known tool used by most penetration testers.
Syntax is easy:

```bash
$ fiesta -h
Usage of fiesta:
  -r string
    	Execute resource file (- for stdin)
```

You can run an interactive shell just typing `fiesta` or you can run a rc file (`fiesta -r my.rc`), which basically types for you a list of FIESTA commands included in the rc file.

## How to use each attack?

Including a detailed description of each attack and writing the `doc.go` files is still a pending task.

## References

* OWASP AppSec EU 2018: [Video](https://www.youtube.com/watch?v=u-drLKiyCSo) [Slides](https://2018.appsec.eu/presos/Hacker_FIESTA_Jose-Selvi_AppSecEU2018.pdf)
* EkoParty 2017: [Video](https://www.youtube.com/watch?v=guZ30nXHkUg)

