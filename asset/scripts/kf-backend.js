var platform = "kf-backend";

var socket = io("ws://goim.smartkf.top:80/", {
    host: "goim.smartkf.top",
    secure: true,
    transports: ['websocket'],
    query: "token=helloworld&platform=kf",
    path: "/socket.io/",
});

socket.on('reply', function (msg) {
    $('#messages').append($('<li>').text("reply-->" + msg));
});

$('form').submit(function () {
    socket.emit('message', JSON.stringify({"msgType": "text", "ip": "127.0.0.1", "content": $("#m").val()}));
    $('#m').val('');
    return false;
});

// 监听断开连接事件
socket.on('disconnect', (reason) => {
    console.log('Disconnected from the server:', reason);
});


$("#push").on("click", function () {
    mockPush();
})

function mockPush() {
    // mock 接收一条消息.
    fetch("http://goim.smartkf.top/api/push", {
        method:"POST",
        headers: {
            "Content-Type":"application/json"
        },
        body: JSON.stringify({
            sessionId: socket.sessionId,
            event: "test",
            data: JSON.stringify({
                msgType: "text",
                content: "hello world",
            })
        })
    })
}