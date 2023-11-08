# go-wechaty

![Go Version](https://img.shields.io/github/go-mod/go-version/wechaty/go-wechaty)
[![Go](https://github.com/wechaty/go-wechaty/workflows/Go/badge.svg)](https://github.com/wechaty/go-wechaty/actions?query=workflow%3AGo)
[![Maintainability](https://api.codeclimate.com/v1/badges/dbae0a43d431b0fccee5/maintainability)](https://codeclimate.com/github/wechaty/go-wechaty/maintainability)

![Go Wechaty](https://wechaty.github.io/go-wechaty/images/go-wechaty.png)

[![Go Wechaty Getting Started](https://img.shields.io/badge/Go%20Wechaty-Getting%20Started-7de)](https://github.com/wechaty/go-wechaty-getting-started)
[![Wechaty in Go](https://img.shields.io/badge/Wechaty-Go-7de)](https://github.com/wechaty/go-wechaty)

## Connecting Chatbots

[![Powered by Wechaty](https://img.shields.io/badge/Powered%20By-Wechaty-brightgreen.svg)](https://github.com/Wechaty/wechaty)

Wechaty is a RPA SDK for Wechat **Individual** Account that can help you create a chatbot in 6 lines of Go.

## Voice of the Developers

> "Wechaty is a great solution, I believe there would be much more users recognize it." [link](https://github.com/Wechaty/wechaty/pull/310#issuecomment-285574472)  
> &mdash; <cite>@Gcaufy, Tencent Engineer, Author of [WePY](https://github.com/Tencent/wepy)</cite>
>
> "太好用，好用的想哭"  
> &mdash; <cite>@xinbenlv, Google Engineer, Founder of HaoShiYou.org</cite>
>
> "最好的微信开发库" [link](http://weibo.com/3296245513/Ec4iNp9Ld?type=comment)  
> &mdash; <cite>@Jarvis, Baidu Engineer</cite>
>
> "Wechaty让运营人员更多的时间思考如何进行活动策划、留存用户，商业变现" [link](http://mp.weixin.qq.com/s/dWHAj8XtiKG-1fIS5Og79g)  
> &mdash; <cite>@lijiarui, Founder & CEO of Juzi.BOT.</cite>
>
> "If you know js ... try Wechaty, it's easy to use."  
> &mdash; <cite>@Urinx Uri Lee, Author of [WeixinBot(Python)](https://github.com/Urinx/WeixinBot)</cite>

See more at [Wiki:Voice Of Developer](https://github.com/Wechaty/wechaty/wiki/Voice%20Of%20Developer)

## Join Us

Wechaty is used in many ChatBot projects by thousands of developers. If you want to talk with other developers, just scan the following QR Code in WeChat with secret code _go wechaty_, join our **Wechaty Go Developers' Home**.

![Wechaty Friday.BOT QR Code](https://wechaty.js.org/img/friday-qrcode.svg)

Scan now, because other Wechaty Go developers want to talk with you too! (secret code: _go wechaty_)

## Usage

```go
package main

import (
	"fmt"
	"github.com/wechaty/go-wechaty/wechaty"
	"github.com/wechaty/go-wechaty/wechaty-puppet/schemas"
	"github.com/wechaty/go-wechaty/wechaty/user"
)

func main() {
	wechaty.NewWechaty().
		OnScan(func(context *wechaty.Context, qrCode string, status schemas.ScanStatus, data string) {
			fmt.Printf("Scan QR Code to login: %s\nhttps://wechaty.github.io/qrcode/%s\n", status, qrCode)
		}).
		OnLogin(func(context *wechaty.Context, user *user.ContactSelf) {
			fmt.Printf("User %s logined\n", user)
		}).
		OnMessage(func(context *wechaty.Context, message *user.Message) {
			fmt.Printf("Message: %s\n", message)
		}).DaemonStart()
}
```

## Requirements

1. Go 1.18+

## Install

```shell
# go get wechaty

go get github.com/wechaty/go-wechaty
```

## Development

```sh
make install
make test
```

## QA
- wechaty-puppet-service: WECHATY_PUPPET_SERVICE_TOKEN not found ? 
  - go-wechaty is the go language implementation of wechaty (TypeScript). Puppet is required to start wechaty, but it is currently known that puppets are written in TypeScript language. In order to enable go-wechaty to use these puppets, we can use wechaty-gateway to convert puppets into grpc service, let go-wechaty connect to the grpc service, go-wechaty -> wechaty-gateway -> puppet, document: https://wechaty.js.org/docs/puppet-services/diy/
  - puppet list: https://wechaty.js.org/docs/puppet-providers/  

## See Also

- [Learn Go in 12 Minutes](https://www.youtube.com/watch?v=C8LgvuEBraI)
- [How to Write Go Code](https://golang.org/doc/code.html)
- [Journey from OO language to Golang - Sergey Kibish @DevFest Switzerland 2018](https://www.youtube.com/watch?v=1ZjvhGfpwJ8)
- [The Go Blog - Publishing Go Modules](https://blog.golang.org/publishing-go-modules)
- [Effective Go](https://golang.org/doc/effective_go.html)

### Golang for Node.js Developer

- [Golang for Node.js Developers - Examples of Golang examples compared to Node.js for learning](https://github.com/miguelmota/golang-for-nodejs-developers)
- [Learning Go as a Node.js Developer](https://nemethgergely.com/learning-go-as-a-nodejs-developer/)
- [Golang Tutorial for Node.js Developers](https://blog.risingstack.com/golang-tutorial-for-nodejs-developers-getting-started/)

## History

### master

### v0.4 (Jun 19, 2020)

Go Wechaty Scala Wechaty **BETA** Released!

Read more from our Multi-language Wechaty Beta Release event from our blog:

- [Multi Language Wechaty Beta Release Announcement!](https://wechaty.js.org/2020/06/19/multi-language-wechaty-beta-release/)

### v0.1 (Apr 03 2020)

1. Welcome our second and third Go Wechaty contributors:
    - Bojie LI (李博杰) <https://github.com/SilkageNet> [#9](https://github.com/wechaty/go-wechaty/pull/9)
    - Chaofei DING (丁超飞) <https://github.com/dchaofei> [#20](https://github.com/wechaty/go-wechaty/pull/20)
1. Enable [GitHub Actions](https://github.com/wechaty/go-wechaty/actions?query=workflow%3AGo)
    1. Enable linting: [golint](https://github.com/golang/lint)
    1. Enable testing: [testing](https://golang.org/pkg/testing/)
1. Add Makefile for easy developing
1. Re-structure module directories: from `src/wechaty` to `wechaty`
1. Rename example bot to `examples/ding-dong-bot.go`

### v0.0.1 (Mar 12, 2020)

1. Project created.
1. Welcome our first Go Wechaty contributor:
    - Xiaoyu DING （丁小雨） <https://github.com/dingdayu> [#2](https://github.com/wechaty/go-wechaty/pull/2)

## Related Projects

- [Wechaty](https://github.com/wechaty/wechaty) - Conversatioanl AI Chatot SDK for Wechaty Individual Accounts (TypeScript)
- [Python Wechaty](https://github.com/wechaty/python-wechaty) - Python WeChaty Conversational AI Chatbot SDK for Wechat Individual Accounts (Python)
- [Go Wechaty](https://github.com/wechaty/go-wechaty) - Go WeChaty Conversational AI Chatbot SDK for Wechat Individual Accounts (Go)
- [Java Wechaty](https://github.com/wechaty/java-wechaty) - Java WeChaty Conversational AI Chatbot SDK for Wechat Individual Accounts (Java)
- [Scala Wechaty](https://github.com/wechaty/scala-wechaty) - Scala WeChaty Conversational AI Chatbot SDK for WechatyIndividual Accounts (Scala)

## Badge

[![Wechaty in Go](https://img.shields.io/badge/Wechaty-Go-7de)](https://github.com/wechaty/go-wechaty)

```md
[![Wechaty in Go](https://img.shields.io/badge/Wechaty-Go-7de)](https://github.com/wechaty/go-wechaty)
```

## Contributors

[![contributor](https://sourcerer.io/fame/huan/wechaty/go-wechaty/images/0)](https://sourcerer.io/fame/huan/wechaty/go-wechaty/links/0)
[![contributor](https://sourcerer.io/fame/huan/wechaty/go-wechaty/images/1)](https://sourcerer.io/fame/huan/wechaty/go-wechaty/links/1)
[![contributor](https://sourcerer.io/fame/huan/wechaty/go-wechaty/images/2)](https://sourcerer.io/fame/huan/wechaty/go-wechaty/links/2)
[![contributor](https://sourcerer.io/fame/huan/wechaty/go-wechaty/images/3)](https://sourcerer.io/fame/huan/wechaty/go-wechaty/links/3)
[![contributor](https://sourcerer.io/fame/huan/wechaty/go-wechaty/images/4)](https://sourcerer.io/fame/huan/wechaty/go-wechaty/links/4)
[![contributor](https://sourcerer.io/fame/huan/wechaty/go-wechaty/images/5)](https://sourcerer.io/fame/huan/wechaty/go-wechaty/links/5)
[![contributor](https://sourcerer.io/fame/huan/wechaty/go-wechaty/images/6)](https://sourcerer.io/fame/huan/wechaty/go-wechaty/links/6)
[![contributor](https://sourcerer.io/fame/huan/wechaty/go-wechaty/images/7)](https://sourcerer.io/fame/huan/wechaty/go-wechaty/links/7)

1. [@SilkageNet](https://github.com/SilkageNet) - Bojie LI (李博杰)
1. [@huan](https://github.com/huan) - Huan LI (李卓桓)

## Creators

- [@dchaofei](https://github.com/dchaofei) - Chaofei DING (丁超飞)
- [@dingdayu](https://github.com/dingdayu) - Xiaoyu DING (丁小雨) 

## Copyright & License

- Code & Docs © 2020 Wechaty Contributors <https://github.com/wechaty>
- Code released under the Apache-2.0 License
- Docs released under Creative Commons

## Thanks
<a href="https://www.jetbrains.com/?from=go-wechaty"><img src="/docs/images/goland.png" width = "75px" height = "75px" alt="goland.png" /></a>
