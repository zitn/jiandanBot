# 煎蛋无聊图bot

## todolist

 - [x] 每二十分钟抓取一次无聊图并推送至channel
 - [x] 将多张图片组合在一起发送
 - [x] 显示吐槽以及吐槽的OOXX
 - [x] 点击更新按钮可以拉取最新吐槽
 - [ ] 将大图以文件形式发送避免压缩
 - [ ] OO XX功能


## config.json 格式
```json
    {
    "Token": "YOUR_TOKEN",
    "AdminID": 000000,
    "ChannelID": 000000,
    "ChannelUsername": "@YOUTR_CHANNEL_NAME",
    "ApiAddress": "jandan.net_api"
    }
```
ApiAddress 可通过抓包官方app得到