# i18n-scanner
A i18n message scanner, scan and parse messages to JSON format.

Well, I can't find a handy i18n-scanner for vue with customized translate message function `t('bla')`, so I wrote this tool.

## Usage

Scan messages in `src`
```
i18n-scanner -d src
```


Use customized translate function name `$t` (Usually in Vue)
```
i18n-scanner -d src -k '\$t'

# i18n-scanner use regular expression to parse message, you have to escape the function name manually
```

Use customized translate function name `__`
```
i18n-scanner -d src -k __
```

Specify output languages to `en`, `zh` and `jp`
```
i18n-scanner -d src -l "en,zh,jp"
```

Use previous saved messages file, and keep translated messages.
```
i18n-scanner -d src -m my-myessages.json
```

## Output Format
```
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

