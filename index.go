package main

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<title>Embed markdown</title>
<meta charset="utf-8">
<meta name="description" content="">
<meta name="author" content="">
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/materialize/0.100.2/css/materialize.min.css">
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/prism/1.8.4/themes/prism.min.css">
</head>
	<body>
		<div class="container">
			<h1 class="header">Embed Markdown</h1>
			<div style="font-size: x-large;">
			Markdown URL to embed:<input type="text" id="url"><br>
			<span id="sync-button" class="btn-large waves-effect waves-light orange">sync</span>&nbsp;
			<span id="async-button" class="btn-large waves-effect waves-light orange">async</span>
			</div>
			<div id="src" style="display: none;">
			<pre class="language-html" id="out"></pre>
			<div id="msg"></id>
			</div>
			<h1>Fork me on github</h1>
			<a href="https://github.com/ledyba/embed-markdown">https://github.com/ledyba/embed-markdown</a>
		</div>
		<script async>
		// <--
		document.addEventListener("DOMContentLoaded", function(event) {
			const outElem = document.getElementById("out");
			const msgElem = document.getElementById("msg");
			const srcElem = document.getElementById("src");

			const syncFn = function() {
				srcElem.style.display="block";

				const url = document.getElementById("url").value;
				const outUrl = location.href+"?"+url;
				outElem.innerText = "<script src=\""+outUrl+"\"><\/script>";

				document.getSelection().selectAllChildren(outElem);
				if(document.execCommand('copy')) {
					msg.innerText = "Copied!";
				}else{
					msg.innerText = "Please copy this snippet.";
				}
			};
			const asyncFn = function() {
				srcElem.style.display="block";

				const url = document.getElementById("url").value;
				const elemId = Math.random().toString(36).slice(-8)+Math.random().toString(36).slice(-8);
				const outUrl = location.href+"async/"+elemId+"?"+url;

				outElem.innerText = "<div id=\""+elemId+"\"><\/div>\n<script src=\""+outUrl+"\" async><\/script>";

				document.getSelection().selectAllChildren(outElem);
				if(document.execCommand('copy')) {
					msg.innerText = "Copied!";
				}else{
					msg.innerText = "Please copy this snippet.";
				}
			};
			document.getElementById("sync-button").addEventListener("click", syncFn);
			document.getElementById("async-button").addEventListener("click", asyncFn);
			});
		// -->
		</script>
		
		<script type="text/javascript" src="https://code.jquery.com/jquery-3.2.1.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/materialize/0.100.2/js/materialize.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.8.4/prism.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.8.4/components/prism-markup.js"></script>
		</body>
</html>
`)
}
