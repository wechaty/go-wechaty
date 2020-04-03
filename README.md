# go-wechaty

![Go Wechaty](https://wechaty.github.io/go-wechaty/images/go-wechaty.png)

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
> &mdash; <cite>@lijiarui, CEO of BotOrange.</cite>
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

	"github.com/wechaty/go-wechaty/src/wechaty"
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
go get github.com/wechaty/go-wechaty
```

## See Also

## History

### master

### v0.0.1 (Mar 12, 2020)

Project created.

## Related Projects

- [Wechaty](https://github.com/wechaty/wechaty) - WeChat Bot SDK for Individual Account in TypeScript
- [Python Wechaty](https://github.com/wechaty/python-wechaty) - Python WeChat Bot SDK for Individual Account.
- [Java Wechaty](https://github.com/wechaty/java-wechaty) - Java WeChat Bot SDK for Individual Account.

## Author

[Huan LI](https://github.com/huan) ([李卓桓](http://linkedin.com/in/zixia)) zixia@zixia.net

[![Profile of Huan LI (李卓桓) on StackOverflow](https://stackexchange.com/users/flair/265499.png)](https://stackexchange.com/users/265499)

## Copyright & License

- Code & Docs © 2020-now Huan LI \<zixia@zixia.net\>
- Code released under the Apache-2.0 License
- Docs released under Creative Commons
