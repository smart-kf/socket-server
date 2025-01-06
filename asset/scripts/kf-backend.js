var platform = "kf-backend";
localStorage.debug = '*'; // 开启 socketio 的 debug 信息.

//  线上
// var socket = io("ws://goim.smartkf.top:80/", {
//     host: "goim.smartkf.top",
//     secure: true,
//     transports: ['websocket'],
//     query: "token=helloworld&platform=kf",
//     path: "/socket.io/",
// });
// 本地.
var socket = io("ws://localhost:9000/", {
    host: "localhost:9000",
    secure: true,
    transports: ['websocket'],
    query: "token=helloworld&platform=kf",
    path: "/socket.io/",
});

socket.on('test', function (msg) {
    $('#messages').append($('<li>').text("reply-->" + msg.content));
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

socket.on("sessionId",function(msg){
    console.log("sessionId-->",msg)
})


function mockPush() {
    // mock 接收一条消息.
    fetch("http://goim.smartkf.top/api/push", {
        method:"POST",
        headers: {
            "Content-Type":"application/json"
        },
        body: JSON.stringify({
            sessionId: socket.io.engine.id,
            event: "test",
            data: JSON.stringify({
                msgType: "text",
                content: "hello world",
            })
        })
    })
}