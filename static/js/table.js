let sheetContent = new Map();
var hot = new Handsontable(document.getElementById('myTable'), {
    data: null,
    colHeaders: true,
    contextMenu:true,
    columnSorting:true,
    filters: true,
    manualColumnResize: true,//列宽自动适应
    manualColumnMove: true,//控制列的移动
    rowHeaders: true, // 默认值是false
    copyable:true,
    // width: window.innerWidth * 0.9, //表宽 多余自动显示滚动条
    mergeCells: [],//数组 格子合并对象{row,col,rowspan,colspan}
    dropdownMenu: true,
    hiddenColumns: {
        indicators: true
    },
    height: window.innerHeight * 0.77,
    multiColumnSorting: true,
    manualRowMove: true,
    afterChange: changeCell,
    afterGetColHeader: alignHeaders,
    afterOnCellMouseDown: changeCheckboxCell,
    beforeRenderer: addClassesToRows,
    licenseKey: 'deee2-f2318-21175-0be43-b6a25'
});

//获取某个sheet
function GetSheetData(filename) {
    let xmlData =  sheetContent.get(filename);
    if (xmlData)
        return parseExcel(xmlData);
}


//获取sheet原始数据
function GetSheetOriginData(filename) {
    let xmlData =  sheetContent.get(filename);
    return xmlData;
}

//保存sheet数据
function SaveSheet(filename,data) {
    sheetContent.set(filename,data);
}


//获取Sheet的Header
function GetSheetHeader(filename) {
    return GetSheetData(filename)[0];
}


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
    if (hot){
        hot.loadData(data.slice(1));
        hot.updateSettings({
            colHeaders:  data[0]
        });
    }
}

// 将表格数据绑定到Handsontable实例中
function clearDataToTable() {
    if (hot)
        hot.loadData(null);
}


// 读取Excel文件并展示解析后的表格数据
function showExcelData(url,filename) {
    let tableData = GetSheetData(filename);
    if(tableData){
        clearDataToTable();
        bindDataToTable(tableData);
        setCurSelect(filename);
        return;
    }

    // 创建一个XMLHttpRequest对象
    let xhr = new XMLHttpRequest();
    // XMLHttpRequest的onreadystatechange事件处理函数
    xhr.onreadystatechange = function() {
        if (xhr.readyState === 4 && xhr.status === 200) {
            // 将Excel文件数据解析为表格数据
            //缓存下数据
            let tableData = parseExcel(xhr.response);
            SaveSheet(filename,xhr.response);

            clearDataToTable();
            // 将表格数据绑定到Handsontable实例中
            bindDataToTable(tableData);
            setCurSelect(filename);
        }
    };
    // 打开一个GET请求，获取Excel文件数据
    xhr.open('GET', url, true);
    xhr.responseType = 'arraybuffer';
    // 发送GET请求
    xhr.send();
}
