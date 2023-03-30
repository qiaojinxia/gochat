// 将Excel文件解析为Handsontable表格数据的函数
function parseExcel(data) {
    // 新建一个Workbook实例
    let workbook = XLSX.read(data, {type: 'binary'});
    // 获取第一个Sheet的数据
    let sheet = workbook.Sheets[workbook.SheetNames[0]];
    // 将Sheet数据解析为二维数组格式
    let sheetData = XLSX.utils.sheet_to_json(sheet, {header: 1});
    // 返回解析后的表格数据
    return sheetData;
}

// 将表格数据绑定到Handsontable实例中
function bindDataToTable(data) {
    // 获取表格容器元素
    let container = document.getElementById('myTable');
    let hot = new Handsontable(container, {
        data: data,
        rowHeaders: true,
        colHeaders: true,
        colWidths: 100,
        readOnly: true,
        contextMenu:true,
        columnSorting:true,
        manualColumnResize: true,//列宽自动适应
        manualColumnMove: true,//控制列的移动
        copyable:true,
        licenseKey: 'deee2-f2318-21175-0be43-b6a25'
    });
}

// 读取Excel文件并展示解析后的表格数据
function showExcelData(filename) {
    // 创建一个XMLHttpRequest对象
    let xhr = new XMLHttpRequest();
    // XMLHttpRequest的onreadystatechange事件处理函数
    xhr.onreadystatechange = function() {
        if (xhr.readyState === 4 && xhr.status === 200) {
            // 将Excel文件数据解析为表格数据
            let tableData = parseExcel(xhr.response);
            // 将表格数据绑定到Handsontable实例中
            bindDataToTable(tableData);
        }
    };
    // 打开一个GET请求，获取Excel文件数据
    xhr.open('GET', filename, true);
    xhr.responseType = 'arraybuffer';
    // 发送GET请求
    xhr.send();
}
