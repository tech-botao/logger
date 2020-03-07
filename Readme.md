# Logger 

> 重复的逻辑不需要再重复开发
> 简单的东西需要模块化

## install

```json
go get github.com/tech-botao/logger
```

## example

```go

export LOG_FILE = /var/logs/{project}.log

import . github.com/tech-botao/logger

Log.debug("message", interface)

```