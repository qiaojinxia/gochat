function toggleSidebar() {
    let sidebar = document.querySelector('.sidebar');
    sidebar.classList.toggle('show'); // 切换show类名
}


let isClick = 0;
// 点击页面任意位置时触发事件
document.addEventListener('click', function(event) {
    // 如果点击的不是侧边栏或按钮，则隐藏侧边栏
    let sidebar = document.getElementById('hbox');
    let toggleButton = document.getElementById('setting');
     if (!sidebar.contains(event.target) && !toggleButton.contains(event.target)) {
         sidebar.classList.remove('show'); // 切换show类名
         sidebar.classList.add('hide');
         isClick = 0;
     }else if (toggleButton.contains(event.target)) {
         if (isClick === 1) {
             isClick = 0;
             sidebar.classList.remove('show');
             sidebar.classList.add('hide');
         } else if(isClick === 0) {
             isClick = 1;
             sidebar.classList.remove('hide');
             sidebar.classList.add('show');
         }
     }

});


function updateRangTemperature(obj) {
    let numStr = obj.value.trim();
    let numFloat = parseFloat(numStr).toFixed(2);
    document.getElementById("rang_label").innerHTML = numFloat;
}

let slider = document.getElementById("tem_slider");
updateRangTemperature(slider);
slider.addEventListener("input", function() {
    updateRangTemperature(slider);
});


const mySwitch = document.getElementById("mySwitch");
// 添加 click 事件监听器
mySwitch.addEventListener("click", function() {
    // 检查开关是否被打开
    if (this.checked) {
        // 如果是，则设置背景颜色为绿色，并更新状态文本为 ON
        const s = document.getElementById("statusColor");
        s.style.backgroundColor = "green";
    } else {
        const s = document.getElementById("statusColor");
        // 否则，将背景颜色重置为默认值（白色），并更新状态文本为 OFF
        s.style.backgroundColor = "gray";
    }
});

// JS Component
class DialogBox {
    constructor(selector,className,imgUrl) {
        this.element = document.querySelector(selector);
        this.header = null;
        this.content = null;
        this.dialogBox = null;
        if (this.element) {
            // 创建头部和内容区域元素
            this.header = document.createElement('div');
            this.header.classList.add('header');
            this.content = document.createElement('div');
            this.content.classList.add('content');
            this.dialogBox = document.createElement('div');
            this.dialogBox.classList.add(className);
            // 添加到容器中
            this.dialogBox.appendChild(this.header);
            this.dialogBox.appendChild(this.content);
            this.element.appendChild(this.dialogBox);
            this.header.style.backgroundImage = "url('" + imgUrl + "')";
            this.header.style.backgroundSize = "cover";
        } else {
            console.error(`Element with selector ${selector} not found.`);
        }
    }
    setTitle(title) {
        if (this.header) {
            // 设置标题文本
            const titleNode = document.createTextNode(title);
            while (this.header.firstChild) {
                // 删除旧节点
                this.header.removeChild(this.header.firstChild);
            }
            this.header.appendChild(titleNode);
        }
    }
    setContent(contentHTML){
        if(this.content){
            this.content.innerHTML= contentHTML;//设置内容区域的 html 内容。
        }
    }
    appendContent(contentHTML){
        if(this.content){
            this.content.innerHTML += contentHTML;//设置内容区域的 html 内容。
        }
    }

    setText(contentText){
            // 创建 Text 节点
            if(this.content){
               let textNode = document.createTextNode(contentText);
               this.p = document.createElement('p');
               this.content.appendChild(this.p);//设置内容区域的 html 内容。
                this.p.appendChild(textNode);//设置内容区域的 html 内容。
            }
    }
    appendText(contentText){
        // 创建 Text 节点
        if(this.content){
            if(!this.p){
                this.p = document.createElement('p');
                let textNode = document.createTextNode(contentText);
                this.content.appendChild(this.p);//设置内容区域的 html 内容。
                this.p.appendChild(textNode);//设置内容区域的 html 内容。
            }else{
                let textNode = this.p.childNodes[0];
                textNode.nodeValue += contentText;
            }
        }
    }

}

function scrollToBottom() {
    //找到消息列表容器
    let messages = $('.messages');

    //将消息列表容器的滚动条滚动到底部位置
    messages.scrollTop(messages.prop("scrollHeight"));
}
async function send() {
    let msg = document.getElementById("chatinput").value;
    let loading = document.getElementById("loading");
    loading.style.display = 'block';

    //检查输入
    if (msg.trim().length === 0) {
        let response = $("<div></div>");
        response.html("🤔<br>Message cannot be null问题不能为空");
        $("#chatlog").append(response);
        loading.style.display = 'none';
        return;
    }
    // 获取元素内容
    let username =  document.getElementById("tt").innerHTML;
    if (username === '' || (username.includes("请输入您的用户名") && username.includes("input") && username.includes("placeholder"))){
        let response = $("<div></div>");
        response.html("🤔<br>Tips: 请先输入用户名");
        $("#chatlog").append(response);
        $('#chatlog').scrollTop($('#chatlog')[0].scrollHeight);
        loading.style.display = 'none';
        return;
    }

    let mySwitch = document.getElementById("mySwitch")

    let openContext = mySwitch.checked === true? '1' : '0';
    let systemRole = roles.find(role => role.name === selectBox.value).content;
    let apiKey = document.getElementById("key").value;

    let dialogBoxYou = new DialogBox('#chatlog','dialog-box1','/static/img/xiaoxin.png');//创建实例对象
    let dialogContentYou=`${msg}`;//构建对话框内容
    dialogBoxYou.setText(dialogContentYou);//将对话框内容添加到组件中显示。
    //对话框拉到底部
    scrollToBottom();

    //对话随机度
    let slider = document.getElementById("tem_slider");
    let numStr = slider.value.trim();
    let temperature = parseFloat(numStr).toFixed(2);


    let url = '/postChatStreamMsg/' + username ;
    let data = {
            msg: msg,
            isOpenContext:openContext,
            systemRole:systemRole,
            apiKey:apiKey,
            temperature:temperature,
    };

    let dialogBoxRobot = new DialogBox('#chatlog', 'dialog-box', '/static/img/chatGpt.png');//创建实例对象
    scrollToBottom();
    await fetch(url, {
        method: "POST",
        keepalive: true,
        body: JSON.stringify(data),
        headers: {
            "Accept-Encoding": "chunked",
            "Content-Type": "application/json"
        },
    }).then(async response => {
        let  reader = response.body.getReader();
        let buffer = '';
        function processChunk(chunk) {
            // 处理当前数据块...
            //获取正常信息时，逐字追加输出
            // let parsedMessageRobot = marked.parse(chunk);

            let dialogContentRobot =`${chunk}`;//构建对话框内容

            dialogBoxRobot.appendContent(dialogContentRobot);//将对话框内容添加到组件中显示。'

            scrollToBottom();

            // 继续读取下一个数据块
            return reader.read().then(processData);
        }

        //Chunk协议解析
        async function processData({done, value}) {
            if (done) return;
            // 将新读取到的字节添加到缓冲区中
            buffer += new TextDecoder().decode(value);
            // 检查缓冲区是否包含完整的分块编码响应。
            while (buffer.includes('\r\n')) {
                let index = buffer.indexOf('\r\n');
                let chunkSizeHex = buffer.slice(0, index);
                let chunkSize = parseInt(chunkSizeHex, 16);
                if (!isNaN(chunkSize)) {
                    let chunkStartIndex = index + '\r\n'.length;
                    let chunkEndIndex = chunkStartIndex + chunkSize;
                    if (buffer.length >= chunkEndIndex + '\r\n'.length) {
                        let dataChunkString = buffer.slice(chunkStartIndex, chunkEndIndex);
                        buffer = buffer.slice(chunkEndIndex + '\r\n'.length);
                        if (buffer.length > 0 && !buffer.endsWith('\r\n')) {
                            console.error('invalid formatter。');
                            await reader.cancel();
                        }
                        processChunk(dataChunkString);
                        continue;
                    }
                }
                break;
            }
            return processData(await reader.read());
        }
        await processData(await reader.read());
    }).catch(
        error => console.error(error)
    );

    let chat_log = document.getElementById("chatlog");
    chat_log.scrollTop = chat_log.scrollHeight;
    loading.style.display = 'none';
}

$(document).ready(function() {
    const input = document.getElementById('username');
    const title = document.getElementById('tt');
    input.addEventListener('blur', function() {
        if (this.value === ''){
            return;
        }
        title.textContent = this.value;
        this.style.display = 'none';
    });

    document.getElementById("chatinput").addEventListener("keyup", function(event) {
        if (event.keyCode === 13) {
            event.preventDefault(); // 阻止默认行为（即提交表单）
            document.getElementById("sendbutton").click();
            this.value = ""; // 清空文本框内容
        }
    });

    $('#sendbutton').click(function() {
        // 获取按钮元素
        let myButton = document.getElementById("sendbutton");
        // 禁用按钮
        myButton.disabled = true;

        send(); // 调用异步函数send()
        //清空对话框
        document.getElementById("chatinput").value = ""
        // 启用按钮
        myButton.disabled = false;
    });

});


function createSelectBox(options, parent) {
    // 创建select和option元素
    const select = document.createElement("select");
    options.forEach((optionText) => {
        const option = document.createElement("option");
        option.text = optionText;
        select.add(option);
    });
    // 添加到指定父元素中
    parent.appendChild(select);
    // 返回select元素
    return select;
}

function createTextarea(parent) {
    // 创建textarea元素
    const textarea = document.createElement("textarea");
    // 添加到指定父元素中
    parent.appendChild(textarea);
    // 返回textarea元素
    return textarea;
}

let roles = [
    { name:"无角色",content: ""},
    { name: "健身教练", content: "我希望你能充当私人教练。我将为你提供一个希望通过体能训练变得更健康、更强壮、更健康的人所需要的所有信息，而你的职责是根据这个人目前的体能水平、目标和生活习惯，为其制定最佳计划。你应该运用你的运动科学知识、营养建议和其他相关因素，以便制定出适合他们的计划。" },
    { name: "外贸开发", content: "我希望你能充当外贸开发。通过tiktok,facebook,Instagram Facebook. WhatsApp 等挖掘潜在客户,并且制定计划推广产品,运用你的知识找到合适的方式渠道推广营销产品,并与客户沟通的建议" },
    { name: "英语翻译", content: "我想让你充当英语翻译员、拼写纠正员和改进员。我会用任何语言与你交谈，你会检测语言，翻译它并用我的文本的更正和改进版本用英语回答。我希望你用更优美优雅的高级英语单词和句子替换我简化的 A0级单词和句子。保持相同的意思，但使它们更文艺。我要你只回复更正、改进，不要写任何解释。" },
    { name: "写作导师", content: "我想让你做一个 AI 写作导师。我将为您提供一名需要帮助改进其写作的学生，您的任务是使用人工智能工具（例如自然语言处理）向学生提供有关如何改进其作文的反馈。您还应该利用您在有效写作技巧方面的修辞知识和经验来建议学生可以更好地以书面形式表达他们的想法和想法的方法。" },
    { name: "催眠治疗师",content: "我希望你能作为一名催眠治疗师。你将帮助病人进入他们的潜意识，并在行为上产生积极的变化，开发技术将客户带入改变的意识状态，使用可视化和放松的方法来引导人们完成强大的治疗体验，并在任何时候都确保病人的安全。"},
    { name: "全栈程序员", content: "我希望你能扮演一个软件开发者的角色。我将提供一些关于网络应用需求的具体信息，而你的工作是提出一个架构和代码，用 Golang 和 Angular 开发安全的应用。" },
    { name: "关系教练" ,content:  "我想让你充当一个关系教练。我将提供一些关于卷入冲突的两个人的细节，而你的工作是提出建议，说明他们如何能够解决使他们分离的问题。这可能包括关于沟通技巧的建议，或改善他们对彼此观点的理解的不同策略。"},
    { name: "内容总结" ,content:  "将以下文字概括为 100 个字，使其易于阅读和理解。避免使用复杂的句子结构或技术术语。"},
    { name :"写作建议" ,content:  "将以下文字概括为 100 个字，使其易于阅读和理解。避免使用复杂的句子结构或技术术语。"},
    {name: "写作标题生成器",content:"我想让你充当书面作品的标题生成器。我将向你提供一篇文章的主题和关键词，你将生成五个吸引人的标题。请保持标题简洁，不超过 20 个字，并确保保持其含义。答复时要利用题目的语言类型。"},
    {name: "写作素材搜集",content:  "生成一份与 [主题] 有关的十大事实、统计数据和趋势的清单，包括其来源"},
    {name: "前端开发",content:"我希望你能担任高级前端开发员。我将描述一个项目的细节，你将用这些工具来编码项目。Create React App, yarn, Ant Design, List, Redux Toolkit, createSlice, thunk, axios. 你应该将文件合并到单一的 index.js 文件中，而不是其他。不要写解释。"},
    {name:"前端：网页设计",content:  "生我希望你能充当网页设计顾问。我将向你提供一个需要协助设计或重新开发网站的组织的相关细节，你的职责是建议最合适的界面和功能，以提高用户体验，同时也满足该公司的业务目标。你应该运用你在 UX/UI 设计原则、编码语言、网站开发工具等方面的知识，为该项目制定一个全面的计划。"},
    {name:"励志教练",content:"我希望你充当一个激励性的教练。我将向你提供一些关于某人的目标和挑战的信息，你的工作是想出可以帮助这个人实现其目标的策略。这可能涉及到提供积极的肯定，给予有用的建议，或建议他们可以做的活动来达到他们的最终目标。"},
    {name:"励志演讲者",content: "我想让你充当一个激励性的演讲者。把激发行动的话语放在一起，让人们感到有能力去做一些超出他们能力的事情。你可以谈论任何话题，但目的是确保你所说的话能引起听众的共鸣，让他们有动力为自己的目标而努力，为更好的可能性而奋斗。"},
    {name:"化妆师",content: "我希望你能成为一名化妆师。你将在客户身上使用化妆品，以增强特征，根据美容和时尚的最新趋势创造外观和风格，提供关于护肤程序的建议，知道如何处理不同质地的肤色，并能够使用传统方法和新技术来应用产品。"},
    {name:"占星家",content: "我希望你能作为一名占星师。你将学习十二星座及其含义，了解行星位置及其对人类生活的影响，能够准确解读星座，并与寻求指导或建议的人分享你的见解。"},
    {name:"医生",content: "我希望你能扮演医生的角色，为疾病想出有创意的治疗方法。你应该能够推荐常规药物、草药疗法和其他自然疗法。在提供建议时，你还需要考虑病人的年龄、生活方式和病史。"},
    {name:"古典音乐作曲家",content: "我想让你充当作曲家。我将提供一首歌的歌词，你将为其创作音乐。这可能包括使用各种乐器或工具，如合成器或采样器，以创造旋律和和声，使歌词变得生动。"},
    {name:"同义词",content: "我希望你能充当同义词提供者。我将告诉你一个词，你将根据我的提示，给我提供一份同义词备选清单。每个提示最多可提供 10 个同义词。如果我想获得更多的同义词，我会用一句话来回答。更多的 x，其中 x 是你寻找的同义词的单词。你将只回复单词列表，而不是其他。词语应该存在。不要写解释。回复 OK 以确认。"},
    {name:"周报生成器",content: "使用下面提供的文本作为中文周报的基础，生成一个简洁的摘要，突出最重要的内容。该报告应以 markdown 格式编写，并应易于阅读和理解，以满足一般受众的需要。特别是要注重提供对利益相关者和决策者有用的见解和分析。你也可以根据需要使用任何额外的信息或来源。"},
    {name:"哲学家",content: "我希望你充当一个哲学家。我将提供一些与哲学研究有关的主题或问题，而你的工作就是深入探讨这些概念。这可能涉及到对各种哲学理论进行研究，提出新的想法，或为解决复杂问题找到创造性的解决方案。"},
    {name:"哲学教师",content: "我希望你充当一名哲学老师。我将提供一些与哲学研究有关的话题，而你的工作是以一种易于理解的方式解释这些概念。这可能包括提供例子，提出问题或将复杂的想法分解成更容易理解的小块。"},
    {name:"商业企划",content: "根据人们的愿望产生数字创业的想法。例如，当我说 [企划目标] 时，你要为数字创业公司生成一份商业计划书，其中包括创意名称、简短的单字、目标用户角色、需要解决的用户痛点、主要价值主张、销售和营销渠道、收入来源、成本结构、关键活动、关键资源、关键合作伙伴、创意验证步骤、预计第一年的运营成本，以及需要寻找的潜在商业挑战。把结果写在一个标记表中。" },
    {name:"小红书风格",content: "请用小红书风格编辑以下中文段落，小红书风格的特点是标题吸引人，每段都有表情符号，并在结尾加上相关标签。请务必保持文本的原始含义。"},
    {name:"小说家",content: "我希望你能作为一个小说家。你要想出有创意的、吸引人的故事，能够长时间吸引读者。你可以选择任何体裁，如幻想、浪漫、历史小说等--但目的是要写出有出色的情节线、引人入胜的人物和意想不到的高潮。我的第一个要求是 小说类型"},
    {name:"广告方案",content: "我想让你充当一个广告商。你将创建一个活动来推广你选择的产品或服务。你将选择一个目标受众，制定关键信息和口号，选择推广的媒体渠道，并决定为达到目标所需的任何额外活动。"},
    {name:"文本情绪分析",content: "文本情绪分析"},
    {name:"新闻评论",content: "我希望你能作为一个评论员。我将为你们提供与新闻有关的故事或话题，你们要写一篇评论文章，对手头的话题提供有见地的评论。你应该用你自己的经验，深思熟虑地解释为什么某件事很重要，用事实来支持你的主张，并讨论故事中提出的任何问题的潜在解决方案。我的第一个要求是 \"新闻评论角度\""},
    {name:"旅游指南",content: "我想让你充当一个旅游向导。我将把我的位置写给你，你将为我的位置附近的一个地方提供参观建议。在某些情况下，我也会给你我要访问的地方的类型。你也将向我推荐靠近我的第一个地点的类似类型的地方。"},

];

const selectWrapper = document.getElementById("select-wrapper");
const selectBox = createSelectBox(roles.map(role => role.name), selectWrapper);

const textareaWrapper = document.getElementById("textarea-wrapper");
const textarea = createTextarea(textareaWrapper);

// 监听select元素的change事件，当选项变化时更新textarea值
selectBox.addEventListener("change", () => {
    let context = roles.find(role => role.name === selectBox.value).content;
    console.log(context);
    textarea.value = context;
});