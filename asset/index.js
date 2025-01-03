var socket = io("/",{
    // hostname: null,
    // port: null,
    // secure: null,
    transports: ['websocket'],
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