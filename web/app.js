'use strict';

const $ = (id) => document.getElementById(id);
const esc = (s) => String(s ?? '').replace(/[&<>"]/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;' }[c]));

let conn = { connected: false };
let view = 'naming';
let running = false;
let wantDownload = false;
let defaults = {};

async function post(path, body) {
  const r = await fetch(path, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body || {}) });
  const j = await r.json().catch(() => ({}));
  if (!r.ok) throw new Error(j.error || `HTTP ${r.status}`);
  return j;
}

let toastTimer = null;
function toast(msg) {
  const t = $('toast');
  t.textContent = msg;
  t.hidden = false;
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => { t.hidden = true; }, 6000);
}

function log(msg) {
  const p = $('pLog');
  const stick = p.scrollTop + p.clientHeight >= p.scrollHeight - 8;
  const line = document.createElement('div');
  const t = new Date().toTimeString().slice(0, 8);
  line.innerHTML = `<span class="t">${t}</span>  ${esc(msg)}`;
  p.appendChild(line);
  while (p.childNodes.length > 5000) p.removeChild(p.firstChild);
  if (stick) p.scrollTop = p.scrollHeight;
}

/* ---------- tables ---------- */

// Both lists of APs behave the same way: filter, select, delete, and a result
// that lands on the row as the run reaches it.
function makeTable(key, cells) {
  const K = key[0].toUpperCase() + key.slice(1);
  const body = $('rows' + K);
  const t = { rows: [], file: '', sel: new Set(), filter: '', chip: null, anchor: null, trs: [] };

  const status = (r) => (!r.res ? 'pending' : r.res.error ? 'failed' : 'ok');
  const rowHtml = (r) => cells(r).map(([text, cls]) => `<td class="${cls || ''}">${esc(text)}</td>`).join('');

  t.visible = () => {
    const f = t.filter.toLowerCase();
    const out = [];
    t.rows.forEach((r, i) => {
      if (t.chip && status(r) !== t.chip) return;
      if (f && !cells(r).some(([text]) => String(text).toLowerCase().includes(f))) return;
      out.push(i);
    });
    return out;
  };

  t.render = () => {
    const vis = t.visible();
    t.trs = [];
    const frag = document.createDocumentFragment();
    for (const i of vis) {
      const tr = document.createElement('tr');
      tr.dataset.i = i;
      if (t.sel.has(i)) tr.className = 'pick';
      tr.innerHTML = rowHtml(t.rows[i]);
      t.trs[i] = tr;
      frag.appendChild(tr);
    }
    body.replaceChildren(frag);
    $('empty' + K).hidden = t.rows.length > 0;
    $('shown' + K).textContent = t.rows.length ? `${vis.length} of ${t.rows.length}` : '';
    $('sel' + K).textContent = t.sel.size ? `${t.sel.size} selected` : '';
    $('remove' + K).disabled = !t.sel.size || running;

    const n = { ok: 0, failed: 0, pending: 0 };
    t.rows.forEach((r) => n[status(r)]++);
    const chips = $('chips' + K);
    chips.innerHTML = t.rows.length && (n.ok || n.failed)
      ? ['ok', 'failed', 'pending'].filter((s) => n[s]).map((s) =>
        `<span class="chip${t.chip === s ? ' on' : ''}" data-s="${s}">${s} ${n[s]}</span>`).join('')
      : '';
    refreshHeader();
  };

  // A result arriving mid-run touches one row, not the whole table.
  t.update = (i) => {
    const tr = t.trs[i];
    if (tr) tr.innerHTML = rowHtml(t.rows[i]);
  };

  t.set = (rows, file) => {
    t.rows = rows; t.file = file || ''; t.sel.clear(); t.chip = null; t.anchor = null;
    t.render();
  };

  t.removeSelected = () => {
    if (running || !t.sel.size) return;
    t.rows = t.rows.filter((_, i) => !t.sel.has(i));
    t.sel.clear(); t.anchor = null;
    t.render();
    tablesChanged();
  };

  t.clearResults = () => { t.rows.forEach((r) => { r.res = null; }); t.chip = null; };

  body.addEventListener('click', (e) => {
    const tr = e.target.closest('tr');
    if (!tr) return;
    const i = +tr.dataset.i;
    if (e.shiftKey && t.anchor !== null) {
      const vis = t.visible();
      const a = vis.indexOf(t.anchor), b = vis.indexOf(i);
      if (a >= 0 && b >= 0) for (let k = Math.min(a, b); k <= Math.max(a, b); k++) t.sel.add(vis[k]);
    } else if (e.ctrlKey || e.metaKey) {
      if (t.sel.has(i)) t.sel.delete(i); else t.sel.add(i);
      t.anchor = i;
    } else {
      const only = t.sel.size === 1 && t.sel.has(i);
      t.sel.clear();
      if (!only) t.sel.add(i);
      t.anchor = i;
    }
    t.render();
  });
  $('chips' + K).addEventListener('click', (e) => {
    const s = e.target.dataset.s;
    if (!s) return;
    t.chip = t.chip === s ? null : s;
    t.render();
  });
  $('filter' + K).addEventListener('input', (e) => { t.filter = e.target.value; t.render(); });
  $('remove' + K).onclick = t.removeSelected;
  return t;
}

const resultCell = (r) => (!r.res ? ['', ''] : r.res.error ? ['ERROR - ' + r.res.error, 'err'] : [r.res.ok || 'OK', 'ok']);

const naming = makeTable('naming', (r) => [
  [r.id], [r.name], [r.description], [r.location], [r.latitude], [r.longitude],
  [r.width24, 'num'], [r.channel24, 'num'], [r.width50, 'num'], [r.channel50, 'num'], resultCell(r),
]);
const overrides = makeTable('overrides', (r) => [
  [r.id], [r.res ? r.res.ok : '', 'ok'], [r.res ? r.res.error : '', 'err'],
]);
const tables = { naming, overrides };

function tablesChanged() {
  const n = naming.rows.length, m = overrides.rows.length;
  $('nRows').textContent = n ? `(${n})` : '';
  $('nMacs').textContent = m ? `(${m})` : '';
  if (!n) $('namingHint').textContent = 'No file loaded.';
  if (!m) $('macsHint').textContent = 'No file loaded.';
  persist();
  refreshHeader();
}

// The original removed selected rows with the Delete key; so does this.
document.addEventListener('keydown', (e) => {
  if (e.key !== 'Delete' && e.key !== 'Backspace') return;
  if (/^(INPUT|TEXTAREA|SELECT)$/.test(e.target.tagName)) return;
  const t = tables[view];
  if (t && t.sel.size) { e.preventDefault(); t.removeSelected(); }
});

/* ---------- files ---------- */

function loadFile(input, kind, done) {
  const f = input.files[0];
  input.value = '';
  if (!f) return;
  const r = new FileReader();
  r.onload = async () => {
    try {
      done(await post('/api/parse', { kind, text: r.result }), f.name);
    } catch (e) { toast(e.message); }
  };
  r.onerror = () => toast('Error accessing file.');
  r.readAsText(f);
}

const plural = (n, w) => `${n} ${w}${n === 1 ? '' : 's'}`;
const skippedNote = (n) => (n ? `; ${plural(n, 'line')} did not start with a MAC address or serial number and ${n === 1 ? 'was' : 'were'} skipped` : '');

$('pickNaming').onclick = () => $('fileNaming').click();
$('fileNaming').onchange = (e) => loadFile(e.target, 'naming', (res, name) => {
  naming.set((res.rows || []).map((r) => ({ ...r, res: null })), name);
  tablesChanged();
  $('namingHint').textContent = `${name}: ${plural(naming.rows.length, 'AP')}${skippedNote(res.skipped)}.`;
  log(`Loaded ${name}: ${plural(naming.rows.length, 'AP')}${skippedNote(res.skipped)}`);
});
$('clearNaming').onclick = () => { if (!running) { naming.set([], ''); tablesChanged(); } };

$('pickMacs').onclick = () => $('fileMacs').click();
$('fileMacs').onchange = (e) => loadFile(e.target, 'ids', (res, name) => {
  overrides.set((res.ids || []).map((id) => ({ id, res: null })), name);
  tablesChanged();
  $('macsHint').textContent = `${name}: ${plural(overrides.rows.length, 'AP')}${skippedNote(res.skipped)}.`;
  log(`Loaded ${name}: ${plural(overrides.rows.length, 'AP')}${skippedNote(res.skipped)}`);
});
$('clearMacs').onclick = () => { if (!running) { overrides.set([], ''); tablesChanged(); } };

/* ---------- settings ---------- */

const kind = () => (document.querySelector('input[name=kind]:checked') || {}).value || 'sz';

// Everything except the secrets is remembered between sessions.
const SAVE = ['host', 'port', 'user', 'r1Region', 'r1Tenant', 'r1ClientId', 'concurrency', 'apiVersion',
  'verifyTls', 'debug', 'radio24', 'radio50', 'ipMode'];

function persist() {
  const s = { kind: kind(), view, zone: $('zone').value, optPicks };
  for (const k of SAVE) {
    const el = $(k);
    s[k] = el.type === 'checkbox' ? el.checked : el.value;
  }
  s.naming = { file: naming.file, rows: naming.rows.map(({ res, ...r }) => r) };
  s.ids = { file: overrides.file, rows: overrides.rows.map((r) => r.id) };
  try { localStorage.setItem('dogtag', JSON.stringify(s)); } catch (e) { /* private mode, or a file too big to keep */ }
  // The secrets must not survive on disk. They are held for this browser
  // session only, so a reload does not empty the fields unnoticed.
  try {
    sessionStorage.setItem('dt-pass', $('pass').value);
    sessionStorage.setItem('dt-secret', $('r1Secret').value);
  } catch (e) { /* ignore */ }
}

let savedZone = '';
let optPicks = {}; // per controller kind: which overrides are ticked

function restore() {
  let s = {};
  try { s = JSON.parse(localStorage.getItem('dogtag') || '{}'); } catch (e) { /* start clean */ }
  for (const k of SAVE) {
    const el = $(k), v = s[k];
    if (v === undefined || v === '') continue;
    if (el.type === 'checkbox') el.checked = !!v; else el.value = v;
  }
  // What the process was started with wins over what was remembered: a flag
  // on the command line is the more recent instruction.
  for (const k of ['host', 'user', 'r1Tenant', 'r1ClientId']) if (defaults[k] && defaults.explicit) $(k).value = defaults[k];
  const k = defaults.r1Flag ? 'r1' : (s.kind || defaults.kind || 'sz');
  (k === 'r1' ? $('kindR1') : $('kindSz')).checked = true;
  savedZone = s.zone || '';
  optPicks = s.optPicks || {};
  // Lists saved by 6.0.0 keyed each AP as "mac"; they are read as they were.
  if (s.naming && s.naming.rows) naming.set(s.naming.rows.map(({ mac, ...r }) => ({ id: mac, ...r, res: null })), s.naming.file);
  const ids = s.ids || s.macs;
  if (ids && ids.rows) overrides.set(ids.rows.map((id) => ({ id, res: null })), ids.file);
  if (naming.rows.length) $('namingHint').textContent = `${naming.file || 'Restored'}: ${plural(naming.rows.length, 'AP')}.`;
  if (overrides.rows.length) $('macsHint').textContent = `${overrides.file || 'Restored'}: ${plural(overrides.rows.length, 'AP')}.`;
  try {
    const p = sessionStorage.getItem('dt-pass'); if (p) $('pass').value = p;
    const q = sessionStorage.getItem('dt-secret'); if (q) $('r1Secret').value = q;
  } catch (e) { /* ignore */ }
  if (['naming', 'overrides', 'channels', 'scaling', 'about'].includes(s.view)) view = s.view;
}

async function loadDefaults() {
  defaults = await (await fetch('/api/defaults')).json();
  $('ver').textContent = defaults.version || '';
  $('aboutVer').textContent = defaults.version || '';
  for (const k of ['host', 'port', 'user', 'r1Region', 'r1Tenant', 'r1ClientId', 'concurrency', 'apiVersion', 'verifyTls', 'debug']) {
    const el = $(k), v = defaults[k];
    if (v === undefined || v === null || v === '') continue;
    if (el.type === 'checkbox') el.checked = !!v; else el.value = v;
  }
  defaults.explicit = !!defaults.host || !!defaults.r1Tenant;
  defaults.r1Flag = defaults.kind === 'r1';
  restore();
  if (defaults.hasPass) $('pass').placeholder = 'given on the command line';
  if (defaults.hasR1Secret) $('r1Secret').placeholder = 'given on the command line';
  if (defaults.debug) $('debug').checked = true;
  kindChanged();
  showView(view);
  tablesChanged();
  applyState(await (await fetch('/api/state')).json());
  listen();
  checkForUpdate();
}

// Asked for after the page has rendered, never before: a blocked or slow
// GitHub must not delay the console.
async function checkForUpdate() {
  try {
    const u = await (await fetch('/api/update')).json();
    if (!u || !u.available) return;
    const el = $('update');
    el.textContent = `v${u.latest} available`;
    el.href = u.url || 'https://github.com/andreacoppini/dogtag/releases/latest';
    el.title = `You are running ${u.current}. Click to open the release page.`;
    el.hidden = false;
  } catch (e) { /* offline or blocked — say nothing */ }
}

for (const k of SAVE) $(k).addEventListener('change', persist);
$('showPass').onchange = () => { $('pass').type = $('showPass').checked ? 'text' : 'password'; };
$('pass').addEventListener('input', persist);
$('r1Secret').addEventListener('input', persist);
document.querySelectorAll('section > h2').forEach((h) => { h.onclick = () => h.parentElement.classList.toggle('collapsed'); });

/* ---------- controller ---------- */

function kindChanged() {
  const r1 = kind() === 'r1';
  $('szOpts').hidden = r1;
  $('r1Opts').hidden = !r1;
  if (!conn.connected) wording(r1);
  persist();
}
$('kindSz').onchange = kindChanged;
$('kindR1').onchange = kindChanged;

// SmartZone groups APs into zones and RUCKUS One into venues; the page says
// whichever the connected controller says.
function wording(r1) {
  $('zoneTitle').textContent = r1 ? 'Venue' : 'Zone';
  document.querySelectorAll('.zoneWord').forEach((el) => { el.textContent = r1 ? 'venue' : 'zone'; });
  $('namingR1Hint').hidden = !r1;
  document.querySelectorAll('.idHead').forEach((el) => { el.textContent = r1 ? 'Serial / MAC' : 'MAC Address'; });
  $('idHint').textContent = r1
    ? 'One AP per row, serial number first (a MAC address works too, once the AP has come online)'
    : 'One AP per row, MAC address first';
  $('filterNaming').placeholder = r1 ? 'Filter by serial, name, result…' : 'Filter by MAC, name, result…';
  $('emptyOverrides').textContent = r1
    ? 'Please load a CSV file — the first column of each row is the AP\'s serial number.'
    : 'Please load a CSV file — the first column of each row is the AP\'s MAC address.';
}

function applyState(st) {
  const was = conn.connected;
  conn = st || { connected: false };
  const on = !!conn.connected;
  $('connDot').classList.toggle('on', on);
  $('connState').textContent = on ? 'Connected' : 'Not connected';
  $('connInfo').textContent = on ? conn.describe : '';
  $('connect').textContent = on ? 'Disconnect' : 'Connect';
  $('connect').classList.toggle('primary', !on);
  for (const id of ['kindSz', 'kindR1', 'host', 'port', 'user', 'pass', 'r1Region', 'r1Tenant', 'r1ClientId', 'r1Secret']) $(id).disabled = on;
  if (on) {
    (conn.kind === 'r1' ? $('kindR1') : $('kindSz')).checked = true;
    $('szOpts').hidden = conn.kind === 'r1';
    $('r1Opts').hidden = conn.kind !== 'r1';
    wording(conn.kind === 'r1');
  }

  const z = $('zone');
  if (on) {
    const keep = z.value || savedZone;
    z.innerHTML = (conn.zones || []).map((x) => `<option value="${esc(x.id)}">${esc(x.name)}</option>`).join('');
    const byName = (n) => (conn.zones || []).find((x) => x.name === n);
    // As the original did: the staging zone if there is one, else the default.
    const pick = (conn.zones || []).find((x) => x.id === keep) || byName('Staging Zone') || byName('Default Zone') || (conn.zones || [])[0];
    if (pick) z.value = pick.id;
  } else if (was) {
    z.innerHTML = '';
  }
  z.disabled = !on || !(conn.zones || []).length;
  renderOptions();
  if (st && st.running !== undefined) setRunning(!!st.running);
  refreshHeader();
}

$('connect').onclick = async () => {
  const btn = $('connect');
  btn.disabled = true;
  try {
    if (conn.connected) {
      applyState(await post('/api/disconnect'));
      return;
    }
    $('connState').textContent = 'Connecting…';
    persist();
    applyState(await post('/api/connect', {
      kind: kind(), host: $('host').value.trim(), port: $('port').value.trim(), user: $('user').value, pass: $('pass').value,
      verifyTls: $('verifyTls').checked, apiVersion: $('apiVersion').value.trim(), debug: $('debug').checked,
      r1Region: $('r1Region').value, r1Tenant: $('r1Tenant').value.trim(), r1ClientId: $('r1ClientId').value.trim(), r1Secret: $('r1Secret').value,
    }));
  } catch (e) {
    applyState({ connected: false });
    toast(e.message);
  } finally { btn.disabled = false; }
};
$('zone').onchange = () => { savedZone = $('zone').value; persist(); refreshHeader(); };
const zoneName = () => { const o = $('zone').selectedOptions[0]; return o ? o.textContent : ''; };

/* ---------- override options ---------- */

function renderOptions() {
  const opts = conn.connected ? (conn.overrides || []) : [];
  const picks = optPicks[conn.kind];
  $('optList').innerHTML = opts.map((o) => {
    const on = picks ? picks.includes(o.id) : o.default;
    const label = o.label === o.id ? esc(o.id) : `${esc(o.label)}`;
    return `<label title="${esc(o.id)}"><input type="checkbox" value="${esc(o.id)}"${on ? ' checked' : ''}> ${label}</label>`;
  }).join('');
  $('optHint').textContent = !conn.connected ? 'Connect to see what this controller offers.'
    : conn.kind === 'r1'
      ? 'Each ticked setting is switched back to "use venue settings" on every AP in the file. The unticked ones exist only on some AP models.'
      : 'Each ticked override is deleted from every AP in the file, so the AP takes that setting from its group or zone again. "picture" deletes the AP\'s photo.';
}
const pickedOptions = () => [...$('optList').querySelectorAll('input:checked')].map((i) => i.value);
function setOptions(fn) {
  $('optList').querySelectorAll('input').forEach((i) => {
    const o = (conn.overrides || []).find((x) => x.id === i.value);
    i.checked = fn(o);
  });
  optPicks[conn.kind] = pickedOptions();
  persist();
}
$('optList').addEventListener('change', () => { optPicks[conn.kind] = pickedOptions(); persist(); });
$('optAll').onclick = () => setOptions(() => true);
$('optNone').onclick = () => setOptions(() => false);
$('optDefault').onclick = () => setOptions((o) => !!(o && o.default));

/* ---------- views ---------- */

function showView(v) {
  view = v;
  document.querySelectorAll('#tools .tab').forEach((t) => t.classList.toggle('on', t.dataset.v === v));
  for (const k of ['naming', 'overrides', 'channels', 'scaling', 'about']) {
    $('v' + k[0].toUpperCase() + k.slice(1)).hidden = k !== v;
  }
  $('namingSection').hidden = v !== 'naming';
  $('overridesSection').hidden = v !== 'overrides';
  $('channelsSection').hidden = v !== 'channels';
  if (tables[v]) tables[v].render();
  refreshHeader();
}
$('tools').addEventListener('click', (e) => {
  const v = e.target.dataset.v;
  if (v) { showView(v); persist(); }
});

function refreshHeader() {
  const run = $('run');
  run.textContent = view === 'channels' ? 'Refresh' : 'Run';
  const t = tables[view];
  const ready = view === 'channels' ? !!$('zone').value : !!(t && t.rows.length);
  run.disabled = running || !conn.connected || !ready;
  run.title = !conn.connected ? 'Connect to a controller first'
    : view === 'channels' ? 'Count the connected APs on each channel'
    : t ? (t.rows.length ? '' : 'Load a CSV file first') : 'Nothing to run on this tab';
  $('exportZone').disabled = running || !conn.connected || !$('zone').value;
  $('expCsv').disabled = !(t ? t.rows.length : view === 'channels' && lastChannels);
  if (!running && t) counters(t.rows.length, t.rows.filter((r) => r.res && !r.res.error).length, t.rows.filter((r) => r.res && r.res.error).length);
}

function counters(total, done, failed) {
  $('cTotal').textContent = total;
  $('cDone').textContent = done;
  $('cFail').textContent = failed;
}

function setRunning(on) {
  running = on;
  $('stop').disabled = !on;
  $('stop').classList.toggle('armed', on);
  document.querySelector('aside').classList.toggle('locked', on);
  $('lockNote').hidden = !on;
  document.querySelectorAll('aside input, aside select, aside button').forEach((el) => {
    if (on) { el.dataset.was = el.disabled ? '1' : ''; el.disabled = true; } else if (el.dataset.was !== undefined) { el.disabled = el.dataset.was === '1'; delete el.dataset.was; }
  });
  if (!on) { applyLocks(); $('bar').style.width = '0'; }
  if (tables[view]) tables[view].render(); else refreshHeader();
}
// After a run the disabled states are recomputed rather than trusted, since
// the connection may have gone while the panel was locked.
function applyLocks() {
  const on = !!conn.connected;
  for (const id of ['kindSz', 'kindR1', 'host', 'port', 'user', 'pass', 'r1Region', 'r1Tenant', 'r1ClientId', 'r1Secret']) $(id).disabled = on;
  $('zone').disabled = !on || !(conn.zones || []).length;
}

/* ---------- running ---------- */

function namingSummary() {
  const n = (f) => naming.rows.filter(f).length;
  const out = [];
  const add = (count, text) => { if (count) out.push(`${text} on ${plural(count, 'AP')}`); };
  add(n((r) => r.name), 'set the name');
  add(n((r) => r.description), 'set the description');
  if (conn.kind !== 'r1') add(n((r) => r.location), 'set the location');
  add(n((r) => r.latitude || r.longitude), 'set the GPS position');
  add(n((r) => r.channel24 || r.width24), 'pin the 2.4 GHz channel');
  add(n((r) => r.channel50 || r.width50), 'pin the 5 GHz channel');
  return out;
}

function overridesSummary() {
  const out = [];
  const picked = pickedOptions();
  if (picked.length) out.push(`hand back ${plural(picked.length, 'setting')}: ${picked.join(', ')}`);
  for (const [id, band] of [['radio24', '2.4 GHz'], ['radio50', '5 GHz']]) {
    if ($(id).value !== 'skip') out.push(`switch the ${band} radio ${$(id).value}`);
  }
  if ($('ipMode').value === 'dynamic') out.push('force the IP address to dynamic (DHCP)');
  return out;
}

function confirmRun(count, items) {
  return new Promise((resolve) => {
    $('cN').textContent = count;
    $('cWhere').textContent = conn.describe || 'the controller';
    $('cList').innerHTML = items.map((i) => `<li>${esc(i)}</li>`).join('');
    $('confirm').hidden = false;
    const close = (v) => { $('confirm').hidden = true; resolve(v); };
    $('cCancel').onclick = () => close(false);
    $('cGo').onclick = () => close(true);
  });
}

async function start(body) {
  try {
    const t = tables[body.job];
    if (t) { t.clearResults(); t.render(); }
    // Marked as running before the request goes out: a short run can finish,
    // and say so over the event stream, before this request has returned.
    const concurrency = +$('concurrency').value || 0;
    setRunning(true);
    $('phase').textContent = 'starting';
    await post('/api/run', { ...body, concurrency });
  } catch (e) {
    setRunning(false);
    $('phase').textContent = 'idle';
    wantDownload = false;
    toast(e.message);
  }
}

$('run').onclick = async () => {
  if (view === 'channels') {
    return start({ job: 'channels', zoneId: $('zone').value, zoneName: zoneName() });
  }
  if (view === 'naming') {
    const items = namingSummary();
    if (!items.length) return toast('Nothing to change: every column after the first is empty.');
    if (!await confirmRun(naming.rows.length, items)) return;
    return start({ job: 'naming', rows: naming.rows.map(({ res, ...r }) => r) });
  }
  if (view === 'overrides') {
    const items = overridesSummary();
    if (!items.length) return toast('Nothing to do: tick at least one setting, or choose a radio or IP change.');
    if (!await confirmRun(overrides.rows.length, items)) return;
    return start({
      job: 'overrides', ids: overrides.rows.map((r) => r.id), overrides: pickedOptions(),
      radio24: $('radio24').value, radio50: $('radio50').value, ipDynamic: $('ipMode').value === 'dynamic',
    });
  }
};

$('exportZone').onclick = () => {
  wantDownload = true;
  start({ job: 'export', zoneId: $('zone').value, zoneName: zoneName() });
};
$('stop').onclick = () => post('/api/stop').catch((e) => toast(e.message));

/* ---------- events ---------- */

function listen() {
  const es = new EventSource('/api/events');
  es.onmessage = (m) => {
    let e;
    try { e = JSON.parse(m.data); } catch (x) { return; }
    onEvent(e);
  };
}

function progress(e) {
  if (!e.total) return;
  counters(e.total, e.done, e.failed);
  $('bar').style.width = `${Math.min(100, 100 * (e.done + e.failed) / e.total)}%`;
}

function onEvent(e) {
  switch (e.kind) {
    case 'log':
      log(e.message);
      break;
    case 'phase':
      $('phase').textContent = e.phase;
      if (e.total) progress(e); else $('bar').style.width = '0';
      break;
    case 'row': {
      const t = tables[e.job];
      const r = t && t.rows[e.row.index];
      // Matched on the AP's ID as well as the position, so results replayed after
      // a reload cannot land on a list that has since been edited.
      if (r && r.id === e.row.id) {
        r.res = { ok: e.row.ok, error: e.row.error };
        t.update(e.row.index);
      }
      if (e.row.error) log(`${e.row.id}  ERROR - ${e.row.error}`);
      progress(e);
      break;
    }
    case 'channels':
      renderChannels(e.channels);
      break;
    case 'done':
      log(e.message);
      $('phase').textContent = e.phase;
      setRunning(false);
      if (e.job === 'export' && wantDownload && e.phase === 'done') window.location = '/api/download';
      wantDownload = false;
      break;
    case 'error':
      log('ERROR - ' + e.message);
      toast(e.message);
      $('phase').textContent = 'failed';
      setRunning(false);
      wantDownload = false;
      break;
    case 'session':
      applyState({ connected: false });
      toast(e.message);
      break;
  }
}

/* ---------- channels ---------- */

let lastChannels = null;

function renderChannels(c) {
  lastChannels = c;
  $('chanInfo').textContent = `${c.zone}: ${c.connected} of ${plural(c.total, 'AP')} connected · counted ${c.ended}`;
  const band = (title, rows) => {
    if (!rows || !rows.length) return title === '6 GHz' ? '' : `<div class="band"><h3>${title}<span>no APs</span></h3></div>`;
    const max = Math.max(...rows.map((r) => r.aps));
    const total = rows.reduce((a, r) => a + r.aps, 0);
    return `<div class="band"><h3>${title}<span>${plural(total, 'AP')}</span></h3><table>
      <thead><tr><th>Channel</th><th>APs</th><th></th></tr></thead><tbody>${rows.map((r) =>
    `<tr><td class="num">${r.channel}</td><td class="num">${r.aps}</td><td class="barcell"><div class="bar" style="width:${100 * r.aps / max}%"></div></td></tr>`).join('')}
      </tbody></table></div>`;
  };
  $('bands').innerHTML = band('2.4 GHz', c.band24) + band('5 GHz', c.band50) + band('6 GHz', c.band6);
  refreshHeader();
}

/* ---------- CSV of what is on screen ---------- */

const csvCell = (v) => { const s = String(v ?? ''); return /[",\r\n]/.test(s) ? `"${s.replace(/"/g, '""')}"` : s; };

function save(name, records) {
  const blob = new Blob(['\ufeff' + records.map((r) => r.map(csvCell).join(',')).join('\r\n') + '\r\n'], { type: 'text/csv' });
  const a = document.createElement('a');
  a.href = URL.createObjectURL(blob);
  a.download = name;
  a.click();
  setTimeout(() => URL.revokeObjectURL(a.href), 1000);
}

$('expCsv').onclick = () => {
  const res = (r) => (!r.res ? '' : r.res.error ? 'ERROR - ' + r.res.error : r.res.ok || 'OK');
  if (view === 'naming') {
    save('dogtag-naming-results.csv', [
      [conn.kind === 'r1' ? 'Serial' : 'MAC', 'Name', 'Description', 'Location', 'GPSLatitude', 'GPSLongitude', '24Width', '24Channel', '50Width', '50Channel', 'Result'],
      ...naming.rows.map((r) => [r.id, r.name, r.description, r.location, r.latitude, r.longitude, r.width24, r.channel24, r.width50, r.channel50, res(r)]),
    ]);
  } else if (view === 'overrides') {
    save('dogtag-overrides-results.csv', [[conn.kind === 'r1' ? 'Serial' : 'MAC', 'OK', 'Error'],
      ...overrides.rows.map((r) => [r.id, r.res ? r.res.ok : '', r.res ? r.res.error : ''])]);
  } else if (view === 'channels' && lastChannels) {
    const c = lastChannels;
    const rows = (b, list) => (list || []).map((r) => [b, r.channel, r.aps]);
    save('dogtag-channels.csv', [['Band', 'Channel', 'APs'], ...rows('2.4GHz', c.band24), ...rows('5GHz', c.band50), ...rows('6GHz', c.band6)]);
  }
};

/* ---------- scaling ---------- */

// SmartZone capacity: how many switches a cluster can carry alongside a given
// number of APs. The figures are the original DogTag's.
function scalCalc() {
  const ver = +$('scalVer').value;
  const nodes = +$('scalNodes').value;
  const high = document.querySelector('input[name=scalDash]:checked').value === 'high';
  const mult = nodes <= 2 ? 1 : nodes - 1; // one node is always held in reserve
  const maxAP = (high ? 10000 : 1024) * mult;
  const ratio = ver === 5.0 ? 20 : high ? 8 : 5;
  const slider = $('scalAP');
  const wasMax = +slider.value >= +slider.max;
  slider.max = maxAP;
  if (wasMax || +slider.value > maxAP) slider.value = maxAP;
  $('scalNodesVal').textContent = nodes;

  if (ver < 5) {
    slider.value = maxAP;
    slider.disabled = true;
    $('scalAPVal').textContent = maxAP;
    $('scalSwVal').textContent = 'N/A';
    $('scalHint').textContent = 'Before 5.0 a SmartZone does not manage switches.';
    return;
  }
  slider.disabled = false;
  const aps = +slider.value;
  $('scalAPVal').textContent = aps;
  $('scalSwVal').textContent = Math.floor(maxAP / ratio - aps / ratio);
  $('scalHint').textContent = `Up to ${maxAP} APs with no switches, or ${Math.floor(maxAP / ratio)} switches with no APs: one switch takes the room of ${ratio} APs.`;
}
for (const id of ['scalVer', 'scalNodes', 'scalAP']) $(id).addEventListener('input', scalCalc);
document.querySelectorAll('input[name=scalDash]').forEach((r) => r.addEventListener('change', scalCalc));
scalCalc();

/* ---------- drawer ---------- */

$('grow').onclick = () => {
  const big = $('drawer').classList.toggle('big');
  $('grow').textContent = big ? 'Shrink' : 'Expand';
};
$('clearLog').onclick = () => { $('pLog').replaceChildren(); };

loadDefaults().catch((e) => toast('Could not reach the console: ' + e.message));
