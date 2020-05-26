# go-wechaty

![Go Version](https://img.shields.io/github/go-mod/go-version/wechaty/go-wechaty)
[![Go](https://github.com/wechaty/go-wechaty/workflows/Go/badge.svg)](https://github.com/wechaty/go-wechaty/actions?query=workflow%3AGo)
[![Maintainability](https://api.codeclimate.com/v1/badges/dbae0a43d431b0fccee5/maintainability)](https://codeclimate.com/github/wechaty/go-wechaty/maintainability)

![Go Wechaty](https://wechaty.github.io/go-wechaty/images/go-wechaty.png)

[![Go Wechaty Getting Started](https://img.shields.io/badge/Go%20Wechaty-Getting%20Started-7de)](https://github.com/wechaty/go-wechaty-getting-started)

## Connecting Chatbots

[![Powered by Wechaty](https://img.shields.io/badge/Powered%20By-Wechaty-brightgreen.svg)](https://github.com/Wechaty/wechaty)

Wechaty is a RPA SDK for Wechat **Individual** Account that can help you create a chatbot in 6 lines of Go.

## WORK IN PROGRESS

Work in progress...

Please come back after 4 weeks...

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

![Wechaty Go Developers' Home](https://wechaty.github.io/wechaty/images/bot-qr-code.png)

Scan now, because other Wechaty Go developers want to talk with you too! (secret code: _go wechaty_)

## The World's Shortest Go ChatBot: 7 lines of Code

```go
package main

import (
	"fmt"

	"github.com/wechaty/go-wechaty/wechaty"
)

func main() {
	_ = wechaty.NewWechaty().
		OnScan(func(qrCode, status string) {
			fmt.Printf("Scan QR Code to login: %s\nhttps://api.qrserver.com/v1/create-qr-code/?data=%s\n", status, qrCode)
		}).
		OnLogin(func(user string) { fmt.Printf("User %s logined\n", user) }).
		OnMessage(func(message string) { fmt.Printf("Message: %s\n", message) }).
		Start()
}
```

## Go Wechaty Developing Plan

We already have Wechaty in TypeScript, It will be not too hard to translate the TypeScript(TS) to Go because [wechaty](https://github.com/wechaty/wechaty) has only 3,000 lines of the TS code, they are well designed and de-coupled by the [wechaty-puppet](https://github.com/wechaty/wechaty-puppet/) abstraction. So after we have translated those 3,000 lines of TypeScript code, we will almost be done.

As we have already a ecosystem of Wechaty in TypeScript, so we will not have to implement everything in Go, especially, in the Feb 2020, we have finished the [@chatie/grpc](https://github.com/chatie/grpc) service abstracting module with the [wechaty-puppet-hostie](https://github.com/wechaty/wechaty-puppet-hostie) implmentation.

The following diagram shows out that we can reuse almost everything in TypeScript, and what we need to do is only the block located at the top right of the diagram: `Wechaty (Go)`.

```ascii
  +--------------------------+ +--------------------------+
  |                          | |                          |
  |   Wechaty (TypeScript)   | |      Wechaty (Go)        |
  |                          | |                          |
  +--------------------------+ +--------------------------+

  +-------------------------------------------------------+
  |                 Wechaty Puppet Hostie                 |
  |                                                       |
  |                (wechaty-puppet-hostie)                |
  +-------------------------------------------------------+

+---------------------  @chatie/grpc  ----------------------+

  +-------------------------------------------------------+
  |                Wechaty Puppet Abstract                |
  |                                                       |
  |                   (wechaty-puppet)                    |
  +-------------------------------------------------------+

  +--------------------------+ +--------------------------+
  |      Pad Protocol        | |      Web Protocol        |
  |                          | |                          |
  | wechaty-puppet-padplus   | |(wechaty-puppet-puppeteer)|
  +--------------------------+ +--------------------------+
  +--------------------------+ +--------------------------+
  |    Windows Protocol      | |       Mac Protocol       |
  |                          | |                          |
  | (wechaty-puppet-windows) | | (wechaty-puppet-macpro)  |
  +--------------------------+ +--------------------------+
```

## Example: How to Translate TypeScript to Go

There's a 100 lines class named `Image` in charge of downloading the WeChat image to different sizes.

It is a great example for demonstrating how do we translate the TypeScript to Go in Wechaty Way:

### Image Class Source Code

- TypeScript: <https://github.com/wechaty/wechaty/blob/master/src/user/image.ts>
- Go: <https://github.com/wechaty/go-wechaty/blob/master/src/wechaty/user/image.go>

If you are interested in the translation and want to look at how it works, it will be a good start from reading and comparing those two `Image` class files in TypeScript and Go at the same time.

## To-do List

- TS: TypeScript
- SLOC: Source Lines Of Code

### Wechaty Internal Modules

1. [ ] Class Wechaty
    - TS SLOC(1160): <https://github.com/wechaty/wechaty/blob/master/src/wechaty.ts>
    - [ ] Code
    - [ ] Unit Tests
    - [ ] Documentation
1. [ ] Class Contact
    - TS SLOC(804): <https://github.com/wechaty/wechaty/blob/master/src/user/contact.ts>
    - [ ] Code
    - [ ] Unit Tests
    - [ ] Documentation
1. [ ] Class ContactSelf
    - TS SLOC(199): <https://github.com/wechaty/wechaty/blob/master/src/user/contact-self.ts>
    - [ ] Code
    - [ ] Unit Tests
    - [ ] Documentation
1. [ ] Class Message
    - TS SLOC(1054): <https://github.com/wechaty/wechaty/blob/master/src/user/message.ts>
    - [ ] Code
    - [ ] Unit Tests
    - [ ] Documentation
1. [ ] Class Room
    - TS SLOC(1194): <https://github.com/wechaty/wechaty/blob/master/src/user/room.ts>
    - [ ] Code
    - [ ] Unit Tests
    - [ ] Documentation
1. [ ] Class Image
    - TS SLOC(60): <https://github.com/wechaty/wechaty/blob/master/src/user/image.ts>
    - [X] Code
    - [ ] Unit Tests
    - [ ] Documentation
1. [ ] Class Accessory
    - TS SLOC(179): <https://github.com/wechaty/wechaty/blob/master/src/accessory.ts>
    - [ ] Code
    - [ ] Unit Tests
    - [ ] Documentation
1. [ ] Class Config
    - TS SLOC(187): <https://github.com/wechaty/wechaty/blob/master/src/config.ts>
    - [ ] Code
    - [ ] Unit Tests
    - [ ] Documentation
1. [ ] Class Favorite
    - TS SLOC(52): <https://github.com/wechaty/wechaty/blob/master/src/user/favorite.ts>
    - [ ] Code
    - [ ] Unit Tests
    - [ ] Documentation
1. [ ] Class Friendship
    - TS SLOC(417): <https://github.com/wechaty/wechaty/blob/master/src/user/friendship.ts>
    - [ ] Code
    - [ ] Unit Tests
    - [ ] Documentation
1. [ ] Class MiniProgram
    - TS SLOC(70): <https://github.com/wechaty/wechaty/blob/master/src/user/mini-program.ts>
    - [ ] Code
    - [ ] Unit Tests
    - [ ] Documentation
1. [ ] Class RoomInvitation
    - TS SLOC(317): <https://github.com/wechaty/wechaty/blob/master/src/user/room-invitation.ts>
    - [ ] Code
    - [ ] Unit Tests
    - [ ] Documentation
1. [ ] Class Tag
    - TS SLOC(190): <https://github.com/wechaty/wechaty/blob/master/src/user/tag.ts>
    - [ ] Code
    - [ ] Unit Tests
    - [ ] Documentation
1. [ ] Class UrlLink
    - TS SLOC(107): <https://github.com/wechaty/wechaty/blob/master/src/user/url-link.ts>
    - [ ] Code
    - [ ] Unit Tests
    - [ ] Documentation

### Wechaty External Modules

1. [ ] Class FileBox
    - TS SLOC(638): <https://github.com/huan/file-box/blob/master/src/file-box.ts>
    - [ ] Code
    - [ ] Unit Tests
    - [ ] Documentation
1. [ ] Class MemoryCard
    - TS SLOC(376): <https://github.com/huan/memory-card/blob/master/src/memory-card.ts>
    - [ ] Code
    - [ ] Unit Tests
    - [ ] Documentation
1. [ ] Class WechatyPuppet
    - TS SLOC(1115): <https://github.com/wechaty/wechaty-puppet/blob/master/src/puppet.ts>
    - [ ] Code
    - [ ] Unit Tests
    - [ ] Documentation
1. [ ] Class WechatyPuppetHostie
    - TS SLOC(909): <https://github.com/wechaty/wechaty-puppet-hostie/blob/master/src/grpc/puppet-client.ts>
    - [ ] Code
    - [ ] Unit Tests
    - [ ] Documentation

## Usage

WIP...

## Requirements

1. Go 1.14+

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

- [Wechaty](https://github.com/wechaty/wechaty) - WeChat Bot SDK for Individual Account in TypeScript (JavaScript)
- [Python Wechaty](https://github.com/wechaty/python-wechaty) - Python WeChat Bot SDK for Individual Account.
- [Java Wechaty](https://github.com/wechaty/java-wechaty) - Java WeChat Bot SDK for Individual Account.

## Contributors

1. [@SilkageNet](https://github.com/SilkageNet) - Bojie LI (李博杰)

## Commiters

- [@dchaofei](https://github.com/dchaofei) - Chaofei DING (丁超飞)
- [@dingdayu](https://github.com/dingdayu) - Xiaoyu DING (丁小雨) 

## Copyright & License

- Code & Docs © 2020 Wechaty Contributors <https://github.com/wechaty>
- Code released under the Apache-2.0 License
- Docs released under Creative Commons
