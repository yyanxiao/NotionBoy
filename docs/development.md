
# Development in local env

1. Run ngrok to reverse proxy local 8001
    ```
    docker run --net=host --rm --name ngrok -it ngrok/ngrok http 8001
    ```

2. 配置[微信测试号](https://mp.weixin.qq.com/debug/cgi-bin/sandbox?t=sandbox/index)
3. 更新 `settings_local.yaml`
4. `make buildd` 构建运行镜像
5. `make rund` 运行镜像


# Code

## configuration

配置使用 viper 读取,
1. 优先级，从上往下，下面的会覆盖上层的配置
   1. 配置文件，配置文件的 `key` 不区分大小写
      - 根目录下的 `settings.[yaml|json|toml|...]` 文件
      - 根目录下的 `settings_local.[yaml|json|toml|...]` 文件

   2. 环境变量, 不区分大小写
     - 环境变量拥有 prefix: `APP_`
     - 环境变量会将 `_` 转换成 `.`, 例如 `APP_TELEGRAM_TOKEN` 对应配置文件中的 `telegram.token`

## logger

使用 zap
- 有两个 Logger:  `Logger` 和 `SugaredLogger`
  - `SugaredLogger` 会自动判断类型
  - `Logger` 需要使用内置的类型，不需要解析类型，所以性能特别好
- 各个级别分别有四个个 methods
  - `Info`: 字符串拼接
  - `Infof`: 带模版的字符串拼接
  - `Infow`: 支持结构化的日志
  - `Infoln`: 字符串拼接，结尾自动换行
