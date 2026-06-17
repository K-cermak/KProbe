// ===== DATA & DOM REFS =====
const table = document.querySelector("#config tbody");
const errorArea = document.querySelector("#errorArea");
const errorList = document.querySelector("#errorList");
const validationStatus = document.querySelector("#validationStatus");

var autoValidateTimer = null;

var modalToDelete = null;
var currentExpressionRow = null;
var expressionSnapshot = null;
var expressionSaved = false;

var data = [
    {
        name: "scan_name",
        type: "ping",
        address: "127.0.0.1",
        timeout: 10,
    }
];


// ===== TABLE RENDERING =====
function renderTable() {
    scheduleAutoValidation();
    table.innerHTML = "";

    for (let i = 0; i < data.length; i++) {
        let row = document.createElement("tr");

        // # index
        let index = document.createElement("td");
        index.innerHTML = "#" + (i + 1);
        row.appendChild(index);

        // Name
        let name = document.createElement("td");
        let nameInput = document.createElement("input");
        nameInput.classList.add("form-control");
        nameInput.dataset.row = i;
        nameInput.dataset.field = "name";
        nameInput.value = data[i].name;
        name.appendChild(nameInput);
        row.appendChild(name);

        // Type
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
        select.dataset.row = i;
        type.appendChild(select);
        row.appendChild(type);

        // Address
        let address = document.createElement("td");
        let addressInput = document.createElement("input");
        addressInput.classList.add("form-control");
        addressInput.dataset.row = i;
        addressInput.dataset.field = "address";
        addressInput.value = data[i].address;
        address.appendChild(addressInput);
        row.appendChild(address);

        // Timeout
        let timeout = document.createElement("td");
        let timeoutInput = document.createElement("input");
        timeoutInput.classList.add("form-control");
        timeoutInput.type = "number";
        timeoutInput.dataset.row = i;
        timeoutInput.dataset.field = "timeout";
        timeoutInput.value = data[i].timeout;
        timeout.appendChild(timeoutInput);
        row.appendChild(timeout);

        // Status Codes (HTTP only)
        let status = document.createElement("td");
        let statusInput = document.createElement("input");
        statusInput.classList.add("form-control");
        statusInput.dataset.row = i;
        statusInput.dataset.field = "status_code";
        if (data[i].type === "http") {
            statusInput.value = data[i].status_code || "";
        } else {
            statusInput.value = "-";
            statusInput.disabled = true;
        }
        status.appendChild(statusInput);
        row.appendChild(status);

        // Expression (HTTP only – button)
        let exprCell = document.createElement("td");
        if (data[i].type === "http") {
            let exprBtn = document.createElement("button");
            exprBtn.classList.add("btn", "btn-sm");

            let hasExpr = data[i].expression && data[i].expression.trim() !== "";
            let hasVars = data[i].variables && data[i].variables.length > 0;

            if (hasExpr || hasVars) {
                exprBtn.classList.add("btn-info");
                let varCount = (data[i].variables || []).length;
                exprBtn.innerHTML = '<i class="bi bi-pencil-square me-1"></i>' + varCount + " var" + (varCount !== 1 ? "s" : "");
            } else {
                exprBtn.classList.add("btn-outline-secondary");
                exprBtn.innerHTML = '<i class="bi bi-plus-lg me-1"></i>Add';
            }

            exprBtn.dataset.row = i;
            exprBtn.addEventListener("click", function () {
                openExpressionModal(parseInt(this.dataset.row));
            });
            exprCell.appendChild(exprBtn);
        } else {
            let dash = document.createElement("span");
            dash.classList.add("text-muted");
            dash.textContent = "-";
            exprCell.appendChild(dash);
        }
        row.appendChild(exprCell);

        // Delete
        let deleteCell = document.createElement("td");
        let button = document.createElement("button");
        button.classList.add("btn", "btn-danger");
        button.innerHTML = "<i class='bi bi-trash3-fill'></i>";
        button.onclick = function () {
            modalToDelete = i;
            genModal(resetFavModal);
        };
        deleteCell.appendChild(button);
        row.appendChild(deleteCell);

        table.appendChild(row);
    }

    setSwitchers();
    setUpdaters();
}


// ===== TYPE SWITCHERS =====
function setSwitchers() {
    let selects = table.querySelectorAll("select");
    for (let i = 0; i < selects.length; i++) {
        selects[i].addEventListener("change", function () {
            let rowIdx = parseInt(this.dataset.row);
            data[rowIdx].type = this.value;
            if (this.value === "ping") {
                delete data[rowIdx].status_code;
                delete data[rowIdx].variables;
                delete data[rowIdx].expression;
            } else {
                data[rowIdx].status_code = "";
            }
            renderTable();
        });
    }
}


// ===== INPUT UPDATERS =====
function setUpdaters() {
    let inputs = table.querySelectorAll("input.form-control");
    for (let i = 0; i < inputs.length; i++) {
        inputs[i].addEventListener("change", function () {
            scheduleAutoValidation();
            let rowIdx = parseInt(this.dataset.row);
            let field = this.dataset.field;

            switch (field) {
                case "name":
                    data[rowIdx].name = this.value;
                    break;
                case "address":
                    data[rowIdx].address = this.value;
                    break;
                case "timeout":
                    data[rowIdx].timeout = parseInt(this.value);
                    break;
                case "status_code":
                    if (data[rowIdx].type === "http") {
                        if (this.value === "") {
                            delete data[rowIdx].status_code;
                        } else {
                            data[rowIdx].status_code = this.value;
                        }
                    }
                    break;
            }
        });
    }
}


// ===== EXPRESSION MODAL =====
function openExpressionModal(rowIndex) {
    currentExpressionRow = rowIndex;
    expressionSaved = false;

    if (!data[rowIndex].variables) {
        data[rowIndex].variables = [];
    }

    // Store snapshot for cancel/restore
    expressionSnapshot = {
        variables: JSON.parse(JSON.stringify(data[rowIndex].variables)),
        expression: data[rowIndex].expression || ""
    };

    renderVariablesInModal();
    document.querySelector("#expressionInput").value = data[rowIndex].expression || "";
    validateExpressionLive();

    var modal = new bootstrap.Modal(document.querySelector("#expressionModal"));
    modal.show();
}

function renderVariablesInModal() {
    let vars = data[currentExpressionRow].variables || [];
    let container = document.querySelector("#variablesContainer");
    container.innerHTML = "";

    if (vars.length === 0) {
        container.innerHTML = '<p class="text-muted mb-0">No variables defined. Click "Add Variable" to create one.</p>';
        return;
    }

    for (let i = 0; i < vars.length; i++) {
        let varRow = document.createElement("div");
        varRow.classList.add("d-flex", "align-items-center", "gap-2", "mb-2", "variable-row");

        // Name input
        let nameInput = document.createElement("input");
        nameInput.classList.add("form-control", "form-control-sm");
        nameInput.placeholder = "var_name";
        nameInput.value = vars[i].name || "";
        nameInput.style.maxWidth = "130px";
        nameInput.dataset.varIndex = i;
        nameInput.addEventListener("input", function () {
            updateVariableInData(parseInt(this.dataset.varIndex), "name", this.value);
        });
        varRow.appendChild(nameInput);

        // Mode select
        let modeSelect = document.createElement("select");
        modeSelect.classList.add("form-select", "form-select-sm");
        modeSelect.style.maxWidth = "95px";
        modeSelect.dataset.varIndex = i;

        let boolOpt = document.createElement("option");
        boolOpt.value = "bool";
        boolOpt.text = "bool";
        let jsonOpt = document.createElement("option");
        jsonOpt.value = "json";
        jsonOpt.text = "json";
        modeSelect.appendChild(boolOpt);
        modeSelect.appendChild(jsonOpt);
        modeSelect.value = vars[i].mode || "bool";

        modeSelect.addEventListener("change", function () {
            let idx = parseInt(this.dataset.varIndex);
            updateVariableInData(idx, "mode", this.value);
            renderVariablesInModal();
            validateExpressionLive();
        });
        varRow.appendChild(modeSelect);

        if (vars[i].mode === "json") {
            // JSON query input
            let queryInput = document.createElement("input");
            queryInput.classList.add("form-control", "form-control-sm");
            queryInput.placeholder = "$.data.temperature";
            queryInput.value = vars[i].query || "";
            queryInput.dataset.varIndex = i;
            queryInput.addEventListener("input", function () {
                updateVariableInData(parseInt(this.dataset.varIndex), "query", this.value);
            });
            varRow.appendChild(queryInput);

            // As select (bool / number)
            let asSelect = document.createElement("select");
            asSelect.classList.add("form-select", "form-select-sm");
            asSelect.style.maxWidth = "115px";
            asSelect.dataset.varIndex = i;

            let asBoolOpt = document.createElement("option");
            asBoolOpt.value = "bool";
            asBoolOpt.text = "as bool";
            let asNumOpt = document.createElement("option");
            asNumOpt.value = "number";
            asNumOpt.text = "as number";
            asSelect.appendChild(asBoolOpt);
            asSelect.appendChild(asNumOpt);
            asSelect.value = vars[i].as || "bool";

            asSelect.addEventListener("change", function () {
                updateVariableInData(parseInt(this.dataset.varIndex), "as", this.value);
            });
            varRow.appendChild(asSelect);
        } else {
            // Bool value input
            let valueInput = document.createElement("input");
            valueInput.classList.add("form-control", "form-control-sm");
            valueInput.placeholder = "Search value";
            valueInput.value = vars[i].value || "";
            valueInput.dataset.varIndex = i;
            valueInput.addEventListener("input", function () {
                updateVariableInData(parseInt(this.dataset.varIndex), "value", this.value);
            });
            varRow.appendChild(valueInput);
        }

        // Delete variable button
        let deleteBtn = document.createElement("button");
        deleteBtn.classList.add("btn", "btn-outline-danger", "btn-sm");
        deleteBtn.innerHTML = '<i class="bi bi-trash3"></i>';
        deleteBtn.dataset.varIndex = i;
        deleteBtn.addEventListener("click", function () {
            removeVariableFromModal(parseInt(this.dataset.varIndex));
        });
        varRow.appendChild(deleteBtn);

        container.appendChild(varRow);
    }
}

function addVariableToModal() {
    if (currentExpressionRow === null) return;
    if (!data[currentExpressionRow].variables) {
        data[currentExpressionRow].variables = [];
    }
    data[currentExpressionRow].variables.push({
        name: "",
        mode: "bool",
        value: ""
    });
    renderVariablesInModal();
    validateExpressionLive();
}

function removeVariableFromModal(index) {
    if (currentExpressionRow === null) return;
    data[currentExpressionRow].variables.splice(index, 1);
    renderVariablesInModal();
    validateExpressionLive();
}

function updateVariableInData(varIndex, field, value) {
    if (currentExpressionRow === null) return;
    let v = data[currentExpressionRow].variables[varIndex];

    if (field === "mode") {
        v.mode = value;
        if (value === "bool") {
            delete v.query;
            delete v.as;
            v.value = v.value || "";
        } else {
            delete v.value;
            v.query = v.query || "";
            v.as = v.as || "bool";
        }
    } else {
        v[field] = value;
    }

    validateExpressionLive();
}

function validateExpressionLive() {
    let preview = document.querySelector("#expressionPreview");
    let exprInput = document.querySelector("#expressionInput");
    let expr = exprInput.value.trim();
    let vars = data[currentExpressionRow].variables || [];

    let issues = [];

    // Validate variable definitions
    let varNames = new Set();
    for (let v of vars) {
        if (!v.name || v.name.trim() === "") {
            issues.push("Variable name cannot be empty.");
        } else if (!/^[a-z][a-z0-9_]*$/.test(v.name)) {
            issues.push('Variable "' + v.name + '" \u2013 name must start with a letter and contain only lowercase letters, digits, and underscores.');
        }
        if (v.name && varNames.has(v.name)) {
            issues.push('Duplicate variable name: "' + v.name + '".');
        }
        varNames.add(v.name);

        if (v.mode === "bool" && (!v.value || v.value.trim() === "")) {
            issues.push('Variable "' + v.name + '" (bool) \u2013 search value is empty.');
        }
        if (v.mode === "json" && (!v.query || v.query.trim() === "")) {
            issues.push('Variable "' + v.name + '" (json) \u2013 JSON query is empty.');
        }
    }

    // Validate expression
    if (expr !== "") {
        let varTypes = buildVarTypeMap(vars);
        let exprResult = validateExpressionSyntax(expr, varNames, varTypes);
        if (exprResult.errors.length > 0) {
            issues.push.apply(issues, exprResult.errors);
        }
    } else if (vars.length > 0) {
        issues.push("Variables defined but no expression set.");
    }

    // Render preview
    if (issues.length === 0 && (expr !== "" || vars.length > 0)) {
        preview.innerHTML = '<span class="text-success"><i class="bi bi-check-circle me-1"></i>Valid</span>';
    } else if (issues.length === 0) {
        preview.innerHTML = '<span class="text-muted">No expression configured.</span>';
    } else {
        preview.innerHTML = issues
            .map(function (msg) {
                return '<span class="text-warning"><i class="bi bi-exclamation-triangle me-1"></i>' + msg + "</span>";
            })
            .join("<br>");
    }
}


// ===== EXPRESSION VALIDATION =====

// Build a map of variable name -> effective type ("bool" or "number")
function buildVarTypeMap(vars) {
    var types = {};
    for (var i = 0; i < vars.length; i++) {
        var v = vars[i];
        if (v.mode === "bool") {
            types[v.name] = "bool";
        } else if (v.mode === "json") {
            types[v.name] = (v.as === "number") ? "number" : "bool";
        }
    }
    return types;
}

function validateExpressionSyntax(expr, definedVars, varTypes) {
    let errors = [];
    varTypes = varTypes || {};

    let result = tokenizeExpression(expr);
    if (result.error) {
        errors.push(result.error);
        return { errors: errors };
    }

    let tokens = result.tokens;
    let comparisonOps = new Set([">", "<", ">=", "<=", "==", "!="]);

    // Check balanced parentheses
    let depth = 0;
    for (let i = 0; i < tokens.length; i++) {
        if (tokens[i] === "(") depth++;
        if (tokens[i] === ")") depth--;
        if (depth < 0) {
            errors.push("Unbalanced parentheses: extra closing ')'.");
            break;
        }
    }
    if (depth > 0) {
        errors.push("Unbalanced parentheses: missing closing ')'.");
    }

    // Check all identifiers reference defined variables
    // and that bool variables are not used with comparison operators
    for (let i = 0; i < tokens.length; i++) {
        let t = tokens[i];
        if (/^[a-z][a-z0-9_]*$/.test(t)) {
            if (!definedVars.has(t)) {
                errors.push('Undefined variable: "' + t + '".');
            } else if (varTypes[t] === "bool") {
                // Bool variables must not be used with comparison operators
                let prevToken = (i > 0) ? tokens[i - 1] : null;
                let nextToken = (i + 1 < tokens.length) ? tokens[i + 1] : null;
                if ((prevToken && comparisonOps.has(prevToken)) || (nextToken && comparisonOps.has(nextToken))) {
                    errors.push('Variable "' + t + '" is bool and cannot be used with comparison operators.');
                }
            } else if (varTypes[t] === "number") {
                // Number variables must be used with a comparison operator
                let prevToken = (i > 0) ? tokens[i - 1] : null;
                let nextToken = (i + 1 < tokens.length) ? tokens[i + 1] : null;
                if (!(prevToken && comparisonOps.has(prevToken)) && !(nextToken && comparisonOps.has(nextToken))) {
                    errors.push('Variable "' + t + '" is a number and must be used with a comparison operator (e.g. ' + t + ' > 0).');
                }
            }
        }
    }

    return { errors: errors };
}

function tokenizeExpression(expr) {
    let tokens = [];
    let i = 0;
    let s = expr.trim();

    while (i < s.length) {
        // Whitespace
        if (/\s/.test(s[i])) {
            i++;
            continue;
        }

        // Parentheses
        if (s[i] === "(" || s[i] === ")") {
            tokens.push(s[i]);
            i++;
            continue;
        }

        // Two-character operators
        if (i + 1 < s.length) {
            let two = s[i] + s[i + 1];
            if (two === "&&" || two === "||" || two === ">=" || two === "<=" || two === "==" || two === "!=") {
                tokens.push(two);
                i += 2;
                continue;
            }
        }

        // Single-character operators
        if (s[i] === "!" || s[i] === ">" || s[i] === "<") {
            tokens.push(s[i]);
            i++;
            continue;
        }

        // Numeric literals (including negative numbers and decimals)
        if (/[0-9]/.test(s[i]) || (s[i] === "-" && i + 1 < s.length && /[0-9]/.test(s[i + 1]))) {
            let num = "";
            if (s[i] === "-") {
                num += "-";
                i++;
            }
            while (i < s.length && /[0-9.]/.test(s[i])) {
                num += s[i];
                i++;
            }
            tokens.push(num);
            continue;
        }

        // Identifiers (variable names)
        if (/[a-z_]/.test(s[i])) {
            let id = "";
            while (i < s.length && /[a-z0-9_]/.test(s[i])) {
                id += s[i];
                i++;
            }
            tokens.push(id);
            continue;
        }

        return { error: "Unexpected character: '" + s[i] + "' at position " + (i + 1) + "." };
    }

    return { tokens: tokens };
}


// ===== VERIFICATION =====
function verifyCheck() {
    let check_ok = true;
    errorArea.classList.add("d-none");
    errorList.innerHTML = "";

    for (let i = 0; i < data.length; i++) {
        // Name validation
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
                if (j === i) continue;
                if (data[j].name === data[i].name) {
                    let newErrorMessage = "<li>Scan name must be unique (scan name: " + data[i].name + ").</li>";
                    if (!errorList.innerHTML.includes(newErrorMessage)) {
                        errorList.innerHTML += newErrorMessage;
                    }
                    check_ok = false;
                    break;
                }
            }
        } while (false);

        // Address validation
        do {
            if (data[i].address === "") {
                errorList.innerHTML += "<li>Address cannot be empty (scan name: " + data[i].name + ").</li>";
                check_ok = false;
                break;
            }
            if (data[i].address.length > 256) {
                errorList.innerHTML += "<li>Address must be less than 256 characters (scan name: " + data[i].name + ").</li>";
                check_ok = false;
                break;
            }
        } while (false);

        // Timeout validation
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

        // Status codes validation (HTTP only)
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

        // Variables & expression validation (HTTP only)
        if (data[i].type === "http") {
            let vars = data[i].variables || [];
            let expr = (data[i].expression || "").trim();
            let varNames = new Set();

            for (let v = 0; v < vars.length; v++) {
                if (!vars[v].name || vars[v].name.trim() === "") {
                    errorList.innerHTML += "<li>Variable name cannot be empty (scan name: " + data[i].name + ").</li>";
                    check_ok = false;
                } else if (!/^[a-z][a-z0-9_]*$/.test(vars[v].name)) {
                    errorList.innerHTML += "<li>Variable name \"" + vars[v].name + "\" is invalid (scan name: " + data[i].name + ").</li>";
                    check_ok = false;
                }
                if (vars[v].name && varNames.has(vars[v].name)) {
                    errorList.innerHTML += "<li>Duplicate variable name: \"" + vars[v].name + "\" (scan name: " + data[i].name + ").</li>";
                    check_ok = false;
                }
                varNames.add(vars[v].name);

                if (vars[v].mode === "bool" && (!vars[v].value || vars[v].value.trim() === "")) {
                    errorList.innerHTML += "<li>Variable \"" + vars[v].name + "\" value is empty (scan name: " + data[i].name + ").</li>";
                    check_ok = false;
                }
                if (vars[v].mode === "json" && (!vars[v].query || vars[v].query.trim() === "")) {
                    errorList.innerHTML += "<li>Variable \"" + vars[v].name + "\" JSON query is empty (scan name: " + data[i].name + ").</li>";
                    check_ok = false;
                }
            }

            if (expr !== "") {
                let varTypes = buildVarTypeMap(vars);
                let result = validateExpressionSyntax(expr, varNames, varTypes);
                for (let e = 0; e < result.errors.length; e++) {
                    errorList.innerHTML += "<li>Expression error: " + result.errors[e] + " (scan name: " + data[i].name + ").</li>";
                    check_ok = false;
                }
            } else if (vars.length > 0) {
                errorList.innerHTML += "<li>Variables defined but no expression set (scan name: " + data[i].name + ").</li>";
                check_ok = false;
            }
        }
    }

    if (check_ok) {
        errorArea.classList.add("d-none");
    } else {
        errorArea.classList.remove("d-none");
    }

    verifyStatusUpdate(check_ok);
    return check_ok;
}

function verifyStatusUpdate(state) {
    if (state) {
        validationStatus.innerHTML = '<span class="text-success"><i class="bi bi-check-circle me-1"></i>Valid</span>';
    } else {
        validationStatus.innerHTML = '<span class="text-warning"><i class="bi bi-exclamation-triangle me-1"></i>Issues found</span>';
    }
}

function scheduleAutoValidation() {
    if (autoValidateTimer) {
        clearTimeout(autoValidateTimer);
    }
    validationStatus.innerHTML = '';
    autoValidateTimer = setTimeout(function () {
        autoValidateTimer = null;
        verifyCheck();
    }, 800);
}


// ===== CONFIG I/O =====

// Split a config line into parts, treating quoted values as single tokens
function splitConfigLine(line) {
    let parts = [];
    let i = 0;

    while (i < line.length) {
        if (/\s/.test(line[i])) {
            i++;
            continue;
        }

        let part = "";
        while (i < line.length && !/\s/.test(line[i])) {
            if (line[i] === '"') {
                // Enter quoted value
                part += '"';
                i++;
                while (i < line.length && line[i] !== '"') {
                    if (line[i] === "\\" && i + 1 < line.length && line[i + 1] === '"') {
                        part += '\\"';
                        i += 2;
                    } else {
                        part += line[i];
                        i++;
                    }
                }
                if (i < line.length) {
                    part += '"';
                    i++;
                }
            } else {
                part += line[i];
                i++;
            }
        }

        if (part !== "") {
            parts.push(part);
        }
    }

    return parts;
}

function parseConfig(content) {
    let lines = content.split("\n").filter(function (l) {
        return l.trim() !== "";
    });
    let result = [];

    for (let li = 0; li < lines.length; li++) {
        let parts = splitConfigLine(lines[li]);

        let item = {
            name: parts[0],
            type: parts[1],
            address: parts[2],
            timeout: 10
        };

        let variables = [];

        for (let p = 3; p < parts.length; p++) {
            let part = parts[p];

            let eqIdx = part.indexOf("=");
            if (eqIdx === -1) continue;

            let key = part.substring(0, eqIdx);
            let value = part.substring(eqIdx + 1);

            // Remove surrounding quotes
            if (value.length >= 2 && value[0] === '"' && value[value.length - 1] === '"') {
                value = value.substring(1, value.length - 1);
            }

            if (key === "timeout") {
                item.timeout = parseInt(value, 10);
            } else if (key === "status_code") {
                item.status_code = value;
            } else if (key === "keyword") {
                // Legacy support: convert keyword to a bool variable + expression
                variables.push({
                    name: "keyword",
                    mode: "bool",
                    value: value
                });
                if (!item.expression) {
                    item.expression = "keyword";
                }
            } else if (key === "expression") {
                item.expression = value;
            } else if (key.startsWith("var_bool:")) {
                let varName = key.substring("var_bool:".length);
                variables.push({
                    name: varName,
                    mode: "bool",
                    value: value
                });
            } else if (key.startsWith("var_json_bool:")) {
                let varName = key.substring("var_json_bool:".length);
                variables.push({
                    name: varName,
                    mode: "json",
                    query: value,
                    as: "bool"
                });
            } else if (key.startsWith("var_json_number:")) {
                let varName = key.substring("var_json_number:".length);
                variables.push({
                    name: varName,
                    mode: "json",
                    query: value,
                    as: "number"
                });
            }
        }

        if (variables.length > 0) {
            item.variables = variables;
        }

        result.push(item);
    }

    return result;
}

function convertToConfig(data) {
    return data
        .map(function (item) {
            let line = item.name + " " + item.type + " " + item.address + " timeout=" + item.timeout;

            if (item.type === "http") {
                if (item.status_code) {
                    line += ' status_code="' + item.status_code + '"';
                }

                if (item.variables) {
                    for (let v = 0; v < item.variables.length; v++) {
                        let variable = item.variables[v];
                        if (variable.mode === "bool") {
                            line += ' var_bool:' + variable.name + '="' + variable.value + '"';
                        } else if (variable.mode === "json") {
                            if (variable.as === "number") {
                                line += ' var_json_number:' + variable.name + '="' + variable.query + '"';
                            } else {
                                line += ' var_json_bool:' + variable.name + '="' + variable.query + '"';
                            }
                        }
                    }
                }

                if (item.expression) {
                    line += ' expression="' + item.expression + '"';
                }
            }

            return line;
        })
        .join("\n");
}

function downloadConfig() {
    var configContent = convertToConfig(data);
    var blob = new Blob([configContent], { type: "text/plain" });
    var link = document.createElement("a");
    link.href = URL.createObjectURL(blob);
    link.download = "config.conf";
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
}


// ===== INITIALIZATION =====
renderTable();

window.onbeforeunload = function () {
    return "Data may be lost if you leave the page, are you sure?";
};

setTimeout(function () {
    var popoverTriggerList = [].slice.call(
        document.querySelectorAll('[data-bs-toggle="popover"]')
    );
    popoverTriggerList.map(function (popoverTriggerEl) {
        return new bootstrap.Popover(popoverTriggerEl);
    });
}, 200);

// Disable focus warning on modal close
document.addEventListener("DOMContentLoaded", function () {
    document.addEventListener("hide.bs.modal", function (event) {
        if (document.activeElement) {
            document.activeElement.blur();
        }
    });
});


// ===== EVENT LISTENERS =====
document.querySelector("#uploadConfig").addEventListener("change", function (event) {
    var file = event.target.files[0];
    if (file) {
        var reader = new FileReader();
        reader.onload = function (e) {
            var content = e.target.result;
            data = parseConfig(content);
            renderTable();
        };
        reader.readAsText(file);
    }
});

document.querySelector("#downloadConfig").addEventListener("click", function () {
    if (autoValidateTimer) {
        clearTimeout(autoValidateTimer);
        autoValidateTimer = null;
    }
    if (verifyCheck()) {
        downloadConfig();
    }
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

// Expression modal listeners
document.querySelector("#addVariableBtn").addEventListener("click", function () {
    addVariableToModal();
});

document.querySelector("#expressionInput").addEventListener("input", function () {
    if (currentExpressionRow !== null) {
        data[currentExpressionRow].expression = this.value;
        validateExpressionLive();
    }
});

document.querySelector("#expressionModal").addEventListener("hidden.bs.modal", function () {
    if (currentExpressionRow !== null) {
        var item = data[currentExpressionRow];

        if (!expressionSaved && expressionSnapshot) {
            // Cancel: restore snapshot
            if (expressionSnapshot.variables.length > 0) {
                item.variables = expressionSnapshot.variables;
            } else {
                delete item.variables;
            }
            if (expressionSnapshot.expression && expressionSnapshot.expression.trim() !== "") {
                item.expression = expressionSnapshot.expression;
            } else {
                delete item.expression;
            }
        } else {
            // Save: clean up empty data
            if (item.variables && item.variables.length === 0) {
                delete item.variables;
            }
            if (item.expression && item.expression.trim() === "") {
                delete item.expression;
            }
        }

        expressionSnapshot = null;
        expressionSaved = false;
        currentExpressionRow = null;
        renderTable();
    }
});

// Save & Close button in expression modal
document.querySelector("#saveExpressionBtn").addEventListener("click", function () {
    expressionSaved = true;
    var modalEl = document.querySelector("#expressionModal");
    var modal = bootstrap.Modal.getInstance(modalEl);
    if (modal) {
        modal.hide();
    }
});
