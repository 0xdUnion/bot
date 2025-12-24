Discord Bot

# 功能

- 监听多个YouTube频道的新视频并在flow发布通知
- 获取硬件占用情况

# 配置

在根目录创建`config.yaml`

~~~yaml
token: abcde123
guild_id: 12345
flow_channel_id: 123456
yt_channels:
  - UCxxxxxxxxxx
  - UCyyyyyyyyyy
~~~

# 服务
~~~toml
[Unit]
Description=Bot Service
After=network.target

[Service]
User=root
WorkingDirectory=/to/home
ExecStart=/to/home/bot
Restart=on-failure


[Install]
WantedBy=multi-user.target
~~~
