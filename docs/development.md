
# Development in local env

1. Run ngrok to reverse proxy local 8001
    ```
    docker run --net=host --rm --name ngrok -it ngrok/ngrok http 8001
    ```

2. 配置[微信测试号](https://mp.weixin.qq.com/debug/cgi-bin/sandbox?t=sandbox/index)
3. 更新 `settings_local.yaml`
4. `make buildd` 构建运行镜像
5. `make rund` 运行镜像
