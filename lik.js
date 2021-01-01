// Это UTF-8

var lik_page = 0;
var lik_trust = 0;
var lik_time = Math.floor(Date.now()/1000);
var likTicker = null;
var tick_total = 0;
var tick_answer = 0;
var tick_step = 250;
var tick_second = 0;
var tick_delta = 0;
var tick_server = 0;
var tick_redraw = 0;
var tick_menuclose = 0;
var tick_bin = 0;
var tick_shift_minute = (new Date()).getTimezoneOffset();
var mouse_locate;
var queue_tick = Array();
var menu_stack = Array();
var menu_active = false;
var menu_locked = false;

var screen_seek = 0;
var screen_width = 800;
var screen_height = 600;

var marshal_step = 0;
var marshal_path = null;
var marshal_tick = 0;

var para_run_once = null;
var pool_step = Array();
var pool_first = Array();
var para_external = null;
var para_to_part = null;
var para_key_press = null;

var tick_delay = 0;
var part_delay = null;

function lik_start() {
	var body = jQuery('body');
    body.mousemove(on_mouse_move);
	body.mousedown(on_mouse_down);
	body.keypress(on_key_press);
    tick_delta = Math.floor(Date.now()/1000)-lik_time;
	likTicker = window.setInterval(run_step,tick_step);
}

function lik_stop() {
	clearInterval(likTicker);
	likTicker = null;
}

/////////////////// Links

function lik_build_url(part) {
	var url = part;
	var parm = "";
	var psp = part.indexOf("?");
	if (psp>=0) {
		parm = part.substr(psp + 1);
		url = part.substr(0, psp);
	}
    var iss = false;
	var isp = false;
    var prms = parm.split("&");
    for (var p=0; p<prms.length; p++) {
    	var prm = prms[p];
    	if (prm.length > 0 && !prm.startsWith("_sw=") && !prm.startsWith("_sh=")) {
            url += (iss) ? "&" : "?";
    		url += prm;
    		iss = true;
            if (prm.startsWith("_sp=")) isp = true;
		}
	}
    if (!isp) {
        url += (iss) ? "&" : "?";
        url += "_sp=" + lik_page;
    }
	if (screen_width!=window.innerWidth) {
		url += '&_sw='+window.innerWidth;
    }
	if (screen_height!=window.innerHeight) {
		url += '&_sh='+window.innerHeight;
    }
	return url;
}

function lik_load_rid_part(rid, part) {
	lik_load_rid_url(rid, lik_build_url(part));
}
function lik_load_rid_url(rid, url) {
	jQuery('[rid='+rid+']').hide(0).load(url, function() { $(this).show(0); });
}

function lik_get_url(url) {
	return jQuery.ajax({ url: url, async: false }).responseText;
}
function lik_get_part(part) {
	return lik_get_url(lik_build_url(part));
}

function lik_go_url(url) {
	history.pushState({lik_page: lik_page}, "", url);
	document.location.href = url;
}

function lik_go_part(part) {
	lik_go_url(lik_build_url(part));
}

function lik_window_part(part) {
	var url = lik_build_url(part);
	window.open(url, '_blank');
}

function lik_reload() {
	document.location.reload();
}

function lik_set_trust(part) {
	if (!lik_trust) {
		history.pushState('', '', lik_build_url(part));
		lik_trust = true;
	}
}

function lik_set_history(part) {
	lik_rule_history(part, false);
}
function lik_push_history(part) {
	lik_rule_history(part, true);
}
function lik_rule_history(part, psh) {
	var title = null;
	if (match = /(.*?)##(.*)/.exec(part)) {
		part = match[1];
		title = match[2];
	}
	if (!part) {
	} else if (psh) {
		history.pushState({lik_page: lik_page}, "", lik_build_url(part));
	} else {
		history.replaceState({lik_page: lik_page}, "", lik_build_url(part));
	}
	if (title) {
		document.title = title;
	}
}

/////////////////// Key+Mouse

function on_key_press(e) {
	var key = e.which;
	if (para_key_press !== null) para_key_press(e);
}
function on_mouse_move(e) {
	mouse_locate = { X: e.pageX, Y: e.pageY };
	//$('#Debug2').text("X:"+e.pageX+", Y:"+e.pageY);
}
function on_mouse_down(e) {
	if (menu_active && !menu_locked && tick_menuclose == 0) {
		tick_menuclose = tick_total + 1000;
	}
}

/////////////////// Tick

function run_step() {
    var tick = Date.now();
    tick_total = tick - tick_shift_minute*60*1000;
    tick_second = Math.floor(tick_total/1000);
    tick_server = Math.floor(tick/1000)-tick_delta;
    if (!tick_answer) tick_answer = tick_total;
	jQuery('#DebTick').text(tick_total);
    if (pool_first != null) {
        for (var p=0; p<pool_first.length; p++) {
            (pool_first[p])();
        }
        pool_first = null;
    }
	if (tick_total >= tick_redraw) {
		tick_redraw = tick_total + 250;
		lik_redraw();
	}
	if (marshal_path && tick_total >= marshal_tick) {
		lik_step_marshal();
	}
	if (menu_active && !menu_locked && tick_menuclose > 0 && tick_total >= tick_menuclose) {
		menu_close_all();
	}
	while (queue_tick.length>0 && tick_total>=queue_tick[0].Tick) {
		var plan = queue_tick.shift();
		var deal = ('Deal' in plan) ? plan.Deal : null;
		if ('deal'+deal in window) {
            var data = ('Data' in plan) ? plan.Data : null;
            window['deal'+deal](data);
        }
	}
	if (pool_step != null) {
		for (var p=0; p<pool_step.length; p++) {
			(pool_step[p])();
        }
    }
	if (tick_delay && part_delay && tick_total >= tick_delay) {
		tick_delay = 0;
		load_to_part(part_delay);
	}
}

function remove_queue(deal) {
	rule_queue_data(-1,deal,null,0);
}
function present_queue(deal) {
	for (var pos=queue_tick.length-1; pos>=0; pos--) {
		var que = queue_tick[pos];
		if (deal==que.Deal) return 1;
	}
	return 0;
}
function before_queue(delay,deal,data) {
	rule_queue_data(delay,deal,data,1);
}
function after_queue(delay,deal,data) {
	rule_queue_data(delay,deal,data,0);
}
function rule_queue_data(delay,deal,data,before) {
	var tick = (delay>=0) ? tick_total+delay : -1;
	var posi = (delay>=0) ? queue_tick.length : -1;
	var poso = -1;
	for (var pos=queue_tick.length-1; pos>=0; pos--) {
		var que = queue_tick[pos];
		if (tick>=0 && tick<que.Tick) {
			posi = pos;
		}
		if (deal==que.Deal) {
			if (poso>=0) {
				if (posi>poso) posi--;
				queue_tick.splice(poso,1);
			}
			poso = pos;
		}
	}
	if (posi>=0 && poso>=0 && (before && posi>=poso || !before && posi<=poso+1)) {
		posi = poso;
		var que = queue_tick[posi];
		if (before && tick<que.Tick) que.Tick = tick;
		else if (!before && tick>que.Tick) que.Tick = tick;
		que.Data = data;
	}
	else {
		if (poso>=0) {
			if (posi>poso) posi--;
			queue_tick.splice(poso,1);
		}
		if (posi>=0) {
			var plan = { Tick: tick, Deal: deal, Data: data };
			queue_tick.splice(posi,0,plan);
		}
	}
}

function lik_set_marshal(step, path) {
	marshal_step = step;
	marshal_path = path;
	marshal_tick = tick_total + marshal_step;
}

/////////////////////////////////

function lik_step_marshal() {
	marshal_tick = tick_total + 10000;
	if (marshal_path) {
		let path = marshal_path;
		let rdr = jQuery('[marshalAnker]');
		let pairs = '';
		if (rdr.size() > 0) {
			rdr.each(function (idx, item) {
				let elm = jQuery(item);
				let anker = elm.attr('marshalAnker');
				let index = elm.attr('marshalIndex');
				if (pairs) pairs += "/";
				pairs += anker + ":" + index;
			});
		}
		if (pairs) {
			path += (/\?/.exec(path)) ? '&' : '?';
			path += 'marshal=' + pairs;
		}
		json_call_part("GET", null, path, null, null);
	}
}

function lik_next_marshal() {
	marshal_tick = tick_total + marshal_step;
}

function lik_force_marshal() {
	marshal_tick = tick_total;
}

////////////////////////////////////

function lik_redraw() {
	jQuery('[topart]').each(function(idx,item) {
		var elm = $(item);
		if (!elm.attr('totick') || parseInt(elm.attr('totick')) <= tick_total) {
			lik_json_prepare(elm);
			var topart = elm.attr('topart');
			json_call_part("GET", null, topart, null, elm);
		}
	});
}

function lik_force_redraw() {
	jQuery('[topart]').each(function(idx,item) {
		if (!elm.hasAttr('totick') || parseInt(elm.attr('totick')) > tick_total) {
			elm.attr('totick', tick_total);
		}
	});
	tick_redraw = tick_total;
}

function hash_change() {
	if (tick_total>tick_hash) {
		if (location.hash.length>1) {
			tick_hash = tick_total+400;
			//link_id_part('Info', location.hash.substr(1));
		}
	}
}

function set_hash(path) {
	hash_used = 1;
	tick_hash = tick_total+400;
	location.hash = '#'+path;
}

//////////////////////// Load

function delay_to_part(delay,part) {
	tick_delay = tick_total+delay;
	part_delay = part;
}

function load_to_part(part) {
    menu_close_all();
	if (para_to_part != null) para_to_part(part);
    else lik_go_part(part);
}

function get_data_part(part) {
	json_call_part("GET", null, part, null, null);
}

function get_data_proc(part,proc,parm) {
	json_call_part("GET", null, part, proc, parm);
}

function post_data_part(part,data) {
	json_call_part("POST", data, part, null, null);
}

function post_data_proc(part,data,proc,parm) {
	json_call_part("POST", data, part, proc, parm);
}

function json_call_part(type, data, part, proc, parm) {
	json_call_url(type, data, lik_build_url(part), proc, parm, part == marshal_path);
}

function json_call_url(type, data, url, proc, parm, ismarshal) {
	jQuery.ajax({
		type: type,
		url: url,
		data: data,
		dataType: "json",
		crossDomain: true,
		timeout: 10000,
		success: function(code) {
			if (ismarshal) lik_next_marshal();
			json_answer_part(code, proc, parm);
		},
		error: function(xhr, status) {
			if (ismarshal) lik_next_marshal();
			json_answer_part(undefined, proc, parm);
		}
	});
}

function json_answer_part(lika, proc, parm) {
	if (lika !== undefined) {
		tick_answer = tick_total;
		marshal_tick = tick_total;
		for (var key in lika) {
			var val = lika[key];
			if (key == "_sw" && val) {
				screen_width = val;
			} else if (key == "_sh" && val) {
				screen_height = val;
			} else if (key == "_sp" && val) {
				lik_page = val;
			} else if (key == "_page" && val) {
				lik_page = val;
				let path = document.location.href;
				if (match = /\/\/[^\/]*(.*)/.exec(path)) {
					path = match[1];
				}
				if (match = /(.*?)\?/.exec(path)) {
					path = match[1];
				}
				lik_go_part(path);
			} else if (key == "_topart" && val) {
				lik_go_part(val);
			} else if (key == "_url" && val) {
				lik_set_history(val);
			} else if (key == "_history" && val) {
				lik_push_history(val);
			} else if (key == "_title" && val) {
				document.title = val;
			} else if (key == "_self" && parm) {
				parm.replaceWith(val);
				lik_json_prepare(parm);
			} else if (key == "_content" && parm) {
				parm.html(val);
				lik_json_prepare(parm);
			} else if (match = /_function_(.+)/.exec(key)) {
				if (match[1] in window) {
					window[match[1]](val);
				}
			} else {
				var elm = jQuery("#" + key);
				if (elm.size() > 0) {
					elm.replaceWith(val);
					lik_json_prepare(elm);
				}
			}
		}
	}
	if (proc) proc(parm, lika);
	return false;
}

function lik_json_prepare(elm) {
	if (elm && elm.size()>0) {
		var freq = (elm.attr('tofreq')) ? parseInt(elm.attr('tofreq')) : 0;
		if (!freq) freq = 1000;
		elm.attr('totick', tick_total + freq);
	}
}

function load_rid(rid,lika) {
	var elm = jQuery('[rid='+rid+']');
	if (elm.size()>0) {
		var topart = ('topart' in lika) ? lika.topart : '';
		if (topart) elm.attr('topart',topart);
		else elm.removeAttr('topart');
		var code = lika.data;
		elm.hide(0);
		elm.empty();
		elm.prepend(code);
		before_queue(100,'Part');
		elm.show(0);
	}
}

///////////////////// Menu

function menu_level() {
	return menu_stack.length;
}

function menu_close_all() {
	menu_active = false;
	menu_locked = false;
	tick_menuclose = 0;
	menu_close_level(0);
}

function menu_close_level(level) {
	while (menu_stack.length>level) {
		menu_stack[menu_stack.length-1].remove();
		menu_stack.splice(menu_stack.length-1,1);
	}
}

function menu_open(elm,data) {
	menu_active = true;
	tick_menuclose = 0;
	elm.append(data);
	menu_stack.push(elm.children(":last-child"));
}

function menu_lock() {
	menu_locked = true;
	tick_menuclose = 0;
}

function menu_free() {
	menu_locked = false;
	tick_menuclose = 0;
}

///////////////////// Coding

function string_to_XS(data) {
	var code = "";
	var str = new String(data);
	var len = str.length;
	for (var i=0; i<len; i++) {
		var ch = str.charCodeAt(i);
		for (var c=0; c<4; c++) {
			var cd = (ch >> 12)&0xF;
			if (cd < 10) code += String.fromCharCode(0x30 + cd);
			else code += String.fromCharCode(0x41 + cd - 10);
			ch <<= 4;
		}
	}
	return code;
}

function string_from_XS(data) {
	var code = "";
	var str = new String(data);
	var len = str.length;
	for (var i=0; i+4<=len; i+=4) {
		var cd = 0;
		for (var c=0; c<4; c++) {
			cd <<= 4;
			var ch = str.charCodeAt(i+c);
			if (ch >= 0x30 && ch <= 0x39) cd |= (ch-0x30);
			else if (ch >= 0x41 && ch <= 0x46) cd |= (ch-0x41+10);
			else if (ch >= 0x61 && ch <= 0x66) cd |= (ch-0x61+10);
			else break;
		}
		code += String.fromCharCode(cd);
	}
	return code;
}

