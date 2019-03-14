var ws;
var peers=[{}];
var localid;




if ("WebSocket" in window) {


				function init(){

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


if(packet.Msg.Type == "id"){
	localid = packet.Msg.Value;
}


if(packet.Msg.Type == "peerlist" ){
var peercnt = packet.Msg.Value.length;
for(cnt=0; cnt < peercnt; ++cnt){
	var peername = packet.Msg.Value[cnt];
	console.log("creating connection for "+peername);

	var peer = new MyPeer( peername, ws );
	peers[cnt]={sessionid: peername, peer: peer};
}
var tableelt = document.getElementById("videotable");
var pcnt=0;
			for (var irow=0; irow<2; irow++) {
				var row = tableelt.insertRow(0);
				for (var icol=0; icol<2; icol++) {
					var vidiv = document.createElement("div");
					vidiv.id = peers[pcnt].sessionid;
					vidiv.innerHTML = vidiv.id+'<br>';
					vidiv.className = 'vidiv';
					++pcnt;
					row.insertCell(0).appendChild(vidiv);
				}
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
}


	};  // end of onmessage

				ws.onclose = function(event) {
			//		statusdiv.innerHTML+="<div class='red'>Closing connection to security server</div>";
				console.log("Websocket closed")

				};

			};  // end of init

			function error(err) { console.log("crap")}



}
