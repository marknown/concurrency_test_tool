# 并发测试工具v1.0.0

## 它能做什么？

本工具通过并发提交从浏览器复制的curl请求，来检测接口的抗并发能力。可用于测试扣减库存/余额提现/活动资格等并发场景。测完可以去redis或者数据库看扣减是否正确。

## 如何使用？

1. 从浏览器中复制你要并发的请求的cURL信息，一般来说在“Copy”下的“Copy of cURL”
2. 双击bin目录下的二进制执行文件 `concurrencyTestTool_`（win,linux,macos)，请根据自己的系统选对应的版本。
3. 会自动打开浏览器，开启一个本地页面。注意：18880端口会被使用到，请在弹出防火墙提示时，允许。
4. 粘贴刚才复制的cURL信息，到第三步页面中的文本框中，点提交。

## 图片示例

![img](./example.png)

## 参数说明

- **并发次数** 同时发几个请求
- **并发间隔** 如果选择了非0的时间，会按此间隔依次按时执行。例如先择30ms，则N次并发请求开始时间分别为0，30ms以后，60ms以后，以此类推。
- **执行时间** 默认是立即执行，如果时间大于当前时间60秒，会后台执行，同时把响应记录到当前二进制文件所在目录，名为 `concurrency_log.txt` 的文件中。
- **提前时间** 如果指定了时间，但又希望第一次执行时，提前一点。可以选择提前时间。这样会在指定时间前N毫秒执行。
