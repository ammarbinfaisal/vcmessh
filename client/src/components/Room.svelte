<script>
    export let roomname;
    import { navigate } from 'svelte-navaid';
    import adapter from "webrtc-adapter";
    import Video from "./Video.svelte";
    
    let password;
    let pcs = {};
    window.pcs = pcs;
    let servers = [{iceServers: [ "stun.l.google.com:19305", "stun1.l.google.com:19305", "stun2.l.google.com:19305", "stun3.l.google.com:19305", "stun4.l.google.com:19305",]}];
    let localStream;
    let remoteStreams = {};
    let size = -1;

    if(!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
        alert("THIS APP WON'T WORK. YOUR BROWSER CAN'T ACCESS YOUR CAMERA");
    }

    const createPeerConn = ({config, websocket, peerId}) => {
        const p = new RTCPeerConnection(servers);
        p.onicecandidate = event => {
            // console.log(event)
            if(event.candidate) {
                websocket.send(JSON.stringify({to: peerId, event: "addIceCandidate", data: JSON.stringify(event.candidate)}));
            }
        };
        
        localStream.getVideoTracks().forEach(track => p.addTrack(track, localStream));

        // p.oniceconnectionstatechange = console.log;
        // p.onsignalingstatechange = console.log
        
        p.ontrack = ev => {
            remoteStreams[peerId] = ev.streams[0];
            console.log(ev);
        }
        return p;
    }

    const manageWebsoocket = ws => {
        ws.onmessage = async ev => {
            try {
                const msg = JSON.parse(ev.data);
                const {from, event, data} = msg;
                console.log(msg);
                if(event === "createOffer") {
                    pcs[from] = createPeerConn({config: servers, peerId: from, websocket: ws});
                    pcs[from].setRemoteDescription(JSON.parse(data));
                    const anwser = await pcs[from].createAnswer();
                    pcs[from].setLocalDescription(anwser);
                    ws.send(JSON.stringify({to: from, event: "createAnswer", data: JSON.stringify(anwser)}));
                } else if(event === "createAnswer") {
                    pcs[from].setRemoteDescription(JSON.parse(data));
                } else if(event === "addIceCandidate") {
                    pcs[from].addIceCandidate(JSON.parse(data));
                }
            } catch(error) {
                console.log(error);
            }
        }
        ws.onopen = async function() {
            if (size > 0) {
                console.log("creating offers")
                for(let i = 0; i < size; ++i) {
                    const p = createPeerConn({config: servers, peerId: i, websocket: ws});
                    const offer = await p.createOffer();
                    p.setLocalDescription(offer);
                    ws.send(JSON.stringify({to: i, event: "createOffer", data: JSON.stringify(offer)}));
                    pcs[i] = p;
                }
                console.log(pcs);
            }
        }
        ws.onclose = console.log
    }

    fetch(`/getRoomSize?roomname=${roomname}`)
        .then(res => res.json())
        .then(async s => {
            size = s.size;
            console.log(`size: ${size}`)
            if(size < 0) {
                alert("Room does not exists.");
                return navigate("/");
            }

            try {
                

            localStream = await navigator.mediaDevices.getUserMedia({
                   audio: true,
                    video: {
                        width: { ideal: 1280 },
                        height: { ideal: 720 }
                    }
                });

                const url = new URL(window.location.href);
                const protocol = url.protocol === "https:" ? "wss" : "ws";
                const ws = new WebSocket(`${protocol}://${url.host}/ws?roomname=${roomname}`);
                manageWebsoocket(ws);
            } catch(error) {
                alert(error);
            }
        }).catch(console.log);
</script>

<div class="h-screen w-screen flex justify-center align-center flex-column">
    <h2>local video</h2>
    <Video stream={localStream} local={true} />
    <h2>foregin videos</h2>
    <div class="flex flex-wrap">
       {#each Object.keys(remoteStreams) as k}
           <Video stream={remoteStreams[k]} />
       {/each} 
    </div>
</div>
