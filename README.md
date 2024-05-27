# 五子棋

![](https://img.shields.io/github/go-mod/go-version/CuteReimu/gobang?filename=go%2Fgo.mod) ![](https://img.shields.io/badge/Java-8-informational)

文字讲解（Java版）请前往[doc目录](doc)查看

Java代码仅支持在Windows运行，程序入口在`ChessBoard.java`中的`main`函数。\
Go代码支持在任何平台运行。

## 注意

如果你的电脑计算比较慢，可以将`maxLevelCount`（思考步数）、`maxCountEachLevel`（每一层最多遍历的节点数）、`maxCheckmateCount`（算杀时最多计算的步数）适当改小一些。
