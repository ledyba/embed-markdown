# EmbedMarkdown

A very simple webapps to embed remote markdown files into your blogs/websites.

# Example

It's used to embed novels hosted in github into [our website](https://hexe.net/).

Examples:
 - [白い狐と黒い猫 (妖精⊸ロケット)](https://hexe.net/2016/11/03222719.php)
 This story is hosted on [github](https://github.com/YorabaTaiju/WhiteFoxAndBlackCat).
 - [インターネットの青い鳥 (妖精⊸ロケット)](https://hexe.net/2017/02/03161419.php) - [github](https://github.com/FairyRockets/Texts/blob/master/blue-bird-in-the-21st-century.md).
 - [キツネとタヌキ (妖精⊸ロケット)](https://hexe.net/2017/06/14220352.php) - [github](https://github.com/FairyRockets/Texts/blob/master/kitsune-to-tanuki.md).

# How to use

## Server

- golang
```bash
$ go get -u "github.com/russross/blackfriday"
$ go get -u "github.com/microcosm-cc/bluemonday"
$ go build -o embed-markdown github.com/ledyba/embed-markdown/...
$ ./embed-markdown -port=8080
```

You can daemonize this process by systemd, upstart, [supervisord](http://supervisord.org/), etc, etc...

## Client

### without async

```html
<body>
...
Please enjoy the story:

<script src="https://<server-location>?<file-to-embed>"></script>

e.g.)
<script src="https://ledyba.org/EmbedMarkdown/?https://raw.githubusercontent.com/YorabaTaiju/WhiteFoxAndBlackCat/master/README.md"></script>
</boby>
```

### async

```html
<body>
...
Please enjoy the story:

<div id="<elem-id>"></div>
<script async src="https://<server-location>/async/<elem-id>?<file-to-embed>"></script>

e.g.)
<div id="Ewm1yvGPjp"></div>
<script src="https://ledyba.org/EmbedMarkdown/async/Ewm1yvGPjp?https://raw.githubusercontent.com/YorabaTaiju/WhiteFoxAndBlackCat/master/README.md"></script>
</boby>
```
