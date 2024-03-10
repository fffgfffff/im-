#5.4 实现发送文字、表情包等

似乎多此一举设置了udp接收，也可以不要？（发送websocket）

前端user1拼接好数据对象Message
msg={id:1,userid:2,dstid:3,cmd:10,media:1,content:txt}
转化成json字符串jsonstr=JSON.stringify(msg)
后端S在recvproc中接收数据data
并做相应的逻辑处理dispatch(data)-转发给user2
user2通过websocker.onmessage收到消息后做解析并显示

#分布式方案
方案一：消息总线（基于redis/kafka）
优点：简单
缺点：不知道节点状态
方案二：局域网通信协议（udp）
优点：简单
缺点：不知道节点状态
方案三：调度应用（http2/tcp）
优点：可靠
缺点：复杂

