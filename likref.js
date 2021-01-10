const likRefId = 'SysNum';
let likRefRegister = {};
let likRefTick = null;

class LikRefTable {
    Table = "";
    Index = 0;
    Elms = {};
    callBacks = [];

    constructor(table) {
        this.Table = table;
    }

    get(callback) {
        if (this.Index > 0) {
            setTimeout(() => callback(this.Elms), 1);
        } else {
            this.callBacks.push(callback);
        }
    }

    set(elms, index) {
        this.Index = index;
        if (elms) {
            for (var id in elms) {
                let elm = elms[id];
                let old = this.Elms[id];
                if (old && (!elm || !elm[likRefId])) {
                    delete this.elms[id];
                } else if (elm && elm[likRefId] == id) {
                    this.Elms[id] = elm;
                }
            }
        }
        while (this.callBacks.length > 0) {
            let callback = this.callBacks.shift();
            this.get(callback);
        }
    }
}

function _lik_ref_tick() {
    if (!likRefTick) {
        likRefTick = setInterval(() => _lik_ref_request_all(), 10000);
        _lik_ref_request_all();
    }
}

function _lik_ref_request_all() {
    for (let table in likRefRegister) {
        _lik_ref_request_table(table, likRefRegister[table].Index);
    }
}

function _lik_ref_request_table(table, index) {
    let url = 'http://localhost:8090/api/' + table + '/' + index;
    jQuery.ajax({
        type: 'GET',
        url: url,
        dataType: "json",
        crossDomain: true,
        timeout: 5000,
        success: function(code) {
            _lik_ref_update(code);
        }
    });
}

function _lik_ref_update(lika) {
    if (lika) {
        const index = lika['index'];
        for (let key in lika) {
            if (key != 'index') {
                let reg = likRefRegister[key];
                if (!reg) {
                    reg = new LikRefTable(key);
                    likRefRegister[key] = reg;
                }
                reg.set(lika[key], index);
            }
        }
    }
}

function likref_get_all(table, callback) {
    _lik_ref_tick();
    let reg = likRefRegister[table];
    if (!reg) {
        reg = new LikRefTable(table);
        likRefRegister[table] = reg;
        _lik_ref_request_table(table, 0);
    }
    reg.get(callback);
}

function likref_get_sys(table, sys, callback) {
    likref_get_all(table, (elms) => {
        callback(elms[sys]);
    });
}

