let state = "start"; // 初始状态为 "start"

function InputAnalysis(msg) {
    let command = ""; // 存储当前命令
    let objects = []; // 存储对象列表
    let content = ""; // 存储内容
    // 循环遍历输入的命令字符串并转移状态
    for (let i = 0; i < msg.length; i++) {
        let char = msg.charAt(i);
        switch (state) {
            case "start":
                if (char === "/") {
                    state = "command";
                }
                break;
            case "command":
                if (char === " ") {
                    state = "objects";
                } else {
                    command += char;
                }
                break;
            case "objects":
                if (char === "+") {
                    objects.push("");
                } else if (char === " ") {
                    state = "content";
                } else {
                    if (objects.length === 0){
                        objects.push("");
                    }
                    objects[objects.length-1] += char;
                }
                break;
            case "content":
                content += char;
                break;
        }
    }
    let res_state = state;
    state = "start";
    return [command,objects,content,res_state];
}

