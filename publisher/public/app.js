$(function () {
    let websocket = new WebSocket("ws://"+window.location.host+"/websocket");
    let room = $("#chat-text");
    let user = "";
    websocket.addEventListener("message",function(e){
        let data = JSON.parse(e.data);
        let chatContent = `<p><strong>${data.username}</strong>: ${data.text}</p>`;
        user = data.username;
        room.append(chatContent);
        room.scrollTop = room.scrollHeight;
    });
    $("#input-form").on("submit",function(event){
        event.preventDefault();
        let username = user;
        let text = $("#input-text")[0].value;
        websocket.send(
            JSON.stringify({
                username:username,
                text:text,
            })
        );
        $("#input-text")[0].value = "";
    });
});