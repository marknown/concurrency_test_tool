# 并发测试工具v1.0.0

## 它能做什么？

本工具通过并发提交从浏览器复制的curl请求，来检测接口的抗并发能力。可用于测试扣减库存/余额提现/活动资格等并发场景。测完可以去redis或者数据库看扣减是否正确。

## 如何使用？

1. 从浏览器中复制你要并发的请求的cURL信息，一般来说在“Copy”下的“Copy of cURL”
2. 双击bin目录下的二进制执行文件 `concurrencyTestTool_`（win,linux,macos)，请根据自己的系统选对应的版本。
3. 会自动打开浏览器，开启一个本地页面。注意：18880端口会被使用到，请在弹出防火墙提示时，允许。
4. 粘贴刚才复制的cURL信息，到第三步页面中的文本框中，点提交。

## 图片示例