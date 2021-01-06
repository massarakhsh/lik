let likRefRegister = {};
let likRefIndex = 0;
let likRefTick = null;

function lik_ref_initialize() {
    likRefRegister = {};
    likRefIndex = 0;
    lik_ref_request();
    if (!likRefTick) {
        likRefTick = setInterval(lik_ref_request, 10000);
    }
}

function lik_ref_request() {
    lik_ref_call("GET", null, 'http://localhost:8090/api/marshal/' + likRefIndex);
}

function lik_ref_call(type, data, url) {
    jQuery.ajax({
        type: type,
        url: url,
        data: data,
        dataType: "json",
        crossDomain: true,
        timeout: 10000,
        success: function(code) {
            lik_ref_update(code);
        }
    });
}

function lik_ref_update(lika) {
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

function lik_ref_get(path) {
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

