//DATA
const table = document.querySelector("#config tbody");
const errorArea = document.querySelector("#errorArea");
const errorList = document.querySelector("#errorList");
const verifyConfig = document.querySelector("#verifyConfig")

var modalToDelete = null;

var data = [
    {
      name: "scan_name",
      type: "ping",
      address: "127.0.0.1",
      timeout: 10,
    }
];

function renderTable() {
    verifyChange(false);
    table.innerHTML = "";

    for (let i = 0; i < data.length; i++) {
        let row = document.createElement("tr");

        let index = document.createElement("td");
        index.innerHTML = "#" + (i + 1);
        row.appendChild(index);

        let name = document.createElement("td");
        let name_input = document.createElement("input");
        name_input.classList.add("form-control");
        name_input.value = data[i].name;
        name.appendChild(name_input);
        row.appendChild(name);

        let type = document.createElement("td");
        let select = document.createElement("select");

        let options = ["ping", "http"];
        for (let j = 0; j < options.length; j++) {
            let option = document.createElement("option");
            option.value = options[j];
            option.text = options[j];
            select.appendChild(option);
        }
        select.value = data[i].type;
        select.classList.add("form-select");
        type.appendChild(select);
        row.appendChild(type);

        let address = document.createElement("td");
        let address_input = document.createElement("input");
        address_input.classList.add("form-control");
        address_input.value = data[i].address;
        address.appendChild(address_input);
        row.appendChild(address);

        let timeout = document.createElement("td");
        let timeout_input = document.createElement("input");
        timeout_input.classList.add("form-control");
        timeout_input.type = "number";
        timeout_input.value = data[i].timeout;
        timeout.appendChild(timeout_input);
        row.appendChild(timeout);

        let status = document.createElement("td");
        let status_input = document.createElement("input");
        status_input.classList.add("form-control");
        if (data[i].type === "http") {
            status_input.value = data[i].status_code || "";
        } else {
            status_input.value = "-";
            status_input.disabled = true;
        }
        status.appendChild(status_input);
        row.appendChild(status);

        let keyword = document.createElement("td");
        let keyword_input = document.createElement("input");
        keyword_input.classList.add("form-control");
        if (data[i].type === "http") {
            keyword_input.value = data[i].keyword || "";
        } else {
            keyword_input.value = "-";
            keyword_input.disabled = true;
        }
        keyword.appendChild(keyword_input);
        row.appendChild(keyword);

        let detele = document.createElement("td");
        let button = document.createElement("button");
        button.classList.add("btn", "btn-danger");
        button.innerHTML = "<i class='bi bi-trash3-fill'></i>";
        button.onclick = function () {
            modalToDelete = i;
            genModal(resetFavModal);
        };
        detele.appendChild(button);
        row.appendChild(detele);
        
        table.appendChild(row);
    }

    setSwitchers();
    setUpdaters();
}

function setSwitchers() {
    let selects = table.querySelectorAll("select");
    for (let i = 0; i < selects.length; i++) {
        selects[i].addEventListener("change", function () {
            data[i].type = this.value;
            if (this.value === "ping") {
                data[i].status_code = "-";
                data[i].keyword = "-";
            } else {
                data[i].status_code = "";
                data[i].keyword = "";
            }

            renderTable();
        });
    }
}

function setUpdaters() {
    let inputs = table.querySelectorAll(".form-control");
    for (let i = 0; i < inputs.length; i++) {
        inputs[i].addEventListener("change", function () {
            verifyChange(false);

            let row = Math.floor(i / 5);
            let col = i % 5;
            switch (col) {
                case 0:
                    data[row].name = this.value;
                    break;
                case 1:
                    data[row].address = this.value;
                    break;
                case 2:
                    data[row].timeout = parseInt(this.value);
                    break;
                case 3:
                    if (data[row].type === "http") {
                        if (this.value === "") {
                            delete data[row].status_code;
                        } else {
                            data[row].status_code = this.value;
                        }
                    }
                    break;
                case 4:
                    if (data[row].type === "http") {
                        if (this.value === "") {
                            delete data[row].keyword;
                        } else {
                            data[row].keyword = this.value;
                        }
                    }
                    break;
            }
        });
    }
}

function verifyCheck() {
    let check_ok = true;
    errorArea.classList.add("d-none");
    errorList.innerHTML = "";

    for (let i = 0; i < data.length; i++) {
        //names
        do {
            if (data[i].name === "") {
                errorList.innerHTML += "<li>Scan name cannot be empty</li>";
                check_ok = false;
                break;
            }

            if (data[i].name.length > 32) {
                errorList.innerHTML += "<li>Scan name must be less than 32 characters (scan name: " + data[i].name + ").</li>";
                check_ok = false;
                break;
            }

            if (!/^[a-z0-9_]+$/g.test(data[i].name)) {
                errorList.innerHTML += "<li>Scan name can only contain lowercase letters, digits and underscores (scan name: " + data[i].name + ").</li>";
                check_ok = false;
                break;
            }

            for (let j = 0; j < data.length; j++) {
                if (j === i) {
                    continue;
                }
                if (data[j].name === data[i].name) {
                    errorList.innerHTML += "<li>Scan name must be unique (scan name: " + data[i].name + ").</li>";
                    check_ok = false;
                    break;
                }
            }
        } while (false);

        //address
        do {
            if (data[i].address === "") {
                errorList.innerHTML += "<li>Address cannot be empty (scan name: " + data[i].name + ").</li>";
                check_ok = false;
                break;
            }

            //if longer than 256 characters
            if (data[i].address.length > 256) {
                errorList.innerHTML += "<li>Address must be less than 256 characters (scan name: " + data[i].name + ").</li>";
                check_ok = false;
                break;
            }

        } while (false);

        //timeout
        do {
            if (data[i].timeout < 0 || data[i].timeout > 30000) {
                errorList.innerHTML += "<li>Timeout must be between 0 and 30000 (scan name: " + data[i].name + ").</li>";
                check_ok = false;
                break;
            }

            if (!Number.isInteger(data[i].timeout)) {
                errorList.innerHTML += "<li>Timeout must be an integer (scan name: " + data[i].name + ").</li>";
                check_ok = false;
                break;
            }
        } while (false);

        //status_code
        do {
            if (data[i].type === "http") {
                if (!data[i].hasOwnProperty("status_code") || data[i].status_code === "") {
                    break;
                }

                if (data[i].status_code.length > 256) {
                    errorList.innerHTML += "<li>Status code must be less than 256 characters (scan name: " + data[i].name + ").</li>";
                    check_ok = false;
                    break;
                }

                if (!/^[0-9,]+$/g.test(data[i].status_code)) {
                    errorList.innerHTML += "<li>Status code can only contain digits and commas (scan name: " + data[i].name + ").</li>";
                    check_ok = false;
                    break;
                }

                let codes = data[i].status_code.split(",");
                for (let j = 0; j < codes.length; j++) {
                    let code = parseInt(codes[j]);
                    if (code < 100 || code > 599 || !Number.isInteger(code)) {
                        errorList.innerHTML += "<li>Status code must be an integer between 100 and 599 (scan name: " + data[i].name + ").</li>";
                        check_ok = false;
                        break;
                    }
                }
            }
        } while (false);
    }

    if (check_ok) {
        errorArea.classList.add("d-none");
    } else {
        errorArea.classList.remove("d-none");
    }
    
    verifyChange(check_ok);
    return check_ok;
}

function verifyChange(state) {
    if (state) {
        verifyConfig.classList.remove("btn-warning");
        verifyConfig.classList.add("btn-success");
        verifyConfig.querySelector("span").innerHTML = "Verification Successful";
        verifyConfig.querySelector("i").classList.remove("bi-exclamation-triangle");
        verifyConfig.querySelector("i").classList.add("bi-check-circle");
    } else {
        verifyConfig.classList.remove("btn-success");
        verifyConfig.classList.add("btn-warning");
        verifyConfig.querySelector("span").innerHTML = "Verify Values";
        verifyConfig.querySelector("i").classList.remove("bi-check-circle");
        verifyConfig.querySelector("i").classList.add("bi-exclamation-triangle");
    }
}

//RUNS
renderTable();

window.onbeforeunload = function() {
    return "Data may be lost if you leave the page, are you sure?";
};

setTimeout(function () {
    var popoverTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="popover"]'))
    popoverTriggerList.map(function (popoverTriggerEl) {
        return new bootstrap.Popover(popoverTriggerEl)
    })
}, 200);

//disable focus warning
document.addEventListener("DOMContentLoaded", function () {
    document.addEventListener('hide.bs.modal', function (event) {
        if (document.activeElement) {
            document.activeElement.blur();
        }
    });
});

//EVENT LISTENERS
document.querySelector("#uploadConfig").addEventListener('change', event => {
    const file = event.target.files[0];
    if (file) {
        const reader = new FileReader();
        reader.onload = (e) => {
            const content = e.target.result;
            data = parseConfig(content);
            renderTable();
        };
        reader.readAsText(file);
    }
});

document.querySelector("#downloadConfig").addEventListener("click", function () {
    if (verifyCheck()) {
        downloadConfig();
    }
});

verifyConfig.addEventListener("click", function () {
    verifyCheck();
});

document.querySelector("#addRow").addEventListener("click", function () {
    data.push({
        name: "new_scan_name",
        type: "ping",
        address: "127.0.0.1",
        timeout: 10,
    });
    renderTable();
});


//UPLOAD
function parseConfig(content) {
    return content.split('\n').map(line => {
        const parts = line.split(' ');
        const item = {
            name: parts[0],
            type: parts[1],
            address: parts[2],
        };
        parts.slice(3).forEach(part => {
            const [key, value] = part.split('=');
            if (key === 'timeout') {
                item.timeout = parseInt(value, 10);
            } else if (key === 'status_code') {
                item.status_code = value.replace(/\"/g, '');
            } else if (key === 'keyword') {
                let position = line.indexOf("keyword");
                item.keyword = line.slice(position + 9, -1);
            }
        });
        return item;
    });
}


//DOWNLOAD
function downloadConfig() {
    const configContent = convertToConfig(data);
    const blob = new Blob([configContent], { type: 'text/plain' });
    const link = document.createElement('a');
    link.href = URL.createObjectURL(blob);
    link.download = 'config.conf';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
}

function convertToConfig(data) {
    return data.map(item => {
        let line = `${item.name} ${item.type} ${item.address} timeout=${item.timeout}`;
        if (item.status_code) {
            line += ` status_code=\"${item.status_code}\"`;
        }
        if (item.keyword) {
            line += ` keyword=\"${item.keyword}\"`;
        }
        return line;
    }).join('\n');
}
