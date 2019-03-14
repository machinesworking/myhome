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

//var shutdowndlg;
//var shutdowndlgbutton;
var nightmodestate = "off";
var options = {
	    mode: 'tree',
	    modes: ['code', 'form', 'text', 'tree', 'view'], // allowed modes
	    error: function (err) {
	      alert(err.toString());
	    }
	  };



if ("WebSocket" in window) {


				function init(){
					editconfiglist = document.getElementById("editConfigSelect");
					configlist = document.getElementById("configSelect");
					container = document.getElementById("jsoneditor");
					msgtext = document.getElementById("msgtext");


					statusdiv = document.getElementById("statusDiv");


					status = document.getElementById("status");

					editdiv = document.getElementById("editDiv");

					editor = new JSONEditor(container, options);


					configlist.addEventListener("change", function() {
						var name = configlist.options[configlist.selectedIndex].text;
						currentConfig = name;


if(name == "note.json")
	msgtext.style.visibility="visible";
	else
	msgtext.style.visibility="hidden";

							var msg = {Type: "command", Msg:{Name: "loadmsg", Value: {filename: name}}};


							ws.send(JSON.stringify(msg));

					});

					editconfiglist.addEventListener("change", function() {
						loadConfig();
					});


					var MESSAGE_SERVER = (location.protocol == 'https:' ? 'wss' : 'ws') + '://'+ document.domain +':8080/ws';
					ws = new WebSocket(MESSAGE_SERVER);
				//	ws = new WebSocket("ws://172.20.3.50:8080/ws");


				ws.onopen = function(event) {
				//msgarea.nodeValue="connected to server ";

				var msg = {Type:"status", Msg:{ Name: "unknown", Type: "text", Desc: location.hostname, State: "online"}};
				ws.send(JSON.stringify(msg));


					};

				ws.onerror = function(event) {
				};

				ws.onmessage = function(event) {


					var packet = JSON.parse(event.data);

					if(packet.Type ==  "control"){
						if(packet.Msg.Type == "notification"){
						status = packet.Msg.Value;

					var div = document.createElement("Div");
			 		div.innerHTML = status;
			 		 statusdiv.appendChild(div);
			 				statusDiv.scrollTop = statusDiv.scrollHeight;
}

if(packet.Msg.Type == "id"){
	localid = packet.Msg.Value;
}


if(packet.Msg.Type == "devicelist" ){


for(cnt=0; cnt < packet.Msg.Value.length; ++cnt){
	var div = document.createElement("Div");
	div.innerHTML = packet.Msg.Value[cnt].Desc +" "+packet.Msg.Value[cnt].State;
	 statusdiv.appendChild(div);
			statusDiv.scrollTop = statusDiv.scrollHeight;
}


}
/*
if(packet.Msg.Type == "peerlist" ){
var peercnt = packet.Msg.Value.length;
for(cnt=0; cnt < peercnt; ++cnt){
	var peername = packet.Msg.Value[cnt];
	console.log("creating connection for "+peername);
	var peer = new MyPeer( peername, ws );
	peers[cnt]={sessionid: peername, peer: peer};
}
}

if(packet.Msg.Type == "answer"){

for(var z=0; z<peers.length; z++){
	if(peers[z].sessionid === packet.Msg.Value.Type){
console.log("got an answer from "+ packet.Msg.Value.Type)

		try {
	   peers[z].peer.setRemoteDescription(new RTCSessionDescription({type: 'answer', sdp: packet.Msg.Value.Value.sdp}))
	  } catch (e) {
	    alert(e)
	  }
	}


	}
}
*/










				if(packet.Msg.Type == "msglist" ){
				while (configlist.firstChild) {
				    configlist.removeChild(configlist.firstChild);

				}



Object.keys(packet.Msg.Value).forEach(function(key){




					var editoption = document.createElement("option");
					var option = document.createElement("option");
					editoption.innerHTML=key;
					option.innerHTML=key;

					if(option.innerHTML == currentConfig)
						option.selected = true;
					configlist.appendChild(option);
					editconfiglist.appendChild(editoption)
				});


			}

				if (packet.Msg.Type == "configuration"){
				editor.set(packet.Msg.Value);

				}
}


	};  // end of onmessage

				ws.onclose = function(event) {
					statusdiv.innerHTML+="<div class='red'>Closing connection to security server</div>";
				};

			};  // end of init

			function error(err) { console.log("crap")}


	function clearLog(){
		while (statusdiv.firstChild) {
				    statusdiv.removeChild(statusdiv.firstChild);

				}
	}
	function saveConfig(){

		var name = editconfiglist.options[editconfiglist.selectedIndex].text;

		if(name=="New.json"){
			do{
			name = prompt("Please Enter A File Name",".json");
		}while(!(name != "New.json" && name != "" && name.lastIndexOf(".json") != name.length-6 && name.length >=7))

		}

		var data = editor.get();

			var msg = {Type: "command", Msg:{Name: "savemsg", Value:{Name: name, Value: data}}};
			ws.send(JSON.stringify(msg));
			}

	function loadConfig(){
	var name = editconfiglist.options[editconfiglist.selectedIndex].text;

	var msg = {Type: "command", Msg:{Name: "loadmsg", Value: {filename: name}}};
	ws.send(JSON.stringify(msg));
  }

	function doneEditing(){
		var editdiv = document.getElementById("editDiv");
		editDiv.style.display="none";
	}

	function deleteConfig(){
		var name = editconfiglist.options[editconfiglist.selectedIndex].text;
		var msg = {Type: "command", Msg:{Name: "deletemsg", Value:{filename: name}}};
		ws.send(JSON.stringify(msg));
	}


function editConfig(){

	editDiv.style.display="block";
}




function sendmsg(){
	var name = configlist.options[configlist.selectedIndex].text;

var msg = editor.get();

if(name == "note.json")
msg.Msg.Value = msgtext.value;


ws.send(JSON.stringify(msg));

}
}
