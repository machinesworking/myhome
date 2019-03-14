"use strict";
function MyPeer(peername, ws){

this.localname = peername;


 let pc = new RTCPeerConnection({
  iceServers: [
    {
      urls: 'stun:stun.l.google.com:19302'
    }
  ]
})
let log = msg => {
//  document.getElementById('vstat').innerHTML += peername+": "+msg + '<br>'
console.log(msg);
}

pc.ontrack = function (event) {
  var el = document.createElement(event.track.kind)
  el.srcObject = event.streams[0]
  el.autoplay = true
  el.controls = true
  el.width = "640"
  el.height = "480"
  document.getElementById(peername).appendChild(el)
}

pc.oniceconnectionstatechange = e => log(pc.iceConnectionState)
pc.onicecandidate = event => {
if(event.candidate === null){
  	var msg = new Object();
   		msg.Type = "control";
   		msg.Targets = [peername];
   		msg.Msg =  {Type: "offer", Value: {from: localid, data: pc.localDescription.sdp}};
      ws.send(JSON.stringify(msg));
}
}
console.log("creating offer")
pc.createOffer({offerToReceiveVideo: true, offerToReceiveAudio: true}).then(d => pc.setLocalDescription(d)).catch(log)
return pc;
/*
pc.startsess = function(d){
  try {
    pc.setRemoteDescription(new RTCSessionDescription({type: 'answer', sdp: d}))
  } catch (e) {
    alert(e)
  }
}
*/
}
