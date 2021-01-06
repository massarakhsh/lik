let likRefRegister = {};
let likRefIndex = 0;
let likRefTick = null;

function lik_ref_initialize() {
    likRefRegister = {};
    likRefIndex = 0;
    requestUpdate();
    if (!likRefTick) {
        likRefTick = setInterval(requestUpdate, 10000);
    }
}

function requestUpdate() {
    //get_data_proc('/api/marshal/' + dataIndex, makeUpdate, null);
    get_data_proc('http://localhost:8090/api/marshal/' + likRefIndex, makeUpdate, null);
}

function makeUpdate(parm, lika) {
    if (lika) {
        for (var key in lika) {
            if (key == 'index') {
                likRefIndex = lika[key];
            } else {
                let tin = lika[key];
                let tis = likRefRegister[key];
                if (!tis) {
                    tis = {};
                    likRefRegister[key] = tis;
                }
                if (tin) {
                    for (var id in tin) {
                        let ein = tin[id];
                        let eis = tis[id];
                        if (eis && (!ein || !ein["SysNum"])) {
                            delete tis[id];
                        } else if (ein && ein["SysNum"] == id) {
                            tis[id] = ein;
                        }
                    }
                }
            }
        }
    }
}

function getRegister(path) {
    let names = path.split('/');
    let loc = likRefRegister;
    for (let p = 0; p < names.length; p++) {
        let name = names[p];
        if (name) {
            let dat = loc[name];
            if (typeof(dat) === 'undefined' || dat === null) return null;
            loc = dat;
        }
    }
    return loc;
}

