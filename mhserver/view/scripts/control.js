var editor=null;
var ws;
var configlist;
var container;
var editor;
var statusdiv;
var editdiv;
var list;
var currentConfig=null;
var msgtext;
var peers=[{}];
var localid;
var alllights = {Type:"light", Name: "light", Desc: "all lights", State: "online", Address:"all"};
var controlmsgs;

//var shutdowndlg;
//var shutdowndlgbutton;



if ("WebSocket" in window) {


				function init(){


					var MESSAGE_SERVER = (location.protocol == 'https:' ? 'wss' : 'ws') + '://'+ document.domain +':8080/ws';
					ws = new WebSocket(MESSAGE_SERVER);
				//	ws = new WebSocket("ws://172.20.3.50:8080/ws");


				ws.onopen = function(event) {
				//msgarea.nodeValue="connected to server ";

				var msg = {Type:"status", Msg:{ Name: "unknown", Type: "text", Desc: location.hostname, State: "online"}};
				ws.send(JSON.stringify(msg));
				var item = document.createElement("div");
				item.innerHTML = "<b>Status sent</b>";
				appendLog(item);
				addctrllist("camera");
				addctrllist("hmi");



					};

				ws.onerror = function(event) {

					var item = document.createElement("div");
					item.innerHTML = "<b>Websocket error</b>";
					appendLog(item);



				};

				ws.onmessage = function(event) {


					var packet = JSON.parse(event.data);

					if(packet.Type ==  "control"){
if(packet.Msg === null)
return;

/*

						if(packet.Msg.Type == "notification"){
						status = packet.Msg.Value;

					var div = document.createElement("Div");
			 		div.innerHTML = status;
			 		 statusdiv.appendChild(div);
			 				statusDiv.scrollTop = statusDiv.scrollHeight;
}
*/
if(packet.Msg.Type == "id"){
	localid = packet.Msg.Value;
	var  msg1 = {Msg: {Name: "listdevices", Value: {Type: "all"}}, Targets: ["hmi"], Type: "command"};
		ws.send(JSON.stringify(msg1));
}


if(packet.Msg.Type == "devicelist" ){


packet.Msg.Value.push({Type:"light", Name: "light", Desc: "all lights", State: "online", Address:"all"});
packet.Msg.Value.push({Type:"camera", Name: "camera", Desc: "all cameras", State: "online", Address:"all"});
packet.Msg.Value.push({Type:"hmi", Name: "hmi", Desc: "all HMIs", State: "online", Address:"all"});




for(cnt=0; cnt < packet.Msg.Value.length; ++cnt){

var dev = packet.Msg.Value[cnt];
var devtype = dev.Type;
if(devtype === "light" || devtype === "camera" || devtype==="hmi"){
adddev(dev);
} else {
	adddevmisc
}

}
}

				if(packet.Msg.Type == "msglist" ){
					controlmsgs = packet.Msg.Value;
					camctrllist = document.getElementById("camerasel");
					hmictrllist = document.getElementById("hmisel");
				while (camctrllist.firstChild) {
						camctrllist.removeChild(camctrllist.firstChild);

				}
				while (hmictrllist.firstChild) {
						hmictrllist.removeChild(hmictrllist.firstChild);

				}


Object.keys(packet.Msg.Value).forEach(function(key){
if(packet.Msg.Value[key].Targets != null){

if(packet.Msg.Value[key].Targets[0]==="camera"){
	var option = document.createElement("option");
					option.innerHTML=key.substring(0, key.indexOf(".json"));
					option.value=key;
					camctrllist.appendChild(option);

}
if(packet.Msg.Value[key].Targets[0]==="hmi"){
	var option = document.createElement("option");
					option.innerHTML=key.substring(0, key.indexOf(".json"));
					option.value=key;
					hmictrllist.appendChild(option);
}
}
				});


			}

				if (packet.Msg.Type == "configuration"){
				//editor.set(packet.Msg.Value);
				// send the control
				}

}


	};  // end of onmessage

				ws.onclose = function(event) {
					var item = document.createElement("div");
					item.innerHTML = "<b>Websocket has been closed</b>";
					appendLog(item);
				};

			};  // end of init

			function error(err) { console.log("crap")}

			function appendLog(item) {
					var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
					log.appendChild(item);
					if (doScroll) {
							log.scrollTop = log.scrollHeight - log.clientHeight;
					}
			}



			function sendSwitchControl(ev){
devicetype = ev.target.dataset.devtype;

				var checked="off";
				if(ev.target.checked)
				checked = "on";
// gang the switches
if(ev.target.id == devicetype){
	devs = document.getElementById(devicetype+"s").querySelectorAll("input")
	Object.keys(devs).forEach(function(key){
devs[key].checked = ev.target.checked;
});


}

switch(ev.target.dataset.devtype){

		case "light":

ws.send(JSON.stringify({Msg:{Name:"mode",Value:checked},Targets:[ev.target.name],Type:"control"}));

break;


case "camera":
/*
sel = document.getElementById("camerasel");
msgname=sel.options[sel.selectedIndex].text;
msg = controlmsgs[msgname+".json"];
msg.Targets=[ev.target.name];
ws.send(JSON.stringify(msg));

devs = document.getElementById(devicetype+"s").querySelectorAll("input")
Object.keys(devs).forEach(function(key){
devs[key].checked = 0;
});
*/
break;

case "hmi":
		//	ws.send(JSON.stringify({Msg:{Name:"mode",Value:checked},Targets:[ev.target.name],Type:"control"}));

break;

case "misc":
		//	ws.send(JSON.stringify({Msg:{Name:"mode",Value:checked},Targets:[ev.target.name],Type:"control"}));

break;


}
			}



function sendSubmitControl(ev){



	switch(ev.target.dataset.devtype){


	case "camera":
	sel = document.getElementById("camerasel");
	msgname=sel.options[sel.selectedIndex].text;
	msg = controlmsgs[msgname+".json"];
devices = [];
	devs = document.getElementById(devicetype+"s").querySelectorAll("input")
	Object.keys(devs).forEach(function(key){
		devname = devs[key].name
	if(devs[key].checked){
		devices.push(devname);
if(devname === "camera")
devices=["camera"];
}

	devs[key].checked = 0;


	});
msg.Targets=devices;
	ws.send(JSON.stringify(msg));

	break;

	case "hmi":
	sel = document.getElementById("hmisel");
	msgname=sel.options[sel.selectedIndex].text;
	msg = controlmsgs[msgname+".json"];
devices = [];
	devs = document.getElementById(devicetype+"s").querySelectorAll("input")
	Object.keys(devs).forEach(function(key){
		devname = devs[key].name
	if(devs[key].checked){
		devices.push(devname);
if(devname === "hmi")
devices=["hmi"];
}

	devs[key].checked = 0;


	});
msg.Targets=devices;
	ws.send(JSON.stringify(msg));

	break;

	case "misc":
			//	ws.send(JSON.stringify({Msg:{Name:"mode",Value:checked},Targets:[ev.target.name],Type:"control"}));

	break;


	}

}


function addctrllist(devtype){
	var dvdiv = document.getElementById(devtype+"s");
		var div = document.createElement("div");
		div.className = "field";
		var inp = document.createElement("select");
		inp.id = devtype+"sel";
		inp.dataset.devtype=devtype;


	var el = document.createElement("button");
	el.dataset.devtype=devtype;
	el.addEventListener("click", sendSubmitControl);
	el.innerHTML="Submit";


		div.appendChild(inp);
		div.appendChild(el);
		dvdiv.appendChild(div);
}



			function adddev(dev){
			var devdiv = document.getElementById(dev.Type+"s");
			if (devdiv === null)
				devdiv = document.getElementById("miscs");
			var div = document.createElement("div");
			div.className = "field";
			var inp = document.createElement("input");
			inp.type = "checkbox";
			inp.id = dev.Name;
			inp.name = dev.Name;
			inp.innerHTML = dev.Desc;
			inp.className = "switch";
			inp.checked = "";
			inp.addEventListener("click", sendSwitchControl);
			inp.dataset.devtype=dev.Type;
			var lbl = document.createElement("label");
			lbl.htmlFor = dev.Name;
			lbl.innerHTML=dev.Desc;
			div.appendChild(inp);
			div.appendChild(lbl);
			devdiv.appendChild(div);
			}



function adddevmisc(dev){}

function sendmsg(){
	var name = configlist.options[configlist.selectedIndex].text;

var msg = editor.get();

if(name == "note.json")
msg.Msg.Value = msgtext.value;


ws.send(JSON.stringify(msg));

}

function loadConfig(){
var name = editconfiglist.options[editconfiglist.selectedIndex].text;

var msg = {Type: "command", Msg:{Name: "loadmsg", Value: {filename: name}}};
ws.send(JSON.stringify(msg));
}




}
