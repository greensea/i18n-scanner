# i18n-scanner
A i18n message scanner, scan and parse messages to JSON format.

Well, I can't find a handy i18n-scanner for vue with customized translate message function `t('bla')`, so I wrote this tool.

## Usage

Scan messages in `src`
```shell
i18n-scanner -d src
```

## Install

Install using go

```shell
go get github.com/greensea/i18n-scanner
```

You can also download the binary at [Release](https://github.com/greensea/i18n-scanner/releases) page.

## More Usage
Use customized translate function name `$t` (Usually in Vue)
```shell
# i18n-scanner use regular expression to parse message, you have to escape the function name manually
i18n-scanner -d src -k '\$t'
```

Use customized translate function name `__`
```shell
i18n-scanner -d src -k __
```

Specify output languages to `en`, `zh` and `jp`
```shell
i18n-scanner -d src -l "en,zh,jp"
```

Use previous saved messages file, and keep translated messages.
```shell
i18n-scanner -d src -m my-myessages.json
```

## Output Format
```json
{
  "en": {
    "Homepage": "",
    "API Reference": ""
  },
  "zh": {
    "Homepage": "",
    "API Reference": ""
  }
}
```

## Why Golang?

emmm...Go is simple and quick?

Well, I am feeling frontend stack is not in the right route, the ideas and code are bulging. I the fronend stack is making simple things complex. I can wrote this tool in an hour, but I afraid I can't do the same job with nodejs.
