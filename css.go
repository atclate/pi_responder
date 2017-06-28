package main

const (
	CSS = `.circle {
	width: 50px;
	height: 50px;
    display: inline;
	-moz-border-radius: 50px;
	-webkit-border-radius: 50px;
	border-radius: 50px;
	float: left;
	margin-top: 5px;
	margin-bottom: 5px;
	margin-left: 2px;
	margin-right: 2px;
}
.not-selected {
	box-shadow: 5px 5px 5px #888888;
}
.selected {
	box-shadow: 0px 0px 0px #888888;
}

.center {
    width: 100%;
    text-align: center;
	margin-left: auto;
    margin-right: auto;
    display: inline-block;
    float: left;
}

/*body {
    text-align: center;
}*/

.button {
    display: inline-block;
    margin: 10px;
    -webkit-border-radius: 8px;
    -moz-border-radius: 8px;
    border-radius: 8px;
    -webkit-box-shadow:    0 8px 0 #777777, 0 15px 20px rgba(0, 0, 0, .35);
    -moz-box-shadow: 0 8px 0 #777777, 0 15px 20px rgba(0, 0, 0, .35);
    box-shadow: 0 8px 0 #777777, 0 15px 20px rgba(0, 0, 0, .35);
    -webkit-transition: -webkit-box-shadow .1s ease-in-out;
    -moz-transition: -moz-box-shadow .1s ease-in-out;
    -o-transition: -o-box-shadow .1s ease-in-out;
    transition: box-shadow .1s ease-in-out;
    font-size: 50px;
    color: #fff;
}

.button span {
    display: inline-block;
    padding: 20px 30px;
    background-image: -webkit-gradient(linear, 0% 0%, 0% 100%, from(hsla(119, 0%, 80%, .8)), to(hsla(119, 0%, 70%, .2)));
    background-image: -webkit-linear-gradient(hsla(119, 0%, 80%, .8), hsla(119, 0%, 70%, .2));
    background-image: -moz-linear-gradient(hsla(119, 0%, 80%, .8), hsla(119, 0%, 70%, .2));
    background-image: -o-linear-gradient(hsla(119, 0%, 80%, .8), hsla(119, 0%, 70%, .2));
    -webkit-border-radius: 8px;
    -moz-border-radius: 8px;
    border-radius: 8px;
    -webkit-box-shadow: inset 0 -1px 1px rgba(255, 255, 255, .15);
    -moz-box-shadow: inset 0 -1px 1px rgba(255, 255, 255, .15);
    box-shadow: inset 0 -1px 1px rgba(255, 255, 255, .15);
    font-family: 'Pacifico', Arial, sans-serif;
    line-height: 1;
    text-shadow: 0 -1px 1px rgba(175, 49, 95, .7);
    -webkit-transition: background-color .2s ease-in-out, -webkit-transform .1s ease-in-out;
    -moz-transition: background-color .2s ease-in-out, -moz-transform .1s ease-in-out;
    -o-transition: background-color .2s ease-in-out, -o-transform .1s ease-in-out;
    transition: background-color .2s ease-in-out, transform .1s ease-in-out;
}

.button:hover span {
    text-shadow: 0 -1px 1px rgba(175, 49, 95, .9), 0 0 5px rgba(255, 255, 255, .8);
}

.button:active, .button:focus {
    -webkit-box-shadow:    0 8px 0 #777777, 0 12px 10px rgba(0, 0, 0, .3);
    -moz-box-shadow: 0 8px 0 #777777, 0 12px 10px rgba(0, 0, 0, .3);
    box-shadow:    0 8px 0 #777777, 0 12px 10px rgba(0, 0, 0, .3);
}

.button:active span {
    -webkit-transform: translate(0, 4px);
    -moz-transform: translate(0, 4px);
    -o-transform: translate(0, 4px);
    transform: translate(0, 4px);
}
.left {
	float: left;
}

#winner_name {
	font-size: 70px;
}

`
)
